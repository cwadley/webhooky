# webhooky
Vomits the request payloads of incoming webhook requests to the console
----
## Usage
To use without compiling, clone this repository, then run:
`go run webhooky.go`

...or to compile then run:
`go build webhooky.go`
`./webhooky`

Command line options:
```
--port 		<port to listen on> (required)
--endpoint 	<endpoint to listen to> (optional, / is default)
--cert	<path to SSL cert> (optional)
--key 	<path to SSL key> (optional)
--body <string containing a custom response body to be returned> (optional)
```
### Docker container
`docker run -it cwadley/webhooky:latest --port <port_num> [--endpoint, --cert, --key, --body]`

If SSL certificates are to be used, they must be volume mounted into the container:
`docker run -it -v $(pwd)/certs:/certs cwadley/webhooky:latest --port 443 --cert /certs/mycert.pem --key /certs/mykey.pem`
