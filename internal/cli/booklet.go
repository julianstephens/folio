package cli

import (
	"github.com/julianstephens/folio/internal/pdf"
)

type Booklet struct {
	InputFileArg
	OutputFileArg
	SheetSize   string `help:"Sheet size for output (letter, a4). If not provided, will be inferred from input." optional:""`
	BindingEdge string `help:"Binding edge orientation (long, short)." enum:"long,short" default:"long"`
	Duplex      string `help:"Duplex printing mode (long, short, none)." enum:"long,short,none" default:"long"`
	Pad         string `help:"Page padding mode (auto, none). Auto pads to multiple of 4." enum:"auto,none" default:"auto"`
}

func (c *Booklet) AfterApply() error {
	return pdf.IsPDFFile(c.InputFileArg.File)
}

func (c *Booklet) Run() error {
	return pdf.CreateBooklet(c.InputFileArg.File, c.OutputFileArg.File, pdf.BookletOptions{
		SheetSize:   c.SheetSize,
		BindingEdge: c.BindingEdge,
		Duplex:      c.Duplex,
		Pad:         c.Pad,
	})
}
