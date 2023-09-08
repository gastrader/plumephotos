package templates

import "embed"

//go:embed *
var FS embed.FS

//will stick all files inside embed.FS file system