package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes"
	versions "github.com/mt-inside/versions-over-ip/api/v1alpha1"
	"github.com/mt-inside/versions-over-ip/cmd/client/versionsclient"
)

func main() {
	ctxt := context.Background()

	client, err := versionsclient.NewVersionsClient(ctxt,
		option.WithGRPCDialOption(grpc.WithInsecure()),
		option.WithoutAuthentication(),
	)
	if err != nil {
		log.Fatalf("Could not make gapic client: %v", err)
	}

	var ss []*versions.Series
	//ss = fetch(client, "zfsonlinux", "zfs", 2, 2)
	//render(ss)
	ss = fetch(client, "golang", "go", 2, 2)
	render(ss)
	ss = fetch(client, "kubernetes", "kubernetes", 2, 5)
	render(ss)
	ss = fetch(client, "helm", "helm", 2, 2)
	render(ss)
	ss = fetch(client, "envoyproxy", "envoy", 2, 2)
	render(ss)
	ss = fetch(client, "istio", "istio", 2, 2)
	render(ss)
	//ss = fetch(client, "linkerd", "linkerd2", 1, 2)
	//render(ss)
	ss = fetch(client, "hashicorp", "terraform", 2, 2)
	render(ss)
}

func fetch(
	client *versionsclient.VersionsClient,
	org string, repo string,
	depth, count int32,
) []*versions.Series {
	ctxt := context.Background()

	resplro, err := client.GetVersions(ctxt, &versions.VersionsRequest{Org: org, Repo: repo, Depth: depth, Count: count})
	if err != nil {
		log.Fatalf("Could not initiate request for versions: %v", err)
	}

	log.Printf("Initiated fetch: %s", resplro.Name())

	value, err := resplro.Wait(ctxt)
	if err != nil {
		log.Fatalf("Could not wait synchronously: %v", err)
	}

	log.Printf("Done: %v", resplro.Done())

	return value.Serieses
}

func render(ss []*versions.Series) {
	for _, s := range ss {
		fmt.Printf("%s: ", s.GetName())
		for _, r := range s.GetReleases() {
			d, _ := ptypes.Timestamp(r.GetDate())

			fmt.Printf("%s %s (%d days ago)", r.GetName(), r.GetVersion(), int(time.Since(d).Hours())/24)
			// TODO: date as first class, calc here days from now()
			fmt.Printf(" | ")
		}
		fmt.Printf("\n")
	}
}
