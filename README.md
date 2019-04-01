# Echo Server
A simple echo server written in Go

## Building and Running
### Docker
If you are using docker, build the server image and run it with
```
docker build -t echo-server .
docker run -d -p 8080:8080 echo-server
```
Note: `echo-server` can be replaced with any tag you like

###  Without Docker
Run
```
dep ensure
go run .
```

