package files

import (
	"syscall/js"
	"time"
)

var (
	blobConstructor       = js.Global().Get("Blob")
	uint8ArrayConstructor = js.Global().Get("Uint8Array")
	fileConstructor       = js.Global().Get("File")
	dateConstructor       = js.Global().Get("Date")
	fileReader            = js.Global().Get("FileReader")
)

/*
	Example:

		var f = files.New("hello.txt", []byte(index), "", time.Now())
		fmt.Println(f.Name())
		fmt.Println(f.Size())
		fmt.Println(f.MimeType())

		rd, err := files.Read(f)
		if err != nil {
			console.Error(err)
			return
		}
		var b = new(bytes.Buffer)
		b.ReadFrom(rd)
		console.Log(b.String())

*/

type File js.Value

func New(name string, data []byte, mimetype string, lastmod time.Time) File {
	var typedArray = uint8ArrayConstructor.New(len(data))
	js.CopyBytesToJS(typedArray, data)
	var date = dateConstructor.New(lastmod.UnixMilli())
	var blob = blobConstructor.New([]interface{}{typedArray}, map[string]interface{}{
		"type": "application/octet-stream",
	})
	var file = fileConstructor.New([]interface{}{blob}, name, map[string]interface{}{
		"type":         mimetype,
		"lastModified": date,
	})
	return File(file)
}

func NewFromJS(file js.Value) File {
	return File(file)
}

func (f File) MarshalJS() js.Value {
	return js.Value(f)
}

func (f File) Name() string {
	return js.Value(f).Get("name").String()
}

func (f File) Size() int {
	return js.Value(f).Get("size").Int()
}

func (f File) MimeType() string {
	return js.Value(f).Get("type").String()
}

func (f File) LastModified() time.Time {
	return time.UnixMilli(int64(js.Value(f).Get("lastModified").Int()))
}
