# Folio Examples

This directory contains sample files and example commands to help you get started with **folio**.

## Sample PDF

`sample.pdf` is an 8-page US Letter document that can be used to test booklet and n-up layouts.

To generate additional test PDFs, see [`../testdata/generate/`](../testdata/generate/).

## Example Commands

All examples assume you have built the `folio` binary and it's available in your PATH or you're running from the project root.

### 1. Get PDF Information

Display information about the sample PDF:

```bash
folio info examples/sample.pdf
```

**Expected output:**

```
File: examples/sample.pdf
Pages: 8
Page size: 612 x 792 pt (US Letter 8.5x11in)
Orientation: portrait
```

### 2. Create a Booklet

Transform the 8-page document into a booklet layout suitable for printing and folding:

```bash
folio booklet examples/sample.pdf -o examples/booklet-output.pdf
```

**What this does:**
- Reorders pages for booklet printing (page order: 8,1,2,7,6,3,4,5)
- Creates a 4-sheet (8-page) output
- Pages are arranged so you can print double-sided and fold in the middle

**To print the booklet:**
1. Print `booklet-output.pdf` double-sided with "flip on short edge"
2. Fold the printed sheets in half
3. Staple along the spine (binding edge)

### 3. Create a 2-up Layout

Arrange two pages per sheet for compact handouts:

```bash
folio nup examples/sample.pdf -o examples/handout-2up.pdf --per-sheet=2
```

**What this does:**
- Places 2 pages side-by-side on each sheet
- Automatically switches to landscape orientation
- Results in 4 sheets (8 pages → 4 sheets in 2-up)

**Common use cases:**
- Handouts for meetings or presentations
- Review copies that save paper
- Quick reference sheets

### 4. Create a 4-up Layout

Arrange four pages per sheet for maximum compactness:

```bash
folio nup examples/sample.pdf -o examples/handout-4up.pdf --per-sheet=4
```

**What this does:**
- Places 4 pages in a 2×2 grid on each sheet
- Keeps portrait orientation
- Results in 2 sheets (8 pages → 2 sheets in 4-up)

**Common use cases:**
- Proof copies for review
- Archive copies to save storage space
- Quick overview of document content

### 5. Booklet with Custom Sheet Size

Create a booklet on A4 paper:

```bash
folio booklet examples/sample.pdf -o examples/booklet-a4.pdf --sheet-size=a4
```

**Note:** This will scale US Letter content to fit A4 sheets.

### 6. N-up with Specific Orientation

Force a specific orientation for the n-up output:

```bash
folio nup examples/sample.pdf -o examples/handout-portrait.pdf --per-sheet=4 --orientation=portrait
```

## Tips

1. **Auto-padding**: By default, folio adds blank pages when needed. For example, a 5-page document becomes 8 pages for booklet printing. Use `--pad=none` to disable this behavior.

2. **Output filenames**: If you don't specify `--output`, the command will fail. Always provide an output filename.

3. **Sheet sizes**: Common values for `--sheet-size` include:
   - `letter` (US Letter: 8.5" × 11")
   - `a4` (A4: 210mm × 297mm)
   - If omitted, the tool uses the source PDF's page size

4. **Binding edge**: For booklets:
   - `--binding-edge=long` (default): Binding along the long edge, suitable for portrait sheets
   - `--binding-edge=short`: Binding along the short edge, suitable for landscape sheets

## Verifying Output

After creating any output file, you can verify it with the `info` command:

```bash
folio info examples/booklet-output.pdf
```

This helps confirm that the page count and dimensions are as expected.

## Cleaning Up

To remove generated output files:

```bash
rm examples/*-output.pdf examples/*-2up.pdf examples/*-4up.pdf examples/*-a4.pdf examples/*-portrait.pdf
```
