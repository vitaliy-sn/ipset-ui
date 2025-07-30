package controller

import (
	"ipset-ui/internal/ipset"
)

type IPSetController struct {
	ipset *ipset.IPSet
}

func NewIPSetController() *IPSetController {
	return &IPSetController{
		ipset: ipset.NewIPSet(),
	}
}
