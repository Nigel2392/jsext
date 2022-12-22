package router_test

import (
	"create_elems/jsext/router"
	"net/url"
	"testing"
)

func TestRouter(t *testing.T) {
	var rt = router.NewRouter().SkipTrailingSlash()

	var homeRoute = rt.Register("home", "/", func(vars router.Vars, u *url.URL) {
		t.Log("Index")
	})
	var innerHome = homeRoute.Register("innerhome", "home/", func(vars router.Vars, u *url.URL) {
		t.Log("Home")
	})
	innerHome.Register("innerinnerhome", "home/", func(vars router.Vars, u *url.URL) {
		t.Log("Home")
	})
	rt.Register("about", "/about/", func(vars router.Vars, u *url.URL) {
		t.Log("About")
	})
	rt.Register("post", "/post/<<post:int>>/<<name:string>>/", func(vars router.Vars, u *url.URL) {
		t.Log(vars.Get("post"))
	})
	rt.Register("postraw", "/post/<<post:raw([0-9]+)>>/", func(vars router.Vars, u *url.URL) {
		t.Log(vars.Get("post"))
	})

	var currentRoute = rt.GetRoute("home:innerhome:innerinnerhome")
	if currentRoute == nil {
		t.Error("Route not found!")
		return
	}
	t.Log(currentRoute)

	currentRoute = rt.GetRoute("post")
	if currentRoute == nil {
		t.Error("Route not found!")
		return
	}
	var path = currentRoute.URL("123", "test")
	if path != "/post/123/test" { // SkipTrailingSlash is true
		t.Error("Wrong path! " + path)
		return
	}
	t.Log("Path formatted: " + path)

	currentRoute, _, found := rt.Match("/")
	if !found {
		t.Error("Route / not found!", rt)
		return
	}
	t.Log(currentRoute)
	currentRoute, _, found = rt.Match("/home/")
	if !found {
		t.Error("Route /home not found!")
		return
	}
	t.Log(currentRoute)
	currentRoute, _, found = rt.Match("/home/home/")
	if !found {
		t.Error("Route /home/home not found!")
		return
	}
	t.Log(currentRoute)

	currentRoute, _, found = rt.Match("/about/")
	if !found {
		t.Error("Route /about not found!")
		return
	}
	t.Log(currentRoute)
	currentRoute, vars, found := rt.Match("/post/123/")
	if !found {
		t.Error("Route /post/123 not found!")
	}
	t.Log("Vars: ", vars)
	post := vars.Get("post")
	if post != "123" {
		t.Error("Post variable not found!")
	}
	t.Log(currentRoute)
	currentRoute, vars, found = rt.Match("/post/123/abc/")
	if !found {
		t.Error("Route /post/123/abc not found!")
	}
	t.Log("Vars: ", vars)
	post = vars.Get("post")
	if post != "123" {
		t.Error("Post variable not found!")
		return
	}
	name := vars.Get("name")
	if name != "abc" {
		t.Error("Name not found!")
		return
	}
	t.Log(currentRoute)
	_, _, found = rt.Match("/post/123a")
	if found {
		t.Error("Route found!")
		return
	}
	_, _, found = rt.Match("/post/123a/abc")
	if found {
		t.Error("Route found!")
		return
	}
	_, _, found = rt.Match("/post")
	if found {
		t.Error("Route found!")
		return
	}

}
