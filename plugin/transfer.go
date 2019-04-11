package plg

import (
	"fmt"

	aqua "github.com/Justyer/Aquaman"
)

type Transferer interface {
	TransferCover()
	TransferVideo()
}

type Transfer struct {
	Middleware
}

func NewTransfer() aqua.MiddleManager {
	return &Transfer{}
}

func (m *Transfer) Run(grt_idx int) {
	fmt.Printf("trans_inchan:\n%#v\n%#v\n%#v\n%#v\n---\n", m, m.InChan, m.InChan.CHL, m.InChan.CHL2)
	m.Pop(func(c *aqua.Carrior) {
		ins := m.Instance(c)
		m.Template(ins)
	}, "trans")
	m.Close()
}
func (m *Transfer) Instance(c *aqua.Carrior) Transferer {
	fmt.Println("transfer", string(c.Data))
	return nil
}

// 中间件处理的业务
func (m *Transfer) Template(d Transferer) {

}
