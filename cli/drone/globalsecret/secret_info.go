package globalsecret

import (
	"html/template"
	"os"

	"github.com/urfave/cli"

	"github.com/woodpecker-ci/woodpecker/cli/drone/internal"
)

var globalSecretInfoCmd = cli.Command{
	Name:   "info",
	Usage:  "display global secret info",
	Action: globalSecretInfo,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "name",
			Usage:    "global secret name",
			Required: true,
		},
		cli.StringFlag{
			Name:   "format",
			Usage:  "format output",
			Value:  tmplSecretList,
			Hidden: true,
		},
	},
}

func globalSecretInfo(c *cli.Context) error {
	var (
		secretName = c.String("name")
		format     = c.String("format") + "\n"
	)
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	secret, err := client.GlobalSecret(secretName)
	if err != nil {
		return err
	}
	tmpl, err := template.New("_").Funcs(secretFuncMap).Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, secret)
}
