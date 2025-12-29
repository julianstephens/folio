package main

import (
	"github.com/alecthomas/kong"

	"github.com/julianstephens/folio/internal/cli"
)

func main() {
	cli := cli.CLI{}

	kongCtx := kong.Parse(
		&cli,
		kong.Name("folio"),
		kong.Description("A CLI for preparing PDFs for printing as booklets or multi-up handouts."),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": "0.1.0",
		},
	)

	err := kongCtx.Run()
	kongCtx.FatalIfErrorf(err)
}
