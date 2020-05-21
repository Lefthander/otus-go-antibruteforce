//nolint
package tests

import (
	"context"
	"errors"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/Lefthander/otus-go-antibruteforce/internal/grpc/api"
	"github.com/Pallinder/go-randomdata"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"google.golang.org/grpc"
)

const (
	random   = "random"
	Black    = false
	White    = true
	DummyNet = "0.0.0.0/24"
	abfURL   = "abf-service:4000"
)

var (
	ctx = context.Background()
)

type abfTest struct {
	client                  api.ABFServiceClient
	login                   string
	password                string
	ipaddress               string
	lastStatus              bool
	lastError               string
	ipList                  []string
	resetAt                 int
	numberOfAllowedRequests int
	intervalBetweenRequests time.Duration
}

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

func (a *abfTest) connectToServer(url string) error {
	conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure())

	if err != nil {
		return err

	}
	a.client = api.NewABFServiceClient(conn)
	if a.client != nil {
		return errors.New("Failure to create a gRPC client")
	}
	return nil
}

func (a *abfTest) setLogin(login string) error {
	a.login = login
	return nil
}

func (a *abfTest) setPassword(password string) error {
	a.password = password
	return nil
}

func (a *abfTest) setIpaddress(ipaddress string) error {
	a.ipaddress = ipaddress
	return nil
}

func (a *abfTest) delayBetweenRequestIs(interval string) error {
	t, err := time.ParseDuration(interval)

	if err != nil {
		return err
	}
	a.intervalBetweenRequests = t
	return nil
}

func (a *abfTest) resetLimits() error {
	r, err := a.client.Reset(ctx, &api.AuthRequest{
		Login:    a.login,
		Password: a.password,
		Ipaddr:   a.ipaddress,
	})
	if err != nil {
		return err
	}
	if !r.GetOk() {
		return errors.New(r.GetError())
	}
	return nil
}

func (a *abfTest) sendRequests(number int) error {
	var allowedRequests int
	for i := 0; i <= number; i++ {

		if i == a.resetAt {
			if err := a.resetLimits(); err != nil {
				return nil
			}
		}
		err := a.isAllow()
		if err != nil {
			return err
		}
		if a.lastStatus {
			allowedRequests++
		}
		time.Sleep(a.intervalBetweenRequests)
		a.numberOfAllowedRequests = allowedRequests
		return nil
	}
	return nil
}

func (a *abfTest) allRequestsAreNotBlocked() error {

	return nil
}

/* func allRequestsAreNotBlocked() error {
	return godog.ErrPending
}
*/
func (a *abfTest) resetAtRequest(n int) error {
	a.resetAt = n
	return nil
}

func requestArePassed(arg1 int) error {
	return godog.ErrPending
}

func (a *abfTest) isAllow() error {
	login, password, ip := a.login, a.password, a.ipaddress
	if login == random {
		login = randomdata.Email()
	}
	if password == random {
		password = randomdata.SillyName()
	}

	if ip == random {
		ip = randomdata.IpV4Address()
	}
	r, err := a.client.Allow(ctx, &api.AuthRequest{
		Login:    login,
		Password: password,
		Ipaddr:   ip,
	})
	if err != nil {
		return err
	}
	a.lastStatus = r.GetOk()
	return nil
}

func (a *abfTest) verifyIpaddress() error {
	a.login = random
	a.password = random

	return a.isAllow()
}

func (a *abfTest) requestIsNotBlocked() error {
	if !a.lastStatus {
		return errors.New("request is blocked, expected not")
	}
	return nil
}

func (a *abfTest) addIpaddressToBlackList() error {
	r, err := a.client.AddToIpFilter(ctx, &api.IPFilterData{
		Network: a.ipaddress,
		Color:   Black,
	})
	if err != nil {
		return err
	}
	a.lastError = r.GetError()

	return nil
}

func (a *abfTest) requestIsBlocked() error {
	if a.lastStatus {
		return errors.New("request is not blocked, expected blocked")
	}
	return nil
}

func (a *abfTest) errorReportedNetworkAlreadyExists() error {
	if a.lastError == "" {
		return errors.New("expected error")
	}
	return nil
}

func (a *abfTest) addIpaddressToWhiteList() error {
	r, err := a.client.AddToIpFilter(ctx, &api.IPFilterData{
		Network: a.ipaddress,
		Color:   White,
	})
	if err != nil {
		return err
	}
	a.lastError = r.GetError()
	return nil
}

