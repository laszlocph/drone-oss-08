package globalsecret

import "github.com/urfave/cli"

// Command exports the secret command.
var Command = cli.Command{
	Name:  "globalsecret",
	Usage: "manage global secrets",
	Subcommands: []cli.Command{
		globalSecretCreateCmd,
		globalSecretDeleteCmd,
		globalSecretUpdateCmd,
		globalSecretInfoCmd,
		globalSecretListCmd,
	},
}
