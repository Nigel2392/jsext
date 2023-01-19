# Errors
Sometimes you may encounter errors, which are not very descriptive.
I will go into detail of some common errors, and how to fix them.

## Error: panic: ValueOf: invalid value
This error is caused by the fact that you are trying to pass a value to a function, which is not supported by the `syscall/js` package.
When using passing any of the jsext wrappers to a function, you will likely need to use the `Value()` function to get the underlying `syscall/js.Value` object.

## I have an error, which is not listed here!
Are any errors occurring to you with the JSExt package, and are they not listed here?
Please consider opening an issue on the [github repository](https://github.com/Nigel2392/jsext/issues/new/choose), and we will try to fix it as soon as possible.
