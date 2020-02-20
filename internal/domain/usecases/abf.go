package usecases

import (
	"context"
	"net"
	"time"

	"github.com/Lefthander/otus-go-antibruteforce/config"
	"github.com/Lefthander/otus-go-antibruteforce/internal/adapters"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/entities"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/interfaces"
	"go.uber.org/zap"
)

// ABFService is a model of high level representation of Antibruteforce service
type ABFService struct {
	ConstraintN     uint32 // Number of Login Attempts per minute
	ConstraintM     uint32 // Number of Password Attempts per minute
	ConstraintK     uint32 // Number of IP attempts per minute
	BucketStorage   interfaces.BucketKeeper
	IPFilterStorage interfaces.FilterKeeper
	loginMap        *adapters.UUIDTable
	passwdMap       *adapters.UUIDTable
	ipMap           *adapters.UUIDTable
	logger          *zap.Logger
	config          *config.ServiceConfig
}

// NewABFService creates a new instance of Antibruteforce service
func NewABFService(numberOfLogin, numberOfPassword, numberOfIP uint32, bucketStorage interfaces.BucketKeeper,
	filterStorage interfaces.FilterKeeper, logger *zap.Logger, config *config.ServiceConfig) *ABFService {
	return &ABFService{
		ConstraintN:     numberOfLogin,
		ConstraintM:     numberOfPassword,
		ConstraintK:     numberOfIP,
		BucketStorage:   bucketStorage,
		IPFilterStorage: filterStorage,
		logger:          logger,
		config:          config,
	}
}

func validateAuthRequest(a entities.AuthenticationRequest) error {
	if a.Login == "" {
		return errors.ErrAuthRequestLoginMissed
	}

	if a.Password == "" {
		return errors.ErrAuthRequestPasswordMissed
	}

	if a.IPAddress == "" {
		return errors.ErrAuthRequestIPMissed
	}

	return nil
}

// IsAuthenticate verifies is allow or not to pass the AuthenticationRequest
func (a *ABFService) IsAuthenticate(ctx context.Context, authRequest entities.AuthenticationRequest) (bool, error) {
	err := validateAuthRequest(authRequest)

	if err != nil {
		return false, err
	}

	//td := time.Duration(time.Minute / time.Duration(a.config.ConstraintN))

	return true, nil
}

// IsIPConform verifies does specified IP included in the filter table either black or white
func (a *ABFService) IsIPConform(ctx context.Context, ip net.IP) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.config.TimeOut)*time.Microsecond)
	defer cancel()

	flag, err := a.IPFilterStorage.IsIPConform(ctx, ip)

	if flag && err == errors.ErrIPFilterMatchedWhiteList {
		return true, nil
	}

	if flag && err == errors.ErrIPFilterMatchedBlackList {
		return false, nil
	}
	return true, nil
}

// AddIPNetwork adds the net to the white or black table
func (a *ABFService) AddIPNetwork(ctx context.Context, net net.IPNet, color bool) error {
	err := a.IPFilterStorage.AddIPNetwork(ctx, net, color)

	if err != nil {
		return err
	}

	return nil
}

// DeleteIPNetwork deletes the specified network from the white or black table
func (a *ABFService) DeleteIPNetwork(ctx context.Context, net net.IPNet, color bool) error {
	err := a.IPFilterStorage.DeleteIPNetwork(ctx, net, color)

	if err != nil {
		return err
	}

	return nil
}
