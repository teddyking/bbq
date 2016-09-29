package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/client"
	"code.cloudfoundry.org/garden/client/connection"
)

const handlePrefix = "bbq-"

func newClient(gardenAddress string) garden.Client {
	return client.New(connection.New("tcp", gardenAddress))
}

func main() {
	var gardenAddress string
	var numContainers int
	var debug bool

	flag.StringVar(&gardenAddress, "gardenAddr", "127.0.0.1:7777", "Garden server address")
	flag.IntVar(&numContainers, "numContainers", 10, "Number of containers")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.Parse()

	gdn := newClient(gardenAddress)

	log.Printf("=== Creating %d containers sequentially ===\n", numContainers)
	for i := 0; i < numContainers; i++ {
		_, err := gdn.Create(garden.ContainerSpec{
			Handle: handle(i),
		})

		if err != nil {
			log.Printf("Error creating container '%s' - %s\n", handle(i), err.Error())
			os.Exit(1)
		}

		fmt.Print(".")
	}
	fmt.Println("")
	log.Printf("=== DONE ===\n")

	log.Printf("=== Retrieving created containers\n")
	containers, err := gdn.Containers(garden.Properties{})
	if err != nil {
		log.Printf("Error getting containers - %s\n", err.Error())
		os.Exit(1)
	}
	log.Printf("=== DONE ===\n")

	if len(containers) != numContainers {
		log.Printf("Expected to find %d containers, but found %d\n", numContainers, len(containers))
		os.Exit(1)
	}

	log.Printf("=== Destroying %d containers sequentially ===\n", numContainers)
	for i := 0; i < numContainers; i++ {
		if err := gdn.Destroy(handle(i)); err != nil {
			log.Printf("Error destroying container '%s' - %s\n", handle(i), err.Error())
			os.Exit(1)
		}

		fmt.Print(".")
	}
	fmt.Println("")
	log.Printf("=== DONE ===\n")
}

func handle(index int) string {
	return fmt.Sprintf("%s%d", handlePrefix, index)
}
