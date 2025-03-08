package main

import (
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/effective-security/x/ctl"
	"github.com/tbilicode/bogclient/internal/cli"
	"github.com/tbilicode/bogclient/internal/cli/account"
	"github.com/tbilicode/bogclient/internal/version"
)

type app struct {
	cli.Cli

	Account account.Cmd `cmd:"" help:"Account operations"`
}

func main() {
	realMain(os.Args, os.Stdout, os.Stderr, os.Exit)
}

func realMain(args []string, out io.Writer, errout io.Writer, exit func(int)) {
	cl := app{
		Cli: cli.Cli{
			Version: ctl.VersionFlag("0.1.1"),
		},
	}
	// cl.Cli.WithErrWriter(errout).
	// 	WithWriter(out)

	parser, err := kong.New(&cl,
		kong.Name("bog"),
		kong.Description("BOG client"),
		//kong.UsageOnError(),
		kong.Writers(out, errout),
		kong.Exit(exit),
		ctl.BoolPtrMapper,
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": version.Current().String(),
		})
	if err != nil {
		panic(err)
	}

	ctx, err := parser.Parse(args[1:])
	parser.FatalIfErrorf(err)

	if ctx != nil {
		err = ctx.Run(&cl.Cli)
		ctx.FatalIfErrorf(err)
	}
}
