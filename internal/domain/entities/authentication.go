package entities

import "net"

// AuthenticationRequest is a data model for Antibruteforce Service Request from remote site
type AuthenticationRequest struct {
	Login     string
	Password  string
	IPAddress net.IP
}
