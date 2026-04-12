//go:build !embed

package static

import "io/fs"

// FS is nil when not built with -tags embed.
// The server falls back to CLIENT env var (proxy or static dir).
var FS fs.FS = nil
