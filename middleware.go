package aqua

type MWHandler func() MiddleManager

// 单个中间件及配置
type MWNode struct {
	Name         string
	Create       MWHandler
	GRT_NUM      int
	CHL_SIZE     int
	ClosedGRTNum int64
	Instances    []MiddleManager
	Next         *MWNode
}

func NewMWNode(mn string, f MWHandler, num int, size int) *MWNode {
	return &MWNode{
		Name:     mn,
		Create:   f,
		GRT_NUM:  num,
		CHL_SIZE: size,
	}
}

func (n *MWNode) NextNode(mwn *MWNode) {
	n.Next = mwn
}
