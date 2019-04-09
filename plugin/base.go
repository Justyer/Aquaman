package plg

import (
	"fmt"
	"log"
	"sync/atomic"

	aqua "github.com/Justyer/Aquaman"
)

// 中间件
type Middleware struct {
	Node            *aqua.MWNode
	IN_CHL          chan *aqua.Carrior
	IN2_CHL         chan *aqua.Carrior
	IN_CHAN_SWITCH  bool
	OUT_CHL         chan *aqua.Carrior
	OUT2_CHL        chan *aqua.Carrior
	MAIN_OUT_CHL    chan *aqua.Carrior
	OUT_CHAN_SWITCH bool
}

func NewMiddleware() aqua.MiddleManager {
	return &Middleware{}
}

func (m *Middleware) SetMWNode(mwn *aqua.MWNode) {
	m.Node = mwn
}

func (m *Middleware) SetInChan(in chan *aqua.Carrior) {
	m.IN_CHL = in
}

func (m *Middleware) SetOutChan(out chan *aqua.Carrior) {
	m.OUT_CHL = out
}

func (m *Middleware) SetIn2Chan(in chan *aqua.Carrior) {
	m.IN2_CHL = in
}

func (m *Middleware) SetOut2Chan(out chan *aqua.Carrior) {
	m.OUT2_CHL = out
}

func (m *Middleware) DefaultChanConfig() {
	m.MAIN_OUT_CHL = m.OUT_CHL
	m.OUT_CHAN_SWITCH = true
	m.IN_CHAN_SWITCH = true
}

func (m *Middleware) GetFreeOutChan() chan *aqua.Carrior {
	if m.OUT_CHAN_SWITCH {
		return m.OUT2_CHL
	} else {
		return m.OUT_CHL
	}
}

func (m *Middleware) SwitchOutChan() {
	if m.OUT_CHAN_SWITCH {
		m.MAIN_OUT_CHL = m.OUT2_CHL
	} else {
		m.MAIN_OUT_CHL = m.OUT_CHL
	}
	m.OUT_CHAN_SWITCH = !m.OUT_CHAN_SWITCH
}

func (m *Middleware) Run(i int) {
	log.Fatal("Must implement Run() func")
}

func (m *Middleware) Close() {
	if m.OUT_CHL == nil {
		return
	}
	closeNum := atomic.AddInt64(&m.Node.ClosedGRTNum, 1)
	if closeNum == int64(m.Node.GRT_NUM) {
		fmt.Printf("Close %p OutChan\n", m.OUT_CHL)
		close(m.OUT_CHL)
	}
}
