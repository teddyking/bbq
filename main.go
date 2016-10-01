package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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
	var enableDiskLimits bool
	var destroyDelay int

	flag.StringVar(&gardenAddress, "gardenAddr", "127.0.0.1:7777", "Garden server address")
	flag.IntVar(&numContainers, "numContainers", 10, "Number of containers")
	flag.IntVar(&destroyDelay, "destroyDelay", 0, "Time in seconds to wait before starting to destroy containers")
	flag.BoolVar(&enableDiskLimits, "enableDiskLimits", false, "Create containers with a disk limit")
	flag.Parse()

	gdn := newClient(gardenAddress)

	log.Printf("=== Creating %d containers sequentially ===\n", numContainers)
	var diskLimits garden.DiskLimits

	if enableDiskLimits {
		diskLimits = garden.DiskLimits{
			ByteSoft: uint64(10000000),
			ByteHard: uint64(20000000),
		}
	}

	for i := 0; i < numContainers; i++ {
		_, err := gdn.Create(garden.ContainerSpec{
			Handle: handle(i),
			Limits: garden.Limits{
				Disk: diskLimits,
			},
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

	if destroyDelay != 0 {
		log.Printf("=== Waiting %d seconds before starting to destroy containers ===\n", destroyDelay)
		time.Sleep(time.Second * time.Duration(destroyDelay))
		log.Printf("=== DONE ===\n")
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
