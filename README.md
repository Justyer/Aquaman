# Aquaman

[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)

中间件管理器

## 简介
集于单个进程中的中间件管理，定义所需中间件，注册到管理器中，即可使用

## 使用

```go
package main

import (
	aqua "github.com/Justyer/Aquaman"
	plg "github.com/Justyer/Aquaman/plugin"
)

func main() {
	// 新建中间件节点
	fetch_node := aqua.NewMWNode("fetch", plg.NewFetch, 1, 100)
	download_node := aqua.NewMWNode("download", plg.NewDownload, 1, 100)
	transfer_node := aqua.NewMWNode("transfer", plg.NewTransfer, 1, 100)

	// 将中间件插入到前一个中间件上
	fetch_node.NextNode(download_node)
	download_node.NextNode(transfer_node)

	// 注册业务线
	mwm.Register(fetch_node)

	// 启动相应业务线
	mwm.Run()
}
```