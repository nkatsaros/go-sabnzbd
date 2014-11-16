package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nkatsaros/go-sabnzbd"
)

func main() {
	var (
		addrFlag   = flag.String("addr", "localhost:8080", "sabnzbd address")
		apikeyFlag = flag.String("apikey", "", "sabnzbd apikey")
		nzbFlag    = flag.String("nzb", "", "nzb file to upload")
	)

	flag.Parse()

	if *apikeyFlag == "" {
		fmt.Println("apikey is required")
		flag.Usage()
		os.Exit(1)
	}

	if *nzbFlag == "" {
		fmt.Println("nzb is required")
		flag.Usage()
		os.Exit(1)
	}

	s, err := sabnzbd.New(sabnzbd.Addr(*addrFlag), sabnzbd.ApikeyAuth(*apikeyFlag))
	if err != nil {
		log.Fatalln("couldn't create sabnzbd:", err)
	}

	auth, err := s.Auth()
	if err != nil {
		log.Fatalln("couldn't get auth type:", err)
	}

	if auth != "apikey" {
		log.Fatalln("sabnzbd instance must be using apikey authentication")
	}

	_, err = s.AddFile(*nzbFlag)
	if err != nil {
		log.Fatalln("failed to upload nzb", err)
	}
}
