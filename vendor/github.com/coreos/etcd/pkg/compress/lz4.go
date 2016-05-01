// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compress

import (
	"io"
	"io/ioutil"

	"github.com/bkaradzic/go-lz4"
)

func NewLz4Compressor() Compressor {
	return lz4Compressor{}
}

type lz4Compressor struct{}

func (_ lz4Compressor) Do(w io.Writer, p []byte) error {
	bts, err := lz4.Encode(nil, p)
	if err != nil {
		return err
	}
	_, err = w.Write(bts)
	return err
}

func (_ lz4Compressor) Type() string {
	return headers[Lz4]
}

func NewLz4Decompressor() Decompressor {
	return lz4Decompressor{}
}

type lz4Decompressor struct{}

func (_ lz4Decompressor) Do(r io.Reader) ([]byte, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return lz4.Decode(nil, src)
}

func (_ lz4Decompressor) Type() string {
	return headers[Lz4]
}
