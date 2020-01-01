// Copyright 2018 tree xie
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

package errorhandler

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
)

type (
	// Config error handler config
	Config struct {
		Skipper      elton.Skipper
		ResponseType string
	}
)

const (
	// ErrCategory error cateogry of error handler
	ErrCategory = "elton-error-handler"
)

// NewDefault create a default error handler
func NewDefault() elton.Handler {
	return New(Config{})
}

// New create a error handler
func New(config Config) elton.Handler {
	skipper := config.Skipper
	if skipper == nil {
		skipper = elton.DefaultSkipper
	}
	return func(c *elton.Context) error {
		if skipper(c) {
			return c.Next()
		}
		err := c.Next()
		// 如果没有出错，直接返回
		if err == nil {
			return nil
		}
		he, ok := err.(*hes.Error)
		if !ok {
			he = hes.Wrap(err)
			he.StatusCode = http.StatusInternalServerError
			he.Exception = true
			he.Category = ErrCategory
		}
		c.StatusCode = he.StatusCode
		if config.ResponseType == "json" ||
			strings.Contains(c.GetRequestHeader("Accept"), "application/json") {
			buf := he.ToJSON()
			c.BodyBuffer = bytes.NewBuffer(buf)
			c.SetHeader(elton.HeaderContentType, elton.MIMEApplicationJSON)
		} else {
			c.BodyBuffer = bytes.NewBufferString(he.Error())
			c.SetHeader(elton.HeaderContentType, elton.MIMETextPlain)
		}

		return nil
	}
}
