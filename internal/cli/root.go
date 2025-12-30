package cli

import "github.com/alecthomas/kong"

type Globals struct {
	Version kong.VersionFlag `name:"version" help:"Show application version."`
}

type InputFileArg struct {
	File string `arg:"" name:"file" help:"Input PDF file." type:"existingfile"`
}

type OutputFileFlag struct {
	File string `flag:"" name:"output" short:"o" help:"Output PDF file." type:"path"`
}

type SheetSizeFlag struct {
	SheetSize string `flag:"" name:"sheet-size" help:"Sheet size (e.g., letter, a4)." optional:""`
}

type PadFlag struct {
	Pad string `flag:"" name:"pad" help:"Page padding mode (auto, none). Auto pads to multiple of 4." enum:"auto,none" default:"auto"`
}

type CLI struct {
	Globals
	Info    Info    `cmd:"" help:"Show information about a PDF file."`
	Nup     Nup     `cmd:"" help:"Apply n-up layout to a PDF file."`
	Booklet Booklet `cmd:"" help:"Create a booklet layout from a PDF file."`
}
