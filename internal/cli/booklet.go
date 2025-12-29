package cli

import (
	"github.com/julianstephens/folio/internal/pdf"
)

type Booklet struct {
	InputFileArg
	OutputFileFlag
	SheetSizeFlag
	PadFlag
	BindingEdge string `help:"Binding edge orientation (long, short)." enum:"long,short" default:"long"`
}

func (c *Booklet) AfterApply() error {
	return pdf.IsPDFFile(c.InputFileArg.File)
}

func (c *Booklet) Run() error {
	return pdf.CreateBooklet(c.InputFileArg.File, c.OutputFileFlag.File, pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: c.SheetSize,
			Pad:       c.Pad,
		},
		BindingEdge: c.BindingEdge,
	})
}
