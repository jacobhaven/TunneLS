package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type goTunneLS struct {
	Nodes        []*node            // slice of nodes to run
	Logfile      string             // path to logfile, use /dev/stdout for standard output and /dev/stderr for standard error
	logInterface chan []interface{} // log channel
}

// read json file into gTLS
func (gTLS *goTunneLS) parseFile(path string) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(raw, gTLS); err != nil {
		log.Fatal(err)
	}
}

// receive on global log channel and append received to logfile
// if logfile doesn't exist, create it, and check continuously
// if it doesn't exist and if so create
func (gTLS *goTunneLS) receiveAndLog() {
	if gTLS.Logfile != "" {
		for {
			logFile, err := os.OpenFile(gTLS.Logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				log.Fatal(err)
			}
			defer logFile.Close()
			log.SetOutput(logFile)
			log.Println("--> global -/ beginning logging")
			for {
				if _, err = os.Stat(gTLS.Logfile); err != nil {
					if os.IsNotExist(err) {
						break
					}
					log.Fatal(err)
				}
				log.Println(<-gTLS.logInterface...)
			}
		}
	}
}

func (gTLS *goTunneLS) log(v ...interface{}) {
	gTLS.logInterface <- append([]interface{}{"--> global -/"}, v...)
}
