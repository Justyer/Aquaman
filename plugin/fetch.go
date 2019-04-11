package plg

import (
	"fmt"
	"time"

	aqua "github.com/Justyer/Aquaman"
)

type Fetcher interface {
	GetData() []int
}

type Fetch struct {
	Middleware
}

func NewFetch() aqua.MiddleManager {
	return &Fetch{}
}

func (m *Fetch) Run(grt_idx int) {
	fmt.Printf("fetch_inchan:\n%#v\n%#v\n---\n", m, m.InChan)
	a := true
	time.AfterFunc(10*time.Second, func() {
		a = false
	})

	for a {
		x := &aqua.Carrior{
			Data: []byte("fetch"),
		}
		m.OutChan.Push(x)
		time.Sleep(time.Second)
	}

	m.Close()

}
func (m *Fetch) Instance(c *aqua.Carrior) Fetcher {
	return nil
}

// 中间件处理的业务
func (m *Fetch) Template(f Fetcher) {

}
