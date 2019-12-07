package main

import (
	"errors"

	"github.com/vicanso/elton"
	errorhandler "github.com/vicanso/elton-error-handler"
)

func main() {
	e := elton.New()

	// 指定出错以json的形式返回
	e.Use(errorhandler.New(errorhandler.Config{
		ResponseType: "json",
	}))

	e.GET("/", func(c *elton.Context) (err error) {
		return errors.New("abcd")
	})

	err := e.ListenAndServe(":3000")
	if err != nil {
		panic(err)
	}
}
