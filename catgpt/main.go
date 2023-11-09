package main

import (
	"context"
	"os"

	"github.com/NickTaporuk/y-y-devops-trainings-cloud-1/pkg/gpt"
)

func main() {
	listenPublic := ":8090"
	if lp := os.Getenv("CATGPT_LISTEN_PUBLIC"); lp != "" {
		listenPublic = lp
	}

	// This listener should not be exposed to the internet.
	listenPrivate := ":9090"
	if lp := os.Getenv("CATGPT_LISTEN_PRIVATE"); lp != "" {
		listenPrivate = lp
	}

	defaultGPT := &gpt.CatGPT{}
	handler := gpt.NewHandler(defaultGPT)
	gpt.Serve(context.Background(), listenPublic, listenPrivate, handler)
}
