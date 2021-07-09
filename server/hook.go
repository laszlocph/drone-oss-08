// Copyright 2018 Drone.IO Inc.
// Copyright 2021 Informatyka Boguslawski sp. z o.o. sp.k., http://www.ib.pl/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// This file has been modified by Informatyka Boguslawski sp. z o.o. sp.k.

package server

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
	"github.com/woodpecker-ci/woodpecker/model"
	"github.com/woodpecker-ci/woodpecker/remote"
	"github.com/woodpecker-ci/woodpecker/shared/httputil"
	"github.com/woodpecker-ci/woodpecker/shared/token"
	"github.com/woodpecker-ci/woodpecker/store"

	"github.com/woodpecker-ci/woodpecker/cncd/pipeline/pipeline/frontend/yaml"
	"github.com/woodpecker-ci/woodpecker/cncd/pipeline/pipeline/rpc"
	"github.com/woodpecker-ci/woodpecker/cncd/pubsub"
	"github.com/woodpecker-ci/woodpecker/cncd/queue"
)

var skipRe = regexp.MustCompile(`\[(?i:ci *skip|skip *ci)\]`)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetQueueInfo(c *gin.Context) {
	c.IndentedJSON(200,
		Config.Services.Queue.Info(c),
	)
}

func PauseQueue(c *gin.Context) {
	Config.Services.Queue.Pause()
	c.Status(http.StatusOK)
}

func ResumeQueue(c *gin.Context) {
	Config.Services.Queue.Resume()
	c.Status(http.StatusOK)
}

func BlockTilQueueHasRunningItem(c *gin.Context) {
	for {
		info := Config.Services.Queue.Info(c)
		if info.Stats.Running == 0 {
			break
		}
	}
	c.Status(http.StatusOK)
}

