package main

import (
	"k8s-demo-emp-api/api"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: k8s-demo-emp-api serve-api|cluster-tests")
	}
	cmd := os.Args[1]
	if cmd == "serve-api" {
		api.ServeApi()
	} else {
		log.Fatalf("Unknown command %s\n", cmd)
	}
}
