# FIXME
No idea where versionsclient is getting `github.com/mt-inside/versions-over-ip/api` from

gogo (server) doesn't support:
* optional (field prescence)
gapic (client) doesn't support:
* oneof - `md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v&%s=%v", "org", url.QueryEscape(req.App.(*apipb.VersionsRequest_Github).Github.GetOrg()), "repo", url.QueryEscape(req.App.(*apipb.VersionsRequest_Github).Github.GetRepo())))`
  * hack: delete md and following line
* bug in golang import
  * hack: rename import alias

# TODO
scripted gen of stubs
* into api
* server - google new proto
* server - gogo
* client - gapic
proto (EN 5 steps for resrouce)
read config from yaml file
be calls can be in series for now
Don't download / vendor api-common-protos - how?
Move all packages to github.com/mt-inside/versions-over-ip
Move most stuff into pgk
.pb.go should be built into api? gprc / gapic stubs should not
make a hyper.js client for the fetch
release dates (make the versions in the proto in to a type rather than string)
gapic / gax: retrying on InvalidArg with it shouldn't, ignoring the timeout

proto pkg needs full name (.../api), build into api.

# Raison d'être
Was born as a GAPIC test, and a NATS test.
Can also be the go-to example for manual proto/grpc generation (cf go-grpc-bazel-example)
* try gapic, gogo, etc
* try proto plugins, extensions, validators

# Architecture
Seems weird to have a client that makes a network call to a server, that then makes another network call, but yanno it's an experiement
Also it allows multiple clients and moves the logic about aggregating versions off them
Also the server could start polling GH, for a cache, push notifications on a change, etc
