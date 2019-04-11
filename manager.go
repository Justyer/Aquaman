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
	TXLS map[string]*MWNode
}

func NewMWManager() *MiddlewareManager {
	return &MiddlewareManager{
		TXLS: make(map[string]*MWNode),
	}
}

// 注册一条业务线
func (mgr *MiddlewareManager) RegisterTXL(n string, mwn *MWNode) {
	mgr.TXLS[n] = mwn
}

// 根据业务线名字执行任务
func (mgr *MiddlewareManager) ExecuteByName(n string) {
	mgr.suffixChannel(n)
}

// 每个中间件创建的channel为前置
// func (mgr *MiddlewareManager) prefixChannel(n string) {
// 	wg := &sync.WaitGroup{}

// 	p := mgr.TXLS[n]
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
func (mgr *MiddlewareManager) suffixChannel(n string) {
	wg := &sync.WaitGroup{}

	p := mgr.TXLS[n]
	var pp *Chan
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
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int, s string) {
				defer wg.Done()
				ins.Run(i)
			}(wg, i, p.Name)
		}
		pp = chl
		p.Instances = mmrs
		p = p.Next
	}

	wg.Wait()
}

// 拔出某个中间件
func (mgr *MiddlewareManager) DropMW(n, mn string) error {
	p := mgr.TXLS[n]
	var pp *MWNode
	var p_in, p_out *Chan
	for {
		if p == nil {
			return errors.New(fmt.Sprintf("have no %s-%s", n, mn))
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

//
func (mgr *MiddlewareManager) ServiceFinder(n, ii string) {
	fmt.Println(ii, "---------")
	p := mgr.TXLS[n]
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
