package globalsecret

import (
	"html/template"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/woodpecker-ci/woodpecker/cli/drone/internal"
)

var globalSecretListCmd = cli.Command{
	Name:   "ls",
	Usage:  "list global secrets",
	Action: globalSecretList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "format",
			Usage:  "format output",
			Value:  tmplSecretList,
			Hidden: true,
		},
	},
}

func globalSecretList(c *cli.Context) error {
	var format = c.String("format") + "\n"
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	list, err := client.GlobalSecretList()
	if err != nil {
		return err
	}
	tmpl, err := template.New("_").Funcs(secretFuncMap).Parse(format)
	if err != nil {
		return err
	}
	for _, registry := range list {
		tmpl.Execute(os.Stdout, registry)
	}
	return nil
}

// template for secret list items
var tmplSecretList = "\x1b[33m{{ .Name }} \x1b[0m" + `
Events: {{ list .Events }}
{{- if .Images }}
Images: {{ list .Images }}
{{- else }}
Images: <any>
{{- end }}
`

var secretFuncMap = template.FuncMap{
	"list": func(s []string) string {
		return strings.Join(s, ", ")
	},
}
