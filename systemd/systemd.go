// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// More details refer to
// https://www.freedesktop.org/software/systemd/man/sd_listen_fds.html

// Package systemd get sockets from systemd manager
package systemd

import (
	"errors"
	"net"
	"os"
	"strconv"
	"syscall"
)

// #define SD_LISTEN_FDS_START 3
const (
	SdListenFdsStart = 3
)

// Internally, sd_listen_fds() checks whether the $LISTEN_PID environment
// variable equals the daemon PID. If not, it returns immediately. Otherwise,
// it parses the number passed in the $LISTEN_FDS environment variable, then
// sets the FD_CLOEXEC flag for the parsed number of file descriptors starting
// from SD_LISTEN_FDS_START. Finally, it returns the parsed number.
// sd_listen_fds_with_names() does the same but also parses $LISTEN_FDNAMES if set.

// Listeners returns sockets installed by .socket unit
func Listeners() ([]net.Listener, error) {
	defer os.Unsetenv("LISTEN_PID")
	defer os.Unsetenv("LISTEN_FDS")

	pid, err := strconv.Atoi(os.Getenv("LISTEN_PID"))
	if err != nil {
		return nil, err
	}
	if pid != os.Getpid() {
		return nil, errors.New("wrong pid from env")
	}

	fds, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
	if err != nil {
		return nil, err
	}
	if fds < 1 {
		return nil, errors.New("socket may be not ready?")
	}

	lns := make([]net.Listener, 0, fds)
	var errCheck error
	for fd := SdListenFdsStart; fd < SdListenFdsStart+fds; fd++ {
		// set the FD_CLOEXEC flag
		syscall.CloseOnExec(fd)
		file := os.NewFile(uintptr(fd), "LISTEN_FDS"+strconv.Itoa(fd))
		ln, err := net.FileListener(file)
		if err != nil {
			errCheck = err
			continue
		}
		lns = append(lns, ln)

		// refer to go doc `func FileListener`
		// FileListener returns a COPY of the network listener corresponding
		// to the open file f. It is the caller's responsibility to close ln
		// when finished. Closing ln does not affect f, and closing f does not
		// affect ln.
		file.Close()
	}
	return lns, errCheck
}
