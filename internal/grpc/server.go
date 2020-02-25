package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/entities"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/usecases"
	"github.com/Lefthander/otus-go-antibruteforce/internal/grpc/api"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ABFServer implemets the AntiBruteForceService Server
type ABFServer struct {
	abfService *usecases.ABFService
}

// NewABFServer creates a new instance of AntiBruteForce Service grpc server
func NewABFServer(abfs *usecases.ABFService) *ABFServer {

	return &ABFServer{
		abfService: abfs,
	}
}

// ListenAndServe runs the GRPC Server
func (a *ABFServer) ListenAndServe(listeningAddr string) error {
	g := grpc.NewServer(grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor))

	// Run wachdog goroutine to handle the system interrupts
	go func() {
		terminate := make(chan os.Signal, 1)
		signal.Notify(terminate, os.Interrupt, syscall.SIGINT)
		<-terminate
		log.Println("Received system interrupt...")
		g.GracefulStop()
	}()

	l, err := net.Listen("tcp", listeningAddr)
	if err != nil {
		return err
	}
	api.RegisterABFServiceServer(g, a)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(g)
	return g.Serve(l)
}

// Allow check authorization of triplet (login,password,ip)
func (a *ABFServer) Allow(ctx context.Context, req *api.AuthRequest) (*api.AuthResponse, error) {
	ipaddress := net.ParseIP(req.GetIpaddr())
	if ipaddress == nil {
		return &api.AuthResponse{
			Response: &api.AuthResponse_Error{
				Error: errors.ErrAuthRequestIPMissed.Error(),
			},
		}, errors.ErrAuthRequestIPMissed
	}
	isOk, err := a.abfService.IsAuthenticate(ctx, entities.AuthenticationRequest{
		Login:     req.GetLogin(),
		Password:  req.GetPassword(),
		IPAddress: req.GetIpaddr(),
	})
	if err != nil {
		log.Println("Errors during IsAUthnenticate", err)
		return &api.AuthResponse{
			Response: &api.AuthResponse_Error{
				Error: status.Error(codes.Internal, err.Error()).Error(),
			},
		}, err
	}
	return &api.AuthResponse{
		Response: &api.AuthResponse_Ok{Ok: isOk},
	}, nil
}

// Reset buckets for (login,ip) to initial state
func (a *ABFServer) Reset(ctx context.Context, req *api.AuthRequest) (*api.AuthResponse, error) {
	ipaddress := net.ParseIP(req.GetIpaddr())

	if ipaddress == nil {
		return &api.AuthResponse{
			Response: &api.AuthResponse_Error{
				Error: errors.ErrAuthRequestIPMissed.Error(),
			},
		}, errors.ErrAuthRequestIPMissed
	}
	err := a.abfService.ResetLimits(ctx, entities.AuthenticationRequest{
		Login:     req.GetLogin(),
		Password:  "",
		IPAddress: req.GetIpaddr(),
	})

	if err != nil {
		log.Println("Error at bucket reset limits", err)
		return &api.AuthResponse{
			Response: &api.AuthResponse_Error{
				Error: status.Error(codes.Internal, err.Error()).Error(),
			},
		}, err
	}

	return &api.AuthResponse{
		Response: &api.AuthResponse_Ok{Ok: true},
	}, nil
}

// AddToIpFilter adds a network to B/W IP Table
func (a *ABFServer) AddToIpFilter(ctx context.Context, req *api.IPFilterData) (*api.IPFilterResponse, error) {
	_, cidr, err := net.ParseCIDR(req.GetNetwork())
	if err != nil {
		return &api.IPFilterResponse{
			Error: status.Error(codes.Internal, err.Error()).Error(),
		}, err
	}

	err = a.abfService.AddIPNetwork(ctx, *cidr, req.GetColor())

	if err != nil {
		return &api.IPFilterResponse{
			Error: status.Error(codes.Internal, err.Error()).Error(),
		}, err
	}

	return &api.IPFilterResponse{
		Error: errors.ErrNoGrpcError.Error(),
	}, nil
}

// DeleteFromIpFilter deletes the network from B/W table
func (a *ABFServer) DeleteFromIpFilter(ctx context.Context, req *api.IPFilterData) (*api.IPFilterResponse, error) {
	_, cidr, err := net.ParseCIDR(req.GetNetwork())
	if err != nil {
		return &api.IPFilterResponse{
			Error: status.Error(codes.Internal, err.Error()).Error(),
		}, err
	}

	err = a.abfService.DeleteIPNetwork(ctx, *cidr, req.GetColor())

	if err != nil {
		return &api.IPFilterResponse{
			Error: status.Error(codes.Internal, err.Error()).Error(),
		}, err
	}

	return &api.IPFilterResponse{
		Error: errors.ErrNoGrpcError.Error(),
	}, nil
}

// GetIpFilters dump all content of specific table Black or White)
func (a *ABFServer) GetIpFilters(ctx context.Context, req *api.IPFilterData) (*api.IPFiltersData, error) {
	ipnets, err := a.abfService.ListIPNetworks(ctx, req.GetColor())

	if err != nil {
		return &api.IPFiltersData{
			Filters: nil,
			Error:   status.Error(codes.Internal, err.Error()).Error(),
		}, err
	}

	return &api.IPFiltersData{
		Filters: ipnets,
		Error:   errors.ErrNoGrpcError.Error(),
	}, nil
}
