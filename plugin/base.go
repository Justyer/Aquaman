package plg

import (
	"fmt"
	"log"
	"sync/atomic"

	aqua "github.com/Justyer/Aquaman"
)

// 中间件
type Middleware struct {
	Node    *aqua.MWNode
	InChan  *aqua.Chan
	OutChan *aqua.Chan
}

func NewMiddleware() aqua.MiddleManager {
	return &Middleware{}
}

func (m *Middleware) SetMWNode(mwn *aqua.MWNode) {
	m.Node = mwn
}

func (m *Middleware) SetInChan(c *aqua.Chan) {
	m.InChan = c
}

func (m *Middleware) SetOutChan(c *aqua.Chan) {
	m.OutChan = c
}

func (m *Middleware) GetInChan() *aqua.Chan {
	return m.InChan
}

func (m *Middleware) GetOutChan() *aqua.Chan {
	return m.OutChan
}

func (m *Middleware) Pop(f func(*aqua.Carrior), nm string) {
	for {
		c := m.InChan
		if c == nil {
			fmt.Println("break2", nm)
			break
		}
		if c.CHL == nil && c.CHL2 == nil {
			fmt.Println("break blank bil", nm)
			break
		}
		var cr *aqua.Carrior
		var isUse bool
		select {
		case cr, isUse = <-c.CHL:
			if !isUse {
				c.CHL = nil
			}
		case cr, isUse = <-c.CHL2:
			if !isUse {
				c.CHL2 = nil
			}
		}

		switch {
		case isUse:
			f(cr)
		case !isUse && (c.CHL == nil && c.CHL2 == nil):
			break
		}
	}
}

func (m *Middleware) Run(i int) {
	log.Fatal("Must implement Run() func")
}

func (m *Middleware) Close() {
	if m.OutChan == nil {
		fmt.Printf("Closedddddd %p OutChan\n", m.OutChan)
		return
	}
	closeNum := atomic.AddInt64(&m.Node.ClosedGRTNum, 1)
	if closeNum == int64(m.Node.GRT_NUM) {
		fmt.Printf("Close %p OutChan\n", m.OutChan)
		m.OutChan.Close()
	}
}
