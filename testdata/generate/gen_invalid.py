#!/usr/bin/env python3

"""
Generate invalid / corrupted / non-PDF files for testing error paths.
"""

import os

OUT_DIR = os.path.join(os.path.dirname(__file__), "..", "pdf")


def ensure_out_dir():
    os.makedirs(OUT_DIR, exist_ok=True)


def make_empty_file():
    path = os.path.join(OUT_DIR, "empty.pdf")
    with open(path, "wb"):
        pass
    print(f"wrote empty file {path}")


def make_not_a_pdf():
    path = os.path.join(OUT_DIR, "not_a_pdf.pdf")
    with open(path, "w") as f:
        f.write("This is definitely not a PDF.")
    print(f"wrote text masquerading as PDF {path}")


def make_truncated_pdf():
    """
    Copy 4p_letter.pdf then truncate it in half.
    Requires that gen_valid.py has already been run.
    """
    src = os.path.join(OUT_DIR, "4p_letter.pdf")
    dst = os.path.join(OUT_DIR, "truncated.pdf")

    if not os.path.exists(src):
        print("missing 4p_letter.pdf; run gen_valid.py first")
        return

    data = None
    with open(src, "rb") as f:
        data = f.read()

    if data is None:
        print("failed to read 4p_letter.pdf")
        return

    half = len(data) // 2
    with open(dst, "wb") as f:
        f.write(data[:half])

    print(f"wrote truncated PDF {dst}")


def main():
    ensure_out_dir()
    make_empty_file()
    make_not_a_pdf()
    make_truncated_pdf()


if __name__ == "__main__":
    main()
