package aqua

import (
	"sync"
)

type MiddleManager interface {
	SetMWNode(mwn *MWNode)
	SetInChan(in chan *Carrior)
	SetOutChan(out chan *Carrior)
	DefaultChanConfig()
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
func (mgr *MiddlewareManager) prefixChannel(n string) {
	wg := &sync.WaitGroup{}

	p := mgr.TXLS[n]
	var pp *MWNode
	for {
		if p == nil {
			break
		}
		var mmrs []MiddleManager
		chl := make(chan *Carrior, p.CHL_SIZE)
		for i := 0; i < p.GRT_NUM; i++ {
			ins := p.Create()
			ins.SetMWNode(p)
			ins.SetInChan(chl)
			mmrs = append(mmrs, ins)
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int) {
				defer wg.Done()
				ins.Run(i)
			}(wg, i)
		}
		p.Instances = mmrs
		if pp != nil {
			for j := 0; j < len(pp.Instances); j++ {
				pp.Instances[j].SetOutChan(chl)
				pp.Instances[j].DefaultChanConfig()
			}
		}
		pp = p
		p = p.Next
	}

	wg.Wait()
}

// 每个中间件创建的channel为后置
func (mgr *MiddlewareManager) suffixChannel(n string) {
	wg := &sync.WaitGroup{}

	p := mgr.TXLS[n]
	var pp chan *Carrior
	for {
		if p == nil {
			break
		}
		var mmrs []MiddleManager
		chl := make(chan *Carrior, p.CHL_SIZE)
		for i := 0; i < p.GRT_NUM; i++ {
			ins := p.Create()
			ins.SetMWNode(p)
			ins.SetOutChan(chl)
			ins.DefaultChanConfig()
			if pp != nil {
				ins.SetInChan(pp)
			}
			mmrs = append(mmrs, ins)
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int) {
				defer wg.Done()
				ins.Run(i)
			}(wg, i)
		}
		pp = chl
		p.Instances = mmrs
		p = p.Next
	}

	wg.Wait()
}

// 切换某个中间件的channel
func (mgr *MiddlewareManager) SwitchChannel(n, mn string) {
	p := mgr.TXLS[n]
	for {
		if p == nil {
			break
		}
		if p.Name == mn {

		}
	}
}

//
