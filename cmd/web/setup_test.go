package main

import (
	"testing"
	"os"
	"net/http"
)

func TestMain(m *testing.M) {



	os.Exit(m.Run())
}

type myhandler struct{}


func (mh *myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// do something
}