package adapters

import (
	"context"
	"net"
	"sync"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/entities"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
)

// IPFilterMemory implements the Black/White IP Filter tables storage
type IPFilterMemory struct {
	BlackIPList entities.IPFilter
	mxBlack     sync.RWMutex
	WhiteIPList entities.IPFilter
	mxWhite     sync.RWMutex
}

// NewIPFilterMemory creates inmemory storage for IP B/W tables
func NewIPFilterMemory() *IPFilterMemory {
	return &IPFilterMemory{
		BlackIPList: entities.IPFilter{Nets: map[string]net.IPNet{}, Color: false},
		mxBlack:     sync.RWMutex{},
		WhiteIPList: entities.IPFilter{Nets: map[string]net.IPNet{}, Color: true},
		mxWhite:     sync.RWMutex{},
	}
}

// IsIPConform verifies does the IP address belongs to White/Black table, if belongs true, error = white/black
func (ipf *IPFilterMemory) IsIPConform(ctx context.Context, ip net.IP) (bool, error) {

	ipf.mxWhite.RLock()
	defer ipf.mxWhite.RUnlock()
	for _, v := range ipf.WhiteIPList.Nets {
		if v.Contains(ip) {
			return true, errors.ErrIPFilterMatchedWhiteList
		}
	}

	ipf.mxBlack.RLock()
	defer ipf.mxBlack.RUnlock()
	for _, v := range ipf.BlackIPList.Nets {
		if v.Contains(ip) {
			return true, errors.ErrIPFilterMatchedBlackList
		}
	}
	return false, errors.ErrIPFilterNoMatch
}

// AddIPNetwork creates a new network in the B/W table, return error in case of network already exists
func (ipf *IPFilterMemory) AddIPNetwork(ctx context.Context, network net.IPNet, color bool) error {

	switch color {
	case true:
		ipf.mxWhite.Lock()
		defer ipf.mxWhite.Unlock()
		if _, ok := ipf.WhiteIPList.Nets[network.String()]; !ok {
			ipf.WhiteIPList.Nets[network.String()] = network
			return nil
		}
	case false:
		ipf.mxBlack.Lock()
		defer ipf.mxBlack.Unlock()
		if _, ok := ipf.BlackIPList.Nets[network.String()]; !ok {
			ipf.BlackIPList.Nets[network.String()] = network
			return nil
		}

	}
	return errors.ErrIPFilterNetworkAlreadyExists
}

// DeleteIPNetwork from the B/W table in accordance with Color Flag, returns error in case of request network not found
func (ipf *IPFilterMemory) DeleteIPNetwork(ctx context.Context, network net.IPNet, color bool) error {

	switch color {
	case true:
		ipf.mxWhite.Lock()
		defer ipf.mxWhite.Unlock()
		if _, ok := ipf.WhiteIPList.Nets[network.String()]; ok {
			delete(ipf.WhiteIPList.Nets, network.String())
			return nil
		}
	case false:
		ipf.mxBlack.Lock()
		defer ipf.mxBlack.Unlock()
		if _, ok := ipf.BlackIPList.Nets[network.String()]; ok {
			delete(ipf.BlackIPList.Nets, network.String())
			return nil
		}

	}
	return errors.ErrIPFilterNetworkNotFound
}

// ListIPNetworks in specified by color (B/W) table
func (ipf *IPFilterMemory) ListIPNetworks(ctx context.Context, color bool) ([]string, error) {

	var iplist []string
	switch color {
	case true:
		ipf.mxWhite.Lock()
		defer ipf.mxWhite.Unlock()
		for k := range ipf.WhiteIPList.Nets {
			iplist = append(iplist, k)
		}
		return iplist, nil
	case false:
		ipf.mxBlack.Lock()
		defer ipf.mxBlack.Unlock()
		for k := range ipf.BlackIPList.Nets {
			iplist = append(iplist, k)
		}
		return iplist, nil
	}
	return nil, nil
}
