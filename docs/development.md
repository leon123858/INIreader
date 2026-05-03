# INIreader development

## Local prerequisites

- Go 1.22 or newer.
- Node.js 22 or newer.
- Wails CLI v2.9.0: `go install github.com/wailsapp/wails/v2/cmd/wails@v2.9.0`
- Linux Wails dependencies for GUI builds on your distro. On Ubuntu 24.04 or current GitHub runners, install `libgtk-3-dev libwebkit2gtk-4.1-dev pkg-config`.

## Common commands

```sh
go test -tags headless ./...
wails dev
wails build -clean
```

The `headless` test tag excludes Wails GUI bindings so backend INI parsing and saving can be validated on Linux without GTK/WebKitGTK installed.

Windows MSI packages are built in GitHub Actions on `windows-latest` with WiX Toolset. Tagged pushes like `v0.1.0` publish the generated MSI to a GitHub Release.
