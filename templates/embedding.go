package templates

import "embed"

//go:embed index.html
var TemplateFS embed.FS
