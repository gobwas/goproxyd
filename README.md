# goproxyd

> Go modules proxy daemon

# About

This is a tiny Go modules proxy server that helps with local builds of Go
programs which need the use of private repos.

Under the hood it just returns **already present** entries in
`$GOPATH/pkg/mod/cache/download` directory.

It is based on awesome [goproxy](https://github.com/goproxy/goproxy)
implementation.

# How to

Essentially this was built to help with local docker builds. So if it is your
case just do the following (assuming that you have downloaded all required
modules on the host machine by running `go mod download` or just by `go
build`):

```
$ docker build -t goproxyd .
$ docker network create goproxyd
$ docker run -it --rm --name goproxyd --network goproxyd --publish 8080:8080 -v $GOPATH/pkg/mod/cache:/cache goproxyd
```

After this you can build you program docker:

```
docker build --network goproxyd --build-arg GOPROXY=http://goproxyd:8080 --build-arg GONOSUMDB=github.com/your-org-here . 
```


