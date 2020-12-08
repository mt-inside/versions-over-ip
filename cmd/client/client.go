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
	ss = fetchLinux(client)
	render(ss)
	//ss = fetch(client, "zfsonlinux", "zfs", 2, 2)
	//render(ss)
	ss = fetchGithub(client, "golang", "go", 2, 2)
	render(ss)
	ss = fetchGithub(client, "kubernetes", "kubernetes", 2, 5)
	render(ss)
	ss = fetchGithub(client, "helm", "helm", 2, 2)
	render(ss)
	ss = fetchGithub(client, "envoyproxy", "envoy", 2, 2)
	render(ss)
	ss = fetchGithub(client, "istio", "istio", 2, 2)
	render(ss)
	//ss = fetchGithub(client, "linkerd", "linkerd2", 1, 2)
	//render(ss)
	ss = fetchGithub(client, "hashicorp", "terraform", 2, 2)
	render(ss)
}

func fetchGithub(
	client *versionsclient.VersionsClient,
	org string, repo string,
	depth, count int32,
) []*versions.Series {
	ctxt := context.Background()

	resplro, err := client.GetVersions(
		ctxt,
		&versions.VersionsRequest{
			App: &versions.VersionsRequest_Github{
				Github: &versions.GithubRepo{
					Org: org, Repo: repo, Depth: depth, Count: count,
				},
			},
		},
	)
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
func fetchLinux(
	client *versionsclient.VersionsClient,
) []*versions.Series {
	ctxt := context.Background()

	resplro, err := client.GetVersions(
		ctxt,
		&versions.VersionsRequest{
			App: &versions.VersionsRequest_Linux{},
		},
	)
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
