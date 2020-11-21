package main

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/grpc"

	versions "github.com/mt-inside/versions-over-ip/api/v1alpha1"
	"github.com/mt-inside/versions-over-ip/cmd/client/versionsclient"
)

func main() {
	ctxt := context.Background()

	client, err := versionsclient.NewVersionsClient(ctxt,
		//option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		option.WithGRPCDialOption(grpc.WithInsecure()),
		option.WithoutAuthentication(),
	)
	if err != nil {
		log.Fatalf("Could not make gapic client: %v", err)
	}

	resplro, err := client.GetVersions(ctxt, &versions.VersionsRequest{Org: "kubernetes", Repo: "kubernetes"})
	if err != nil {
		log.Fatalf("Could not initiate request for versions: %v", err)
	}

	log.Printf("Initiated fetch: %s", resplro.Name())

	value, err := resplro.Wait(ctxt)
	if err != nil {
		log.Fatalf("Could not wait synchronously: %v", err)
	}

	log.Printf("Done: %v", resplro.Done())
	log.Printf("%v", value)
}
