[ ! -d api-common-protos ] && git clone https://github.com/googleapis/api-common-protos

# protoc (OS dependant)
go get -v google.golang.org/protobuf/cmd/protoc-gen-go@v1.26.0
go get -v google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get -v github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic@v0.20.1 # versions of google.golang.org/api and protoc-gen-go_gapic are sensitive to each other
