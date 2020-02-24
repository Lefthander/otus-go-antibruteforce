package grpc

import (
	"context"
	"net"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/entities"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/usecases"
	"github.com/Lefthander/otus-go-antibruteforce/internal/grpc/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ABFServer implemets the AntiBruteForceService Server
type ABFServer struct {
	abfService *usecases.ABFService
}

// Allow check authorization of triplet (login,password,ip)
func (a *ABFServer) Allow(ctx context.Context, req *api.AuthRequest) *api.AuthResponse {
	ipaddress := net.ParseIP(req.GetIpaddr())
    if ipaddress == nil {
		return &api.AuthResponse{
			 Response: &api.AuthResponse_Error{
			    Error: errors.ErrAuthRequestIPMissed.Error(),
		     },
	    },nil			
		}
	isOk, err := a.abfService.IsAuthenticate(ctx,entities.AuthenticationRequest{
		Login: req.GetLogin(),
		Password: req.GetPassword(),
		IPAddress: req.GetIpaddr(),
	})
	if err != nil {
		log.Println("Errors during IsAUthnenticate",err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.AuthResponse{
		Response: &&api.AuthResponse_Ok{Ok:isOk},
	},nil
}

// Reset buckets for (login,ip) to initial state
func (a *ABFServer) Reset(ctx context.Context, req *api.AuthRequest) *api.AuthResponse {
	ipaddress := net.ParseIP(req.GetIpaddr())
    if ipaddress == nil {
		return &api.AuthResponse{
			 Response: &api.AuthResponse_Error{
			    Error: errors.ErrAuthRequestIPMissed.Error(),
		     },
		},nil
	}
	err:= a.abfService.		

}

// AddToIPFilter adds a network to B/W IP Table
func (a *ABFServer) AddToIPFilter(ctx context.Context, req *api.IPFilterData) *api.IPFilterResponse {

}

// DeleteFromIPFilter deletes the network from B/W table
func (a *ABFServer) DeleteFromIPFilter(ctx context.Context, req *api.IPFilterData) *api.IPFilterResponse {

}

// GetIPFilters dump all content of specific table Black or White)
func (a *ABFServer) GetIPFilters(ctx context.Context, req *api.IPFilterData) *api.IPFiltersData {

}
