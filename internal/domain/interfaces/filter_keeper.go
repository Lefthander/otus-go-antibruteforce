package interfaces

import (
	"context"
	"net"
)

// FilterKeeper is interface which implement Black/White filters storage
type FilterKeeper interface {
	AddIPNetwork(ctx context.Context, network net.IPNet, color bool) error    // Add to IPFilter table color: true - white , false - black
	DeleteIPNetwork(ctx context.Context, network net.IPNet, color bool) error // Delete network or IP from the IPFilter table, color: true - white, false - black
	IsIPConform(ctx context.Context, ip net.IP) (bool, error)                 // Verify does the IP address belongs to White/Black table, if belongs true, error = white/black
	ListIPNetworks(ctx context.Context, color bool) ([]string, error)         // Dump all records in the table , white table - color = true, black color = false
}
