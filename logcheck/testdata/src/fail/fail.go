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
	"context"

	"github.com/scionproto/scion/go/lib/log"
)

const (
	untyped     = "untyped_key"
	typed   key = "typed_key"
)

var value = 1

func validParity() {
	log.Trace("message")
	log.Debug("message")
	log.Info("message")
	log.Warn("message")
	log.Error("message")
	log.Crit("message")

	log.Trace("message")
	log.Debug("message")
	log.Info("message")
	log.Warn("message")
	log.Error("message")
	log.Crit("message")

	log.Trace("message", "key", value)
	log.Debug("message", "key", value)
	log.Info("message", "key", value)
	log.Warn("message", "key", value)
	log.Error("message", "key", value)
	log.Crit("message", "key", value)

	log.Trace("message", "key", value, "key", value)
	log.Debug("message", "key", value, "key", value)
	log.Info("message", "key", value, "key", value)
	log.Warn("message", "key", value, "key", value)
	log.Error("message", "key", value, "key", value)
	log.Crit("message", "key", value, "key", value)
}

func validTypes() {
	log.Debug("message", "key", value)
	log.Debug("message", untyped, value)
	log.Debug("message", typed, value)
}

func invalidParity() {
	log.Trace("message", "key") // want `context should be even: len=1 ctx=\["key"\]`
	log.Debug("message", "key") // want `context should be even: len=1 ctx=\["key"\]`
	log.Info("message", "key")  // want `context should be even: len=1 ctx=\["key"\]`
	log.Warn("message", "key")  // want `context should be even: len=1 ctx=\["key"\]`
	log.Error("message", "key") // want `context should be even: len=1 ctx=\["key"\]`
	log.Crit("message", "key")  // want `context should be even: len=1 ctx=\["key"\]`

	log.Trace("message", "key", value, "key") // want `context should be even: len=3 ctx=\["key",value,"key"\]`
	log.Debug("message", "key", value, "key") // want `context should be even: len=3 ctx=\["key",value,"key"\]`
	log.Info("message", "key", value, "key")  // want `context should be even: len=3 ctx=\["key",value,"key"\]`
	log.Warn("message", "key", value, "key")  // want `context should be even: len=3 ctx=\["key",value,"key"\]`
	log.Error("message", "key", value, "key") // want `context should be even: len=3 ctx=\["key",value,"key"\]`
	log.Crit("message", "key", value, "key")  // want `context should be even: len=3 ctx=\["key",value,"key"\]`
}

func invalidType() {
	log.Info("message", value, value) // want `key should be string: type="int" name="value"`
}

func logger() {
	logger := log.FromCtx(context.Background())
	loggerN := log.New()
	loggerR := log.Root()
	logger.Info("message", "key")  // want `context should be even: len=1 ctx=\["key"\]`
	loggerN.Info("message", "key") // want `context should be even: len=1 ctx=\["key"\]`
	loggerR.Info("message", "key") // want `context should be even: len=1 ctx=\["key"\]`

	log.FromCtx(context.Background()).Info("message", "key") // want `context should be even: len=1 ctx=\["key"\]`
	log.New().Info("message", "key")                         // want `context should be even: len=1 ctx=\["key"\]`
	log.Root().Info("message", "key")                        // want `context should be even: len=1 ctx=\["key"\]`
}

type key string
