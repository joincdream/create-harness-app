package templates

import (
	"embed"
)

// EmbeddedFS holds default blueprint templates
//go:embed all:default all:antigravity
var EmbeddedFS embed.FS
