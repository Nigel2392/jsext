//go:build skipimports
// +build skipimports

package try

import (
	"syscall/js"
)

func getRandString(n int) string {
	var js = `(function(len){
		function dec2hex (dec) {
			return dec.toString(16).padStart(2, "0")
		}

		var arr = new Uint8Array((len || 40) / 2)
		window.crypto.getRandomValues(arr)
		return Array.from(arr, dec2hex).join('')
	})`
	var funcRandString = js.Global().Call("eval", js)
	return ("f_" + funcRandString.Invoke(n).String())
}
