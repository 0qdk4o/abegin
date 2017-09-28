// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package domain

import (
	"bytes"
	"crypto/sha1"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestJSONReadWrite(t *testing.T) {
	jsonbuff, err := ioutil.ReadFile(JSONFile)
	if err != nil {
		t.Fatal(err)
	}
	OriginFileSHA1 := sha1.Sum(jsonbuff)
	t.Logf("OriginFileSHA1: %v", OriginFileSHA1)

	ptemp, err := ioutil.TempFile("", "TestJSONReadWrite")
	if err != nil {
		t.Fatal(err)
	}
	DestFile := ptemp.Name()
	ptemp.Close()
	defer os.Remove(DestFile)

	t.Logf("Temp File: %s\n", DestFile)
	err = WriteRegistrarsInfo2File(ptemp.Name())
	if err != nil {
		t.Fatal(err)
	}

	afterbuff, err := ioutil.ReadFile(JSONFile)
	if err != nil {
		t.Fatal(err)
	}
	WriteFileSHA1 := sha1.Sum(afterbuff)
	t.Logf("WriteFileSHA1 : %v", WriteFileSHA1)
	if bytes.Compare(OriginFileSHA1[:], WriteFileSHA1[:]) != 0 {
		t.Fatal("SHA1 compare Not match")
	}
}

type TestCompose struct {
	name      string
	actoin    Handler
	domains   []string
	expectURL string
}

var testPair = []TestCompose{
	TestCompose{
		"Us2",
		&CHKAVAI,
		[]string{"google.com", "jack.name", "isscape.net", "hello.cc"},
		"https://httpapi.com/api/domains/available.json?auth-userid=000000&api-key=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&domain-name=google&tlds=com&domain-name=jack&tlds=name&domain-name=isscape&tlds=net&domain-name=hello&tlds=cc",
	},
	TestCompose{
		"Us2",
		&CHKAVAI,
		[]string{"google.com"},
		"https://httpapi.com/api/domains/available.json?auth-userid=000000&api-key=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&domain-name=google&tlds=com",
	},
	TestCompose{
		"Us2Demo",
		&CHKAVAI,
		[]string{"google.com", "jack.name", "isscape.net", "hello.cc"},
		"https://test.httpapi.com/api/domains/available.json?auth-userid=000000&api-key=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&domain-name=google&tlds=com&domain-name=jack&tlds=name&domain-name=isscape&tlds=net&domain-name=hello&tlds=cc",
	},
	TestCompose{
		"Us2Demo",
		&CHKAVAI,
		[]string{"google.com"},
		"https://test.httpapi.com/api/domains/available.json?auth-userid=000000&api-key=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&domain-name=google&tlds=com",
	},
}

func TestSplicingURL(t *testing.T) {
	for _, p := range registrars {
		for _, v := range testPair {
			if strings.Compare(v.name, p.Name) == 0 {
				retURL := v.actoin.SplicingURL(p, v.domains)
				if strings.Compare(retURL, v.expectURL) != 0 {
					t.Fatalf("(%s) Return: %v, Want %v\n", v.name, retURL, v.expectURL)
				}
			}
		}
	}
}

func TestHTTPGet(t *testing.T) {
	for _, v := range testPair {
		for _, p := range registrars {
			if strings.Compare(p.Name, v.name) == 0 {
				httpData, err := httpGet(v.actoin.SplicingURL(p, v.domains))
				if err != nil {
					t.Fatalf("httpGet() return: %v, want nil\n", err)
				}
				t.Logf("[%s][%s]%s\n", p.Name, v.actoin, httpData)
				time.Sleep(time.Second * 1)
			}
		}
	}
}