func PostHook(c *gin.Context) {
	remote_ := remote.FromContext(c)

	tmprepo, build, err := remote_.Hook(c.Request)
	if err != nil {
		logrus.Errorf("failure to parse hook. %s", err)
		c.AbortWithError(400, err)
		return
	}
	if build == nil {
		c.Writer.WriteHeader(200)
		return
	}
	if tmprepo == nil {
		logrus.Errorf("failure to ascertain repo from hook.")
		c.Writer.WriteHeader(400)
		return
	}

	// skip the build if any case-insensitive combination of the words "skip" and "ci"
	// wrapped in square brackets appear in the commit message
	skipMatch := skipRe.FindString(build.Message)
	if len(skipMatch) > 0 {
		logrus.Infof("ignoring hook. %s found in %s", skipMatch, build.Commit)
		c.Writer.WriteHeader(204)
		return
	}

	repo, err := store.GetRepoOwnerName(c, tmprepo.Owner, tmprepo.Name)
	if err != nil {
		logrus.Errorf("failure to find repo %s/%s from hook. %s", tmprepo.Owner, tmprepo.Name, err)
		c.AbortWithError(404, err)
		return
	}
	if !repo.IsActive {
		logrus.Errorf("ignoring hook. %s/%s is inactive.", tmprepo.Owner, tmprepo.Name)
		c.AbortWithError(204, err)
		return
	}

	// get the token and verify the hook is authorized
	parsed, err := token.ParseRequest(c.Request, func(t *token.Token) (string, error) {
		return repo.Hash, nil
	})
	if err != nil {
		logrus.Errorf("failure to parse token from hook for %s. %s", repo.FullName, err)
		c.AbortWithError(400, err)
		return
	}
	if parsed.Text != repo.FullName {
		logrus.Errorf("failure to verify token from hook. Expected %s, got %s", repo.FullName, parsed.Text)
		c.AbortWithStatus(403)
		return
	}

	if repo.UserID == 0 {
		logrus.Warnf("ignoring hook. repo %s has no owner.", repo.FullName)
		c.Writer.WriteHeader(204)
		return
	}
	var skipped = true
	if (build.Event == model.EventPush && repo.AllowPush) ||
		(build.Event == model.EventPull && repo.AllowPull) ||
		(build.Event == model.EventDeploy && repo.AllowDeploy) ||
		(build.Event == model.EventTag && repo.AllowTag) {
		skipped = false
	}

	if skipped {
		logrus.Infof("ignoring hook. repo %s is disabled for %s events.", repo.FullName, build.Event)
		c.Writer.WriteHeader(204)
		return
	}

	user, err := store.GetUser(c, repo.UserID)
	if err != nil {
		logrus.Errorf("failure to find repo owner %s. %s", repo.FullName, err)
		c.AbortWithError(500, err)
		return
	}

	// if the remote has a refresh token, the current access token
	// may be stale. Therefore, we should refresh prior to dispatching
	// the build.
	if refresher, ok := remote_.(remote.Refresher); ok {
		ok, _ := refresher.Refresh(user)
		if ok {
			store.UpdateUser(c, user)
		}
	}

	// fetch the build file from the remote
	configFetcher := &configFetcher{remote_: remote_, user: user, repo: repo, build: build}
	remoteYamlConfigs, err := configFetcher.Fetch()
	if err != nil {
		logrus.Errorf("error: %s: cannot find %s in %s: %s", repo.FullName, repo.Config, build.Ref, err)
		c.AbortWithError(404, err)
		return
	}

	filtered, err := branchFiltered(build, remoteYamlConfigs)
	if err != nil {
		logrus.Errorf("failure to parse yaml from hook for %s. %s", repo.FullName, err)
		c.AbortWithError(400, err)
	}
	if filtered {
		c.String(200, "Branch does not match restrictions defined in yaml")
		return
	}

	if zeroSteps(build, remoteYamlConfigs) {
		c.String(200, "Step conditions yield zero runnable steps")
		return
	}

	// update some build fields
	build.RepoID = repo.ID
	build.Verified = true
	build.Status = model.StatusPending

	if repo.IsGated && build.Sender != user.Login {
		build.Status = model.StatusBlocked
	}

	err = store.CreateBuild(c, build, build.Procs...)
	if err != nil {
		logrus.Errorf("failure to save commit for %s. %s", repo.FullName, err)
		c.AbortWithError(500, err)
		return
	}

	// persist the build config for historical correctness, restarts, etc
	for _, remoteYamlConfig := range remoteYamlConfigs {
		_, err := findOrPersistPipelineConfig(build, remoteYamlConfig)
		if err != nil {
			logrus.Errorf("failure to find or persist build config for %s. %s", repo.FullName, err)
			c.AbortWithError(500, err)
			return
		}
	}

	c.JSON(200, build)

	if build.Status == model.StatusBlocked {
		return
	}

	netrc, err := remote_.Netrc(user, repo)
	if err != nil {
		c.String(500, "Failed to generate netrc file. %s", err)
		return
	}

	envs := map[string]string{}
	if Config.Services.Environ != nil {
		globals, _ := Config.Services.Environ.EnvironList(repo)
		for _, global := range globals {
			envs[global.Name] = global.Value
		}
	}

	secs, err := Config.Services.Secrets.SecretListBuild(repo, build)
	if err != nil {
		logrus.Debugf("Error getting secrets for %s#%d. %s", repo.FullName, build.Number, err)
	}

	globalsecs, err := Config.Services.GlobalSecrets.GlobalSecretListBuild(build)
	if err != nil {
		logrus.Debugf("Error getting global secrets for #%d. %s", build.Number, err)
	}

	regs, err := Config.Services.Registries.RegistryList(repo)
	if err != nil {
		logrus.Debugf("Error getting registry credentials for %s#%d. %s", repo.FullName, build.Number, err)
	}

	// get the previous build so that we can send status change notifications
	last, _ := store.GetBuildLastBefore(c, repo, build.Branch, build.ID)

	b := procBuilder{
		Repo:       repo,
		Curr:       build,
		Last:       last,
		Netrc:      netrc,
		Secs:       secs,
		GlobalSecs: globalsecs,
		Regs:       regs,
		Envs:       envs,
		Link:       httputil.GetURL(c.Request),
		Yamls:      remoteYamlConfigs,
	}
	buildItems, err := b.Build()
	if err != nil {
		if _, err = UpdateToStatusError(store.FromContext(c), *build, err); err != nil {
			logrus.Errorf("Error setting error status of build for %s#%d. %s", repo.FullName, build.Number, err)
		}
		return
	}
	build = setBuildStepsOnBuild(b.Curr, buildItems)

	err = store.FromContext(c).ProcCreate(build.Procs)
	if err != nil {
		logrus.Errorf("error persisting procs %s/%d: %s", repo.FullName, build.Number, err)
	}

	defer func() {
		for _, item := range buildItems {
			uri := fmt.Sprintf("%s/%s/%d", Config.Server.Host, repo.FullName, build.Number)
			if len(buildItems) > 1 {
				err = remote_.Status(user, repo, build, uri, item.Proc)
			} else {
				err = remote_.Status(user, repo, build, uri, nil)
			}
			if err != nil {
				logrus.Errorf("error setting commit status for %s/%d: %v", repo.FullName, build.Number, err)
			}
		}
	}()

	publishToTopic(c, build, repo, model.Enqueued)
	queueBuild(build, repo, buildItems)
}

