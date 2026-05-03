.PHONY: test build dev prepare-frontend

# Detect OS: set platform-specific variables
ifeq ($(OS),Windows_NT)
    MKDIR = if not exist "frontend\dist" mkdir "frontend\dist"
    CP    = copy /Y index.html "frontend\dist\index.html"
    CP_ICON = if exist "iniicon.png" copy /Y iniicon.png "build\iniicon.png"
    CP_FAVICON = if exist "iniicon.png" copy /Y iniicon.png "frontend\dist\iniicon.png"
    WAILS_TAGS =
else
    MKDIR = mkdir -p frontend/dist
    CP    = cp index.html frontend/dist/index.html
    CP_ICON = [ -f iniicon.png ] && cp iniicon.png build/iniicon.png || true
    CP_FAVICON = [ -f iniicon.png ] && cp iniicon.png frontend/dist/iniicon.png || true
    WAILS_TAGS = -tags webkit2_41
endif

test:
	go test -tags headless ./...

prepare-frontend:
	$(MKDIR)
	$(CP)
	$(CP_ICON)
	$(CP_FAVICON)

build: prepare-frontend
	wails build $(WAILS_TAGS) -clean

dev: prepare-frontend
	wails dev $(WAILS_TAGS)
