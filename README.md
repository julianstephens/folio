# folio

A command-line tool for preparing PDFs for printing as booklets or multi-up handouts.

## Overview

**folio** simplifies the process of transforming regular PDF documents into printer-ready formats:

- **Booklet layout**: Reorder and arrange pages for saddle-stitch or folded booklet printing
- **N-up layout**: Place multiple pages on a single sheet (2-up or 4-up) for handouts or review copies
- **PDF info**: Quickly inspect page count, dimensions, and orientation

Ideal for personal printing projects, zine creation, or preparing conference handouts.

## Features

- **Booklet creation**: Automatically reorders pages for booklet printing with configurable binding edge
- **N-up layouts**: Arrange 2 or 4 pages per sheet with automatic orientation detection
- **Flexible sheet sizes**: Supports common paper sizes (letter, A4) or auto-detection from source
- **Auto-padding**: Adds blank pages when needed to meet layout requirements
- **Simple CLI**: Straightforward commands with sensible defaults

## Installation

### From Source

**Prerequisites:**
- Go 1.25 or later

**Steps:**

```bash
# Clone the repository
git clone https://github.com/julianstephens/folio.git
cd folio

# Build the binary
make build

# The binary is now available at ./bin/folio
# Optionally, move it to your PATH
cp bin/folio /usr/local/bin/
```

## Usage

### Basic Commands

#### Get PDF Information

Display page count, dimensions, and orientation of a PDF file:

```bash
folio info <input.pdf>
```

**Example:**

```bash
folio info document.pdf
```

**Output:**

```
File: document.pdf
Pages: 8
Page size: 612 x 792 pt (US Letter 8.5x11in)
Orientation: portrait
```

#### Create a Booklet

Transform a PDF into a booklet layout suitable for saddle-stitch binding:

```bash
folio booklet <input.pdf> -o <output.pdf>
```

**Example:**

```bash
folio booklet document.pdf -o booklet.pdf
```

**Options:**
- `-o, --output`: Output file path (required)
- `--sheet-size`: Target sheet size (e.g., `letter`, `a4`)
- `--pad`: Page padding mode (`auto` or `none`, default: `auto`)
- `--binding-edge`: Binding edge orientation (`long` or `short`, default: `long`)

**Example with options:**

```bash
folio booklet document.pdf -o booklet.pdf --sheet-size=letter --binding-edge=long
```

#### Create N-up Layout

Arrange multiple pages per sheet for compact printing:

```bash
folio nup <input.pdf> -o <output.pdf> --per-sheet=2
```

**Example (2-up layout):**

```bash
folio nup document.pdf -o handout-2up.pdf --per-sheet=2
```

**Example (4-up layout):**

```bash
folio nup document.pdf -o handout-4up.pdf --per-sheet=4
```

**Options:**
- `-o, --output`: Output file path (required)
- `--per-sheet`: Number of pages per sheet (`2` or `4`, default: `2`)
- `--sheet-size`: Target sheet size (e.g., `letter`, `a4`)
- `--orientation`: Output orientation (`portrait`, `landscape`, or empty for auto-detection)
- `--pad`: Page padding mode (`auto` or `none`, default: `auto`)

### Complete Examples

See the [`examples/`](examples/) directory for sample PDFs and step-by-step command examples.

## Limitations and Known Caveats

- **PDF version support**: Works with standard PDF files; some features may not work with encrypted or heavily compressed PDFs
- **Page reordering only**: Does not modify page content, only arrangement and layout
- **Auto-padding**: When `--pad=auto` (default), blank pages are added to meet layout requirements (e.g., booklets need multiples of 4 pages)
- **Sheet sizes**: Limited to common paper sizes; custom dimensions are not yet supported
- **Binding edges**: Booklet binding edge options are `long` (default, for portrait sheets) and `short` (for landscape sheets)
- **No duplex configuration**: Output assumes your printer's duplex settings are correct for the layout (e.g., flip on short edge for booklets)

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Acknowledgments

- Built with [pdfcpu](https://github.com/pdfcpu/pdfcpu) for PDF processing
- CLI powered by [kong](https://github.com/alecthomas/kong)