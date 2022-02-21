package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var mh myHandler
	h := NoSurf(&mh)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("Expected type http.Handler, instead got %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var mh myHandler
	h := SessionLoad(&mh)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("Expected type http.Handler, instead got %T", v))
	}
}
