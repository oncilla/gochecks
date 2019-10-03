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

package fail

import (
	"github.com/scionproto/scion/go/lib/serrors"
)

const (
	untyped     = "untyped_key"
	typed   key = "typed_key"
)

var (
	errWrap = serrors.New("wrap")
	errBase = serrors.New("base")
	value   = 1
)

func validParity() {
	serrors.New("some error")
	serrors.Wrap(errWrap, errBase)
	serrors.WrapStr("wrap", errBase)

	serrors.New("some error", "key", value)
	serrors.WithCtx(errBase, "key", value)
	serrors.Wrap(errWrap, errBase, "key", value)
	serrors.WrapStr("wrap", errBase, "key", value)

	serrors.New("some error", "key", value, "key", value)
	serrors.WithCtx(errBase, "key", value, "key", value)
	serrors.Wrap(errWrap, errBase, "key", value, "key", value)
	serrors.WrapStr("wrap", errBase, "key", value, "key", value)
}

func validTypes() {
	serrors.New("some error", "key", value)
	serrors.New("some error", untyped, value)
	serrors.New("some error", typed, value)
}

func invalidParity() {
	serrors.New("some error", "key")        // want `context should be even: len=1 ctx=\["key"\]`
	serrors.WithCtx(errBase, "key")         // want `context should be even: len=1 ctx=\["key"\]`
	serrors.Wrap(errWrap, errBase, "key")   // want `context should be even: len=1 ctx=\["key"\]`
	serrors.WrapStr("wrap", errBase, "key") // want `context should be even: len=1 ctx=\["key"\]`

	serrors.New("some error", "key", value, "key")        // want `context should be even: len=3 ctx=\["key",value,"key"\]`
	serrors.WithCtx(errBase, "key", value, "key")         // want `context should be even: len=3 ctx=\["key",value,"key"\]`
	serrors.Wrap(errWrap, errBase, "key", value, "key")   // want `context should be even: len=3 ctx=\["key",value,"key"\]`
	serrors.WrapStr("wrap", errBase, "key", value, "key") // want `context should be even: len=3 ctx=\["key",value,"key"\]`
}

func invalidType() {
	serrors.New("some error", value, value)        // want `key should be string: type="int" name="value"`
	serrors.WithCtx(errBase, value, value)         // want `key should be string: type="int" name="value"`
	serrors.Wrap(errWrap, errBase, value, value)   // want `key should be string: type="int" name="value"`
	serrors.WrapStr("wrap", errBase, value, value) // want `key should be string: type="int" name="value"`
}

func noCtx() {
	serrors.WithCtx(errBase) // want `should have context:.*`
}

type key string

type multiKey key
