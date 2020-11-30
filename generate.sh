# --include_imports - copy imports to descriptor
# go_out talks to protoc-gen-go, which is used for server stubs
# go_gapic_[out,opt] talks to protoc-gen-go_gapic, which is used for client stubs

# TODO split invocations client/server
# This is all messy cause server and client share the same code base. Ideally you'd have an API repo, and then separate client and server codebases, which would vendor snapshots of those protos (yay forwards/backwards compat)
# - We currently build proto + grpc (client and server) stubs to /api, the server uses them, the client ignores them but the gapic stubs it calls instead use them

# go_out - versions.pb.go - the structs
# go-grpc_out - versions_grpc.pb.go - gRPC server and client (we only use the server)
# go_gapic_out - versionsclient/ - gAPIC super client based on gax
protoc \
    -I "api-common-protos" \
    -I "api/v1alpha1" \
    --experimental_allow_proto3_optional \
    --descriptor_set_out="api/v1alpha1/versions.proto.pb" \
    --include_imports \
    --go_out="api/v1alpha1" \
    --go_opt=paths=source_relative \
    --go-grpc_out="api/v1alpha1" \
    --go-grpc_opt=paths=source_relative \
    --go_gapic_out="cmd/client" \
    --go_gapic_opt="go-gapic-package=versionsclient;versionsclient" \
    --go_gapic_opt="grpc-service-config=cmd/client/versions_grpc_service_config.json" \
    "api/v1alpha1/versions.proto"

#go-gapic-package=github.com/mt-inside/versions-over-ip/cmd/client/versionsclient;versionsclient,module=github.com/mt-inside/versions-over-ip
#go-gapic-package=path/to/out;pkg,module=path,transport=rest+grpc,gapic-service-config=gapic_cfg.json,release-level=alpha
