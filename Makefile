.PHONY: test build dev prepare-frontend

# Detect OS: set platform-specific variables
ifeq ($(OS),Windows_NT)
    MKDIR = if not exist "frontend\dist" mkdir "frontend\dist"
    MKDIR_BUILD = if not exist "build" mkdir "build"
    CP    = copy /Y index.html "frontend\dist\index.html"
    CP_ICON = if exist "appicon.png" copy /Y appicon.png "build\appicon.png"
    CP_FAVICON = if exist "appicon.png" copy /Y appicon.png "frontend\dist\appicon.png"
    WAILS_TAGS =
else
    MKDIR = mkdir -p frontend/dist
    MKDIR_BUILD = mkdir -p build
    CP    = cp index.html frontend/dist/index.html
    CP_ICON = [ -f appicon.png ] && cp appicon.png build/appicon.png || true
    CP_FAVICON = [ -f appicon.png ] && cp appicon.png frontend/dist/appicon.png || true
    WAILS_TAGS = -tags webkit2_41
endif

test:
	go test -tags headless ./...

prepare-frontend:
	$(MKDIR)
	$(MKDIR_BUILD)
	$(CP)
	$(CP_ICON)
	$(CP_FAVICON)

build: prepare-frontend
	wails build $(WAILS_TAGS) -clean

dev: prepare-frontend
	wails dev $(WAILS_TAGS)
