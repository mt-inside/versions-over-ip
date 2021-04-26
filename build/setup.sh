rm -rf api-common-protos && git clone git@github.com:googleapis/api-common-protos.git

# protoc (OS dependant)
go get -v google.golang.org/protobuf/cmd/protoc-gen-go@v1.26.0
go get -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.37.0
go get -v github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic@v0.19.0 # versions of google.golang.org/api and protoc-gen-go_gapic are sensitive to each other
