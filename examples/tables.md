```go
	type Friend struct {
		Name string `table:"Friend Name" width:"100px"` // This is the column name
	}

	type Person struct {
		Name   string                         `table:"Name" width:"100px"` // This is the column name
		Age    int                            `table:"Age" width:"100px"`
		Friend `table:"Friend" width:"100px"` // This is the column name
	}
    
    var people = []Person{
		{Name: "John", Age: 20, Friend: Friend{
			Name: "Jane",
		}},
		{Name: "Jane", Age: 21, Friend: Friend{
			Name: "Jack",
		}},
		{Name: "Jack", Age: 22, Friend: Friend{
			Name: "Jill",
		}},
		{Name: "Jill", Age: 23, Friend: Friend{
			Name: "Joe",
		}},
		{Name: "Joe", Age: 24, Friend: Friend{
			Name: "Jen",
		}},
		{Name: "Jen", Age: 25, Friend: Friend{
			Name: "Jenny",
		}},
		{Name: "Jenny", Age: 26, Friend: Friend{
			Name: "John",
		}},
	}
    
    var div = elements.Div()
	var Table = table.NewFromStruct("100%", people, nil)
	div.Append(Table.Run())
```