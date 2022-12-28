//go:build js && wasm
// +build js,wasm

package forms_test

import (
	"testing"

	"github.com/Nigel2392/jsext/components/forms"
)

func TestForms(t *testing.T) {
	var data = make(map[string]string)
	data["Name"] = "Jane"
	// data["time"] = time.Now().Format("2006-01-02T15:04")
	data["Age"] = "22"
	data["Friend_Age"] = "25"    //TODO TODO TODO
	data["Friend_Name"] = "John" //TODO TODO TODO
	data["Friend_Other_Age"] = "25"
	data["Friend_Other_Name"] = "Jippy"

	type Me struct {
		Name   string
		Age    int
		Friend struct {
			Name  string
			Age   int
			Other *struct {
				Name string
				Age  int
			}
		}
	}
	var testForm Me
	forms.FormDataToStruct(data, &testForm)
	if testForm.Name != "Jane" {
		t.Error("Name is not Jane")
	}
	if testForm.Age != 22 {
		t.Error("Age is not 22")
	}
	if testForm.Friend.Name != "John" {
		t.Error("Friend name is not John")
	}
	if testForm.Friend.Age != 25 {
		t.Error("Friend age is not 25")
	}
	if testForm.Friend.Other == nil {
		t.Error("Other is nil")
	}
	if testForm.Friend.Other.Name != "Jippy" {
		t.Error("Friend other name is not Jippy")
	}
	if testForm.Friend.Other.Age != 25 {
		t.Error("Friend other age is not 25")
	}
	t.Log(testForm)
}
