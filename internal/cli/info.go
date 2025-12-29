package cli

import (
	"errors"
	"fmt"

	"github.com/julianstephens/folio/internal/pdf"
	pdftype "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

type Info struct {
	InputFileArg
}

func (c *Info) AfterApply() error {
	valid, err := pdf.IsPDFFile(c.File)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("file %q is not a valid PDF", c.File)
	}

	return nil
}

func (c *Info) Run() error {
	return nil
}
