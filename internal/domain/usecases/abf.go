package usecases

import (
	"context"
	"net"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/entities"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/interfaces"
)

// ABFService is a model of high level representation of Antibruteforce service
type ABFService struct {
	ConstraintN     uint32 // Number of Login Attempts per minute
	ConstraintM     uint32 // Number of Password Attempts per minute
	ConstraintK     uint32 // Number of IP attempts per minute
	BucketStorage   interfaces.BucketKeeper
	IPFilterStorage interfaces.FilterKeeper
}

// NewABFService creates a new instance of Antibruteforce service
func NewABFService(numberOfLogin, numberOfPassword, numberOfIP uint32, bucketStorage interfaces.BucketKeeper, filterStorage interfaces.FilterKeeper) *ABFService {
	return &ABFService{
		ConstraintN:     numberOfLogin,
		ConstraintM:     numberOfPassword,
		ConstraintK:     numberOfIP,
		BucketStorage:   bucketStorage,
		IPFilterStorage: filterStorage,
	}

}

// IsAuthenticate verifies is allow or not to pass the AuthenticationRequest
func (a *ABFService) IsAuthenticate(ctx context.Context, authRequest entities.AuthenticationRequest) (bool, error) {

	// TODO: Fill the content

	return true, nil

}

// IsIPConform verifies does specified IP included in the filter table either black or white
func (a *ABFService) IsIPConform(ctx context.Context, ip net.IP) (bool, error) {
	// TODO: Write some code here...
	return true, nil
}

// AddIPNetwork adds the net to the white or black table
func (a *ABFService) AddIPNetwork(ctx context.Context, net net.IPNet, color bool) error {

	// TODO: Do some code )
	return nil
}

// DeleteIPNetwork deletes the specified network from the white or black table
func (a *ABFService) DeleteIPNetwork(ctx context.Context, net net.IPNet, color bool) error {

	// TODO: Do some code here)
	return nil

}
