module github.com/mt-inside/versions-over-ip

go 1.15

// versions of google.golang.org/api and protoc-gen-go_gapic are sensitive to each other
require (
	cloud.google.com/go v0.72.0
	github.com/fatih/color v1.10.0
	github.com/go-logr/logr v0.4.0
	github.com/golang-commonmark/html v0.0.0-20180910111043-7d7c804e1d46 // indirect
	github.com/golang-commonmark/linkify v0.0.0-20180910111149-f05efb453a0e // indirect
	github.com/golang-commonmark/markdown v0.0.0-20180910011815-a8f139058164 // indirect
	github.com/golang-commonmark/mdurl v0.0.0-20180910110917-8d018c6567d6 // indirect
	github.com/golang-commonmark/puny v0.0.0-20180910110745-050be392d8b8 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.1.2
	github.com/googleapis/gapic-generator-go v0.19.0 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/hashicorp/go-version v1.2.1
	github.com/mt-inside/go-usvc v0.0.4
	github.com/opennota/wd v0.0.0-20180911144301-b446539ab1e7 // indirect
	golang.org/x/net v0.0.0-20210414194228-064579744ee0 // indirect
	golang.org/x/sys v0.0.0-20210415045647-66c3f260301c // indirect
	google.golang.org/api v0.36.0
	google.golang.org/genproto v0.0.0-20210414175830-92282443c685
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)
