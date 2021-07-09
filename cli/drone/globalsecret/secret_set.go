package globalsecret

import (
	"io/ioutil"
	"strings"

	"github.com/woodpecker-ci/woodpecker/cli/drone/internal"
	"github.com/woodpecker-ci/woodpecker/drone-go/drone"

	"github.com/urfave/cli"
)

var globalSecretUpdateCmd = cli.Command{
	Name:   "update",
	Usage:  "update a global secret",
	Action: globalSecretUpdate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "global secret name",
		},
		cli.StringFlag{
			Name:  "value",
			Usage: "global secret value",
		},
		cli.StringSliceFlag{
			Name:  "event",
			Usage: "global secret limited to these events",
		},
		cli.StringSliceFlag{
			Name:  "image",
			Usage: "global secret limited to these images",
		},
	},
}

func globalSecretUpdate(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	secret := &drone.Secret{
		Name:   c.String("name"),
		Value:  c.String("value"),
		Images: c.StringSlice("image"),
		Events: c.StringSlice("event"),
	}
	if strings.HasPrefix(secret.Value, "@") {
		path := strings.TrimPrefix(secret.Value, "@")
		out, ferr := ioutil.ReadFile(path)
		if ferr != nil {
			return ferr
		}
		secret.Value = string(out)
	}
	_, err = client.GlobalSecretUpdate(secret)
	return err
}
