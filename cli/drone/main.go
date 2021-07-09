package main

import (
	"fmt"
	"os"

	"github.com/woodpecker-ci/woodpecker/cli/drone/build"
	"github.com/woodpecker-ci/woodpecker/cli/drone/deploy"
	"github.com/woodpecker-ci/woodpecker/cli/drone/exec"
	"github.com/woodpecker-ci/woodpecker/cli/drone/globalsecret"
	"github.com/woodpecker-ci/woodpecker/cli/drone/info"
	"github.com/woodpecker-ci/woodpecker/cli/drone/log"
	"github.com/woodpecker-ci/woodpecker/cli/drone/registry"
	"github.com/woodpecker-ci/woodpecker/cli/drone/repo"
	"github.com/woodpecker-ci/woodpecker/cli/drone/secret"
	"github.com/woodpecker-ci/woodpecker/cli/drone/user"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// drone version number
var version string

func main() {
	app := cli.NewApp()
	app.Name = "drone"
	app.Version = version
	app.Usage = "command line utility"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "t, token",
			Usage:  "server auth token",
			EnvVar: "DRONE_TOKEN,WOODPECKER_TOKEN",
		},

		cli.StringFlag{
			Name:   "s, server",
			Usage:  "server address",
			EnvVar: "DRONE_SERVER,WOODPECKER_SERVER",
		},
		cli.BoolFlag{
			Name:   "skip-verify",
			Usage:  "skip ssl verfification",
			EnvVar: "DRONE_SKIP_VERIFY,WOODPECKER_SKIP_VERIFY",
			Hidden: true,
		},
		cli.StringFlag{
			Name:   "socks-proxy",
			Usage:  "socks proxy address",
			EnvVar: "SOCKS_PROXY",
			Hidden: true,
		},
		cli.BoolFlag{
			Name:   "socks-proxy-off",
			Usage:  "socks proxy ignored",
			EnvVar: "SOCKS_PROXY_OFF",
			Hidden: true,
		},
	}
	app.Commands = []cli.Command{
		build.Command,
		log.Command,
		deploy.Command,
		exec.Command,
		info.Command,
		registry.Command,
		secret.Command,
		globalsecret.Command,
		repo.Command,
		user.Command,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
