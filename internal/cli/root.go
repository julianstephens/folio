package cli

import "github.com/alecthomas/kong"

type Globals struct {
	Version kong.VersionFlag `name:"version" help:"Show application version."`
}

type InputFileArg struct {
	File string `arg:"" name:"file" help:"Input PDF file." type:"existingfile"`
}

type OutputFileArg struct {
	File string `arg:"" name:"output" help:"Output PDF file." type:"path"`
}

type CLI struct {
	Globals
	Info    Info    `cmd:"" help:"Show information about a PDF file."`
	Nup     Nup     `cmd:"" help:"Apply n-up layout to a PDF file."`
	Booklet Booklet `cmd:"" help:"Create a booklet layout from a PDF file."`
}
