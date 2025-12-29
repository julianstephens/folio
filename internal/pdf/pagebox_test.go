package pdf_test

import (
	"testing"

	"github.com/julianstephens/folio/internal/pdf"
)

func TestDetectKnownSize(t *testing.T) {
	tests := []struct {
		widthPt  float64
		heightPt float64
		expected pdf.KnownSize
	}{
		{612, 792, pdf.KnownSize_.Letter}, // 8.5 x 11 in
		{595, 842, pdf.KnownSize_.A4},     // 210 x 297 mm
		{0, 0, pdf.KnownSize_.Custom},     // Invalid size
		{500, 700, pdf.KnownSize_.Custom}, // Non-standard size
	}

	for _, tt := range tests {
		sizeName := pdf.DetectKnownSize(tt.widthPt, tt.heightPt)
		if sizeName != tt.expected {
			t.Errorf("DetectKnownSize(%f, %f) = %q; want %q", tt.widthPt, tt.heightPt, sizeName, tt.expected)
		}
	}
}

func TestDetermineOrientation(t *testing.T) {
	tests := []struct {
		widthPt  float64
		heightPt float64
		expected pdf.Orientation
	}{
		{612, 792, pdf.Orientation_.Portrait},  // Portrait
		{792, 612, pdf.Orientation_.Landscape}, // Landscape
		{600, 600, pdf.Orientation_.Portrait},  // Perfect square treated as Portrait
	}

	for _, tt := range tests {
		orientation := pdf.DetectOrientation(tt.widthPt, tt.heightPt)
		if orientation != tt.expected {
			t.Errorf("DetermineOrientation(%f, %f) = %q; want %q", tt.widthPt, tt.heightPt, orientation, tt.expected)
		}
	}
}
