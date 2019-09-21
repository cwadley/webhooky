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
--sslCert	<path to SSL cert> (optional)
--sslKey 	<path to SSL key> (optional)
```
### Docker container
`docker run -it cwadley/webhooky:latest --port <port_num> [--endpoint, --sslCert, --sslKey]`

If SSL certificates are to be used, they must be volume mounted into the container:
`docker run -it -v $(pwd)/certs:/certs cwadley/webhooky:latest --port 443 --sslCert /certs/mycert.pem --sslKey /certs/mykey.pem`
