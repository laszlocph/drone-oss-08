// Copyright 2019 Laszlo Fogas
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
	"time"

	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "DRONE_DEBUG,WOODPECKER_DEBUG",
		Name:   "debug",
		Usage:  "enable server debug mode",
	},
	cli.StringFlag{
		EnvVar: "DRONE_SERVER_HOST,DRONE_HOST,WOODPECKER_SERVER_HOST,WOODPECKER_HOST",
		Name:   "server-host",
		Usage:  "server fully qualified url (<scheme>://<host>)",
	},
	cli.StringFlag{
		EnvVar: "WOODPECKER_HOST_INTERNAL",
		Name:   "server-host-internal",
		Usage:  "server internal fully qualified url (<scheme>://<host>)",
	},
	cli.StringFlag{
		EnvVar: "DRONE_SERVER_ADDR,WOODPECKER_SERVER_ADDR",
		Name:   "server-addr",
		Usage:  "server address",
		Value:  ":8000",
	},
	cli.StringFlag{
		EnvVar: "DRONE_SERVER_CERT,WOODPECKER_SERVER_CERT",
		Name:   "server-cert",
		Usage:  "server ssl cert path",
	},
	cli.StringFlag{
		EnvVar: "DRONE_SERVER_KEY,WOODPECKER_SERVER_KEY",
		Name:   "server-key",
		Usage:  "server ssl key path",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_LETS_ENCRYPT,WOODPECKER_LETS_ENCRYPT",
		Name:   "lets-encrypt",
		Usage:  "enable let's encrypt",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_QUIC,WOODPECKER_QUIC",
		Name:   "quic",
		Usage:  "enable quic",
	},
	cli.StringFlag{
		EnvVar: "DRONE_WWW,WOODPECKER_WWW",
		Name:   "www",
		Usage:  "serve the website from disk",
		Hidden: true,
	},
	cli.StringSliceFlag{
		EnvVar: "DRONE_ADMIN,WOODPECKER_ADMIN",
		Name:   "admin",
		Usage:  "list of admin users",
	},
	cli.StringSliceFlag{
		EnvVar: "DRONE_ORGS,WOODPECKER_ORGS",
		Name:   "orgs",
		Usage:  "list of approved organizations",
	},
	cli.StringSliceFlag{
		EnvVar: "DRONE_REPO_OWNERS,WOODPECKER_REPO_OWNERS",
		Name:   "repo-owners",
		Usage:  "List of syncable repo owners",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_OPEN,WOODPECKER_OPEN",
		Name:   "open",
		Usage:  "enable open user registration",
	},
	cli.StringFlag{
		EnvVar: "DRONE_REPO_CONFIG,WOODPECKER_REPO_CONFIG",
		Name:   "repo-config",
		Usage:  "file path for the drone config",
		Value:  ".drone.yml",
	},
	cli.StringFlag{
		EnvVar: "DRONE_DOCS,WOODPECKER_DOCS",
		Name:   "docs",
		Usage:  "link to user documentation",
		Value:  "https://woodpecker.laszlo.cloud",
	},
	cli.DurationFlag{
		EnvVar: "DRONE_SESSION_EXPIRES,WOODPECKER_SESSION_EXPIRES",
		Name:   "session-expires",
		Usage:  "session expiration time",
		Value:  time.Hour * 72,
	},
	cli.StringSliceFlag{
		EnvVar: "DRONE_ESCALATE,WOODPECKER_ESCALATE",
		Name:   "escalate",
		Usage:  "images to run in privileged mode",
		Value: &cli.StringSlice{
			"plugins/docker",
			"plugins/gcr",
			"plugins/ecr",
		},
	},
	cli.StringSliceFlag{
		EnvVar: "DRONE_VOLUME,WOODPECKER_VOLUME",
		Name:   "volume",
	},
	cli.StringFlag{
		EnvVar: "DRONE_DOCKER_CONFIG,WOODPECKER_DOCKER_CONFIG",
		Name:   "docker-config",
	},
	cli.StringSliceFlag{
		EnvVar: "DRONE_ENVIRONMENT,WOODPECKER_ENVIRONMENT",
		Name:   "environment",
	},
	cli.StringSliceFlag{
		EnvVar: "DRONE_NETWORK,WOODPECKER_NETWORK",
		Name:   "network",
	},
	cli.StringFlag{
		EnvVar: "DRONE_AGENT_SECRET,DRONE_SECRET,WOODPECKER_AGENT_SECRET,WOODPECKER_SECRET",
		Name:   "agent-secret",
		Usage:  "server-agent shared password",
	},
	cli.StringFlag{
		EnvVar: "DRONE_SECRET_ENDPOINT,WOODPECKER_SECRET_ENDPOINT",
		Name:   "secret-service",
		Usage:  "secret plugin endpoint",
	},
	cli.StringFlag{
		EnvVar: "DRONE_REGISTRY_ENDPOINT,WOODPECKER_REGISTRY_ENDPOINT",
		Name:   "registry-service",
		Usage:  "registry plugin endpoint",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GATEKEEPER_ENDPOINT,WOODPECKER_GATEKEEPER_ENDPOINT",
		Name:   "gating-service",
		Usage:  "gated build endpoint",
	},
	cli.StringFlag{
		EnvVar: "DRONE_DATABASE_DRIVER,DATABASE_DRIVER,WOODPECKER_DATABASE_DRIVER,DATABASE_DRIVER",
		Name:   "driver",
		Usage:  "database driver",
		Value:  "sqlite3",
	},
	cli.StringFlag{
		EnvVar: "DRONE_DATABASE_DATASOURCE,DATABASE_CONFIG,WOODPECKER_DATABASE_DATASOURCE,DATABASE_CONFIG",
		Name:   "datasource",
		Usage:  "database driver configuration string",
		Value:  "drone.sqlite",
	},
	cli.StringFlag{
		EnvVar: "DRONE_PROMETHEUS_AUTH_TOKEN,WOODPECKER_PROMETHEUS_AUTH_TOKEN",
		Name:   "prometheus-auth-token",
		Usage:  "token to secure prometheus metrics endpoint",
		Value:  "",
	},
	//
	// resource limit parameters
	//
	cli.Int64Flag{
		EnvVar: "DRONE_LIMIT_MEM_SWAP,WOODPECKER_LIMIT_MEM_SWAP",
		Name:   "limit-mem-swap",
		Usage:  "maximum swappable memory allowed in bytes",
	},
	cli.Int64Flag{
		EnvVar: "DRONE_LIMIT_MEM,WOODPECKER_LIMIT_MEM",
		Name:   "limit-mem",
		Usage:  "maximum memory allowed in bytes",
	},
	cli.Int64Flag{
		EnvVar: "DRONE_LIMIT_SHM_SIZE,WOODPECKER_LIMIT_SHM_SIZE",
		Name:   "limit-shm-size",
		Usage:  "docker compose /dev/shm allowed in bytes",
	},
	cli.Int64Flag{
		EnvVar: "DRONE_LIMIT_CPU_QUOTA,WOODPECKER_LIMIT_CPU_QUOTA",
		Name:   "limit-cpu-quota",
		Usage:  "impose a cpu quota",
	},
	cli.Int64Flag{
		EnvVar: "DRONE_LIMIT_CPU_SHARES,WOODPECKER_LIMIT_CPU_SHARES",
		Name:   "limit-cpu-shares",
		Usage:  "change the cpu shares",
	},
	cli.StringFlag{
		EnvVar: "DRONE_LIMIT_CPU_SET,WOODPECKER_LIMIT_CPU_SET",
		Name:   "limit-cpu-set",
		Usage:  "set the cpus allowed to execute containers",
	},
	//
	// remote parameters
	//
	cli.BoolFlag{
		EnvVar: "DRONE_GITHUB,WOODPECKER_GITHUB",
		Name:   "github",
		Usage:  "github driver is enabled",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITHUB_URL,WOODPECKER_GITHUB_URL",
		Name:   "github-server",
		Usage:  "github server address",
		Value:  "https://github.com",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITHUB_CONTEXT,WOODPECKER_GITHUB_CONTEXT",
		Name:   "github-context",
		Usage:  "github status context",
		Value:  "continuous-integration/drone",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITHUB_CLIENT,WOODPECKER_GITHUB_CLIENT",
		Name:   "github-client",
		Usage:  "github oauth2 client id",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITHUB_SECRET,WOODPECKER_GITHUB_SECRET",
		Name:   "github-secret",
		Usage:  "github oauth2 client secret",
	},
	cli.StringSliceFlag{
		EnvVar: "DRONE_GITHUB_SCOPE,WOODPECKER_GITHUB_SCOPE",
		Name:   "github-scope",
		Usage:  "github oauth scope",
		Value: &cli.StringSlice{
			"repo",
			"repo:status",
			"user:email",
			"read:org",
		},
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITHUB_GIT_USERNAME,WOODPECKER_GITHUB_GIT_USERNAME",
		Name:   "github-git-username",
		Usage:  "github machine user username",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITHUB_GIT_PASSWORD,WOODPECKER_GITHUB_GIT_PASSWORD",
		Name:   "github-git-password",
		Usage:  "github machine user password",
	},
	cli.BoolTFlag{
		EnvVar: "DRONE_GITHUB_MERGE_REF,WOODPECKER_GITHUB_MERGE_REF",
		Name:   "github-merge-ref",
		Usage:  "github pull requests use merge ref",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GITHUB_PRIVATE_MODE,WOODPECKER_GITHUB_PRIVATE_MODE",
		Name:   "github-private-mode",
		Usage:  "github is running in private mode",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GITHUB_SKIP_VERIFY,WOODPECKER_GITHUB_SKIP_VERIFY",
		Name:   "github-skip-verify",
		Usage:  "github skip ssl verification",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GOGS,WOODPECKER_GOGS",
		Name:   "gogs",
		Usage:  "gogs driver is enabled",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GOGS_URL,WOODPECKER_GOGS_URL",
		Name:   "gogs-server",
		Usage:  "gogs server address",
		Value:  "https://github.com",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GOGS_GIT_USERNAME,WOODPECKER_GOGS_GIT_USERNAME",
		Name:   "gogs-git-username",
		Usage:  "gogs service account username",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GOGS_GIT_PASSWORD,WOODPECKER_GOGS_GIT_PASSWORD",
		Name:   "gogs-git-password",
		Usage:  "gogs service account password",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GOGS_PRIVATE_MODE,WOODPECKER_GOGS_PRIVATE_MODE",
		Name:   "gogs-private-mode",
		Usage:  "gogs private mode enabled",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GOGS_SKIP_VERIFY,WOODPECKER_GOGS_SKIP_VERIFY",
		Name:   "gogs-skip-verify",
		Usage:  "gogs skip ssl verification",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GITEA,WOODPECKER_GITEA",
		Name:   "gitea",
		Usage:  "gitea driver is enabled",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITEA_URL,WOODPECKER_GITEA_URL",
		Name:   "gitea-server",
		Usage:  "gitea server address",
		Value:  "https://try.gitea.io",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITEA_CLIENT,WOODPECKER_GITEA_CLIENT",
		Name:   "gitea-client",
		Usage:  "gitea oauth2 client id",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITEA_SECRET,WOODPECKER_GITEA_SECRET",
		Name:   "gitea-secret",
		Usage:  "gitea oauth2 client secret",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITEA_CONTEXT,WOODPECKER_GITEA_CONTEXT",
		Name:   "gitea-context",
		Usage:  "gitea status context",
		Value:  "continuous-integration/drone",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITEA_GIT_USERNAME,WOODPECKER_GITEA_GIT_USERNAME",
		Name:   "gitea-git-username",
		Usage:  "gitea service account username",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITEA_GIT_PASSWORD,WOODPECKER_GITEA_GIT_PASSWORD",
		Name:   "gitea-git-password",
		Usage:  "gitea service account password",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GITEA_PRIVATE_MODE,WOODPECKER_GITEA_PRIVATE_MODE",
		Name:   "gitea-private-mode",
		Usage:  "gitea private mode enabled",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GITEA_SKIP_VERIFY,WOODPECKER_GITEA_SKIP_VERIFY",
		Name:   "gitea-skip-verify",
		Usage:  "gitea skip ssl verification",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_BITBUCKET,WOODPECKER_BITBUCKET",
		Name:   "bitbucket",
		Usage:  "bitbucket driver is enabled",
	},
	cli.StringFlag{
		EnvVar: "DRONE_BITBUCKET_CLIENT,WOODPECKER_BITBUCKET_CLIENT",
		Name:   "bitbucket-client",
		Usage:  "bitbucket oauth2 client id",
	},
	cli.StringFlag{
		EnvVar: "DRONE_BITBUCKET_SECRET,WOODPECKER_BITBUCKET_SECRET",
		Name:   "bitbucket-secret",
		Usage:  "bitbucket oauth2 client secret",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GITLAB,WOODPECKER_GITLAB",
		Name:   "gitlab",
		Usage:  "gitlab driver is enabled",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITLAB_URL,WOODPECKER_GITLAB_URL",
		Name:   "gitlab-server",
		Usage:  "gitlab server address",
		Value:  "https://gitlab.com",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITLAB_CLIENT,WOODPECKER_GITLAB_CLIENT",
		Name:   "gitlab-client",
		Usage:  "gitlab oauth2 client id",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITLAB_SECRET,WOODPECKER_GITLAB_SECRET",
		Name:   "gitlab-secret",
		Usage:  "gitlab oauth2 client secret",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITLAB_GIT_USERNAME,WOODPECKER_GITLAB_GIT_USERNAME",
		Name:   "gitlab-git-username",
		Usage:  "gitlab service account username",
	},
	cli.StringFlag{
		EnvVar: "DRONE_GITLAB_GIT_PASSWORD,WOODPECKER_GITLAB_GIT_PASSWORD",
		Name:   "gitlab-git-password",
		Usage:  "gitlab service account password",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GITLAB_SKIP_VERIFY,WOODPECKER_GITLAB_SKIP_VERIFY",
		Name:   "gitlab-skip-verify",
		Usage:  "gitlab skip ssl verification",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GITLAB_PRIVATE_MODE,WOODPECKER_GITLAB_PRIVATE_MODE",
		Name:   "gitlab-private-mode",
		Usage:  "gitlab is running in private mode",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_GITLAB_V3_API,WOODPECKER_GITLAB_V3_API",
		Name:   "gitlab-v3-api",
		Usage:  "gitlab is running the v3 api",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_STASH,WOODPECKER_STASH",
		Name:   "stash",
		Usage:  "stash driver is enabled",
	},
	cli.StringFlag{
		EnvVar: "DRONE_STASH_URL,WOODPECKER_STASH_URL",
		Name:   "stash-server",
		Usage:  "stash server address",
	},
	cli.StringFlag{
		EnvVar: "DRONE_STASH_CONSUMER_KEY,WOODPECKER_STASH_CONSUMER_KEY",
		Name:   "stash-consumer-key",
		Usage:  "stash oauth1 consumer key",
	},
	cli.StringFlag{
		EnvVar: "DRONE_STASH_CONSUMER_RSA,WOODPECKER_STASH_CONSUMER_RSA",
		Name:   "stash-consumer-rsa",
		Usage:  "stash oauth1 private key file",
	},
	cli.StringFlag{
		EnvVar: "DRONE_STASH_CONSUMER_RSA_STRING,WOODPECKER_STASH_CONSUMER_RSA_STRING",
		Name:   "stash-consumer-rsa-string",
		Usage:  "stash oauth1 private key string",
	},
	cli.StringFlag{
		EnvVar: "DRONE_STASH_GIT_USERNAME,WOODPECKER_STASH_GIT_USERNAME",
		Name:   "stash-git-username",
		Usage:  "stash service account username",
	},
	cli.StringFlag{
		EnvVar: "DRONE_STASH_GIT_PASSWORD,WOODPECKER_STASH_GIT_PASSWORD",
		Name:   "stash-git-password",
		Usage:  "stash service account password",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_STASH_SKIP_VERIFY,WOODPECKER_STASH_SKIP_VERIFY",
		Name:   "stash-skip-verify",
		Usage:  "stash skip ssl verification",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_CODING,WOODPECKER_CODING",
		Name:   "coding",
		Usage:  "coding driver is enabled",
	},
	cli.StringFlag{
		EnvVar: "DRONE_CODING_URL,WOODPECKER_CODING_URL",
		Name:   "coding-server",
		Usage:  "coding server address",
		Value:  "https://coding.net",
	},
	cli.StringFlag{
		EnvVar: "DRONE_CODING_CLIENT,WOODPECKER_CODING_CLIENT",
		Name:   "coding-client",
		Usage:  "coding oauth2 client id",
	},
	cli.StringFlag{
		EnvVar: "DRONE_CODING_SECRET,WOODPECKER_CODING_SECRET",
		Name:   "coding-secret",
		Usage:  "coding oauth2 client secret",
	},
	cli.StringSliceFlag{
		EnvVar: "DRONE_CODING_SCOPE,WOODPECKER_CODING_SCOPE",
		Name:   "coding-scope",
		Usage:  "coding oauth scope",
		Value: &cli.StringSlice{
			"user",
			"project",
			"project:depot",
		},
	},
	cli.StringFlag{
		EnvVar: "DRONE_CODING_GIT_MACHINE,WOODPECKER_CODING_GIT_MACHINE",
		Name:   "coding-git-machine",
		Usage:  "coding machine name",
		Value:  "git.coding.net",
	},
	cli.StringFlag{
		EnvVar: "DRONE_CODING_GIT_USERNAME,WOODPECKER_CODING_GIT_USERNAME",
		Name:   "coding-git-username",
		Usage:  "coding machine user username",
	},
	cli.StringFlag{
		EnvVar: "DRONE_CODING_GIT_PASSWORD,WOODPECKER_CODING_GIT_PASSWORD",
		Name:   "coding-git-password",
		Usage:  "coding machine user password",
	},
	cli.BoolFlag{
		EnvVar: "DRONE_CODING_SKIP_VERIFY,WOODPECKER_CODING_SKIP_VERIFY",
		Name:   "coding-skip-verify",
		Usage:  "coding skip ssl verification",
	},
	cli.DurationFlag{
		EnvVar: "DRONE_KEEPALIVE_MIN_TIME,WOODPECKER_KEEPALIVE_MIN_TIME",
		Name:   "keepalive-min-time",
		Usage:  "server-side enforcement policy on the minimum amount of time a client should wait before sending a keepalive ping.",
	},
	cli.BoolFlag{
		EnvVar: "WOODPECKER_GITEA_REV_PROXY_AUTH",
		Name:   "gitea-rev-proxy-auth",
		Usage:  "enable gitea authentication using HTTP header specified in WOODPECKER_GITEA_REV_PROXY_AUTH_HEADER",
	},
	cli.StringFlag{
		EnvVar: "WOODPECKER_GITEA_REV_PROXY_AUTH_HEADER",
		Name:   "gitea-rev-proxy-auth-header",
		Usage:  "HTTP header with gitea authenticated user login",
	},
}
