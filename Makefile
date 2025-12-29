
.PHONY: build	

build:
	go build -o bin/folio ./cmd/folio
	chmod +x bin/folio
