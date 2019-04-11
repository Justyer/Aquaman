package aqua

import "fmt"

type Chan struct {
	CHL     chan *Carrior
	CHL2    chan *Carrior
	CHL_USE bool
}

func NewChan() *Chan {
	return &Chan{}
}

func (c *Chan) Init(chl_size int) {
	chl := make(chan *Carrior, chl_size)
	c.CHL = chl
	c.CHL_USE = true
}

func (c *Chan) Push(e *Carrior) {
	c.Active() <- e
}

func (c *Chan) Pop(f func(*Carrior)) {
	for {
		var cr *Carrior
		isUse, ok := false, false
		select {
		case cr, ok = <-c.CHL:
			isUse = ok || isUse
		case cr, ok = <-c.CHL2:
			isUse = ok || isUse
		}
		if isUse {
			f(cr)
		} else {
			fmt.Println("break")
			break
		}
	}
}

func (c *Chan) Switch() {
	c.CHL_USE = !c.CHL_USE
}

func (c *Chan) Active() chan *Carrior {
	if c.CHL_USE {
		return c.CHL
	} else {
		return c.CHL2
	}
}

func (c *Chan) SetFree(cr chan *Carrior) {
	if c.CHL_USE {
		c.CHL2 = cr
	} else {
		c.CHL = cr
	}
}

func (c *Chan) Free() chan *Carrior {
	if c.CHL_USE {
		return c.CHL2
	} else {
		return c.CHL
	}
}

func (c *Chan) Close() {
	close(c.Active())
}
