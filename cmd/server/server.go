package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mt-inside/go-usvc"
	versionspb "github.com/mt-inside/versions-over-ip/api/v1alpha1"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/mt-inside/versions-over-ip/pkg/fetch"
)

type workItem struct {
	ID  string
	Req *versionspb.VersionsRequest
}

var (
	grpcAddr = ":50051"

	workItems map[string]workItem = make(map[string]workItem)
)

type versionsServer struct {
	versionspb.UnimplementedVersionsServer

	log logr.Logger
}

func (s *versionsServer) GetVersions(ctxt context.Context, in *versionspb.VersionsRequest) (*lropb.Operation, error) {
	uid, _ := uuid.NewUUID()
	work := workItem{ID: uid.String(), Req: in}
	workItems[uid.String()] = work

	res := &lropb.Operation{
		Name: uid.String(),
		Done: false,
	}

	s.log.V(1).Info("Long-running operation started", "work item id", uid.String())

	return res, nil
}

/* See:
* * https://github.com/googleapis/gapic-showcase/blob/master/server/services/operations_service.go
* * https://github.com/googleapis/googleapis/blob/master/google/longrunning/operations.proto
* Semantics:
* * I *think* the semantics are: Get() returns right away (done & result / not done), Wait() blocks and returns (done & result / not done after timeout)
*   * However, calling Wait() in the client does not result in a call to WaitOperation() here...
* * The docs talk about "polling", and indeed Get() is a short-poll, but Wait() is a long poll / stream so it's ok?
* * Get() and Wait() are meant to talk to the same underlying data source (set of long-running ops that may be done or not done) - the examples are confusing cause they don't
* * This API is a pretty dump wrapper over that map.
* * The map is meant to hold onto stuff even after it's done, so that it can be queried
* * Get() - instant - is it done? If so, result. Sample code indicates this should delete the work item
* * Wait() - blocks - wait until it's done or timeout. If it becomes done in time, result
* * List() - optional - dump work queue (including done)
* * Cancel() - optional - best-effort, async cancel (use Get() to check). When/if cancelled, keep in the map and set error to CANCELLED
* * Delete() - optional - remove from map, ie client no longer interested in the *result*. Does not cancel it from happening
 */

type lroServer struct {
	log logr.Logger
}

func (s *lroServer) ListOperations(ctx context.Context, in *lropb.ListOperationsRequest) (*lropb.ListOperationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "List is not implemented")
}

func (s *lroServer) GetOperation(ctxt context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
	s.log.V(1).Info("GetOperation", "name", in.Name)

	notDoneResponse := &lropb.Operation{
		Name: in.Name,
		Done: false,
	}

	item, ok := workItems[in.Name]
	if !ok {
		return notDoneResponse, status.Errorf(codes.NotFound, "Versions request %s not found.", in.Name)
	}

	// lol hack - sync innit

	var ss []fetch.Series
	var err error
	switch app := item.Req.GetApp().(type) {
	case *versionspb.VersionsRequest_Github:
		gh := app.Github
		s.log.V(1).Info("GetVersions", "type", "github", "address", gh.GetOrg()+"/"+gh.GetRepo())

		if gh.GetRepo() == "go" || gh.GetRepo() == "pfsense" || gh.GetRepo() == "nginx" {
			// FIXME hack. The GH website shows Releases for golang/go, but the API returns an empty list...
			// fix by falling back to tags if release names (?) are an empty list
			ss, err = fetch.Github(s.log, gh.GetOrg(), gh.GetRepo(), true, gh.GetDepth(), gh.GetCount())
		} else {
			ss, err = fetch.Github(s.log, gh.GetOrg(), gh.GetRepo(), false, gh.GetDepth(), gh.GetCount())
		}
	case *versionspb.VersionsRequest_Linux:
		s.log.V(1).Info("GetVersions", "type", "linux")

		ss, err = fetch.Linux()
	}
	if err != nil {
		s.log.V(1).Error(err, "can't fetch")
		// TODO: translate errors, eg dealine exceeded (DeadlineExceeded) vs json parse error (InvalidArgument), etc
		return notDoneResponse, status.Errorf(codes.Unknown, "Error fetching versions for request %s: %v.", in.Name, err)
	}

	value, _ := anypb.New(
		&versionspb.VersionsResponse{Serieses: series2proto(ss)},
	)
	resp := &lropb.Operation{
		Name:   in.Name,
		Done:   true,
		Result: &lropb.Operation_Response{Response: value},
	}

	s.log.Info("LRO Complete", "name", in.Name)
	delete(workItems, in.Name)

	return resp, nil
}

func series2proto(ss []fetch.Series) (res []*versionspb.Series) {
	res = make([]*versionspb.Series, len(ss))

	for i, s := range ss {
		rs := []*versionspb.Release{}
		for n, r := range s.Releases {
			d := timestamppb.New(r.Date)
			rs = append(rs, &versionspb.Release{Name: n, Version: r.Version.String(), Date: d})
		}

		res[i] = &versionspb.Series{
			Name:     s.Name,
			Releases: rs,
		}
	}
	return

}

func (s *lroServer) WaitOperation(ctxt context.Context, in *lropb.WaitOperationRequest) (*lropb.Operation, error) {
	s.log.V(1).Info("WaitOperation", "name", in.Name)

	var resp *lropb.Operation
	var err error

	if _, ok := workItems[in.Name]; ok {
		// Stuff an empty op in here? Is that why Get delete()s?
		resp = &lropb.Operation{
			Name:   in.Name,
			Done:   true,
			Result: &lropb.Operation_Response{},
		}
		err = nil
	} else {
		resp = &lropb.Operation{
			Name: in.Name,
			Done: false,
		}
		err = status.Errorf(codes.NotFound, "Versions request %s not found.", in.Name)
	}

	return resp, err
}

func (s *lroServer) CancelOperation(ctx context.Context, in *lropb.CancelOperationRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "Cancel is not implemented")
}
func (s *lroServer) DeleteOperation(ctx context.Context, in *lropb.DeleteOperationRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "Delete is not implemented")
}

type healthServer struct {
	log logr.Logger
}

func (s *healthServer) Check(ctxt context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	s.log.V(2).Info("grpc Health::Check()")
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}
func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

func main() {
	log := usvc.GetLogger(true, 0)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Error(err, "failed to listen", "addr", grpcAddr)
		os.Exit(1)
	}

	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}
	s := grpc.NewServer(sopts...)

	versionspb.RegisterVersionsServer(s, &versionsServer{log: log})
	healthpb.RegisterHealthServer(s, &healthServer{log})
	lropb.RegisterOperationsServer(s, &lroServer{log})

	// TODO: argh what the shutdown?
	// TODO: add grpc server to golang-graceful-shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		s.GracefulStop()
		log.Info("Shutting down...")
		close(idleConnsClosed)
	}()

	log.Info("Starting gRPC server", "addr", grpcAddr)
	_ = s.Serve(lis)
}
