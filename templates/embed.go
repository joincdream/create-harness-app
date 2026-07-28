package templates

import (
	"embed"
)

// EmbeddedFS holds default blueprint templates
//go:embed default/*
var EmbeddedFS embed.FS
