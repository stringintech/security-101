package auth

import "net/http"

type VirtualFilterChain struct {
	chain    []Filter
	position int // Current position in filter chain
}

func NewVirtualFilterChain(c []Filter) *VirtualFilterChain {
	return &VirtualFilterChain{
		chain:    c,
		position: -1,
	}
}

func (c *VirtualFilterChain) DoFilter(w http.ResponseWriter, r *http.Request) {
	c.position++

	// If we have more filters, continue the chain
	if c.position < len(c.chain) {
		c.chain[c.position].DoFilter(w, r, c)
	}
}
