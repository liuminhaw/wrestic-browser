package static

import "embed"

//go:embed *.css
//go:embed js/*.js
var FS embed.FS
