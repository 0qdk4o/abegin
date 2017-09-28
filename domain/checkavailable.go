// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package domain

import (
	"bytes"
	"encoding/json"
	"strings"
)

// Checkavailable represents operation type
type Checkavailable Cmd

// SplicingURL splice pieces of variable to long url
func (a *Checkavailable) SplicingURL(p *Registrar, domains []string) string {
	if len(domains) < 1 || p == nil {
		return ""
	}
	var t = Cmd(*a)
	url := bytes.NewBufferString((&t).SplicingURL(p, domains))

	uniqToplv := make(map[string]bool)
	uniqSecondlv := make(map[string]bool)
	for _, v := range domains {
		namesSlice := strings.Split(v, ".")
		namesLen := len(namesSlice)
		if namesLen < 2 {
			return ""
		}

		if !uniqSecondlv[namesSlice[namesLen-2]] {
			url.WriteString("&domain-name=")
			url.WriteString(namesSlice[namesLen-2])
			uniqSecondlv[namesSlice[namesLen-2]] = true
		}

		if !uniqToplv[namesSlice[namesLen-1]] {
			url.WriteString("&tlds=")
			url.WriteString(namesSlice[namesLen-1])
			uniqToplv[namesSlice[namesLen-1]] = true
		}
	}

	return url.String()
}

// ProcessHTTPData should implemented by derived struct
func (a *Checkavailable) ProcessHTTPData(b []byte) (*Result, error) {
	var res Result
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	_, err := decoder.Token() // return '{'
	if err != nil {
		return nil, err
	}

	st, err := decoder.Token() // check if status field is output in the json
	if err != nil {
		return nil, err
	}

	if strings.Compare(st.(string), "status") == 0 {
		err = json.Unmarshal(b, &res)
	} else {
		err = json.Unmarshal(b, &res.NameRes)
	}

	return &res, err
}

func (a *Checkavailable) String() string {
	return string(*a)
}
