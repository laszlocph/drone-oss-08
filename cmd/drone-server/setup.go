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

package main

import (
	"fmt"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/woodpecker-ci/woodpecker/cncd/queue"
	"github.com/woodpecker-ci/woodpecker/model"
	"github.com/woodpecker-ci/woodpecker/plugins/environments"
	"github.com/woodpecker-ci/woodpecker/plugins/registry"
	"github.com/woodpecker-ci/woodpecker/plugins/secrets"
	"github.com/woodpecker-ci/woodpecker/remote"
	"github.com/woodpecker-ci/woodpecker/remote/bitbucket"
	"github.com/woodpecker-ci/woodpecker/remote/bitbucketserver"
	"github.com/woodpecker-ci/woodpecker/remote/coding"
	"github.com/woodpecker-ci/woodpecker/remote/gitea"
	"github.com/woodpecker-ci/woodpecker/remote/github"
	"github.com/woodpecker-ci/woodpecker/remote/gitlab"
	"github.com/woodpecker-ci/woodpecker/remote/gitlab3"
	"github.com/woodpecker-ci/woodpecker/remote/gogs"
	droneserver "github.com/woodpecker-ci/woodpecker/server"
	"github.com/woodpecker-ci/woodpecker/server/web"
	"github.com/woodpecker-ci/woodpecker/store"
	"github.com/woodpecker-ci/woodpecker/store/datastore"
	"golang.org/x/sync/errgroup"

	"github.com/urfave/cli"
)

func setupStore(c *cli.Context) store.Store {
	return datastore.New(
		c.String("driver"),
		c.String("datasource"),
	)
}

func setupQueue(c *cli.Context, s store.Store) queue.Queue {
	return model.WithTaskStore(queue.New(), s)
}

func setupSecretService(c *cli.Context, s store.Store) model.SecretService {
	return secrets.New(s)
}

func setupRegistryService(c *cli.Context, s store.Store) model.RegistryService {
	if c.String("docker-config") != "" {
		return registry.Combined(
			registry.New(s),
			registry.Filesystem(c.String("docker-config")),
		)
	} else {
		return registry.New(s)
	}
}

func setupEnvironService(c *cli.Context, s store.Store) model.EnvironService {
	return environments.Filesystem(c.StringSlice("environment"))
}

// helper function to setup the remote from the CLI arguments.
func SetupRemote(c *cli.Context) (remote.Remote, error) {
	switch {
	case c.Bool("github"):
		return setupGithub(c)
	case c.Bool("gitlab"):
		return setupGitlab(c)
	case c.Bool("bitbucket"):
		return setupBitbucket(c)
	case c.Bool("stash"):
		return setupStash(c)
	case c.Bool("gogs"):
		return setupGogs(c)
	case c.Bool("gitea"):
		return setupGitea(c)
	case c.Bool("coding"):
		return setupCoding(c)
	default:
		return nil, fmt.Errorf("version control system not configured")
	}
}

// helper function to setup the Bitbucket remote from the CLI arguments.
func setupBitbucket(c *cli.Context) (remote.Remote, error) {
	return bitbucket.New(
		c.String("bitbucket-client"),
		c.String("bitbucket-secret"),
	), nil
}

// helper function to setup the Gogs remote from the CLI arguments.
func setupGogs(c *cli.Context) (remote.Remote, error) {
	return gogs.New(gogs.Opts{
		URL:         c.String("gogs-server"),
		Username:    c.String("gogs-git-username"),
		Password:    c.String("gogs-git-password"),
		PrivateMode: c.Bool("gogs-private-mode"),
		SkipVerify:  c.Bool("gogs-skip-verify"),
	})
}

// helper function to setup the Gitea remote from the CLI arguments.
func setupGitea(c *cli.Context) (remote.Remote, error) {
	if !c.IsSet("gitea-client") {
		return gitea.New(gitea.Opts{
			URL:                c.String("gitea-server"),
			Context:            c.String("gitea-context"),
			Username:           c.String("gitea-git-username"),
			Password:           c.String("gitea-git-password"),
			PrivateMode:        c.Bool("gitea-private-mode"),
			SkipVerify:         c.Bool("gitea-skip-verify"),
			RevProxyAuth:       c.Bool("gitea-rev-proxy-auth"),
			RevProxyAuthHeader: c.String("gitea-rev-proxy-auth-header"),
		})
	}
	return gitea.NewOauth(gitea.Opts{
		URL:         c.String("gitea-server"),
		Context:     c.String("gitea-context"),
		Username:    c.String("gitea-git-username"),
		Password:    c.String("gitea-git-password"),
		Client:      c.String("gitea-client"),
		Secret:      c.String("gitea-secret"),
		PrivateMode: c.Bool("gitea-private-mode"),
		SkipVerify:  c.Bool("gitea-skip-verify"),
	})
}

