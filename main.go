package main

import (
	"fmt"
	"jview/engine"
	"jview/platform/darwin"
	"jview/transport"
	"log"
	"os"
	"runtime"
)

func main() {
	// macOS requires the main thread for AppKit
	runtime.LockOSThread()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: jview <file.jsonl>\n")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Initialize platform
	darwin.AppInit()
	disp := darwin.NewDispatcher()
	rend := darwin.NewRenderer()

	// Create session
	sess := engine.NewSession(rend, disp)

	// Create file transport
	ft := transport.NewFileTransport(filePath)

	// Process messages in a goroutine
	go func() {
		ft.Start()

		for {
			select {
			case msg, ok := <-ft.Messages():
				if !ok {
					log.Println("main: transport closed")
					return
				}
				sess.HandleMessage(msg)

			case err, ok := <-ft.Errors():
				if !ok {
					return
				}
				log.Printf("main: transport error: %v", err)
				return
			}
		}
	}()

	// Run the macOS event loop (blocks forever)
	darwin.AppRun()
}
