# elton-error-handler

[![Build Status](https://img.shields.io/travis/vicanso/elton-error-handler.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-error-handler)

Error handler for elton, it convert error to json response(NewDefault). Suggest to use `hes.Error` for custom error.

```go
package main

import (
	"errors"

	"github.com/vicanso/elton"

	errorhandler "github.com/vicanso/elton-error-handler"
)

func main() {

	d := elton.New()
	d.Use(errorhandler.NewDefault())

	d.GET("/", func(c *elton.Context) (err error) {
		err = errors.New("abcd")
		return
	})

	d.ListenAndServe(":7001")
}
```