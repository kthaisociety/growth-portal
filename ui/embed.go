// Package ui embeds the built frontend (ui/dist) into the Go binary (docs/adr/0001). Run `npm run
// build` in ui/ first; until then only dist/.gitkeep is embedded and the server answers "frontend not
// built".
package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

// Dist is the built frontend, rooted at dist/.
func Dist() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err) // dist is embedded at compile time, so this can't fail
	}
	return sub
}
