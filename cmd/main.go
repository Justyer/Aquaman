package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	aqua "github.com/Justyer/Aquaman"
	plg "github.com/Justyer/Aquaman/plugin"
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

	mwm := aqua.NewMWManager()

	// 新建中间件节点
	fetch_node := aqua.NewMWNode("fetch", plg.NewFetch, 1, 100)
	download_node := aqua.NewMWNode("download", plg.NewDownload, 1, 100)
	transfer_node := aqua.NewMWNode("transfer", plg.NewTransfer, 1, 100)

	// 将中间件插入到前一个中间件上
	fetch_node.NextNode(download_node)
	download_node.NextNode(transfer_node)

	// 注册业务线
	mwm.Register(fetch_node)

	// time.AfterFunc(3*time.Second, func() {
	// 	mwm.MWIter("1")
	// 	err := mwm.DropMW("download")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	mwm.MWIter("2")
	// })

	time.AfterFunc(3*time.Second, func() {
		mwm.MWIter("1")
		record_node := aqua.NewMWNode("record", plg.NewStorage, 1, 100)
		err := mwm.InsertMWBack("transfer", record_node)
		if err != nil {
			fmt.Println(err)
		}
		mwm.MWIter("2")
	})

	// 启动相应业务线
	mwm.ExecuteByName()
}