func branchFiltered(build *model.Build, remoteYamlConfigs []*remote.FileMeta) (bool, error) {
	for _, remoteYamlConfig := range remoteYamlConfigs {
		parsedPipelineConfig, err := yaml.ParseString(string(remoteYamlConfig.Data))
		if err != nil {
			return false, err
		}

		if !parsedPipelineConfig.Branches.Match(build.Branch) && build.Event != model.EventTag && build.Event != model.EventDeploy {
		} else {
			return false, nil
		}
	}
	return true, nil
}

func zeroSteps(build *model.Build, remoteYamlConfigs []*remote.FileMeta) bool {
	b := procBuilder{
		Repo:       &model.Repo{},
		Curr:       build,
		Last:       &model.Build{},
		Netrc:      &model.Netrc{},
		Secs:       []*model.Secret{},
		GlobalSecs: []*model.GlobalSecret{},
		Regs:       []*model.Registry{},
		Link:       "",
		Yamls:      remoteYamlConfigs,
	}

	buildItems, err := b.Build()
	if err != nil {
		return false
	}
	if len(buildItems) == 0 {
		return true
	}

	return false
}

func findOrPersistPipelineConfig(build *model.Build, remoteYamlConfig *remote.FileMeta) (*model.Config, error) {
	sha := shasum(remoteYamlConfig.Data)
	conf, err := Config.Storage.Config.ConfigFindIdentical(build.RepoID, sha)
	if err != nil {
		conf = &model.Config{
			RepoID: build.RepoID,
			Data:   string(remoteYamlConfig.Data),
			Hash:   sha,
			Name:   sanitizePath(remoteYamlConfig.Name),
		}
		err = Config.Storage.Config.ConfigCreate(conf)
		if err != nil {
			// retry in case we receive two hooks at the same time
			conf, err = Config.Storage.Config.ConfigFindIdentical(build.RepoID, sha)
			if err != nil {
				return nil, err
			}
		}
	}

	buildConfig := &model.BuildConfig{
		ConfigID: conf.ID,
		BuildID:  build.ID,
	}
	err = Config.Storage.Config.BuildConfigCreate(buildConfig)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// publishes message to UI clients
func publishToTopic(c *gin.Context, build *model.Build, repo *model.Repo, event model.EventType) {
	message := pubsub.Message{
		Labels: map[string]string{
			"repo":    repo.FullName,
			"private": strconv.FormatBool(repo.IsPrivate),
		},
	}
	buildCopy := *build
	buildCopy.Procs = model.Tree(buildCopy.Procs)
	message.Data, _ = json.Marshal(model.Event{
		Type:  model.Enqueued,
		Repo:  *repo,
		Build: buildCopy,
	})
	Config.Services.Pubsub.Publish(c, "topic/events", message)
}

func queueBuild(build *model.Build, repo *model.Repo, buildItems []*buildItem) {
	var tasks []*queue.Task
	for _, item := range buildItems {
		if item.Proc.State == model.StatusSkipped {
			continue
		}
		task := new(queue.Task)
		task.ID = fmt.Sprint(item.Proc.ID)
		task.Labels = map[string]string{}
		for k, v := range item.Labels {
			task.Labels[k] = v
		}
		task.Labels["platform"] = item.Platform
		task.Labels["repo"] = repo.FullName
		task.Dependencies = taskIds(item.DependsOn, buildItems)
		task.RunOn = item.RunsOn
		task.DepStatus = make(map[string]string)

		task.Data, _ = json.Marshal(rpc.Pipeline{
			ID:      fmt.Sprint(item.Proc.ID),
			Config:  item.Config,
			Timeout: repo.Timeout,
		})

		Config.Services.Logs.Open(context.Background(), task.ID)
		tasks = append(tasks, task)
	}
	Config.Services.Queue.PushAtOnce(context.Background(), tasks)
}

func taskIds(dependsOn []string, buildItems []*buildItem) []string {
	taskIds := []string{}
	for _, dep := range dependsOn {
		for _, buildItem := range buildItems {
			if buildItem.Proc.Name == dep {
				taskIds = append(taskIds, fmt.Sprint(buildItem.Proc.ID))
			}
		}
	}
	return taskIds
}

func shasum(raw []byte) string {
	sum := sha256.Sum256(raw)
	return fmt.Sprintf("%x", sum)
}
