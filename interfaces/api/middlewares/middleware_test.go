package middlewares_test

import (
	"fmt"
	"net/http"
	"testing"
)

type Middleware http.HandlerFunc

func (m Middleware) Add(f http.HandlerFunc) Middleware {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
		m(w, r)
	}
}

func handlerA(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A")
}

func handlerB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("B")
}

func handlerC(w http.ResponseWriter, r *http.Request) {
	fmt.Println("C")
}

func TestMiddleware(t *testing.T) {
	m := Middleware(handlerA).Add(handlerB).Add(handlerC)
	var w http.ResponseWriter
	var r *http.Request
	m(w, r)
}

type MiddlewareGroup struct {
	group []Middleware
}

func NewMiddlewareGroup() *MiddlewareGroup {
	return &MiddlewareGroup{
		group: make([]Middleware, 0),
	}
}

func (mg *MiddlewareGroup) Add(m ...Middleware) *MiddlewareGroup {
	mg.group = append(mg.group, m...)
	return mg
}

func (mg *MiddlewareGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, m := range mg.group {
		m(w, r)
	}
}

func TestMiddlewareGroup(t *testing.T) {
	mg := NewMiddlewareGroup()
	mg.Add(handlerA, handlerB, handlerC)
	var w http.ResponseWriter
	var r *http.Request
	mg.ServeHTTP(w, r)
}
