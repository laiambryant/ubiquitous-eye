# ubiquitous-eye

A small portfolio application, used to test opentelemetry with go.

## Running in server mode

To run the application you can either run it in the docker container through the 


``` bash
docker compose up 
```

command while in the root directory. When running containerized version of the app you will have access to the prometheus and Jaeger UIs reachable at the following addresses:

- http://localhost:16686/ Jaeger UI
- http://localhost:9090/ Prometheus UI

If you have go installed on your machine run:

``` bash
go run main.go server-mode
```

while in the src folder.

## Running in batch mode

The program can also be run in batch mode to produce the static website that will be served by github pages.

Use the following command:

``` bash
go run main.go
```

Running in batch mode will produce an index.html file that will then be placed in the docs folder and served by github pages.