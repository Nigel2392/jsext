package plugins

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/errs"
)

type Plugin struct {
	ExportPath   string              `js:"exportPath"`
	Exports      map[string]js.Value `js:"exports"`
	Module       js.Value            `js:"module"`
	Instance     js.Value            `js:"instance"`
	ExportObject js.Value            `js:"exportObject"`
	GoObject     js.Value            `js:"goObject"`
	Flags        uint8
}

func (e *Plugin) Call(varName string, args ...any) (js.Value, error) {
	var v, err = e.Get(varName)
	if err != nil {
		return js.Undefined(), err
	}
	return v.Invoke(args...), nil
}

func (e *Plugin) Get(varName string) (js.Value, error) {
	var v = e.Exports[varName]
	if v.IsUndefined() {
		return js.Undefined(), fmt.Errorf("variable %s not found", varName)
	}
	return v, nil
}

type Plugins map[string]*Plugin

func New() Plugins {
	var m = make(map[string]*Plugin)
	var instantiateStreaming = webAssembly.Get("instantiateStreaming")
	if instantiateStreaming.IsUndefined() {
		// polyfill
		instantiateStreaming = js.Global().Call("eval", `
		WebAssembly.instantiateStreaming = async (resp, importObject) => {
			const source = await (await resp).arrayBuffer();
			return await WebAssembly.instantiate(source, importObject);
		};`)
		webAssembly.Set("instantiateStreaming", instantiateStreaming)
	}
	return m
}

var (
	webAssembly = js.Global().Get("WebAssembly")
	fetch       = js.Global().Get("fetch")
	globalGo    = js.Global().Get("Go")
	obj         = js.Global().Get("Object")
)

const (
	F_IMPORT_FROM_GO_OBJECT uint8 = 1 << iota
	F_IMPORT_FROM_MODULE
	F_IMPORT_FROM_INSTANCE
	F_IMPORT_FROM_GLOBAL
)

func (e *Plugins) NewPlugin(name, url, pathToExports string, flag uint8) (*Plugin, error) {
	if e == nil {
		*e = make(map[string]*Plugin)
	}

	var wait chan any = make(chan any)
	var goObj = globalGo.New()
	var importObject = goObj.Get("importObject")
	var instantiateStreaming = webAssembly.Get("instantiateStreaming")
	var module, instance js.Value
	var err error
	instantiateStreaming.Invoke(fetch.Invoke(url), importObject).Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		module = args[0].Get("module")
		instance = args[0].Get("instance")
		goObj.Call("run", instance)
		wait <- true
		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) any {
		var jsErr = args[0]
		err = errs.Error(jsErr.Get("message").String())
		wait <- true
		return nil
	}))
	<-wait
	if err != nil {
		return nil, err
	}
	var exportObject js.Value
	var split = strings.Split(pathToExports, ".")
	switch {
	case flag&F_IMPORT_FROM_GO_OBJECT != 0:
		exportObject = goObj
	case flag&F_IMPORT_FROM_MODULE != 0:
		exportObject = module
	case flag&F_IMPORT_FROM_INSTANCE != 0:
		exportObject = instance
	case flag&F_IMPORT_FROM_GLOBAL != 0:
		exportObject = js.Global()
	default:
		return nil, fmt.Errorf("invalid flag %d", flag)
	}
	for i := 0; i < len(split); i++ {
		if strings.HasPrefix(split[i], "[") {
			var index = strings.TrimPrefix(split[i], "[")
			index = strings.TrimSuffix(index, "]")
			split[i] = index
		}
		exportObject = exportObject.Get(split[i])
		if exportObject.IsUndefined() {
			return nil, fmt.Errorf("export %s not found", pathToExports)
		}
	}
	var exports = make(map[string]js.Value)
	var keys = obj.Call("keys", exportObject)
	for i := 0; i < keys.Length(); i++ {
		var key = keys.Index(i).String()
		exports[key] = exportObject.Get(key)
	}

	var exp = &Plugin{
		ExportPath:   pathToExports,
		Exports:      exports,
		Module:       module,
		Instance:     instance,
		ExportObject: exportObject,
		GoObject:     goObj,
		Flags:        flag,
	}

	(*e)[name] = exp
	return exp, nil
}

func (e *Plugins) Call(pluginName, varName string, args ...any) (js.Value, error) {
	var v, err = e.Get(pluginName, varName)
	if err != nil {
		return js.Undefined(), err
	}
	return v.Invoke(args...), nil
}

func (e *Plugins) Get(pluginName, varName string) (js.Value, error) {
	var export = (*e)[pluginName]
	if export == nil {
		return js.Undefined(), fmt.Errorf("export %s not found", pluginName)
	}
	var v = export.Exports[varName]
	if v.IsUndefined() {
		return js.Undefined(), fmt.Errorf("variable %s not found", varName)
	}
	return v, nil
}
