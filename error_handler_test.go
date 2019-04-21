package errorhandler

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/cod"
)

func TestSkipAndNoError(t *testing.T) {
	fn := NewDefault()
	t.Run("skip", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/me", nil)
		resp := httptest.NewRecorder()
		c := cod.NewContext(resp, req)
		c.Committed = true
		c.Next = func() error {
			return nil
		}
		err := fn(c)
		if err != nil || c.BodyBuffer != nil {
			t.Fatalf("skip error handler fail, %v", err)
		}
	})

	t.Run("no error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/me", nil)
		resp := httptest.NewRecorder()
		c := cod.NewContext(resp, req)
		c.Next = func() error {
			return nil
		}
		err := fn(c)
		if err != nil || c.BodyBuffer != nil {
			t.Fatalf("no error handler fail, %v", err)
		}
	})
}

func TestErrorHandler(t *testing.T) {
	t.Run("json type", func(t *testing.T) {
		fn := NewDefault()
		req := httptest.NewRequest("GET", "/users/me", nil)
		resp := httptest.NewRecorder()
		c := cod.NewContext(resp, req)
		c.Next = func() error {
			return errors.New("abcd")
		}
		c.CacheMaxAge("5m")
		err := fn(c)
		if err != nil {
			t.Fatalf("error handler fail, %v", err)
		}
		if c.GetHeader(cod.HeaderCacheControl) != "public, max-age=300" {
			t.Fatalf("cache control field is invalid")
		}
		ct := c.GetHeader(cod.HeaderContentType)
		if c.BodyBuffer.String() != `{"statusCode":500,"category":"cod-error-handler","message":"abcd","exception":true}` ||
			ct != "application/json; charset=UTF-8" {
			t.Fatalf("error handler fail")
		}
	})

	t.Run("text type", func(t *testing.T) {
		fn := New(Config{
			ResponseType: "text",
		})
		req := httptest.NewRequest("GET", "/users/me", nil)
		resp := httptest.NewRecorder()
		c := cod.NewContext(resp, req)
		c.Next = func() error {
			return errors.New("abcd")
		}
		c.CacheMaxAge("5m")
		err := fn(c)
		if err != nil {
			t.Fatalf("error handler fail, %v", err)
		}
		if c.GetHeader(cod.HeaderCacheControl) != "public, max-age=300" {
			t.Fatalf("cache control field is invalid")
		}
		ct := c.GetHeader(cod.HeaderContentType)
		if c.BodyBuffer.String() != "category=cod-error-handler, message=abcd" ||
			ct != "text/plain; charset=UTF-8" {
			t.Fatalf("error handler fail")
		}
	})
}