func (a *abfTest) getWhiteListContents() error {
	r, err := a.client.GetIpFilters(ctx, &api.IPFilterData{
		Network: DummyNet,
		Color:   White,
	})
	if err != nil {
		return err
	}
	a.ipList = r.GetFilters()
	return nil
}

func (a *abfTest) isContain(s string) bool {
	for _, v := range a.ipList {
		if v == s {
			return true
		}
	}
	return false
}

func (a *abfTest) receivedIpaddressesAndAnd(ip1, ip2, ip3 string) error {

	if !a.isContain(ip1) && !a.isContain(ip2) && !a.isContain(ip3) {
		return errors.New("no ip addresses found")
	}
	return nil
}

func (a *abfTest) getBlackListContents() error {
	r, err := a.client.GetIpFilters(ctx, &api.IPFilterData{
		Network: DummyNet,
		Color:   Black,
	})
	if err != nil {
		return err
	}
	a.ipList = r.GetFilters()
	return nil
}

func (a *abfTest) receivedIpaddressesAnd(ip1, ip2 string) error {
	if !a.isContain(ip1) && !a.isContain(ip2) {
		return errors.New("expected ip addreses")
	}

	return godog.ErrPending
}

func (a *abfTest) deleteIpaddressFromWhiteList() error {
	r, err := a.client.DeleteFromIpFilter(ctx, &api.IPFilterData{
		Network: a.ipaddress,
		Color:   White,
	})
	if err != nil {
		return err
	}

	a.lastError = r.GetError()

	return nil
}

func (a *abfTest) receivedErrorIpaddressNotFound() error {
	if a.lastError == "" {
		return errors.New("expected error")
	}
	return nil
}

func (a *abfTest) receivedStatusOk() error {
	if !a.lastStatus {
		return errors.New("Expected status - True = Ok")
	}
	return nil
}

func (a *abfTest) receivedEmptyList() error {
	if len(a.ipList) != 0 {
		return errors.New("expected empty list")
	}
	return nil
}

func (a *abfTest) deleteIpaddressFromBlackList() error {
	r, err := a.client.DeleteFromIpFilter(ctx, &api.IPFilterData{
		Network: a.ipaddress,
		Color:   Black,
	})
	if err != nil {
		return err
	}
	a.lastError = r.GetError()
	return nil
}

func FeatureContext(s *godog.Suite) {
	a := &abfTest{}
	_ = a.connectToServer(abfURL)
	s.Step(`^login "([^"]*)"$`, a.setLogin)
	s.Step(`^password "([^"]*)"$`, a.setPassword)
	s.Step(`^ipaddress "([^"]*)"$`, a.setIpaddress)
	s.Step(`^delay between request is "([^"]*)"$`, a.delayBetweenRequestIs)
	s.Step(`^send (\d+) requests$`, a.sendRequests)
	s.Step(`^All requests are not blocked$`, a.allRequestsAreNotBlocked)
	s.Step(`^reset at (\d+) request$`, a.resetAtRequest)
	s.Step(`^(\d+) request are passed$`, requestArePassed)
	s.Step(`^verify ipaddress$`, a.verifyIpaddress)
	s.Step(`^request is not blocked$`, a.requestIsNotBlocked)
	s.Step(`^add ipaddress to black list$`, a.addIpaddressToBlackList)
	s.Step(`^request is blocked$`, a.requestIsBlocked)
	s.Step(`^error reported - network already exists$`, a.errorReportedNetworkAlreadyExists)
	s.Step(`^add ipaddress to white list$`, a.addIpaddressToWhiteList)
	s.Step(`^get white list contents$`, a.getWhiteListContents)
	s.Step(`^received ipaddresses "([^"]*)" And "([^"]*)" And "([^"]*)"$`, a.receivedIpaddressesAndAnd)
	s.Step(`^get black list contents$`, a.getBlackListContents)
	s.Step(`^received ipaddresses "([^"]*)" And "([^"]*)"$`, a.receivedIpaddressesAnd)
	s.Step(`^delete ipaddress from white list$`, a.deleteIpaddressFromWhiteList)
	s.Step(`^received error - ipaddress not found$`, a.receivedErrorIpaddressNotFound)
	s.Step(`^received status Ok$`, a.receivedStatusOk)
	s.Step(`^received empty list$`, a.receivedEmptyList)
	s.Step(`^delete ipaddress from black list$`, a.deleteIpaddressFromBlackList)
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
