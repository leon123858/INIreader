.PHONY: test build dev prepare-frontend

# Detect OS: set platform-specific variables
ifeq ($(OS),Windows_NT)
    SHELL := cmd.exe
    MKDIR = if not exist "frontend\dist" mkdir "frontend\dist"
    MKDIR_BUILD = if not exist "build" mkdir "build"
    MKDIR_WIX = if not exist "build\windows\installer" mkdir "build\windows\installer"
    CP    = copy /Y index.html "frontend\dist\index.html"
    CP_ICON = if exist "appicon.png" copy /Y appicon.png "build\appicon.png"
    CP_FAVICON = if exist "appicon.png" copy /Y appicon.png "frontend\dist\appicon.png"
    CP_WIX = if exist "Product.wxs" copy /Y Product.wxs "build\windows\installer\Product.wxs"
    WAILS_TAGS =
else
    MKDIR = mkdir -p frontend/dist
    MKDIR_BUILD = mkdir -p build
    MKDIR_WIX = mkdir -p build/windows/installer
    CP    = cp index.html frontend/dist/index.html
    CP_ICON = [ -f appicon.png ] && cp appicon.png build/appicon.png || true
    CP_FAVICON = [ -f appicon.png ] && cp appicon.png frontend/dist/appicon.png || true
    CP_WIX = [ -f Product.wxs ] && cp Product.wxs build/windows/installer/Product.wxs || true
    WAILS_TAGS = -tags webkit2_41
endif

test:
	go test -tags headless ./...

prepare-frontend:
	$(MKDIR)
	$(MKDIR_BUILD)
	$(MKDIR_WIX)
	$(CP)
	$(CP_ICON)
	$(CP_FAVICON)
	$(CP_WIX)

build: prepare-frontend
	wails build $(WAILS_TAGS) -clean

dev: prepare-frontend
	wails dev $(WAILS_TAGS)

release:
	git tag v0.0.2
	git push origin v0.0.2
