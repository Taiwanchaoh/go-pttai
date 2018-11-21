// Copyright 2018 The go-pttai Authors
// This file is part of the go-pttai library.
//
// The go-pttai library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-pttai library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-pttai library. If not, see <http://www.gnu.org/licenses/>.

package ptthttp

import "regexp"

const (
	MaxUploadSize = 10000000 // 10MB
)

// re

var (
	reRPCPort   *regexp.Regexp
	reHTTPAddr  *regexp.Regexp
	extRPCPort  []byte
	extHTTPAddr []byte
)

func init() {
	reRPCPort, _ = regexp.Compile("localhost:14779")
	reHTTPAddr, _ = regexp.Compile("localhost:9774")
}