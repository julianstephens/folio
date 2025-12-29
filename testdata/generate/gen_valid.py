#!/usr/bin/env python3

"""
Generate valid, clean PDFs with uniform page sizes and simple content.
Uses ReportLab to produce deterministic, well-formed PDFs.
"""

import os

from reportlab.lib.pagesizes import A4, letter
from reportlab.pdfgen import canvas

OUT_DIR = os.path.join(os.path.dirname(__file__), "..", "pdf")


def ensure_out_dir():
    os.makedirs(OUT_DIR, exist_ok=True)


def make_pdf(path, pages, pagesize, label):
    c = canvas.Canvas(path, pagesize=pagesize)
    width, height = pagesize
    for i in range(1, pages + 1):
        c.drawString(72, height - 72, f"{label} — Page {i}")
        c.showPage()
    c.save()
    print(f"wrote {path}")


def main():
    ensure_out_dir()

    make_pdf(
        os.path.join(OUT_DIR, "4p_letter.pdf"),
        pages=4,
        pagesize=letter,
        label="4p letter",
    )

    make_pdf(
        os.path.join(OUT_DIR, "5p_letter.pdf"),
        pages=5,
        pagesize=letter,
        label="5p letter",
    )

    make_pdf(
        os.path.join(OUT_DIR, "8p_letter.pdf"),
        pages=8,
        pagesize=letter,
        label="8p letter",
    )

    make_pdf(os.path.join(OUT_DIR, "4p_a4.pdf"), pages=4, pagesize=A4, label="4p A4")


if __name__ == "__main__":
    main()
