package try

import (
	"strings"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/errs"
	"github.com/Nigel2392/jsext/v2/jsc"
	"github.com/Nigel2392/jsext/v2/jsrand"
)

// Catch is a function that takes a js.Func and any number of arguments.
//
// If the js.Func throws/raises an error, it will return the error.
//
// It does so by effectively shelling out to the following javascript:
//
//	(function(arg1, arg2, ...){
//	    try{
//	        return func(arg1, arg2, ...)
//	    }catch(e){
//	        return e
//	    }
//	})
//
// The returned value is the value which the js.Func returns.
//
// If an error occurs, the js.Value will be js.Null() and the error will be non-nil.
func Catch(f js.Func, args ...interface{}) (js.Value, error) {
	var argsJS, err = jsc.ValuesOf(args...)
	if err != nil {
		return js.Null(), err
	}
	// generate a random name for each argument
	var argnames = make([]string, len(args))
	args = make([]interface{}, len(args))
	for i, arg := range argsJS {
		argnames[i] = jsrand.String(16)
		args[i] = arg
	}
	// set the function to the global scope
	var funcName = jsrand.String(32)
	js.Global().Set(funcName, f)

	// delete the function after we're done
	defer js.Global().Delete(funcName)

	// create the function
	var b strings.Builder
	b.WriteString("(function(")
	for i, arg := range argnames {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(arg)
	}
	b.WriteString("){")
	b.WriteString("try{")
	b.WriteString("return ")
	b.WriteString(funcName)
	b.WriteString("(")
	for i, arg := range argnames {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(arg)
	}
	b.WriteString(")")
	b.WriteString("}catch(e){return e}})")

	// call the function
	var tryFunc = js.Global().Call("eval", b.String())
	var ret = tryFunc.Invoke(args...)

	// check if the return value is an error
	if ret.InstanceOf(js.Global().Get("Error")) {
		return js.Null(), errs.Error(ret.Get("message").String())
	}
	return ret, nil
}
