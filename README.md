# JSExt!

Golang Wasm extension for easy JS interop.

A simple library for more easily interacting with Javascript like you are used to.

## Binary sizes

When compiling with the regular go compiler, the binary sizes can get pretty big. We do implement a page loader, but this may not be your way to go.
This is why we recomment starting out with TinyGO, as the binary sizes are much smaller, but even that could get pretty large without any optimizations.
If you do suffer binary size issues, you can try to use the `-ldflags="-s -w"` flag when compiling with the regular go compiler.
When compiling with TinyGO, you can try to use the `-no-debug` flag, which will remove debug information from the binary.
There are also some great optimization tools for WebAssembly, such as WasmOpt, which can be used to optimize the binary size, or speed.
You can find more information about ***wasm-opt** [here](https://github.com/WebAssembly/binaryen)*.
