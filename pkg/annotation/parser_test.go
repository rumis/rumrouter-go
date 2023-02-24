package annotation

import (
	"testing"
)

func TestParseRouterGroup(t *testing.T) {
	g := "// @RouterGroup(Middleware=\"m2,m1\",prefix=\"/user\",namespace=\"app\")"
	ga, err := ParseRouterGroup(g)
	if err != nil {
		t.Error(err)
	}
	if ga.Middleware != "m2,m1" {
		t.Error("routergroup middleware parse failed")
	}
	if ga.Prefix != "/user" {
		t.Error("routergroup prefix parse failed")
	}
	if ga.Namespace != "app" {
		t.Error("routergroup namespace parse failed")
	}
}

func TestParseRouter(t *testing.T) {
	r1 := "// @Router(method=\"GET,POST\",path=\"/getconfig\",namespace=\"app,api\",context=\"x=3,y=4\")"
	r1a, err := ParseRouter(r1)
	if err != nil {
		t.Error(err)
	}
	if r1a.Method != "GET,POST" {
		t.Error("router method parse failed")
	}
	if r1a.Path != "/getconfig" {
		t.Error("router path parse failed")
	}
	if r1a.Namespace != "app,api" {
		t.Error("router namespace parse failed")
	}

	if r1a.Context != "x=3,y=4" {
		t.Error("router context parse failed")
	}

	r2 := "// @Router(method=\"options\",path=\"/getconfig\",middleware=\"auth\")"
	r2a, err := ParseRouter(r2)
	if err != nil {
		t.Error(err)
	}
	if r2a.Method != "options" {
		t.Error("router method parse failed")
	}
	if r2a.Path != "/getconfig" {
		t.Error("router path parse failed")
	}
	if r2a.Middleware != "auth" {
		t.Error("router middleware parse failed")
	}
}

func TestParseMiddleware(t *testing.T) {

	m := "@Middleware(name=\"m1\")"
	ma, err := ParseMiddleware(m)
	if err != nil {
		t.Error(err)
	}
	if ma.Name != "m1" {
		t.Error("middleware name parse failed")
	}
}
