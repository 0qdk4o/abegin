// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/abegin/domain"
)

func main() {
	r := domain.NewRegistrar("Us2Demo")
	if r == nil {
		fmt.Println("no Us2Demo registrar info")
	}

	ret, err := r.DoRequest(&domain.CHKAVAI, []string{"google.com", "163.com"})
	if err != nil {
		fmt.Printf("[Debug]Unmarshal: %s\n", err)
		return
	}
	if ret == nil {
		panic("Unreachable branch")
	}
	if ret.Status != nil {
		fmt.Printf("[Status]: %s, [Message]: %s\n", *ret.Status, *ret.Message)
		return
	}
	if ret.NameRes == nil {
		fmt.Fprintf(os.Stderr, "no data\n")
	}
	fmt.Printf("%v\n", ret.NameRes)
}
