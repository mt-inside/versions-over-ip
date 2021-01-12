rm -rf api-common-protos && git clone git@github.com:googleapis/api-common-protos.git

# protoc?
go get -u -v google.golang.org/protobuf/cmd/protoc-gen-go
go get -u -v google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get -u -v github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic
