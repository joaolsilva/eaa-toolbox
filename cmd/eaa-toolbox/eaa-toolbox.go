package main

import (
	"log"
	"r2discover.com/go/eaa-toolbox/pkg/toolbox"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	t := toolbox.NewToolbox()
	t.Start()
}
