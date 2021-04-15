module github.com/mt-inside/versions-over-ip

go 1.15

// versions of google.golang.org/api and protoc-gen-go_gapic are sensitive to each other
require (
	cloud.google.com/go v0.72.0
	github.com/fatih/color v1.10.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.1.2
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/hashicorp/go-version v1.2.1
	golang.org/x/net v0.0.0-20210414194228-064579744ee0 // indirect
	golang.org/x/sys v0.0.0-20210415045647-66c3f260301c // indirect
	google.golang.org/api v0.36.0
	google.golang.org/genproto v0.0.0-20210414175830-92282443c685
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)
