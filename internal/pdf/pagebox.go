package pdf

import "math"

const tolerance = 2

const LetterWidth = 612
const LetterHeight = 792
const A4Width = 595
const A4Height = 842

type Orientation string

var Orientation_ = struct {
	Landscape Orientation
	Portrait  Orientation
}{
	Landscape: "landscape",
	Portrait:  "portrait",
}

type KnownSize string

var KnownSize_ = struct {
	Letter KnownSize
	A4     KnownSize
	Custom KnownSize
}{
	Letter: "US Letter 8.5x11in",
	A4:     "A4 210x297mm",
	Custom: "Custom",
}

func (ks KnownSize) Dimensions() (width, height float64) {
	switch ks {
	case KnownSize_.Letter:
		return LetterWidth, LetterHeight
	case KnownSize_.A4:
		return A4Width, A4Height
	default:
		return 0, 0
	}
}

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
	switch {
	case approxEqual(width, LetterWidth, tolerance) && approxEqual(height, LetterHeight, tolerance):
		return KnownSize_.Letter

	case approxEqual(width, A4Width, tolerance) && approxEqual(height, A4Height, tolerance):
		return KnownSize_.A4

	default:
		return KnownSize_.Custom
	}
}

// DetectOrientation determines the orientation based on width and height.
func DetectOrientation(width, height float64) Orientation {
	if width > height {
		return Orientation_.Landscape
	}
	return Orientation_.Portrait
}
