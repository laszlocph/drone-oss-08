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

package gitea

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"code.gitea.io/sdk/gitea"
	"github.com/woodpecker-ci/woodpecker/model"
	"github.com/woodpecker-ci/woodpecker/remote"
)

// Opts defines configuration options.
type Opts struct {
	URL                string // Gitea server url.
	Context            string // Context to display in status check
	Client             string // OAuth2 Client ID
	Secret             string // OAuth2 Client Secret
	Username           string // Optional machine account username.
	Password           string // Optional machine account password.
	PrivateMode        bool   // Gitea is running in private mode.
	SkipVerify         bool   // Skip ssl verification.
	RevProxyAuth       bool   // Enable reverse proxy authentication using RevProxyAuthHeader.
	RevProxyAuthHeader string // Name of HTTP header with username for reverse proxy authentication.
}

type client struct {
	URL                     string
	Context                 string
	Machine                 string
	Username                string
	Password                string
	PrivateMode             bool
	SkipVerify              bool
	RevProxyAuth            bool
	RevProxyAuthHeader      string
	RevProxyAuthHeaderValue string
}

const (
	DescPending  = "the build is pending"
	DescRunning  = "the build is running"
	DescSuccess  = "the build was successful"
	DescFailure  = "the build failed"
	DescCanceled = "the build canceled"
	DescBlocked  = "the build is pending approval"
	DescDeclined = "the build was rejected"
)

// getStatus is a helper function that converts a Drone
// status to a Gitea status.
func getStatus(status string) gitea.StatusState {
	switch status {
	case model.StatusPending, model.StatusBlocked:
		return gitea.StatusPending
	case model.StatusRunning:
		return gitea.StatusPending
	case model.StatusSuccess:
		return gitea.StatusSuccess
	case model.StatusFailure, model.StatusError:
		return gitea.StatusFailure
	case model.StatusKilled:
		return gitea.StatusFailure
	case model.StatusDeclined:
		return gitea.StatusWarning
	default:
		return gitea.StatusFailure
	}
}

// getDesc is a helper function that generates a description
// message for the build based on the status.
func getDesc(status string) string {
	switch status {
	case model.StatusPending:
		return DescPending
	case model.StatusRunning:
		return DescRunning
	case model.StatusSuccess:
		return DescSuccess
	case model.StatusFailure, model.StatusError:
		return DescFailure
	case model.StatusKilled:
		return DescCanceled
	case model.StatusBlocked:
		return DescBlocked
	case model.StatusDeclined:
		return DescDeclined
	default:
		return DescFailure
	}
}

// New returns a Remote implementation that integrates with Gitea, an open
// source Git service written in Go. See https://gitea.io/
func New(opts Opts) (remote.Remote, error) {
	url, err := url.Parse(opts.URL)
	if err != nil {
		return nil, err
	}
	host, _, err := net.SplitHostPort(url.Host)
	if err == nil {
		url.Host = host
	}
	return &client{
		URL:                opts.URL,
		Context:            opts.Context,
		Machine:            url.Host,
		Username:           opts.Username,
		Password:           opts.Password,
		PrivateMode:        opts.PrivateMode,
		SkipVerify:         opts.SkipVerify,
		RevProxyAuth:       opts.RevProxyAuth,
		RevProxyAuthHeader: opts.RevProxyAuthHeader,
	}, nil
}

// Login authenticates an account with Gitea using basic authentication. The
// Gitea account details are returned when the user is successfully authenticated.
func (c *client) Login(res http.ResponseWriter, req *http.Request) (*model.User, error) {

	var username string
	var password string

	if c.RevProxyAuth {
		username = req.Header.Get(c.RevProxyAuthHeader)
		c.RevProxyAuthHeaderValue = username
		password = ""
	} else {
		username = req.FormValue("username")
		password = req.FormValue("password")
	}

	// Redirect to the login screen if user credentials are incomplete.
	if len(username) == 0 || (!c.RevProxyAuth && len(password) == 0) {
		http.Redirect(res, req, "/login/form", http.StatusSeeOther)
		return nil, nil
	}

	// Create Client with Basic Auth
	client, err := c.newClientBasicAuth(username, password)
	if err != nil {
		return nil, err
	}

	// since api does not return token secret, if drone token exists create new one
	resp, err := client.DeleteAccessToken("drone")
	if err != nil && !(resp != nil && resp.StatusCode == 404) {
		return nil, err
	}

	token, _, terr := client.CreateAccessToken(
		gitea.CreateAccessTokenOption{Name: "drone"},
	)
	if terr != nil {
		return nil, terr
	}
	accessToken := token.Token

	client, err = c.newClientToken(accessToken)
	if err != nil {
		return nil, err
	}
	account, _, err := client.GetUserInfo(username)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Token:  accessToken,
		Login:  account.UserName,
		Email:  account.Email,
		Avatar: expandAvatar(c.URL, account.AvatarURL),
	}, nil
}

