#!/usr/bin/env python3

"""
Generate edge-case PDFs that are syntactically valid but should
trigger specific behavior in folio (e.g., rejection for booklet mode).
"""

import os

from reportlab.lib.pagesizes import A4, letter
from reportlab.pdfgen import canvas

OUT_DIR = os.path.join(os.path.dirname(__file__), "..", "pdf")


def ensure_out_dir():
    os.makedirs(OUT_DIR, exist_ok=True)


def blank_one_page_pdf():
    path = os.path.join(OUT_DIR, "blank_one_page.pdf")
    c = canvas.Canvas(path, pagesize=letter)
    c.showPage()
    c.save()
    print(f"wrote {path}")


def mixed_sizes_pdf():
    path = os.path.join(OUT_DIR, "mixed_sizes.pdf")
    c = canvas.Canvas(path, pagesize=letter)

    # Page 1 — Letter
    c.drawString(72, 720, "Letter-sized page")
    c.showPage()

    # Page 2 — A4
    c.setPageSize(A4)
    c.drawString(72, 800, "A4-sized page")
    c.showPage()

    c.save()
    print(f"wrote {path}")


def main():
    ensure_out_dir()
    blank_one_page_pdf()
    mixed_sizes_pdf()


if __name__ == "__main__":
    main()
