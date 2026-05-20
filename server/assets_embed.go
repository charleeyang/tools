//go:build embed

package main

import (
	"embed"
	"io/fs"
)

// 发布构建 (go build -tags embed): 内嵌前端资源, 生成自包含单文件
// webembed/ 由 `make release` 在构建前从 ../web 复制生成
//
//go:embed all:webembed
var webembed embed.FS

func embeddedFS() fs.FS {
	sub, err := fs.Sub(webembed, "webembed")
	if err != nil {
		return nil
	}
	return sub
}
