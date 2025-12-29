package cli

import (
	"errors"
	"fmt"

	pdftype "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"

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

	var dim *pdftype.Dim
	for d := range info.PageDimensions {
		if d.Height > 0 && d.Width > 0 {
			dim = &d
			break
		}
	}
	if dim == nil {
		return errors.New("could not determine page dimensions")
	}

	fmt.Println("PDF Information:")
	fmt.Printf("Source: %s\n", c.File)
	fmt.Printf("Page Count: %d\n", info.PageCount)
	fmt.Printf("Page Size: %.2fx%.2f\n", dim.Width, dim.Height)
	orientation := "Portrait"
	if dim.Landscape() {
		orientation = "Landscape"
	}
	fmt.Printf("Orientation: %s\n", orientation)

	return nil
}
