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
	fmt.Printf("fetch_inchan:\n%#v\n%#v\n%d\n%d\n---\n", m, m.IN_CHL, len(m.IN_CHL), cap(m.IN_CHL))
	a := true
	time.AfterFunc(time.Second, func() {
		a = false
	})

	for a {
		for i := 0; i < 3; i++ {
			x := &aqua.Carrior{
				Data: []byte(fmt.Sprintf("%d", i)),
			}
			m.MAIN_OUT_CHL <- x
		}
	}

	m.Close()

}
func (m *Fetch) Instance(c *aqua.Carrior) Fetcher {
	return nil
}

// 中间件处理的业务
func (m *Fetch) Template(f Fetcher) {

}
