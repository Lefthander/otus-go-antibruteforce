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
	"github.com/Lefthander/otus-go-antibruteforce/internal/tokenbucket"
	"go.uber.org/zap"
)

// 2DO: Implement the mechanism of clearance of longtime unused buckets

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
	abf := &ABFService{
		ConstraintN:     numberOfLogin,
		ConstraintM:     numberOfPassword,
		ConstraintK:     numberOfIP,
		BucketStorage:   bucketStorage,
		IPFilterStorage: filterStorage,
		logger:          logger,
		config:          config,
	}
	abf.loginMap = adapters.NewUUIDTable()
	abf.passwdMap = adapters.NewUUIDTable()
	abf.ipMap = adapters.NewUUIDTable()

	return abf
}

// validateAuthRequest is a service function just to verify the consistency of AuthRequest
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

// CheckBuckets verify the Authentication Request againsts token buckets, create them if necessary
func (a *ABFService) CheckBuckets(ctx context.Context, request entities.AuthenticationRequest) (bool, error) {

	loginID := a.loginMap.AddToTable(ctx, request.Login)
	passwdID := a.passwdMap.AddToTable(ctx, request.Password)
	ipID := a.ipMap.AddToTable(ctx, request.IPAddress)

	_, err := a.BucketStorage.GetBucket(ctx, loginID)

	if err == errors.ErrTokenBucketNotFound {
		tb, err := tokenbucket.NewTokenBucket(ctx, 1, time.Minute/time.Duration(a.config.ConstraintN))

		if err != nil {
			return false, err
		}

		err = a.BucketStorage.CreateBucket(ctx, loginID, tb)

		if err != nil {
			return false, err
		}
	}

	loginbucket, err := a.BucketStorage.GetBucket(ctx, loginID)

	if err != nil {
		return false, err
	}

	isloginOK := loginbucket.Allow(ctx)

	_, err = a.BucketStorage.GetBucket(ctx, passwdID)

	if err == errors.ErrTokenBucketNotFound {
		tb, err := tokenbucket.NewTokenBucket(ctx, 1, time.Minute/time.Duration(a.config.ConstraintM))

		if err != nil {
			return false, err
		}

		err = a.BucketStorage.CreateBucket(ctx, passwdID, tb)

		if err != nil {
			return false, err
		}
	}

	passwdbucket, err := a.BucketStorage.GetBucket(ctx, passwdID)

	if err != nil {
		return false, err
	}

	ispasswdOK := passwdbucket.Allow(ctx)

	_, err = a.BucketStorage.GetBucket(ctx, ipID)

	if err == errors.ErrTokenBucketNotFound {
		tb, err := tokenbucket.NewTokenBucket(ctx, 1, time.Minute/time.Duration(a.config.ConstraintK))

		if err != nil {
			return false, err
		}

		err = a.BucketStorage.CreateBucket(ctx, ipID, tb)

		if err != nil {
			return false, err
		}
	}

	ipbucket, err := a.BucketStorage.GetBucket(ctx, ipID)

	if err != nil {
		return false, err
	}

	isipOK := ipbucket.Allow(ctx)

	return isloginOK && ispasswdOK && isipOK, nil
}

// IsAuthenticate verifies is allow or not to pass the AuthenticationRequest
func (a *ABFService) IsAuthenticate(ctx context.Context, authRequest entities.AuthenticationRequest) (bool, error) {
	err := validateAuthRequest(authRequest)

	if err != nil {
		return false, err
	}

	flag, err := a.IsIPConform(ctx, net.ParseIP(authRequest.IPAddress))

	if flag && err == errors.ErrIPFilterMatchedWhiteList {
		return true, nil
	}

	if flag && err == errors.ErrIPFilterMatchedBlackList {
		return false, errors.ErrIPFilterMatchedBlackList
	}

	flag, err = a.CheckBuckets(ctx, authRequest)

	if !flag {
		return false, err
	}

	return flag, nil
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

// ResetLimits clear corresponding buckets for pair of login && ip
func (a *ABFService) ResetLimits(ctx context.Context, request entities.AuthenticationRequest) error {
	loginID := a.loginMap.AddToTable(ctx, request.Login)
	ipID := a.ipMap.AddToTable(ctx, request.IPAddress)
	loginbucket, err := a.BucketStorage.GetBucket(ctx, loginID)

	if err != nil {
		return err
	}

	loginbucket.Reset(ctx)

	ipbucket, err := a.BucketStorage.GetBucket(ctx, ipID)

	if err != nil {
		return err
	}

	ipbucket.Reset(ctx)

	return nil
}
