//go:build !embed

package main

import "io/fs"

// 默认(开发)构建: 不内嵌前端, 从磁盘 -web 目录提供
func embeddedFS() fs.FS { return nil }
