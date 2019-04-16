# cod-error-handler

[![Build Status](https://img.shields.io/travis/vicanso/cod-error-handler.svg?label=linux+build)](https://travis-ci.org/vicanso/cod-error-handler)

Error handler for cod, it convert error to json response(NewDefault). Suggest to use `hes.Error` for custom error.

```go
package main

import (
	"errors"

	"github.com/vicanso/cod"

	errorhandler "github.com/vicanso/cod-error-handler"
)

func main() {

	d := cod.New()
	d.Use(errorhandler.NewDefault())

	d.GET("/", func(c *cod.Context) (err error) {
		err = errors.New("abcd")
		return
	})

	d.ListenAndServe(":7001")
}
```