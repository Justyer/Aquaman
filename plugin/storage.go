package plg

import (
	"fmt"

	aqua "github.com/Justyer/Aquaman"
)

type Storager interface {
	StorageCover()
	StorageVideo()
}

type Storage struct {
	Middleware
}

func NewStorage() aqua.MiddleManager {
	return &Storage{}
}

func (m *Storage) Run(grt_idx int) {
	fmt.Printf("store_inchan:\n%#v\n%#v\n%#v\n%#v\n---\n", m, m.InChan, m.InChan.CHL, m.InChan.CHL2)
	m.Pop(func(c *aqua.Carrior) {
		ins := m.Instance(c)
		m.Template(ins)
	}, "store")
	m.Close()
}
func (m *Storage) Instance(c *aqua.Carrior) Storager {
	fmt.Println("storage", string(c.Data))
	return nil
}

// 中间件处理的业务
func (m *Storage) Template(d Storager) {

}
