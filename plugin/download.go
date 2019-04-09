package mw

import (
	mw "Gamora/src/Aquaman/base/middleware"
	"fmt"
)

type Downloader interface {
	DownloadCover()
	DownloadVideo()
}

type Download struct {
	Middleware
}

func NewDownload() mw.MiddleManager {
	return &Download{}
}

func (m *Download) Run(grt_idx int) {
	fmt.Printf("down_inchan:\n%#v\n%#v\n%d\n%d\n---\n", m, m.IN_CHL, len(m.IN_CHL), cap(m.IN_CHL))
	for {
		var c *mw.Carrior
		isUse, ok := false, false
		select {
		case c, ok = <-m.IN_CHL:
			isUse = ok || isUse
		case c, ok = <-m.IN2_CHL:
			isUse = ok || isUse
		}
		if isUse {
			ins := m.Instance(c)
			m.Template(ins)
		} else {
			fmt.Println("break")
			break
		}

	}
	m.Close()
}
func (m *Download) Instance(c *mw.Carrior) Downloader {
	return nil
}

// 中间件处理的业务
func (m *Download) Template(d Downloader) {

}