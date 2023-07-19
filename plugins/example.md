# Add and import a plugin in a javascript runtime.

To get started creating a plugin, lets first create a main.go file which can be compiled to webassembly.

In the following example, we will set our exports to the global variable `my.exports`.

```go
package main

import "syscall/js"

var myObj = js.Global().Get("Object").New()
var exportObj = js.Global().Get("Object").New()

func AddStrings(this js.Value, args []js.Value) interface{} {
	var sum = ""
	for _, v := range args {
		sum += v.String()
	}
	return sum
}

var waiter = make(chan struct{})

func main() {
	js.Global().Set("my", myObj)
	myObj.Set("exports", exportObj)
	exportObj.Set("addStrings", js.FuncOf(AddStrings))
	<-waiter
}

```

After compiling the above go file, we can get started importing the the plugin.

In another webassembly environment, you could call the following functions to retrieve and call the "addStrings" function./*

```go
var exp = plugins.New()
var plugin, err = exp.NewPlugin("plugins", "/static/plugins.wasm", "my.exports", plugins.F_IMPORT_FROM_GLOBAL)
if err != nil {
	console.Error(err.Error())
	return
}
for k, v := range plugin.Exports {
	fmt.Printf("%s: %s\n", k, v.String())
}

v, err := plugin.Call("addStrings", "Hello", " World", "!")
if err != nil {
	console.Error(err.Error())
	return
}

fmt.Println("addStrings:", v.String())
```
