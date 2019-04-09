package mw

import (
	mw "Gamora/src/Aquaman/base/middleware"
	"fmt"
	"log"
	"sync/atomic"
)

// 中间件
type Middleware struct {
	Node            *mw.MWNode
	IN_CHL          chan *mw.Carrior
	IN2_CHL         chan *mw.Carrior
	IN_CHAN_SWITCH  bool
	OUT_CHL         chan *mw.Carrior
	OUT2_CHL        chan *mw.Carrior
	MAIN_OUT_CHL    chan *mw.Carrior
	OUT_CHAN_SWITCH bool
}

func NewMiddleware() mw.MiddleManager {
	return &Middleware{}
}

func (m *Middleware) SetMWNode(mwn *mw.MWNode) {
	m.Node = mwn
}

func (m *Middleware) SetInChan(in chan *mw.Carrior) {
	m.IN_CHL = in
}

func (m *Middleware) SetOutChan(out chan *mw.Carrior) {
	m.OUT_CHL = out
}

func (m *Middleware) SetIn2Chan(in chan *mw.Carrior) {
	m.IN2_CHL = in
}

func (m *Middleware) SetOut2Chan(out chan *mw.Carrior) {
	m.OUT2_CHL = out
}

func (m *Middleware) DefaultChanConfig() {
	m.MAIN_OUT_CHL = m.OUT_CHL
	m.OUT_CHAN_SWITCH = true
	m.IN_CHAN_SWITCH = true
}

func (m *Middleware) GetFreeOutChan() chan *mw.Carrior {
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
