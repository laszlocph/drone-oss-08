package drone

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	pathSelf           = "%s/api/user"
	pathFeed           = "%s/api/user/feed"
	pathRepos          = "%s/api/user/repos"
	pathRepo           = "%s/api/repos/%s/%s"
	pathRepoMove       = "%s/api/repos/%s/%s/move?to=%s"
	pathChown          = "%s/api/repos/%s/%s/chown"
	pathRepair         = "%s/api/repos/%s/%s/repair"
	pathBuilds         = "%s/api/repos/%s/%s/builds"
	pathBuild          = "%s/api/repos/%s/%s/builds/%v"
	pathApprove        = "%s/api/repos/%s/%s/builds/%d/approve"
	pathDecline        = "%s/api/repos/%s/%s/builds/%d/decline"
	pathJob            = "%s/api/repos/%s/%s/builds/%d/%d"
	pathLog            = "%s/api/repos/%s/%s/logs/%d/%d"
	pathLogPurge       = "%s/api/repos/%s/%s/logs/%d"
	pathRepoSecrets    = "%s/api/repos/%s/%s/secrets"
	pathRepoSecret     = "%s/api/repos/%s/%s/secrets/%s"
	pathRepoRegistries = "%s/api/repos/%s/%s/registry"
	pathRepoRegistry   = "%s/api/repos/%s/%s/registry/%s"
	pathUsers          = "%s/api/users"
	pathUser           = "%s/api/users/%s"
	pathBuildQueue     = "%s/api/builds"
	pathQueue          = "%s/api/queue"
	pathVersion        = "%s/version"
	pathGlobalSecrets  = "%s/api/globalsecrets"
	pathGlobalSecret   = "%s/api/globalsecrets/%s"
)

type client struct {
	client *http.Client
	addr   string
}

// New returns a client at the specified url.
func New(uri string) Client {
	return &client{http.DefaultClient, strings.TrimSuffix(uri, "/")}
}

// NewClient returns a client at the specified url.
func NewClient(uri string, cli *http.Client) Client {
	return &client{cli, strings.TrimSuffix(uri, "/")}
}

// SetClient sets the http.Client.
func (c *client) SetClient(client *http.Client) {
	c.client = client
}

// SetAddress sets the server address.
func (c *client) SetAddress(addr string) {
	c.addr = addr
}

// Self returns the currently authenticated user.
func (c *client) Self() (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathSelf, c.addr)
	err := c.get(uri, out)
	return out, err
}

// User returns a user by login.
func (c *client) User(login string) (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathUser, c.addr, login)
	err := c.get(uri, out)
	return out, err
}

// UserList returns a list of all registered users.
func (c *client) UserList() ([]*User, error) {
	var out []*User
	uri := fmt.Sprintf(pathUsers, c.addr)
	err := c.get(uri, &out)
	return out, err
}

// UserPost creates a new user account.
func (c *client) UserPost(in *User) (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathUsers, c.addr)
	err := c.post(uri, in, out)
	return out, err
}

// UserPatch updates a user account.
func (c *client) UserPatch(in *User) (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathUser, c.addr, in.Login)
	err := c.patch(uri, in, out)
	return out, err
}

// UserDel deletes a user account.
func (c *client) UserDel(login string) error {
	uri := fmt.Sprintf(pathUser, c.addr, login)
	err := c.delete(uri)
	return err
}

// Repo returns a repository by name.
func (c *client) Repo(owner string, name string) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathRepo, c.addr, owner, name)
	err := c.get(uri, out)
	return out, err
}

// RepoList returns a list of all repositories to which
// the user has explicit access in the host system.
func (c *client) RepoList() ([]*Repo, error) {
	var out []*Repo
	uri := fmt.Sprintf(pathRepos, c.addr)
	err := c.get(uri, &out)
	return out, err
}

// RepoListOpts returns a list of all repositories to which
// the user has explicit access in the host system.
func (c *client) RepoListOpts(sync, all bool) ([]*Repo, error) {
	var out []*Repo
	uri := fmt.Sprintf(pathRepos+"?flush=%v&all=%v", c.addr, sync, all)
	err := c.get(uri, &out)
	return out, err
}

// RepoPost activates a repository.
func (c *client) RepoPost(owner string, name string) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathRepo, c.addr, owner, name)
	err := c.post(uri, nil, out)
	return out, err
}

// RepoChown updates a repository owner.
func (c *client) RepoChown(owner string, name string) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathChown, c.addr, owner, name)
	err := c.post(uri, nil, out)
	return out, err
}

// RepoRepair repais the repository hooks.
func (c *client) RepoRepair(owner string, name string) error {
	uri := fmt.Sprintf(pathRepair, c.addr, owner, name)
	return c.post(uri, nil, nil)
}

// RepoPatch updates a repository.
func (c *client) RepoPatch(owner, name string, in *RepoPatch) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathRepo, c.addr, owner, name)
	err := c.patch(uri, in, out)
	return out, err
}

