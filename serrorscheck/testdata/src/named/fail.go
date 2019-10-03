// MIT License
//
// Copyright (c) 2019 Oncilla
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package named

import (
	serrors "fake/serrors"

	errors "github.com/scionproto/scion/go/lib/serrors"
)

var (
	errWrap = errors.New("wrap")
	errBase = errors.New("base")
	value   = 1
)

func fakeImportIgnored() {
	serrors.New("some error", "key")
	serrors.Wrap(errWrap, errBase, "key")
}

func valid() {
	errors.New("some error", "key", value, "key", value)
	errors.WithCtx(errBase, "key", value, "key", value)
	errors.Wrap(errWrap, errBase, "key", value, "key", value)
	errors.WrapStr("wrap", errBase, "key", value, "key", value)
}
