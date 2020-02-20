package interfaces

import (
	"context"
	"net"
)

// FilterKeeper is interface which implement Black/White filters storage
type FilterKeeper interface {
	// Add to IPFilter table color: true - white , false - black
	AddIPNetwork(ctx context.Context, network net.IPNet, color bool) error
	// Delete network or IP from the IPFilter table, color: true - white, false - black
	DeleteIPNetwork(ctx context.Context, network net.IPNet, color bool) error
	// Verify does the IP address belongs to White/Black table, if belongs true, error = white/black
	IsIPConform(ctx context.Context, ip net.IP) (bool, error)
	// Dump all records in the table , white table - color = true, black color = false
	ListIPNetworks(ctx context.Context, color bool) ([]string, error)
}
