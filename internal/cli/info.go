package cli

import (
	"errors"
	"fmt"

	"github.com/julianstephens/folio/internal/pdf"
)

type Info struct {
	InputFileArg
}

func (c *Info) AfterApply() error {
	return pdf.IsPDFFile(c.File)
}

func (c *Info) Run() error {
	reader, err := pdf.GetPDFReader(c.File)
	if err != nil {
		return err
	}

	info, err := pdf.GetPDFInfo(reader, c.File)
	if err != nil {
		return err
	}

	if len(info.PageBoundaries) != info.PageCount {
		return errors.New("page boundaries count does not match page count")
	}

	firstPageBoundaries := info.PageBoundaries[0]

	box, err := pdf.ExtractPageBox(firstPageBoundaries)
	if err != nil {
		return err
	}

	fmt.Printf("File: %s\n", c.File)
	fmt.Printf("Pages: %d\n", info.PageCount)
	fmt.Printf("Page size: %.0f x %.0f pt (%s) \n", box.Width, box.Height, box.Size)
	fmt.Printf("Orientation: %s\n", box.Orientation)

	return nil
}
