package globalsecret

import (
	"github.com/urfave/cli"

	"github.com/woodpecker-ci/woodpecker/cli/drone/internal"
)

var globalSecretDeleteCmd = cli.Command{
	Name:   "rm",
	Usage:  "remove a global secret",
	Action: globalSecretDelete,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "global secret name",
		},
	},
}

func globalSecretDelete(c *cli.Context) error {
	var secret = c.String("name")
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	return client.GlobalSecretDelete(secret)
}
