package tests

const PROJECT_ROOT = "../../"

const TESTDATA_DIR = PROJECT_ROOT + "testdata/pdf"

var (
	VALID_PDFS = []string{
		"4p_a4.pdf",
		"4p_letter.pdf",
		"5p_letter.pdf",
		"8p_letter.pdf",
		"blank_one_page.pdf",
	}

	INVALID_PDFS = []string{
		"empty.pdf",
		"encrypted.pdf",
		"mixed_sizes.pdf",
		"not_a_pdf.pdf",
		"truncated.pdf",
	}
)
