# Current workarounds

gapic (client)
* oneof - `md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v&%s=%v", "org", url.QueryEscape(req.App.(*apipb.VersionsRequest_Github).Github.GetOrg()), "repo", url.QueryEscape(req.App.(*apipb.VersionsRequest_Github).Github.GetRepo())))`
  * hack: delete md and following line
* bug in golang import
  * hack: rename import alias


# Raison d'être
Was born as a GAPIC test, and a NATS test.
Can also be the go-to example for manual proto/grpc generation (cf go-grpc-bazel-example)
* try gapic, gogo, etc
* try proto plugins, extensions, validators

## Experience Report
Dropping gogo cause it doesn't support new features and looks to be dying


# Architecture
Seems weird to have a client that makes a network call to a server, that then makes another network call, but yanno it's an experiement
Also it allows multiple clients and moves the logic about aggregating versions off them
Also the server could start polling GH, for a cache, push notifications on a change, etc
