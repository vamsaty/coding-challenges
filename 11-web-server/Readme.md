# Write Your Own "Web Server"

## Description
The challenge was to use build your HTTP server. 
The tool should be able to start a server and listen on a port.
It reads the files under the `www` directory and serves them as static files.
It supports a registry (mux) to register handlers for specific paths and methods.

NOTE: The mux matches the prefix of the request path. If there are paths starting
with the same prefix, the first one is returned. if "/" is registered and "/api" is
requested, "/" is returned. This can lead to only "/" being returned.

## Usage

Steps to build the binary and execute it -
1. `go build -o webserver .`
2. `./webserver --port <port> --host <host>`

## Flags

| Flag  | Description              | Default   |
|-------|--------------------------|-----------|
| -host | Host of the HTTP server  | localhost |
| -port | Port used by HTTP server | 9999      |

NOTE: 
