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
		{612, 792, pdf.KnownSizes.Letter}, // 8.5 x 11 in
		{792, 612, pdf.KnownSizes.Letter}, // Landscape Letter
		{595, 842, pdf.KnownSizes.A4},     // 210 x 297 mm
		{842, 595, pdf.KnownSizes.A4},     // Landscape A4
		{0, 0, pdf.KnownSizes.Custom},     // Invalid size
		{500, 700, pdf.KnownSizes.Custom}, // Non-standard size
	}

	for _, tt := range tests {
		sizeName := pdf.DetectKnownSize(tt.widthPt, tt.heightPt)
		if sizeName != tt.expected {
			t.Errorf("DetectKnownSize(%f, %f) = %q; want %q", tt.widthPt, tt.heightPt, sizeName, tt.expected)
		}
	}
}

func TestDetectOrientation(t *testing.T) {
	tests := []struct {
		widthPt  float64
		heightPt float64
		expected pdf.Orientation
	}{
		{612, 792, pdf.Orientations.Portrait},  // Portrait
		{792, 612, pdf.Orientations.Landscape}, // Landscape
		{600, 600, pdf.Orientations.Portrait},  // Perfect square treated as Portrait
	}

	for _, tt := range tests {
		orientation := pdf.DetectOrientation(tt.widthPt, tt.heightPt)
		if orientation != tt.expected {
			t.Errorf("DetectOrientation(%f, %f) = %q; want %q", tt.widthPt, tt.heightPt, orientation, tt.expected)
		}
	}
}
