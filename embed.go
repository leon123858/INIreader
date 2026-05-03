//go:build !headless

package main

import "embed"

// assets holds our frontend distribution files.  When running `wails build`,
// the frontend is built into the `frontend/dist` directory and then embedded
// into this binary.  The embed directive below matches all files under
// frontend/dist.

//go:embed all:frontend/dist
var assets embed.FS