// helper function to setup the Stash remote from the CLI arguments.
func setupStash(c *cli.Context) (remote.Remote, error) {
	return bitbucketserver.New(bitbucketserver.Opts{
		URL:               c.String("stash-server"),
		Username:          c.String("stash-git-username"),
		Password:          c.String("stash-git-password"),
		ConsumerKey:       c.String("stash-consumer-key"),
		ConsumerRSA:       c.String("stash-consumer-rsa"),
		ConsumerRSAString: c.String("stash-consumer-rsa-string"),
		SkipVerify:        c.Bool("stash-skip-verify"),
	})
}

// helper function to setup the Gitlab remote from the CLI arguments.
func setupGitlab(c *cli.Context) (remote.Remote, error) {
	if c.Bool("gitlab-v3-api") {
		return gitlab3.New(gitlab3.Opts{
			URL:         c.String("gitlab-server"),
			Client:      c.String("gitlab-client"),
			Secret:      c.String("gitlab-secret"),
			Username:    c.String("gitlab-git-username"),
			Password:    c.String("gitlab-git-password"),
			PrivateMode: c.Bool("gitlab-private-mode"),
			SkipVerify:  c.Bool("gitlab-skip-verify"),
		})
	}
	return gitlab.New(gitlab.Opts{
		URL:         c.String("gitlab-server"),
		Client:      c.String("gitlab-client"),
		Secret:      c.String("gitlab-secret"),
		Username:    c.String("gitlab-git-username"),
		Password:    c.String("gitlab-git-password"),
		PrivateMode: c.Bool("gitlab-private-mode"),
		SkipVerify:  c.Bool("gitlab-skip-verify"),
	})
}

// helper function to setup the GitHub remote from the CLI arguments.
func setupGithub(c *cli.Context) (remote.Remote, error) {
	return github.New(github.Opts{
		URL:         c.String("github-server"),
		Context:     c.String("github-context"),
		Client:      c.String("github-client"),
		Secret:      c.String("github-secret"),
		Scopes:      c.StringSlice("github-scope"),
		Username:    c.String("github-git-username"),
		Password:    c.String("github-git-password"),
		PrivateMode: c.Bool("github-private-mode"),
		SkipVerify:  c.Bool("github-skip-verify"),
		MergeRef:    c.BoolT("github-merge-ref"),
	})
}

// helper function to setup the Coding remote from the CLI arguments.
func setupCoding(c *cli.Context) (remote.Remote, error) {
	return coding.New(coding.Opts{
		URL:        c.String("coding-server"),
		Client:     c.String("coding-client"),
		Secret:     c.String("coding-secret"),
		Scopes:     c.StringSlice("coding-scope"),
		Machine:    c.String("coding-git-machine"),
		Username:   c.String("coding-git-username"),
		Password:   c.String("coding-git-password"),
		SkipVerify: c.Bool("coding-skip-verify"),
	})
}

func setupTree(c *cli.Context) *httptreemux.ContextMux {
	tree := httptreemux.NewContextMux()
	web.New(
		web.WithDir(c.String("www")),
		web.WithSync(time.Hour*72),
		web.WithDocs(c.String("docs")),
	).Register(tree)
	return tree
}

func before(c *cli.Context) error { return nil }

func setupMetrics(g *errgroup.Group, store_ store.Store) {
	pendingJobs := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "drone",
		Name:      "pending_jobs",
		Help:      "Total number of pending build processes.",
	})
	waitingJobs := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "drone",
		Name:      "waiting_jobs",
		Help:      "Total number of builds waiting on deps.",
	})
	runningJobs := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "drone",
		Name:      "running_jobs",
		Help:      "Total number of running build processes.",
	})
	workers := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "drone",
		Name:      "worker_count",
		Help:      "Total number of workers.",
	})
	builds := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "drone",
		Name:      "build_total_count",
		Help:      "Total number of builds.",
	})
	users := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "drone",
		Name:      "user_count",
		Help:      "Total number of users.",
	})
	repos := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "drone",
		Name:      "repo_count",
		Help:      "Total number of repos.",
	})

	g.Go(func() error {
		for {
			stats := droneserver.Config.Services.Queue.Info(nil)
			pendingJobs.Set(float64(stats.Stats.Pending))
			waitingJobs.Set(float64(stats.Stats.WaitingOnDeps))
			runningJobs.Set(float64(stats.Stats.Running))
			workers.Set(float64(stats.Stats.Workers))
			time.Sleep(500 * time.Millisecond)
		}
	})
	g.Go(func() error {
		for {
			repoCount, _ := store_.GetRepoCount()
			userCount, _ := store_.GetUserCount()
			buildCount, _ := store_.GetBuildCount()
			builds.Set(float64(buildCount))
			users.Set(float64(userCount))
			repos.Set(float64(repoCount))
			time.Sleep(10 * time.Second)
		}
	})
}
