package integration

import (
	"github.com/Lefthander/otus-go-antibruteforce/internal/grpc/api"
	"github.com/cucumber/godog"
)

type ABFAPI struct {
	verifyAuthRequest *api.AuthRequest
	api.IPFilterData
}

const (
	abfURL = "http://localhost:4000"
)

func ipaddress(arg1 string) error {
	return godog.ErrPending
}

func checkAddress() error {
	return godog.ErrPending
}

func requestIsNotBlocked() error {
	return godog.ErrPending
}

func network(arg1 string) error {
	return godog.ErrPending
}

func addNetworkToBlackList() error {
	return godog.ErrPending
}

/* func addNetworkToBlackList() error {
	return godog.ErrPending
}
*/
func errorNetworkAlreadyExists() error {
	return godog.ErrPending
}

func addNetworkToWhiteList() error {
	return godog.ErrPending
}

func checkIpaddress() error {
	return godog.ErrPending
}

func requestIsBlocked() error {
	return godog.ErrPending
}

func deleteNetworkFromBlackList() error {
	return godog.ErrPending
}

func errorNetworkNotFound() error {
	return godog.ErrPending
}

func deleteNetworkFromWhiteList() error {
	return godog.ErrPending
}

func login(arg1 string) error {
	return godog.ErrPending
}

func password(arg1 string) error {
	return godog.ErrPending
}

func iP(arg1 string) error {
	return godog.ErrPending
}

func delayBetweenRequestIs(arg1 string) error {
	return godog.ErrPending
}

func sendRequests(arg1 int) error {
	return godog.ErrPending
}

func requestsAreNotBlocked(arg1 int) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^ipaddress "([^"]*)"$`, ipaddress)
	s.Step(`^check address$`, checkAddress)
	s.Step(`^request is not blocked$`, requestIsNotBlocked)
	s.Step(`^network "([^"]*)"$`, network)
	s.Step(`^add network to black List$`, addNetworkToBlackList)
	s.Step(`^add network to black list$`, addNetworkToBlackList)
	s.Step(`^error network already exists$`, errorNetworkAlreadyExists)
	s.Step(`^add network to white List$`, addNetworkToWhiteList)
	s.Step(`^check ipaddress$`, checkIpaddress)
	s.Step(`^request is blocked$`, requestIsBlocked)
	s.Step(`^delete network from black list$`, deleteNetworkFromBlackList)
	s.Step(`^error network not found$`, errorNetworkNotFound)
	s.Step(`^delete network from white list$`, deleteNetworkFromWhiteList)
	s.Step(`^login "([^"]*)"$`, login)
	s.Step(`^password "([^"]*)"$`, password)
	s.Step(`^IP "([^"]*)"$`, iP)
	s.Step(`^delay between request is "([^"]*)"$`, delayBetweenRequestIs)
	s.Step(`^send (\d+) requests$`, sendRequests)
	s.Step(`^(\d+) requests are not blocked$`, requestsAreNotBlocked)
}
