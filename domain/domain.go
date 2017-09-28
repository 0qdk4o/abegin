// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// JSONFile specify registrars config file
var (
	JSONFile = "registrars.json"
	CHKAVAI  = Checkavailable("CheckAvailable")

	registrars = make([]*Registrar, 0, 8)
)

// Cmd is the request type
type Cmd string

// A Handler is the interface
type Handler interface {
	SplicingURL(p *Registrar, domains []string) string
	ProcessHTTPData(b []byte) (*Result, error)
}

// SplicingURL splice pieces of variable to long url
func (a *Cmd) SplicingURL(p *Registrar, domains []string) string {
	if len(domains) < 1 || p == nil || p.Host == "" ||
		p.Paths[(*a)] == "" || p.Authid == "" || p.APIkey == "" {
		return ""
	}
	return fmt.Sprintf("%s%s?auth-userid=%s&api-key=%s",
		p.Host, p.Paths[(*a)], p.Authid, p.APIkey)
}

// ProcessHTTPData should implemented by derived struct
func (a *Cmd) ProcessHTTPData(b []byte) (*Result, error) {
	panic("should not call here")
}

// Registrar represent registrar data
type Registrar struct {
	// name is Registrar name such as namechecapCom, nameCom, dynadotCom
	Name string `json:"name"`

	// Host specify registrar host to where request send to.
	// example: http://test.hostapi.com:666 or https://test.hostapi.com
	Host string `json:"host"`

	// APIkey is authenticated by registrar server host when querying
	APIkey string `json:"apikey"`

	// Authid is the userID to send the request to registrar server
	Authid string `json:"authid"`

	// Paths map operation name to script path at the registrar host
	// example checkavailable => /api/domains/available.json
	Paths map[Cmd]string `json:"paths"`
}

// Result represents the error result
type Result struct {
	Status  *string `json:"status"`
	Message *string `json:"message"`
	NameRes
}

// NameRes represents the success query returned result
type NameRes map[string]struct {
	Classkey string `json:"classkey"`
	Status   string `json:"status"`
}

// NewRegistrar returns default registrar pointer to which provide service
func NewRegistrar(name string) *Registrar {
	if len(registrars) < 1 {
		return nil
	}

	for _, p := range registrars {
		if strings.Compare(strings.ToUpper(p.Name), strings.ToUpper(name)) == 0 {
			return p
		}
	}
	return nil
}

// DoRequest query registrar host
func (p *Registrar) DoRequest(h Handler, domains []string) (*Result, error) {
	url := h.SplicingURL(p, domains)
	if url == "" {
		return nil, errors.New("splice url error")
	}

	httpData, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	return h.ProcessHTTPData(httpData)
}

func init() {
	err := loadRegistrarInfo(JSONFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}
}

func loadRegistrarInfo(file string) error {
	buff, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewBuffer(buff))
	_, err = decoder.Token()
	if err != nil {
		return err
	}
	for decoder.More() {
		var reg Registrar
		err := decoder.Decode(&reg)
		if err != nil {
			registrars = registrars[0:0]
			return err
		}
		registrars = append(registrars, &reg)
	}

	_, err = decoder.Token()
	if err != nil {
		registrars = registrars[0:0]
		panic(err)
	}
	return nil
}

// PrintRegs is used to debug
func PrintRegs() {
	if len(registrars) < 1 {
		fmt.Printf("no items in Registrars.")
	}
	for _, v := range registrars {
		fmt.Printf("%s => %v\n", CHKAVAI, (&CHKAVAI).SplicingURL(v, []string{"gupu.com"}))
		fmt.Printf("name: %s\nHost: %s\nAuthid: %s\n", v.Name, v.Host, v.Authid)
		fmt.Printf("APIkey: %s {\n", v.APIkey)
		for k := range v.Paths {
			fmt.Printf("\t%s(%T): %s\n", k, k, v.Paths[k])
		}

		fmt.Println("}")
	}
}

// WriteRegistrarsInfo2File encodes the regs and write to file
func WriteRegistrarsInfo2File(file string) error {
	jsonOut, err := json.MarshalIndent(registrars, "", "\t")
	if err != nil {
		return err
	}
	pfile, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer pfile.Close()
	_, err = pfile.Write(jsonOut)
	if err != nil {
		return err
	}
	return nil
}

func httpGet(url string) ([]byte, error) {
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return ret, nil
}
