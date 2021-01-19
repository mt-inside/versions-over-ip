module github.com/mt-inside/versions-over-ip

go 1.15

// versions of google.golang.org/api and protoc-gen-go_gapic are sensitive to each other
require (
	cloud.google.com/go v0.72.0
	github.com/fatih/color v1.10.0
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/hashicorp/go-version v1.2.1
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/sys v0.0.0-20210113181707-4bcb84eeeb78 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/api v0.36.0
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.1-0.20200805231151-a709e31e5d12
)
