# PDF test corpus generator

This directory contains small Python scripts for generating
test PDFs for `folio`.

Generated files are written to: `testdata/pdf/`

## Requirements

- Python 3
- `reportlab` (`pip install reportlab`)
- optional: `qpdf` for encrypted PDFs

## Usage

```
make # generate all PDFs
make valid # only valid PDFs
make special # mixed sizes / blank page
make invalid # empty, truncated, not-a-pdf
```


## Generated files

### Valid PDFs

- `4p_letter.pdf`
- `5p_letter.pdf`
- `8p_letter.pdf`
- `4p_a4.pdf`

### Valid but edge case

- `blank_one_page.pdf`
- `mixed_sizes.pdf`

### Invalid

- `empty.pdf`
- `not_a_pdf.pdf`
- `truncated.pdf`

### Optional

- `encrypted.pdf` (requires qpdf)
