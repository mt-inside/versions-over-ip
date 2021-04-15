package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	versions "github.com/mt-inside/versions-over-ip/api/v1alpha1"
	"github.com/mt-inside/versions-over-ip/api/v1alpha1/client"
)

const trace_enabled = false

func trace(fmt string, args ...interface{}) {
	if trace_enabled { // from args of etcd
		log.Printf(fmt, args...)
	}
}

func main() {
	ctxt := context.Background()

	client, err := client.NewVersionsClient(ctxt,
		option.WithGRPCDialOption(grpc.WithInsecure()),
		option.WithoutAuthentication(),
	)
	if err != nil {
		log.Fatalf("Could not make gapic client: %v", err)
	}

	var ss []*versions.Series
	ss = fetchLinux(client)
	render("Linux", ss)
	//ss = fetch(client, "zfsonlinux", "zfs", 2, 2)
	//render(ss)
	ss = fetchGithub(client, "golang", "go", 2, 2)
	render("golang", ss)
	ss = fetchGithub(client, "kubernetes", "kubernetes", 2, 5)
	render("kubernetes", ss)
	ss = fetchGithub(client, "helm", "helm", 2, 2)
	render("helm", ss)
	ss = fetchGithub(client, "envoyproxy", "envoy", 2, 2)
	render("envoy", ss)
	ss = fetchGithub(client, "istio", "istio", 2, 2)
	render("istio", ss)
	//ss = fetchGithub(client, "linkerd", "linkerd2", 1, 2)
	//render(ss)
	ss = fetchGithub(client, "hashicorp", "terraform", 2, 2)
	render("terraform", ss)
}

func fetchGithub(
	client *client.VersionsClient,
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

	trace("Initiated fetch: %s", resplro.Name())

	value, err := resplro.Wait(ctxt)
	if err != nil {
		log.Fatalf("Could not wait synchronously: %v", err)
	}

	trace("Done: %v", resplro.Done())

	return value.Serieses
}
func fetchLinux(
	client *client.VersionsClient,
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

	trace("Initiated fetch: %s", resplro.Name())

	value, err := resplro.Wait(ctxt)
	if err != nil {
		log.Fatalf("Could not wait synchronously: %v", err)
	}

	trace("Done: %v", resplro.Done())

	return value.Serieses
}

func render(name string, ss []*versions.Series) {
	c := color.New(color.FgHiWhite).Add(color.Bold)
	c.Printf("== %s ==\n", name)

	for _, s := range ss {
		fmt.Printf("%s: ", s.GetName())
		for _, r := range s.GetReleases() {
			d := r.GetDate().AsTime().Local()

			if isPreReleaseRelease(r.GetName()) || isPreReleaseSeries(s.GetName()) {
				color.Set(color.FgHiBlack)
			} else if isLTSRelease(r.GetName()) || isLTSSeries(s.GetName()) {
				color.Set(color.FgBlue)
			}
			fmt.Printf("%s %s (%d days ago)", r.GetName(), r.GetVersion(), int(time.Since(d).Hours())/24)
			color.Unset()
			fmt.Printf("  ")
		}
		fmt.Println()
	}

	fmt.Println()
}

// TODO: hardcoding bad. Fetchers should register their appropriate strings? Or yanno, implement an interface (on the release objects they return), where we can dispatch to this
// TODO: pre-release-ness can be an attribute of the release (GH) or the series (linux). The linux release object should have a pointer to its series, and the method on that struct should answer it using that
func isPreReleaseRelease(name string) bool {
	return name == "pre"
}

func isLTSRelease(name string) bool {
	return false
}

func isPreReleaseSeries(name string) bool {
	return name == "mainline"
}

func isLTSSeries(name string) bool {
	return name == "longterm"
}
