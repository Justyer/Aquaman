package main

import (
	mw "Gamora/src/Aquaman/base/middleware"
	mwp "Gamora/src/Aquaman/base/middleware/plugin"
	"log"
	"os"

	"runtime/pprof"
)

func main() {
	f, err := os.Create("cpu")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	f2, err := os.Create("heap")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f2)

	mwm := mw.NewMiddlewareManager()

	// 新建中间件节点
	fetch_node := mw.NewMWNode("fetch", mwp.NewFetch, 2, 100)
	download_node := mw.NewMWNode("download", mwp.NewDownload, 2, 100)

	// 将中间件插入到前一个中间件上
	fetch_node.NextNode(download_node)

	// 注册业务线
	mwm.RegisterTXL("downloader", fetch_node)

	// 启动相应业务线
	mwm.ExecuteByName("downloader")
}
