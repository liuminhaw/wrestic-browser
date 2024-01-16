package static

import "embed"

//go:embed *.css
//go:embed js/*.js
//go:embed config/*.json
var FS embed.FS
