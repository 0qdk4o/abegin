// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package Base58 implements base58 encode and decode
package Base58

import (
	"errors"
	"math/big"
)

// BASE58CHARS is all alphanumeric characters except for "0", "I", "O", and "l"
const (
	BASE58CHARS = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

var char2int = make(map[byte]int64)

func init() {
	for i, v := range BASE58CHARS {
		char2int[byte(v)] = int64(i)
	}
}

// EncodeBase58 represents my version base58 encode
func EncodeBase58(b []byte) string {
	n := &big.Int{}
	n.SetBytes(b)

	zero := big.NewInt(0)
	div := big.NewInt(58)
	mod := big.NewInt(0)

	var buf = make([]byte, 0, 64)
	for n.Cmp(zero) > 0 {
		n.DivMod(n, div, mod)
		buf = append(buf, BASE58CHARS[mod.Int64()])
	}

	for _, v := range b {
		if v != 0x00 {
			break
		}
		buf = append(buf, BASE58CHARS[0])
	}

	return reverse(buf)
}

func reverse(b []byte) string {
	lenb := len(b)
	for i, j := 0, lenb-1; i < lenb/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

// DecodeBase58 decode string s to bytes
func DecodeBase58(s string) ([]byte, error) {
	lens := len(s)
	n := big.NewInt(0)
	b := big.NewInt(58) // Base
	tmp := big.NewInt(0)

	var res = make([]byte, 0, 64)
	for i := 0; i < lens; i++ {
		if s[i] == BASE58CHARS[0] {
			res = append(res, 0x00)
			continue
		}
		idx, err := char2int[s[i]]
		if !err {
			return nil, errors.New("Invalid decoded string")
		}
		tmp.SetInt64(idx)
		n.Add(n.Mul(n, b), tmp)
	}
	res = append(res, n.Bytes()...)
	return res, nil
}