// Auth is not supported by the Gitea driver.
func (c *client) Auth(token, secret string) (string, error) {
	return "", fmt.Errorf("Not Implemented")
}

// Teams is supported by the Gitea driver.
func (c *client) Teams(u *model.User) ([]*model.Team, error) {
	client, err := c.newClientToken(u.Token)
	if err != nil {
		return nil, err
	}

	orgs, _, err := client.ListMyOrgs(gitea.ListOrgsOptions{})
	if err != nil {
		return nil, err
	}

	var teams []*model.Team
	for _, org := range orgs {
		teams = append(teams, toTeam(org, c.URL))
	}
	return teams, nil
}

// TeamPerm is not supported by the Gitea driver.
func (c *client) TeamPerm(u *model.User, org string) (*model.Perm, error) {
	return nil, nil
}

// Repo returns the named Gitea repository.
func (c *client) Repo(u *model.User, owner, name string) (*model.Repo, error) {
	client, err := c.newClientToken(u.Token)
	if err != nil {
		return nil, err
	}

	repo, _, err := client.GetRepo(owner, name)
	if err != nil {
		return nil, err
	}
	if c.PrivateMode {
		repo.Private = true
	}
	return toRepo(repo, c.PrivateMode), nil
}

// Repos returns a list of all repositories for the Gitea account, including
// organization repositories.
func (c *client) Repos(u *model.User) ([]*model.Repo, error) {
	repos := []*model.Repo{}

	client, err := c.newClientToken(u.Token)
	if err != nil {
		return nil, err
	}

	// Gitea SDK forces us to read repo list paginated.
	var page int = 1
	for {
		all, _, err := client.ListMyRepos(
			gitea.ListReposOptions{
				ListOptions: gitea.ListOptions{
					Page:     page,
					PageSize: 50, // Gitea SDK limit per page.
				},
			},
		)

		// Gitea SDK does not return error when asking for
		// non existing repos page (empty list is returned)
		// so this should be safe.
		if err != nil {
			return repos, err
		}

		for _, repo := range all {
			repos = append(repos, toRepo(repo, c.PrivateMode))
		}

		// Check if no more repos are available; we don't test len(all) < 50
		// because of Gitea SDK bug https://gitea.com/gitea/go-sdk/issues/507.
		if len(all) == 0 {
			// Empty page returned - finish loop.
			break
		} else {
			// Last page was not empty so more repos may be available - continue loop.
			page = page + 1
		}
	}

	return repos, nil
}

// Perm returns the user permissions for the named Gitea repository.
func (c *client) Perm(u *model.User, owner, name string) (*model.Perm, error) {
	client, err := c.newClientToken(u.Token)
	if err != nil {
		return nil, err
	}

	repo, _, err := client.GetRepo(owner, name)
	if err != nil {
		return nil, err
	}
	return toPerm(repo.Permissions), nil
}

// File fetches the file from the Gitea repository and returns its contents.
func (c *client) File(u *model.User, r *model.Repo, b *model.Build, f string) ([]byte, error) {
	client, err := c.newClientToken(u.Token)
	if err != nil {
		return nil, err
	}

	cfg, _, err := client.GetFile(r.Owner, r.Name, b.Commit, f)
	return cfg, err
}

func (c *client) Dir(u *model.User, r *model.Repo, b *model.Build, f string) ([]*remote.FileMeta, error) {
	var configs []*remote.FileMeta

	client, err := c.newClientToken(u.Token)
	if err != nil {
		return nil, err
	}

	// List files in repository. Path from root
	tree, _, err := client.GetTrees(r.Owner, r.Name, b.Commit, true)
	if err != nil {
		return nil, err
	}

	f = path.Clean(f) // We clean path and remove trailing slash
	f += "/" + "*"    // construct pattern for match i.e. file in subdir
	for _, e := range tree.Entries {
		// Filter path matching pattern and type file (blob)
		if m, _ := filepath.Match(f, e.Path); m && e.Type == "blob" {
			data, err := c.File(u, r, b, e.Path)
			if err != nil {
				return nil, fmt.Errorf("multi-pipeline cannot get %s: %s", e.Path, err)
			}

			configs = append(configs, &remote.FileMeta{
				Name: e.Path,
				Data: data,
			})
		}
	}

	return configs, nil
}

