package main

import (
	"log"
	"os"
	"sync"
)

// read nodes.json file into the global variable gTLS
// then begin logging on gTLS.log channel
// then start each node and wait until they all exit
func main() {
	log.SetPrefix("goTunneLS: ")
	gTLS := new(goTunneLS)
	if err := os.Chdir("/etc/goTunneLS"); err != nil {
		log.Fatal(err)
	}
	gTLS.parseFile("nodes.json")
	if gTLS.Logfile != "" {
		gTLS.logInterface = make(chan []interface{})
	}
	go gTLS.receiveAndLog()
	var nodeWG sync.WaitGroup
	nodeWG.Add(len(gTLS.Nodes))
	for _, n := range gTLS.Nodes {
		gTLS.log("--> initalizing", n.Mode, "node", n.Name)
		// prepend space to name in named nodes to separate mode in logging
		if n.Name != "" {
			n.Name = " " + n.Name
		}
		n.logInterface = gTLS.logInterface
		n.nodeWG = nodeWG
		gTLS.log("--> starting", n.Mode, "node"+n.Name)
		go n.run()
	}
	nodeWG.Wait()
	gTLS.log("--> exiting")
}
