package main

import (
	"fmt"
	"log"

	"github.com/timtosi/fishfinger"
)

// -----------------------------------------------------------------------------

func main() {
	// First of all, let's create a new Compose file handler.
	c, err := fishfinger.NewCompose("./docker-compose.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Time to start services. As no argument is provided, all services will be
	// started.
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	// Now let's say you have a neat redis service running. The Redis client
	// port you previously set in your image is the `6379` but Docker port
	// forwarding will take another one available on your host machine.
	// The `Port` function will return the full address in the following
	// form `<host>:<port>`.
	addr, err := c.Port("redis", "6379/tcp")
	if err != nil {
		log.Fatal(err)
	}

	// Now you can instanciate your redis client and do whatever stuff you want.
	fmt.Printf("My redis is reachable at %s !\n", addr)

	// Finally, you maybe would stop and remove created containers.
	// As no argument is provided, all services will be stopped.
	if err := c.Stop(); err != nil {
		log.Fatal(err)
	}
}
