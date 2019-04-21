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

	jsoniter "github.com/json-iterator/go"

	"github.com/vicanso/cod"
	"github.com/vicanso/hes"
)

type (
	// Config error handler config
	Config struct {
		Skipper      cod.Skipper
		ResponseType string
	}
)

const (
	errErrorHandlerCategory = "cod-error-handler"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// NewDefault create a default error handler
func NewDefault() cod.Handler {
	return New(Config{
		ResponseType: "json",
	})
}

// New create a error handler
func New(config Config) cod.Handler {
	skipper := config.Skipper
	if skipper == nil {
		skipper = cod.DefaultSkipper
	}
	return func(c *cod.Context) error {
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
			he.Category = errErrorHandlerCategory
		}
		c.StatusCode = he.StatusCode
		if config.ResponseType == "json" {
			buf, e := json.Marshal(he)
			if e != nil {
				return e
			}
			c.BodyBuffer = bytes.NewBuffer(buf)
			c.SetHeader(cod.HeaderContentType, cod.MIMEApplicationJSON)
		} else {
			c.BodyBuffer = bytes.NewBufferString(he.Error())
			c.SetHeader(cod.HeaderContentType, cod.MIMETextPlain)
		}

		return nil
	}
}
