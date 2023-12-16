package web

import "embed"

//go:embed index.html ovpnmgr.ico css img js
var Static embed.FS