// RepoDel deletes a repository.
func (c *client) RepoDel(owner, name string) error {
	uri := fmt.Sprintf(pathRepo, c.addr, owner, name)
	err := c.delete(uri)
	return err
}

// RepoMove moves a repository
func (c *client) RepoMove(owner, name, newFullName string) error {
	uri := fmt.Sprintf(pathRepoMove, c.addr, owner, name, newFullName)
	return c.post(uri, nil, nil)
}

// Build returns a repository build by number.
func (c *client) Build(owner, name string, num int) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathBuild, c.addr, owner, name, num)
	err := c.get(uri, out)
	return out, err
}

// Build returns the latest repository build by branch.
func (c *client) BuildLast(owner, name, branch string) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathBuild, c.addr, owner, name, "latest")
	if len(branch) != 0 {
		uri += "?branch=" + branch
	}
	err := c.get(uri, out)
	return out, err
}

// BuildList returns a list of recent builds for the
// the specified repository.
func (c *client) BuildList(owner, name string) ([]*Build, error) {
	var out []*Build
	uri := fmt.Sprintf(pathBuilds, c.addr, owner, name)
	err := c.get(uri, &out)
	return out, err
}

// BuildQueue returns a list of enqueued builds.
func (c *client) BuildQueue() ([]*Activity, error) {
	var out []*Activity
	uri := fmt.Sprintf(pathBuildQueue, c.addr)
	err := c.get(uri, &out)
	return out, err
}

// BuildStart re-starts a stopped build.
func (c *client) BuildStart(owner, name string, num int, params map[string]string) (*Build, error) {
	out := new(Build)
	val := mapValues(params)
	uri := fmt.Sprintf(pathBuild, c.addr, owner, name, num)
	err := c.post(uri+"?"+val.Encode(), nil, out)
	return out, err
}

// BuildStop cancels the running job.
func (c *client) BuildStop(owner, name string, num, job int) error {
	uri := fmt.Sprintf(pathJob, c.addr, owner, name, num, job)
	err := c.delete(uri)
	return err
}

// BuildApprove approves a blocked build.
func (c *client) BuildApprove(owner, name string, num int) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathApprove, c.addr, owner, name, num)
	err := c.post(uri, nil, out)
	return out, err
}

// BuildDecline declines a blocked build.
func (c *client) BuildDecline(owner, name string, num int) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathDecline, c.addr, owner, name, num)
	err := c.post(uri, nil, out)
	return out, err
}

// BuildKill force kills the running build.
func (c *client) BuildKill(owner, name string, num int) error {
	uri := fmt.Sprintf(pathBuild, c.addr, owner, name, num)
	err := c.delete(uri)
	return err
}

// BuildLogs returns the build logs for the specified job.
func (c *client) BuildLogs(owner, name string, num, job int) (io.ReadCloser, error) {
	return nil, errors.New("Method not implemented")
}

// Deploy triggers a deployment for an existing build using the
// specified target environment.
func (c *client) Deploy(owner, name string, num int, env string, params map[string]string) (*Build, error) {
	out := new(Build)
	val := mapValues(params)
	val.Set("event", "deployment")
	val.Set("deploy_to", env)
	uri := fmt.Sprintf(pathBuild, c.addr, owner, name, num)
	err := c.post(uri+"?"+val.Encode(), nil, out)
	return out, err
}

// LogsPurge purges the build logs for the specified build.
func (c *client) LogsPurge(owner, name string, num int) error {
	uri := fmt.Sprintf(pathLogPurge, c.addr, owner, name, num)
	err := c.delete(uri)
	return err
}

// Registry returns a registry by hostname.
func (c *client) Registry(owner, name, hostname string) (*Registry, error) {
	out := new(Registry)
	uri := fmt.Sprintf(pathRepoRegistry, c.addr, owner, name, hostname)
	err := c.get(uri, out)
	return out, err
}

// RegistryList returns a list of all repository registries.
func (c *client) RegistryList(owner string, name string) ([]*Registry, error) {
	var out []*Registry
	uri := fmt.Sprintf(pathRepoRegistries, c.addr, owner, name)
	err := c.get(uri, &out)
	return out, err
}

// RegistryCreate creates a registry.
func (c *client) RegistryCreate(owner, name string, in *Registry) (*Registry, error) {
	out := new(Registry)
	uri := fmt.Sprintf(pathRepoRegistries, c.addr, owner, name)
	err := c.post(uri, in, out)
	return out, err
}

// RegistryUpdate updates a registry.
func (c *client) RegistryUpdate(owner, name string, in *Registry) (*Registry, error) {
	out := new(Registry)
	uri := fmt.Sprintf(pathRepoRegistry, c.addr, owner, name, in.Address)
	err := c.patch(uri, in, out)
	return out, err
}

