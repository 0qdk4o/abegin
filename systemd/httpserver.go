// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/abegin/systemd"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	heihei := fmt.Sprintf("Request proto %v from client %v\n", req.Proto, req.RemoteAddr)
	io.WriteString(w, "hello socket activated world!\n")
	io.WriteString(w, heihei+"\n")
}

func StartupTLS(ln net.Listener, c chan error) {
	mysrv := http.Server{Handler: http.HandlerFunc(HelloServer)}
	cer, err := tls.LoadX509KeyPair("/home/go/bin/server.crt",
		"/home/go/bin/server.key")
	if err != nil {
		c <- err
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	fmt.Println("Starting Https service...")
	err = mysrv.Serve(tls.NewListener(ln, config))
	if err != nil {
		c <- err
	}
}

func StartupHttp(ln net.Listener, c chan error) {
	mysrv := &http.Server{Handler: http.HandlerFunc(HelloServer)}
	fmt.Println("Starting Http service...")
	err := mysrv.Serve(ln)
	if err != nil {
		c <- err
	}
}

func RunServer(lnhttp, lnhttps net.Listener) chan error {
	cherr := make(chan error)

	go StartupTLS(lnhttps, cherr)
	go StartupHttp(lnhttp, cherr)

	return cherr
}

func main() {
	listeners, err := systemd.Listeners()
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: %v\n", err)
	}

	if len(listeners) != 2 {
		panic("Unexpected number of socket activation fds")
	}

	cherr := RunServer(listeners[0], listeners[1])

	select {
	case err := <-cherr:
		fmt.Fprintf(os.Stderr, "Could not start service %v\n", err)
		os.Exit(1)
	}
}
