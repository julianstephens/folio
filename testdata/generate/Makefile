.PHONY: all valid special invalid clean

PYTHON = python3
OUTDIR = ../pdf

all: valid special invalid encrypted

valid:
	$(PYTHON) gen_valid.py

special:
	$(PYTHON) gen_special.py

invalid:
	$(PYTHON) gen_invalid.py

encrypted:
	qpdf --encrypt user owner 128 -- ../pdf/4p_letter.pdf ../pdf/encrypted.pdf

clean:
	rm -f $(OUTDIR)/*.pdf
