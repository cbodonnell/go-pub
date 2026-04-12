//go:build embed

package static

import (
	"embed"
	"io/fs"
)

//go:embed dist
var distFS embed.FS

// FS is the embedded client dist directory.
// Only populated when built with -tags embed (i.e. production/Docker builds).
var FS fs.FS

func init() {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic("static: failed to sub dist: " + err.Error())
	}
	FS = sub
}
