rm -rf api-common-protos && git clone git@github.com:googleapis/api-common-protos.git

# protoc?
go get -u -v google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
go get -u -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.35.0
go get -u -v github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic@v0.17.0 # versions of google.golang.org/api and protoc-gen-go_gapic are sensitive to each other
