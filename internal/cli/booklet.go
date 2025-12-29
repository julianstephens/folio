package cli

import (
	"github.com/julianstephens/folio/internal/pdf"
)

type Booklet struct {
	InputFileArg
	OutputFileArg
}

func (c *Booklet) AfterApply() error {
	return pdf.IsPDFFile(c.InputFileArg.File)
}

func (c *Booklet) Run() error {
	return nil
}
