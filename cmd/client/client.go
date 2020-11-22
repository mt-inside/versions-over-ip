package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/grpc"

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
		fmt.Printf("%s: ", s.GetPrefix())
		if s.GetStable() != "" { // "" is how optionals are presented, not nil
			fmt.Printf("%s", s.GetStable())
		}
		if s.GetPrerelease() != "" /* && s.Prerelease.GreaterThan(s.Stable) */ {
			fmt.Printf(" (PRE %s)", s.GetPrerelease())
		}
		fmt.Printf("\n")
	}
}
