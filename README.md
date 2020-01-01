# elton-error-handler

[![Build Status](https://img.shields.io/travis/vicanso/elton-error-handler.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-error-handler)

Error handler for elton, it convert error to json/text response(NewDefault). Suggest to use `hes.Error` for custom error.

```go
package main

import (
	"errors"

	"github.com/vicanso/elton"

	errorhandler "github.com/vicanso/elton-error-handler"
)

func main() {

	e := elton.New()
	e.Use(errorhandler.NewDefault())

	e.GET("/", func(c *elton.Context) (err error) {
		err = errors.New("abcd")
		return
	})

	err := e.ListenAndServe(":3000")
	if err != nil {
		panic(err)
	})
}
```