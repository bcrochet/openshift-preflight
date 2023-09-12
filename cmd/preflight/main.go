package main

import (
	"context"
	"log"

	"github.com/redhat-openshift-ecosystem/preflight/cmd/preflight/root"
)

func main() {
	entrypoint := root.NewCommand()

	if err := entrypoint.ExecuteContext(context.Background()); err != nil {
		log.Fatal(err)
	}
}
