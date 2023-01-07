package middleware

import (
	"net/url"
	"runtime/debug"

	"github.com/Nigel2392/jsext/framework/router/rterr"
	"github.com/Nigel2392/jsext/framework/router/vars"
)

func Recoverer(varMap vars.Vars, path *url.URL, err rterr.ErrorThrower) bool {
	defer func() {
		if r := recover(); r != nil {
			println(string(debug.Stack()))
			switch newR := r.(type) {
			case error:
				err.Error(500, newR.Error())
			case string:
				err.Error(500, newR)
			default:
				err.Error(500, "Unknown error occurred.")
			}
		}
	}()
	return true
}
