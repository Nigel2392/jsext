# TinyGO HTTP client

This example shows how to use the APIClient inside of a TinyGO application.
The client is essentially just a wrapper around the javascript fetch API.

```go
	var dataMap = make(map[string]interface{})
	dataMap["name"] = "John Doe"
	dataMap["age"] = 22
	dataMap["height"] = 1.84
	dataMap["weight"] = 80.0
	dataMap["Gender"] = "Male"
	dataMap["isHuman"] = true
	dataMap["isAlive"] = true

	var client = requester.NewAPIClient().Post("https://httpbin.org/post")
	client.WithData(dataMap, requester.JSON)
	var f, err = client.Do()
	if err != nil {
		panic(err)
	}
	println(string(f.Body))
	for key, value := range f.Headers {
		println(fmt.Sprintf("%v: %v", key, value))
	}
```