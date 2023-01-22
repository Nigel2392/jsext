package encoders

import (
	"bytes"
	"io"
	"mime/multipart"
)

// Define a type for multipart files
type File struct {
	FileName  string    // Name of the file
	FieldName string    // Name of the field
	Reader    io.Reader // Reader of the file
}

func ToMultipart(data map[string]any, files ...File) []byte {
	bod := &bytes.Buffer{}
	writer := multipart.NewWriter(bod)
	for _, f := range files {
		part, err := writer.CreateFormFile(f.FieldName, f.FileName)
		if err != nil {
			panic(err)
		}
		io.Copy(part, f.Reader)
	}
	for k, v := range data {
		writer.WriteField(k, toString(v))
	}
	writer.Close()
	return bod.Bytes()
}
