package aqua

import (
	"errors"
	"fmt"
	"sync"
)

type MiddleManager interface {
	SetMWNode(mwn *MWNode)
	SetInChan(in *Chan)
	SetOutChan(out *Chan)
	GetInChan() *Chan
	GetOutChan() *Chan
	Run(int)
}

// 中间件管理器
type MiddlewareManager struct {
	TXL *MWNode
	wg  *sync.WaitGroup
}

func NewMWManager() *MiddlewareManager {
	return &MiddlewareManager{
		wg: &sync.WaitGroup{},
	}
}

// 注册业务线
func (mgr *MiddlewareManager) Register(mwn *MWNode) {
	mgr.TXL = mwn
}

// 根据业务线名字执行任务
func (mgr *MiddlewareManager) Run() {
	mgr.suffixChannel()
}

// 每个中间件创建的channel为前置
// func (mgr *MiddlewareManager) prefixChannel(n string) {
// 	wg := &sync.WaitGroup{}
// 	p := mgr.TXL[n]
// 	var pp *MWNode
// 	for {
// 		if p == nil {
// 			break
// 		}
// 		var mmrs []MiddleManager
// 		chl := NewChan()
// 		chl.Init(p.CHL_SIZE)
// 		for i := 0; i < p.GRT_NUM; i++ {
// 			ins := p.Create()
// 			ins.SetMWNode(p)
// 			ins.SetInChan(chl)
// 			mmrs = append(mmrs, ins)
// 			wg.Add(1)
// 			go func(wg *sync.WaitGroup, i int) {
// 				defer wg.Done()
// 				ins.Run(i)
// 			}(wg, i)
// 		}
// 		p.Instances = mmrs
// 		if pp != nil {
// 			for j := 0; j < len(pp.Instances); j++ {
// 				pp.Instances[j].SetOutChan(chl)
// 				pp.Instances[j].DefaultChanConfig()
// 			}
// 		}
// 		pp = p
// 		p = p.Next
// 	}

// 	wg.Wait()
// }

// 每个中间件创建的channel为后置
func (mgr *MiddlewareManager) suffixChannel() {
	p := mgr.TXL
	var pp *Chan
	var exe_order []MiddleManager
	var exe_order_idx []int
	for {
		if p == nil {
			break
		}
		var mmrs []MiddleManager
		chl := NewChan()
		chl.Init(p.CHL_SIZE)
		for i := 0; i < p.GRT_NUM; i++ {
			ins := p.Create()
			ins.SetMWNode(p)
			ins.SetOutChan(chl)
			if pp != nil {
				ins.SetInChan(pp)
			}
			mmrs = append(mmrs, ins)
			exe_order = append(exe_order, ins)
			exe_order_idx = append(exe_order_idx, i)
		}
		pp = chl
		p.Instances = mmrs
		p = p.Next
	}

	for i := len(exe_order) - 1; i >= 0; i-- {
		mgr.wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			exe_order[i].Run(exe_order_idx[i])
		}(mgr.wg, i)
	}

	mgr.wg.Wait()
}

// 修改缓冲区大小
func (mgr *MiddlewareManager) ChangeChanBufferSize(mn string, chl_size int) error {
	p := mgr.TXL
	for {
		if p == nil {
			return errors.New(fmt.Sprintf("have no mw: %s", mn))
		}
		if p.Name == mn {
			out_chl := p.Instances[0].GetOutChan()
			if out_chl.Free() == nil {
				chl := make(chan *Carrior, chl_size)
				out_chl.SetFree(chl)
				out_chl.Switch()
				close(out_chl.Free())
				break
			} else {
				return errors.New(fmt.Sprintf("mw %s: have no free chan", mn))
			}
		}
		p = p.Next
	}
	return nil
}

// 在某个中间件的前面插入添加的中间件
func (mgr *MiddlewareManager) InsertMWBack(mn string, nn *MWNode) error {
	p := mgr.TXL
	for {
		if p == nil {
			return errors.New(fmt.Sprintf("have no mw: %s", mn))
		}
		if p.Name == mn {
			var mmrs []MiddleManager
			chl := NewChan()
			chl.Init(nn.CHL_SIZE)
			for i := 0; i < nn.GRT_NUM; i++ {
				ins := nn.Create()
				ins.SetMWNode(nn)
				ins.SetOutChan(chl)
				ins.SetInChan(p.Instances[0].GetOutChan())

				mmrs = append(mmrs, ins)
				mgr.wg.Add(1)
				go func(wg *sync.WaitGroup, i int) {
					defer wg.Done()
					ins.Run(i)
				}(mgr.wg, i)
			}
			nn.Instances = mmrs
			if p.Next != nil {
				for i := 0; i < len(p.Next.Instances); i++ {
					p.Next.Instances[i].SetInChan(chl)
				}
			}
			nn.Next = p.Next
			p.Next = nn
			break
		}
		p = p.Next
	}
	return nil
}

// 拔出某个中间件
func (mgr *MiddlewareManager) DropMW(mn string) error {
	p := mgr.TXL
	var pp *MWNode
	var p_in, p_out *Chan
	for {
		if p == nil {
			return errors.New(fmt.Sprintf("have no mw: %s", mn))
		}
		if p.Name == mn {
			p_in = p.Instances[0].GetInChan()
			p_out = p.Instances[0].GetOutChan()
			for i := 0; i < len(p.Instances); i++ {
				p.Instances[i].SetInChan(nil)
			}
			break
		}
		pp = p
		p = p.Next
	}
	pn := p.Next
	if pn != nil {
		for i := 0; i < len(pn.Instances); i++ {
			pn.Instances[i].SetInChan(p_in)
		}
		p_in.SetFree(p_out.Active())
	}
	pp.Next = p.Next
	return nil
}

// 遍历当前业务线的中间件组成
func (mgr *MiddlewareManager) MWIter(ii string) {
	fmt.Println(ii, "---------")
	p := mgr.TXL
	for {
		if p == nil {
			break
		}
		for i := 0; i < len(p.Instances); i++ {
			in := p.Instances[i].GetInChan()
			out := p.Instances[i].GetOutChan()
			fmt.Printf("%s:\n%#v\n", p.Name, p.Instances[i])
			if in != nil {
				fmt.Printf("%#v %#v\n", in.CHL, in.CHL2)
			}
			fmt.Printf("%#v %#v\n---\n", out.CHL, out.CHL2)
		}
		p = p.Next
	}
	fmt.Println("---------")
}
