# JSExt!

Golang Wasm extension for easy JS interop.

## Installation
Easily install  with the following command:
```
go get github.com/Nigel2392/jsext
```

## Examples
Examples can be found in the examples folder.
There are only a few examples for now, but we will add more in the future.

## TinyGO support
JSExt is fully supported by TinyGO, and can be used with the TinyGO compiler.
Some features may however work differently with the TinyGO compiler, and the regular go compiler.
This is due to the fact that TinyGO does not support the reflect package, HTTP, and some other packages.

To compile the package with TinyGO support, you need to add the following build tags to your project:
```
tinygo build -tags=tinygo -o <output file> <input file>
```

### Making requests using TinyGo
When making requests when the project is compiled with TinyGO, (when compiling with the `tinygo` tag), 
we will automatically resort to using the Javascript `fetch` api.
When the request is made, we wait for the promise in a goroutine, and return the response when it is done using a channel.

### Limitations of making requests using TinyGo
Some features are not implemented yet when making requests using TinyGO.
These functions mostly have to do with encoding.
We can encode certain items, such as slices and maps. 
This is done using our custom encoder located at `github.com/Nigel2392/jsext/requester/fetch/dirtyjson.go`.
*This encoder however is not fully tested, and may very well break on certain inputs.*

## Creating a project
To easily create projects, it is best to install the jsext cli tool:
```
$ go install github.com/Nigel2392/jsexttool
```
Following that, you can easily create an example project with the following command, which will contain a basic example of some jsext functionality:
```
$ jsexttool -init -n <project name> (optional: -vscode for auto creation of vscode config.)
```
There is also an option to create a plain project, this will create a project with no example code, but still some basic setup:
```
$ jsexttool -plain -n <project name> (optional: -vscode for auto creation of vscode config.)
```
### Note:
Both of the jsexttool options provide a some powershell build scripts, which can be used to quickly build the project with tinygo, or the regular go compiler.
If the chosen compiler is the normal go compiler, you will need to edit the index.html file to import go's default wasm_exec.js file, instead of the tinygo one.
The jsexttool automatically uses the tinygo compiler on setup.

The before mentioned build flags are already set in the build scripts.

There is also a server.go file when initializing the project with the `-init` flag. This file can be used to quickly serve the project for development.

## Binary sizes
When compiling with the regular go compiler, the binary sizes can get pretty big. We do implement a page loader, but this may not be your way to go.
This is why we recomment starting out with TinyGO, as the binary sizes are much smaller, but even that could get up to 2MB+ without any optimizations.
If you do suffer binary size issues, you can try to use the `-ldflags="-s -w"` flag when compiling with the regular go compiler.
When compiling with TinyGO, you can try to use the `-no-debug` flag, which will remove debug information from the binary.
There are also some great optimization tools for WebAssembly, such as WasmOpt, which can be used to optimize the binary size, or speed.
You can find more information about ***wasm-opt** [here](https://github.com/WebAssembly/binaryen)*.
wasm-opt .\static\main.wasm -o=".\static\main.wasm" -Oz --shrink-level=3 --optimize-level=3
