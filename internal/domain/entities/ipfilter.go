package entities

import (
	"net"
)

// IPFilter is a model of IP Filter list
type IPFilter struct {
	Nets  map[string]net.IPNet // Table of networks
	Color bool                 // Color indicate the purpose of the table: white - true, black - false
}
