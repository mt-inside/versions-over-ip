# TODO
Don't download / vendor api-common-protos - how?
Move all packages to github.com/mt-inside/versions-over-ip
Move most stuff into pgk
.pb.go should be built into api? gprc / gapic stubs should not
make a hyper.js client for the fetch

proto pkg needs full name (.../api), build into api.

# Architecture
Seems weird to have a client that makes a network call to a server, that then makes another network call, but yanno it's an experiement
Also it allows multiple clients and moves the logic about aggregating versions off them
Also the server could start polling GH, for a cache, push notifications on a change, etc
