package cli

import (
	"fmt"

	"github.com/julianstephens/folio/internal/pdf"
)

type Booklet struct {
	InputFileArg
	OutputFileArg
}

func (c *Booklet) AfterApply() error {
	valid, err := pdf.IsPDFFile(c.InputFileArg.File)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("file %q is not a valid PDF", c.InputFileArg.File)
	}
	return nil
}

func (c *Booklet) Run() error {
	return nil
}
