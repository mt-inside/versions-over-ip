# TODO
_see GH_

Direction:
* Maybe forget the NATS part; that was meant to be for scatter/gather in-line queries, whereas now it's going to be a monitor. Should still offer LRO API, and should do requests in-line (on a background thread). Always runs in k8s; CLI can hit it any time. Polling client runs as k8s CronJob, asks it about each repo in turn, reports to slack.
  * should report them all in one go, so should serialise its progress in case of crashes and restart on the last repo it was attempting - learn to do this
  * something should track difference-since-last-time and be able to tell you what's new. Since the CLI wants to do this too, feels like the server should do it.
    * "reverse event source" - every time the server is asked to look at anything, it compares to the last version it saw, records an event if it's changed, and sticks that somewhere - kafka/nats? Each event gets an ID, and clients can ask to check a set of repos and get updates since their last seen ID

# Current workarounds

gapic (client)
* oneof - `md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v&%s=%v", "org", url.QueryEscape(req.App.(*apipb.VersionsRequest_Github).Github.GetOrg()), "repo", url.QueryEscape(req.App.(*apipb.VersionsRequest_Github).Github.GetRepo())))`
  * reason: gapic doesn't support OneOf yet (think there's a GH issue for it)
  * hack: delete md and following line


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
