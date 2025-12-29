package cli

import (
	"github.com/julianstephens/folio/internal/pdf"
)

type Booklet struct {
	InputFileArg
	OutputFileArg
	SheetSizeArg
	PadArg
	BindingEdge string `help:"Binding edge orientation (long, short)." enum:"long,short" default:"long"`
}

func (c *Booklet) AfterApply() error {
	return pdf.IsPDFFile(c.InputFileArg.File)
}

func (c *Booklet) Run() error {
	return pdf.CreateBooklet(c.InputFileArg.File, c.OutputFileArg.File, pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: c.SheetSizeArg.SheetSize,
			Pad:       c.PadArg.Pad,
		},
		BindingEdge: c.BindingEdge,
	})
}
