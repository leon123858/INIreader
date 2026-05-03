.PHONY: test build dev prepare-frontend

GOCACHE ?= $(CURDIR)/.cache/go-build
export GOCACHE

test:
	go test -tags headless ./...

prepare-frontend:
	mkdir -p frontend/dist
	cp index.html frontend/dist/index.html

build: prepare-frontend
	wails build -tags webkit2_41 -clean

dev: prepare-frontend
	wails dev -tags webkit2_41
