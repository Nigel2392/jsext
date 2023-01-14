# Example for misc.JobTree

This example shows how to use the JobTree component.

```go
Application.Render(JobTree(&misc.Jobtree{
	Items: []misc.JobtreeItem{
		{
			Company:   "Company 1",
			Title:     "Title 1",
			Tags:      []string{"Tag 1", "Tag 2", "Tag 3"},
			StartDate: "2020-01-01",
			EndDate:   "2020-01-01",
		},
		{
			Company:   "Company 2",
			Title:     "Title 2",
			Tags:      []string{"Tag 1", "Tag 2", "Tag 3"},
			StartDate: "2020-01-01",
			EndDate:   "2020-01-01",
		},
		{
			Company:   "Company 3",
			Title:     "Title 3",
			Tags:      []string{"Tag 1", "Tag 2", "Tag 3"},
			StartDate: "2020-01-01",
			EndDate:   "2020-01-01",
		},
	},
	TagBackgroundColors: []string{
		"#ff0000",
		"#00ff00",
		"#0000ff",
	},
	Width: "100%",
}))
```