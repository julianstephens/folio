package pdf

import "math"

// Tolerance in points for size detection
const tolerance = 2

// Standard page sizes in points (1 point = 1/72 inch)
const LetterWidth = 612
const LetterHeight = 792
const A4Width = 595
const A4Height = 842

// Orientation indicates whether the page is in portrait or landscape mode.
type Orientation string

var Orientations = struct {
	Landscape Orientation
	Portrait  Orientation
}{
	Landscape: "landscape",
	Portrait:  "portrait",
}

// KnownSize represents standard paper sizes.
type KnownSize string

var KnownSizes = struct {
	Letter KnownSize
	A4     KnownSize
	Custom KnownSize
}{
	Letter: "US Letter 8.5x11in",
	A4:     "A4 210x297mm",
	Custom: "Custom",
}

// Dimensions returns the width and height in points for the KnownSize.
func (ks KnownSize) Dimensions() (width, height float64) {
	switch ks {
	case KnownSizes.Letter:
		return LetterWidth, LetterHeight
	case KnownSizes.A4:
		return A4Width, A4Height
	default:
		return 0, 0
	}
}

// PageBox holds information about a PDF page's dimensions and size.
type PageBox struct {
	Width       float64
	Height      float64
	Size        KnownSize
	Orientation Orientation
}

func approxEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// DetectKnownSize checks if the given width and height correspond to a known size.
func DetectKnownSize(width, height float64) KnownSize {
	// normalize to portrait for comparison
	w, h := width, height
	if w > h {
		w, h = h, w
	}

	switch {
	case approxEqual(w, LetterWidth, tolerance) && approxEqual(h, LetterHeight, tolerance):
		return KnownSizes.Letter

	case approxEqual(w, A4Width, tolerance) && approxEqual(h, A4Height, tolerance):
		return KnownSizes.A4

	default:
		return KnownSizes.Custom
	}
}

// DetectOrientation determines the orientation based on width and height.
func DetectOrientation(width, height float64) Orientation {
	if width > height {
		return Orientations.Landscape
	}
	return Orientations.Portrait
}
