# webhooky
Vomits the request payloads of incoming webhook requests to the console
----

To use, clone this repository, then run:

go run webhooky.go

Command line options:
```
--port 		<port to listen on> (required)
--endpoint 	<endpoint to listen to> (optional, / is default)
--sslCert	<path to SSL cert> (optional)
--sslKey 	<path to SSL key> (optional)
```