package templates

import "embed"

//go:embed *.html
var Staticfiles embed.FS