// Status is supported by the Gitea driver.
func (c *client) Status(u *model.User, r *model.Repo, b *model.Build, link string, proc *model.Proc) error {
	client, err := c.newClientToken(u.Token)
	if err != nil {
		return err
	}

	status := getStatus(b.Status)
	desc := getDesc(b.Status)

	_, _, err = client.CreateStatus(
		r.Owner,
		r.Name,
		b.Commit,
		gitea.CreateStatusOption{
			State:       status,
			TargetURL:   link,
			Description: desc,
			Context:     c.Context,
		},
	)

	return err
}

// Netrc returns a netrc file capable of authenticating Gitea requests and
// cloning Gitea repositories. The netrc will use the global machine account
// when configured.
func (c *client) Netrc(u *model.User, r *model.Repo) (*model.Netrc, error) {
	if c.Password != "" {
		return &model.Netrc{
			Login:    c.Username,
			Password: c.Password,
			Machine:  c.Machine,
		}, nil
	}
	return &model.Netrc{
		Login:    u.Login,
		Password: u.Token,
		Machine:  c.Machine,
	}, nil
}

// Activate activates the repository by registering post-commit hooks with
// the Gitea repository.
func (c *client) Activate(u *model.User, r *model.Repo, link string) error {
	config := map[string]string{
		"url":          link,
		"secret":       r.Hash,
		"content_type": "json",
	}
	hook := gitea.CreateHookOption{
		Type:   "gitea",
		Config: config,
		Events: []string{"push", "create", "pull_request"},
		Active: true,
	}

	client, err := c.newClientToken(u.Token)
	if err != nil {
		return err
	}
	_, _, err = client.CreateRepoHook(r.Owner, r.Name, hook)
	return err
}

// Deactivate deactives the repository be removing repository push hooks from
// the Gitea repository.
func (c *client) Deactivate(u *model.User, r *model.Repo, link string) error {
	client, err := c.newClientToken(u.Token)
	if err != nil {
		return err
	}

	hooks, _, err := client.ListRepoHooks(r.Owner, r.Name, gitea.ListHooksOptions{})
	if err != nil {
		return err
	}

	hook := matchingHooks(hooks, link)
	if hook != nil {
		_, err := client.DeleteRepoHook(r.Owner, r.Name, hook.ID)
		return err
	}

	return nil
}

// Hook parses the incoming Gitea hook and returns the Repository and Build
// details. If the hook is unsupported nil values are returned.
func (c *client) Hook(r *http.Request) (*model.Repo, *model.Build, error) {
	return parseHook(r)
}

// authTransport forwards authentication HTTP header to gitea.
type authTransport struct {
	headerName          string
	headerValue         string
	underlyingTransport http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(t.headerName, t.headerValue)
	if t.underlyingTransport != nil {
		return t.underlyingTransport.RoundTrip(req)
	} else {
		return http.DefaultTransport.RoundTrip(req)
	}
}

// helper function to return the Gitea client with Token
func (c *client) newClientToken(token string) (*gitea.Client, error) {
	httpClient := &http.Client{}
	if c.SkipVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	// Forward authentication header in every HTTP request to Gitea
	// in reverse proxy authentication mode.
	if c.RevProxyAuth {
		httpClient.Transport = &authTransport{
			headerName:          c.RevProxyAuthHeader,
			headerValue:         c.RevProxyAuthHeaderValue,
			underlyingTransport: httpClient.Transport,
		}
	}
	return gitea.NewClient(c.URL, gitea.SetToken(token), gitea.SetHTTPClient(httpClient))
}

// helper function to return the Gitea client with Basic Auth
func (c *client) newClientBasicAuth(username, password string) (*gitea.Client, error) {
	httpClient := &http.Client{}
	if c.SkipVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return gitea.NewClient(c.URL, gitea.SetBasicAuth(username, password), gitea.SetHTTPClient(httpClient))
}

// helper function to return matching hooks.
func matchingHooks(hooks []*gitea.Hook, rawurl string) *gitea.Hook {
	link, err := url.Parse(rawurl)
	if err != nil {
		return nil
	}
	for _, hook := range hooks {
		if val, ok := hook.Config["url"]; ok {
			hookurl, err := url.Parse(val)
			if err == nil && hookurl.Host == link.Host {
				return hook
			}
		}
	}
	return nil
}
