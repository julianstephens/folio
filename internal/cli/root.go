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

type PerSheetFlag struct {
	PerSheet int `flag:"" name:"per-sheet" help:"Number of pages per sheet (2 or 4)." enum:"2,4" default:"2"`
}

type OrientationFlag struct {
	Orientation string `flag:"" name:"orientation" help:"Output orientation (portrait, landscape, or empty for auto-detection)." default:""`
}

type CLI struct {
	Globals
	Info    Info    `cmd:"" help:"Show information about a PDF file."`
	Nup     Nup     `cmd:"" help:"Apply n-up layout to a PDF file."`
	Booklet Booklet `cmd:"" help:"Create a booklet layout from a PDF file."`
}
