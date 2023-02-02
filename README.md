# hookinterfaces

Finds foreman hosts that doesn't have an associated subnet for interfaces, searches for a fitting subnet
based on the interfaces IP address and assigns it.

## Building

Run the following to build the tool

    docker run --rm -e GOOS=<os> -e GOARCH=<arch> -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.19-alpine go build cmd/hookinterfaces.go

Replace os and arch with the os and architecture you're running on. For macOS on ARM, use e.g.:

    docker run --rm -e GOOS=darwin -e GOARCH=arm64 -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.19-alpine go build cmd/hookinterfaces.go

## Running

Run `hookinterfaces --help` to see all available options.