// RegistryDelete deletes a registry.
func (c *client) RegistryDelete(owner, name, hostname string) error {
	uri := fmt.Sprintf(pathRepoRegistry, c.addr, owner, name, hostname)
	return c.delete(uri)
}

// Secret returns a secret by name.
func (c *client) Secret(owner, name, secret string) (*Secret, error) {
	out := new(Secret)
	uri := fmt.Sprintf(pathRepoSecret, c.addr, owner, name, secret)
	err := c.get(uri, out)
	return out, err
}

// SecretList returns a list of all repository secrets.
func (c *client) SecretList(owner string, name string) ([]*Secret, error) {
	var out []*Secret
	uri := fmt.Sprintf(pathRepoSecrets, c.addr, owner, name)
	err := c.get(uri, &out)
	return out, err
}

// SecretCreate creates a secret.
func (c *client) SecretCreate(owner, name string, in *Secret) (*Secret, error) {
	out := new(Secret)
	uri := fmt.Sprintf(pathRepoSecrets, c.addr, owner, name)
	err := c.post(uri, in, out)
	return out, err
}

// SecretUpdate updates a secret.
func (c *client) SecretUpdate(owner, name string, in *Secret) (*Secret, error) {
	out := new(Secret)
	uri := fmt.Sprintf(pathRepoSecret, c.addr, owner, name, in.Name)
	err := c.patch(uri, in, out)
	return out, err
}

// SecretDelete deletes a secret.
func (c *client) SecretDelete(owner, name, secret string) error {
	uri := fmt.Sprintf(pathRepoSecret, c.addr, owner, name, secret)
	return c.delete(uri)
}

func (c *client) QueueInfo() (*Info, error) {
	out := new(Info)
	uri := fmt.Sprintf(pathQueue+"/info", c.addr)
	err := c.get(uri, out)
	return out, err
}

// GlobalSecret returns a secret by name.
func (c *client) GlobalSecret(secret string) (*Secret, error) {
	out := new(Secret)
	uri := fmt.Sprintf(pathGlobalSecret, c.addr, secret)
	err := c.get(uri, out)
	return out, err
}

// GlobalSecretList returns a list of all global secrets.
func (c *client) GlobalSecretList() ([]*Secret, error) {
	var out []*Secret
	uri := fmt.Sprintf(pathGlobalSecrets, c.addr)
	err := c.get(uri, &out)
	return out, err
}

// GlobalSecretCreate creates a secret.
func (c *client) GlobalSecretCreate(in *Secret) (*Secret, error) {
	out := new(Secret)
	uri := fmt.Sprintf(pathGlobalSecrets, c.addr)
	err := c.post(uri, in, out)
	return out, err
}

// GlobalSecretUpdate updates a secret.
func (c *client) GlobalSecretUpdate(in *Secret) (*Secret, error) {
	out := new(Secret)
	uri := fmt.Sprintf(pathGlobalSecret, c.addr, in.Name)
	err := c.patch(uri, in, out)
	return out, err
}

// GlobalSecretDelete deletes a secret.
func (c *client) GlobalSecretDelete(secret string) error {
	uri := fmt.Sprintf(pathGlobalSecret, c.addr, secret)
	return c.delete(uri)
}

//
// http request helper functions
//

// helper function for making an http GET request.
func (c *client) get(rawurl string, out interface{}) error {
	return c.do(rawurl, "GET", nil, out)
}

// helper function for making an http POST request.
func (c *client) post(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "POST", in, out)
}

// helper function for making an http PUT request.
func (c *client) put(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "PUT", in, out)
}

// helper function for making an http PATCH request.
func (c *client) patch(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "PATCH", in, out)
}

// helper function for making an http DELETE request.
func (c *client) delete(rawurl string) error {
	return c.do(rawurl, "DELETE", nil, nil)
}

// helper function to make an http request
func (c *client) do(rawurl, method string, in, out interface{}) error {
	body, err := c.open(rawurl, method, in, out)
	if err != nil {
		return err
	}
	defer body.Close()
	if out != nil {
		return json.NewDecoder(body).Decode(out)
	}
	return nil
}

// helper function to open an http request
func (c *client) open(rawurl, method string, in, out interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return nil, err
	}
	if in != nil {
		decoded, derr := json.Marshal(in)
		if derr != nil {
			return nil, derr
		}
		buf := bytes.NewBuffer(decoded)
		req.Body = ioutil.NopCloser(buf)
		req.ContentLength = int64(len(decoded))
		req.Header.Set("Content-Length", strconv.Itoa(len(decoded)))
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > http.StatusPartialContent {
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("client error %d: %s", resp.StatusCode, string(out))
	}
	return resp.Body, nil
}

// mapValues converts a map to url.Values
func mapValues(params map[string]string) url.Values {
	values := url.Values{}
	for key, val := range params {
		values.Add(key, val)
	}
	return values
}
