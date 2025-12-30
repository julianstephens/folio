package cli

import (
	"github.com/julianstephens/folio/internal/pdf"
)

type Nup struct {
	InputFileArg
	OutputFileFlag
}

func (c *Nup) AfterApply() error {
	return pdf.IsPDFFile(c.InputFileArg.File)
}

func (c *Nup) Run() error {
	return nil
}
