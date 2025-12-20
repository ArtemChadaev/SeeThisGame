package migrate

import "embed"

// FS экспортирует все SQL файлы в этой папке
//
//go:embed *.sql
var FS embed.FS
