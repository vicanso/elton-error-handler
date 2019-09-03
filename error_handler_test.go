package errorhandler

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func TestSkipAndNoError(t *testing.T) {
	fn := NewDefault()
	t.Run("skip", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/users/me", nil)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.Committed = true
		c.Next = func() error {
			return nil
		}
		err := fn(c)
		assert.Nil(err)
		assert.Nil(c.BodyBuffer)
	})

	t.Run("no error", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/users/me", nil)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.Next = func() error {
			return nil
		}
		err := fn(c)
		assert.Nil(err)
		assert.Nil(c.BodyBuffer)
	})
}

func TestErrorHandler(t *testing.T) {
	t.Run("json type", func(t *testing.T) {
		assert := assert.New(t)
		fn := NewDefault()
		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set("Accept", "application/json, text/plain, */*")
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.Next = func() error {
			return errors.New("abcd")
		}
		c.CacheMaxAge("5m")
		err := fn(c)
		assert.Nil(err)
		assert.Equal("public, max-age=300", c.GetHeader(elton.HeaderCacheControl))
		assert.True(strings.HasSuffix(c.BodyBuffer.String(), `"statusCode":500,"category":"elton-error-handler","message":"abcd","exception":true}`))
		assert.Equal("application/json; charset=UTF-8", c.GetHeader(elton.HeaderContentType))
	})

	t.Run("text type", func(t *testing.T) {
		assert := assert.New(t)
		fn := New(Config{
			ResponseType: "text",
		})
		req := httptest.NewRequest("GET", "/users/me", nil)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.Next = func() error {
			return errors.New("abcd")
		}
		c.CacheMaxAge("5m")
		err := fn(c)
		assert.Nil(err)
		assert.Equal("public, max-age=300", c.GetHeader(elton.HeaderCacheControl))
		ct := c.GetHeader(elton.HeaderContentType)
		assert.Equal("category=elton-error-handler, message=abcd", c.BodyBuffer.String())
		assert.Equal("text/plain; charset=UTF-8", ct)
	})
}

// https://stackoverflow.com/questions/50120427/fail-unit-tests-if-coverage-is-below-certain-percentage
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.9 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
