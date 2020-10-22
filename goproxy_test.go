package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/net/proxy"
)

var proxyClient *http.Client

func init() {
	go main()
	dialSocksProxy, err := proxy.SOCKS5("tcp", "localhost:1080", nil, proxy.Direct)
	if err != nil {
		log.Fatal("Error connecting to proxy:", err)
	}
	proxyClient = &http.Client{
		Transport: &http.Transport{
			Dial: dialSocksProxy.Dial,
		},
	}
}

func TestInitialization(t *testing.T) {
	select {
	case <-time.After(30 * time.Second):
		t.Fatal("Initialization timed out")
	case <-GetInitChannel():
	}
}

func TestProxyAllow(t *testing.T) {
	testMessage := "Hello, Client"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testMessage)
	}))
	defer ts.Close()

	resp, err := proxyClient.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Unexpected status " + resp.Status)
	}
	greeting, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(greeting) != testMessage {
		t.Error("Unexpected message: " + string(greeting))
	}
}

func TestProxyBlock(t *testing.T) {
	block("127.0.0.1")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Test message")
	}))
	defer ts.Close()

	_, err := proxyClient.Get(ts.URL)
	if err == nil {
		t.Error("Not expected to allow the connection")
	}
}
