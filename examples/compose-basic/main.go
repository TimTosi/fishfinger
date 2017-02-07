package main

import (
	"fmt"
	"log"

	"github.com/timtosi/fishfinger"
)

// -----------------------------------------------------------------------------

func main() {
	// First of all, let's create a new Compose file handler by using the
	// `fishfinger.NewCompose` function that takes the path to your Compose file
	// as an argument.
	c, err := fishfinger.NewCompose("./docker-compose.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Time to start services. As no argument is provided, all services will be
	// started, but `Compose.Start` function accepts a variable number of
	// service name in argument and will start them in the order specified.
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	// At this stage, you should have a neat redis cluster running. The Redis
	// client ports you previously set for each of your services are the `6379`
	// but you had to specify different available ports on your host machine.
	// You can hard code the different ports OR you can use the `Compose.Port`
	// function to find the correct port to use from a service name and the
	// combination of the port exposed and the protocol used such as `6379/tcp`.
	addr, err := c.Port("redis-01", "6379/tcp")
	if err != nil {
		log.Fatal(err)
	}

	// The `Compose.Port` function will return the full address in the following
	// form `<host>:<port>`. Now you can instanciate your redis client and do
	// whatever stuff you want.
	fmt.Printf("My redis is reachable at %s !\n", addr)

	// In the end, you maybe want to clean the containers created. Just use the
	// `Compose.Stop` function just as you would use the `Compose.Start`
	// function. As no argument is provided, all services will be stopped.
	if err := c.Stop(); err != nil {
		log.Fatal(err)
	}
}
