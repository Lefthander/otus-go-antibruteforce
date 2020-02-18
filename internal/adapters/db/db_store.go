package db

import (
	"context"
	"database/sql"
	"net"

	"github.com/jmoiron/sqlx"
)

/* type FilterKeeper interface {
	AddIPNetwork(ctx context.Context, network net.IPNet, color bool) error    // Add to IPFilter table color: true - white , false - black
	DeleteIPNetwork(ctx context.Context, network net.IPNet, color bool) error // Delete network or IP from the IPFilter table, color: true - white, false - black
	IsIPConform(ctx context.Context, ip net.IP) (bool, error)                 // Verify does the IP address belongs to White/Black table, if belongs true, error = white/black
	ListIPNetworks(ctx context.Context, color bool) ([]string, error)         // Dump all records in the table , white table - color = true, black color = false
} */

// IPFilterDB implements interface FilterKeeper to store IP B/W tables in the SQL DB (Postgress)
type IPFilterDB struct {
	*sqlx.DB
}

// NewDBStore create a new instance of IPFilterDB
func NewDBStore(db *sqlx.DB) *IPFilterDB {
	return &IPFilterDB{DB: db}
}

// AddIPNetwork creates a new network in the B/W table, return error in case of network already exists
func (d *IPFilterDB) AddIPNetwork(ctx context.Context, network net.IPNet, color bool) error {
	switch color {
	case true:
		request := `INSERT INTO ip_white_list (network) VALUES ($1)`
		_, err := d.DB.ExecContext(ctx, request, network.String())
		if err != nil {
			return err
		}

	case false:
		request := `INSERT INTO ip_black_list (network) VALUES ($1)`
		_, err := d.DB.ExecContext(ctx, request, network.String())
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteIPNetwork from the B/W table in accordance with Color Flag, returns error in case of request network not found
func (d *IPFilterDB) DeleteIPNetwork(ctx context.Context, network net.IPNet, color bool) error {
	switch color {
	case true:
		request := `DELETE FROM ip_white_list WHERE network=$1`
		_, err := d.DB.ExecContext(ctx, request, network.String())
		if err != nil {
			return err
		}

	case false:
		request := `DELETE FROM ip_black_list WHERE network=$1`
		_, err := d.DB.ExecContext(ctx, request, network.String())
		if err != nil {
			return err
		}

	}
	return nil
}

// IsIPConform verifies does the IP address belongs to White/Black table, if belongs true, error = white/black
func (d *IPFilterDB) IsIPConform(ctx context.Context, ip net.IP) (bool, error) {

	// TODO:
	return true, nil
}

// ListIPNetworks in specified by color (B/W) table
func (d *IPFilterDB) ListIPNetworks(ctx context.Context, color bool) ([]string, error) {
	// TODO:
	var result sql.Result
	switch color {
	case true:
		request := `SELECT * FROM ip_white_list`
		result, err := d.DB.ExecContext(ctx, request)
		if err != nil {
			return nil, err
		}
	case false:
		request := `SELECT * FROM ip_black_list`
		result, err := d.DB.ExecContext(ctx, request)
		if err != nil {
			return nil, err
		}

	}
	return nil, nil
}
