package closer

import "log"

type Closer struct {
	closeList []func()
}

func New() *Closer {
	return &Closer{}
}

func (c *Closer) AddToList(fn ...func()) {
	for _, f := range fn {
		c.closeList = append(c.closeList, f)
	}
}

func (c *Closer) Shutdown() {
	for _, f := range c.closeList {
		f()
	}
	log.Print("closer finished")
}
