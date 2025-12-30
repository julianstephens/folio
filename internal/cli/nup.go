package cli

import (
	"fmt"
	"strings"

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

	// Validate orientation if provided (case-insensitive)
	if c.Orientation != "" {
		orientLower := strings.ToLower(c.Orientation)
		if orientLower != "portrait" && orientLower != "landscape" {
			return fmt.Errorf("orientation must be 'portrait', 'landscape', or empty (for auto-detection), got: %q", c.Orientation)
		}
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
