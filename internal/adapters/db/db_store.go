package db

import (
	"context"
	"database/sql"
	"net"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/jmoiron/sqlx"
)

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
	//requestWhite := `SELECT * FROM ip_white_list WHERE ipaddr=$1 << ANY (network)`
	requestWhite := `SELECT * FROM ip_white_list WHERE $1 << (network)`
	//wnets := make([]string, 0)
	wnets := make([]struct {
		Id      int64
		Network string
	}, 0)

	err := d.DB.SelectContext(ctx, &wnets, requestWhite, ip.String())

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if len(wnets) != 0 {
		return true, errors.ErrIPFilterMatchedWhiteList
	}

	//bnets := make([]string, 0)

	bnets := make([]struct {
		Id      int64
		Network string
	}, 0)

	//select * from ip_white_list where cidr '10.10.0.0/32' << network(network);

	//requestBlack := `SELECT * FROM ip_black_list WHERE ipaddr=$1 << ANY (network)`
	requestBlack := `SELECT * FROM ip_black_list WHERE $1 << (network)`
	err = d.DB.SelectContext(ctx, &bnets, requestBlack, ip.String())

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if len(bnets) != 0 {
		return true, errors.ErrIPFilterMatchedBlackList
	}

	return false, errors.ErrIPFilterNoMatch
}

// ListIPNetworks in specified by color (B/W) table
func (d *IPFilterDB) ListIPNetworks(ctx context.Context, color bool) ([]string, error) {
	//nets := make([]string, 0)
	nets := []struct {
		Id      int64 // nolint
		Network string
	}{}

	switch color {
	case true:
		request := `SELECT * FROM ip_white_list`
		err := d.DB.SelectContext(ctx, &nets, request)

		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	case false:
		request := `SELECT * FROM ip_black_list`
		err := d.DB.SelectContext(ctx, &nets, request)

		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	}
	var result []string

	for _, v := range nets {
		result = append(result, v.Network)
	}

	return result, nil
}
