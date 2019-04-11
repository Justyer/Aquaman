package main

import (
	"fmt"
	"time"

	aqua "github.com/Justyer/Aquaman"
	plg "github.com/Justyer/Aquaman/plugin"
)

func main() {
	// f, err := os.Create("cpu")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	// f2, err := os.Create("heap")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.WriteHeapProfile(f2)

	mwm := aqua.NewMWManager()

	// 新建中间件节点
	fetch_node := aqua.NewMWNode("fetch", plg.NewFetch, 1, 100)
	download_node := aqua.NewMWNode("download", plg.NewDownload, 1, 100)
	transfer_node := aqua.NewMWNode("transfer", plg.NewTransfer, 1, 100)

	// 将中间件插入到前一个中间件上
	fetch_node.NextNode(download_node)
	download_node.NextNode(transfer_node)

	// 注册业务线
	mwm.RegisterTXL("downloader", fetch_node)

	time.AfterFunc(3*time.Second, func() {
		fmt.Println("\n\n\n\n-----------")
		mwm.ServiceFinder("downloader", "1")
		fmt.Println("-----------")
		mwm.DropMW("downloader", "download")
		fmt.Println("-----------")
		mwm.ServiceFinder("downloader", "2")
		fmt.Printf("-----------\n\n\n\n")
	})

	// 启动相应业务线
	mwm.ExecuteByName("downloader")
}
