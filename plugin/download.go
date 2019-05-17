package plg

import (
	"fmt"

	aqua "github.com/Justyer/Aquaman"
)

type Downloader interface {
	DownloadCover()
	DownloadVideo()
}

type Download struct {
	Middleware
}

func NewDownload() aqua.MiddleManager {
	return &Download{}
}

func (m *Download) Run(grt_idx int) {
	fmt.Printf("down_inchan:\n%#v\n%#v\n%#v\n%#v\n---\n", m, m.InChan, m.InChan.CHL, m.InChan.CHL2)
	m.Pop(func(c *aqua.Carrior) {
		ins := m.Instance(c)
		m.Template(ins)
		m.OutChan.Push(&aqua.Carrior{
			Data: []byte("download"),
		})
	})

	m.Close()
}
func (m *Download) Instance(c *aqua.Carrior) Downloader {
	fmt.Println("download", string(c.Data))
	return nil
}

// 中间件处理的业务
func (m *Download) Template(d Downloader) {

}
