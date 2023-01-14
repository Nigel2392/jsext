# Example for misc.JobTree

This example shows how to use the RoadMap component.

```go
Application.Render(misc.RoadMap(&misc.RoadMapOptions{
	Items: []misc.RoadMapItem{
		{
			Name:        "Company 1",
			Title:       "Title 1",
			Description: "Description 1",
			Tags:        []string{"Tag 1", "Tag 2", "Tag 3"},
			StartDate:   "2020-01-01",
			EndDate:     "2020-01-01",
		},
		{
			Name:        "Company 2",
			Title:       "Title 2",
			Description: "Description 2",
			Tags:        []string{"Tag 1", "Tag 2", "Tag 3"},
			StartDate:   "2020-01-01",
			EndDate:     "2020-01-01",
		},
		{
			Name:        "Company 3",
			Title:       "Title 3",
			Description: "Description 3",
			Tags:        []string{"Tag 1", "Tag 2", "Tag 3"},
			StartDate:   "2020-01-01",
			EndDate:     "2020-01-01",
		},
	},
	TagBackgroundColors: []string{
		"#ff0000",
		"#00ff00",
		"#0000ff",
	},
	Width:        "100%",
	Style:        misc.RoadMapStyleTwo,
	CardMargin:   "30px",
	DivisorWidth: "4px",
}))
```