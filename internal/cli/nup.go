package cli

import (
	"fmt"

	"github.com/julianstephens/folio/internal/pdf"
)

type Nup struct {
	InputFileArg
	OutputFileFlag
	SheetSizeFlag
	OrientationFlag
	PerSheetFlag
	PadFlag
}

func (c *Nup) AfterApply() error {
	if err := pdf.IsPDFFile(c.InputFileArg.File); err != nil {
		return err
	}
	
	// Validate orientation if provided
	if c.Orientation != "" && c.Orientation != "portrait" && c.Orientation != "landscape" {
		return fmt.Errorf("orientation must be 'portrait', 'landscape', or empty (for auto-detection), got: %q", c.Orientation)
	}
	
	return nil
}

func (c *Nup) Run() error {
	return pdf.CreateNup(c.InputFileArg.File, c.OutputFileFlag.File, pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: c.SheetSize,
			Pad:       c.Pad,
		},
		PerSheet:    c.PerSheet,
		Orientation: c.Orientation,
	})
}
