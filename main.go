package main

import (
	"log"

	"github.com/nkbhasker/go-auth-starter/cmd"
)

var version string = "dev"

func main() {
	cmd.Version(version)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
