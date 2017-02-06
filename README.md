![FishFinger](assets/fishfinger-logo-640-218.png)

[![GoDoc](https://godoc.org/github.com/timtosi/fishfinger?status.svg)](https://godoc.org/github.com/timtosi/fishfinger)
[![Go Report Card](https://goreportcard.com/badge/github.com/timtosi/fishfinger)](https://goreportcard.com/report/github.com/timtosi/fishfinger)

## What is FishFinger ?
FishFinger is a Dockerfile and Docker-Compose lightweight programmatic library
written in Go. This project provide easy to use abstractions of official Docker
libraries with no other dependency than [official Go library](https://github.com/golang/go)
and [official Docker libraries](https://github.com/docker/).

FishFinger also provide a one possible solution to deal with [container initialization delay](https://github.com/docker/compose/issues/374)
and makes room for further improvement through user-defined solutions.

This project is fully tested and provide an [extensive documentation](https://godoc.org/github.com/timtosi/fishfinger).

## Quickstart

First of all, you need to install and configure the latest [Docker version](https://docs.docker.com/engine/installation/)
on your machine.

Then go get this repository:
```golang
go get -d github.com/timtosi/fishfinger
```

Finally, check how to use the components provided:

* [Use one or more services from a Compose file](examples/compose-basic/main.go)
* [Use one or more services from a Compose file with a Backoff](examples/compose-backoff/main.go)

## Use a Compose File

The `fishfinger.Compose` component represent an abstraction of the [libcompose](https://github.com/docker/libcompose)
library. It provides functions allowing the user to use a [Compose file](https://docs.docker.com/compose/compose-file/)
programmatically.

An extensive documentation of this component is available [here](https://godoc.org/github.com/TimTosi/fishfinger).

### Basic Compose Usage

Let's say you have a basic Compose file where you defined N services.

// HERE COMPOSE FilE example

For instance, you could want to be able to write a test suite that will directly
make use of your containers instead of simply mocking their behaviour through
burdensome to write functions.

First of all, let's create a new Compose file handler.

```go
func main() {
	c, err := fishfinger.NewCompose("./docker-compose.yaml")
	if err != nil {
		log.Fatal(err)
	}
}
```

The `fishfinger.NewCompose` function takes the path to your Compose file as an
argument.

Time to start services.

```go
func main() {
	c, err := fishfinger.NewCompose("./docker-compose.yaml")
	if err != nil {
		log.Fatal(err)
	}
    if err := c.Start(); err != nil {
		log.Fatal(err)
	}
}
```

As no argument is provided, all services will be started, but `Compose.Start`
function accepts a variable number of service name in argument and will start
them in the order specified.

At this stage, you should have a neat redis service running. The Redis client
port you previously set in your image is the `6379` but Docker port forwarding
will take another one available on your host machine. The `Compose.Port` function
is here to find the correct port to use from a service name and the combination
of the port exposed and the protocol used such as `6379/tcp`.

```go
func main() {
	c, err := fishfinger.NewCompose("./docker-compose.yaml")
	if err != nil {
		log.Fatal(err)
	}
    if err := c.Start(); err != nil {
		log.Fatal(err)
	}
    addr, err := c.Port("redis", "6379/tcp")
	if err != nil {
		log.Fatal(err)
	}
}
```

The `Port` function will return the full address in the following form
`<host>:<port>`. Now you can instanciate your redis client and do whatever stuff
you want.

In the end, you maybe want to clean the containers created. Just use the
`Compose.Stop` function just as you would use the `Compose.Start` function.

```go
func main() {
	c, err := fishfinger.NewCompose("./docker-compose.yaml")
	if err != nil {
		log.Fatal(err)
	}
    if err := c.Start(); err != nil {
		log.Fatal(err)
	}
    addr, err := c.Port("redis", "6379/tcp")
	if err != nil {
		log.Fatal(err)
	}
    if err := c.Stop(); err != nil {
		log.Fatal(err)
	}
}
```
As no argument is provided, all services will be stopped.

Complete working code can be found [here](examples/compose-basic/main.go).

## Dockerfile

/!\ Coming soon /!\

## License
Every file provided here is available under the [MIT License](http://opensource.org/licenses/MIT).

## Not Good Enough ?
If you encouter any issue by using what is provided here, please
[let me know](https://github.com/TimTosi/fishfinger/issues) ! 
Help me to improve by sending your thoughts to timothee.tosi@gmail.com !
