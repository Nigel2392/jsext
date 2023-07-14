package try

import (
	"bytes"
	"errors"
	"math/rand"
	"syscall/js"
	"unsafe"

	"github.com/Nigel2392/jsext/v2/jsc"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(0<<63 - 1)

// https://stackoverflow.com/a/31832326/18020941
func randStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

type argTuple struct {
	name string
	val  js.Value
}

// Catch is a function that takes a js.Func and any number of arguments.
//
// It will call the js.Func with the arguments and return the result.
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
func Catch(f js.Func, args ...interface{}) error {
	var argsJS, err = jsc.ValuesOf(args...)
	if err != nil {
		return err
	}
	// generate a random name for each argument
	var argsTuple = make([]argTuple, len(args))
	for i, arg := range argsJS {
		argsTuple[i] = argTuple{
			name: randStringBytesMaskImprSrcUnsafe(16),
			val:  arg,
		}
	}
	// set the function to the global scope
	var funcName = randStringBytesMaskImprSrcUnsafe(32)
	js.Global().Set(funcName, f)

	// delete the function after we're done
	defer js.Global().Delete(funcName)

	// create the function
	var b bytes.Buffer
	b.WriteString("(function(")
	for i, arg := range argsTuple {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(arg.name)
	}
	b.WriteString("){")
	b.WriteString("try{")
	b.WriteString("return ")
	b.WriteString(funcName)
	b.WriteString("(")
	for i, arg := range argsTuple {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(arg.name)
	}
	b.WriteString(")")
	b.WriteString("}catch(e){return e}})")

	// call the function
	var tryFunc = js.Global().Call("eval", b.String())
	var ret = tryFunc.Invoke(args...)

	// check if the return value is an error
	if ret.InstanceOf(js.Global().Get("Error")) {
		return errors.New(ret.Get("message").String())
	}
	return nil
}
