// Copyright (c) 2015-2021, NVIDIA CORPORATION.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"

	"golang.org/x/sys/unix"

	"github.com/NVIDIA/proxyfs/ramswift"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("no .conf file specified")
	}

	doneChan := make(chan bool, 1) // Must be buffered to avoid race

	go ramswift.Daemon(os.Args[1], os.Args[2:], nil, doneChan, unix.SIGINT, unix.SIGTERM, unix.SIGHUP)

	_ = <-doneChan
}
