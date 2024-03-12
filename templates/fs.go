package templates

import "embed"

//go:embed *.gohtml
//go:embed */*.html
//go:embed */*.gohtml
var FS embed.FS
