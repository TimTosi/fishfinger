![FishFinger](assets/fishfinger-logo-640-218.png)

[![GoDoc](https://godoc.org/github.com/timtosi/fishfinger?status.svg)](https://godoc.org/github.com/timtosi/fishfinger)
[![Go Report Card](https://goreportcard.com/badge/github.com/timtosi/fishfinger)](https://goreportcard.com/report/github.com/timtosi/fishfinger)

## What is FishFinger ?
FishFinger is a Docker and Docker-Compose lightweight programmatic library
written in Go. This project provide easy to use abstractions of official Docker
libraries with no other dependency than [official Go library](https://github.com/golang/go)
and [official Docker libraries](https://github.com/docker/).

Finally, this project is fully tested and provide an extensive documentation.

## Quickstart

Install the latest [Docker version](https://docs.docker.com/engine/installation/)
then go get this repository:
```go get -d github.com/timtosi/fishfinger```

Finally check how to use the components provided:

* [Use one or more services from a Compose file](examples/compose-basic/main.go)
* [Use one or more services from a Compose file with a Backoff](examples/compose-backoff/main.go)

## Compose

The `fishfinger.Compose` component represent an abstraction of the [libcompose](https://github.com/docker/libcompose)
library. It provides functions allowing the user to use a [Compose file](https://docs.docker.com/compose/compose-file/)
programmatically.

An extensive documentation of this component is available [here](https://godoc.org/github.com/TimTosi/fishfinger).

## Docker

/!\ Coming soon /!\

## License
Every file provided here is available under the [MIT License](http://opensource.org/licenses/MIT).

## Not Good Enough ?
If you encouter any issue by using what is provided here, please
[let me know](https://github.com/TimTosi/fishfinger/issues) ! 
Help me to improve by sending your thoughts to timothee.tosi@gmail.com !
