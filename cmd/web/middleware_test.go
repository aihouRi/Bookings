package main

import (
	"fmt"
	"testing"
	"net/http"
)

func TestNoSurf(t *testing.T) {
	var myH myhandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing; test passed
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myhandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing; test passed
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}