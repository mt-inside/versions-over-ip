package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

type versionsServer struct{}

func (s *versionsServer) GetVersions(ctxt context.Context, in *versionspb.VersionsRequest) (*lropb.Operation, error) {
	log.Printf("GetVersions(%s/%s)", in.Org, in.Repo)

	uid, _ := uuid.NewUUID()
	work := workItem{ID: uid.String(), Req: in}
	workItems[uid.String()] = work

	res := &lropb.Operation{
		Name: uid.String(),
		Done: false,
	}

	log.Printf("Long-running operation started: %s", uid.String())

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

type lroServer struct{}

func (s *lroServer) ListOperations(ctx context.Context, in *lropb.ListOperationsRequest) (*lropb.ListOperationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "List is not implemented")
}

func (s *lroServer) GetOperation(ctxt context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
	log.Printf("GetOperation: %s", in.Name)

	var resp *lropb.Operation
	var err error

	if item, ok := workItems[in.Name]; ok {
		// lol hack
		ss, err := fetch.Github(item.Req.Org, item.Req.Repo, item.Req.Depth, item.Req.Count)
		log.Printf("%v", ss)
		if err != nil {
			resp = &lropb.Operation{
				Name: in.Name,
				Done: false,
			}
			log.Printf("error fetching: %v", err)
			err = status.Errorf(codes.InvalidArgument, "Error fetching versions for request %s: %v.", in.Name, err)
		} else {
			value, _ := ptypes.MarshalAny(
				&versionspb.VersionsResponse{Serieses: series2proto(ss)},
			)
			resp = &lropb.Operation{
				Name:   in.Name,
				Done:   true,
				Result: &lropb.Operation_Response{Response: value},
			}
			err = nil

			log.Printf("LRO Complete: %s", in.Name)
			delete(workItems, in.Name)
		}
	} else {
		resp = &lropb.Operation{
			Name: in.Name,
			Done: false,
		}
		err = status.Errorf(codes.NotFound, "Versions request %s not found.", in.Name)
	}

	return resp, err
}

func series2proto(ss []fetch.Series) (res []*versionspb.Series) {
	res = make([]*versionspb.Series, len(ss))
	for i, s := range ss {
		res[i] = &versionspb.Series{
			Prefix:     *mapVersionString(s.Prefix),
			Stable:     mapVersionString(s.Stable),
			Prerelease: mapVersionString(s.Prerelease),
		}
	}
	return
}

func mapVersionString(v *version.Version) *string {
	if v == nil {
		return nil
	} else {
		/* ffs golang */
		s := v.String()
		return &s
	}
}

func (s *lroServer) WaitOperation(ctxt context.Context, in *lropb.WaitOperationRequest) (*lropb.Operation, error) {
	log.Printf("WaitOperation: %s", in.Name)

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

type healthServer struct{}

func (s *healthServer) Check(ctxt context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	log.Printf("grpc Health::Check()")
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}
func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

func main() {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}

	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}
	s := grpc.NewServer(sopts...)

	versionspb.RegisterVersionsServer(s, &versionsServer{})
	healthpb.RegisterHealthServer(s, &healthServer{})
	lropb.RegisterOperationsServer(s, &lroServer{})

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		s.GracefulStop()
		log.Printf("Shutting down...")
		close(idleConnsClosed)
	}()

	log.Printf("Starting gRPC server on %s", grpcAddr)
	s.Serve(lis)
}
