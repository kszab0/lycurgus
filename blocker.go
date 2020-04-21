package main

import (
	"log"
	"net/http"

	"gopkg.in/elazarl/goproxy.v1"
)

// Blocker blocks HTTP requests based on different rules
type Blocker struct {
	enabled bool

	proxy     *goproxy.ProxyHttpServer
	blocklist Matcher
	blacklist Matcher
	whitelist Matcher
}

// NewBlocker creates and initializes a Blocker
func NewBlocker(enabled bool) *Blocker {
	b := &Blocker{
		enabled: enabled,
	}
	b.initProxy()
	return b
}

// Toggle toggles the enabled state
func (b *Blocker) Toggle() {
	b.enabled = !b.enabled
}

// ServeHTTP implements the http.Handler interface
func (b *Blocker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.proxy.ServeHTTP(w, r)
}

func (b *Blocker) initProxy() {
	b.proxy = goproxy.NewProxyHttpServer()
	b.proxy.OnRequest().HandleConnectFunc(b.handleConnect)
}

func (b *Blocker) handleConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	if !b.enabled {
		log.Printf("Host accepted (proxy disabled): %s\n", host)
		return goproxy.OkConnect, host
	}
	if b.whitelist != nil && b.whitelist.Match(host) {
		log.Printf("Host accepted (whitelist): %s\n", host)
		return goproxy.OkConnect, host
	}
	if b.blocklist != nil && b.blocklist.Match(host) {
		log.Printf("Host rejected (blocklist): %s\n", host)
		return goproxy.RejectConnect, host
	}
	if b.blacklist != nil && b.blacklist.Match(host) {
		log.Printf("Host rejected (blacklist): %s\n", host)
		return goproxy.RejectConnect, host
	}
	log.Printf("Host accepted: %s\n", host)
	return goproxy.OkConnect, host
}
