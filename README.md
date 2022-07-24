![Coverage](https://img.shields.io/badge/Coverage-85.7%25-brightgreen)
# gostein

[Stein](https://steinhq.com/) API wrapper for Go.

## Installation

```go
go get github.com/pararang/gostein
```

## Usage
### Create a new client
```go
import "github.com/pararang/gostein"
...
steinClient = gostein.New("http://yourstein.host/v1/storage/your-api-id", nil, nil)

// with basic auth
steinClient = gostein.New(
    "http://yourstein.host/v1/storage/your-api-id", 
    nil, 
    &AuthParams{
		Username: "pararang",
		Password: "pararang123",
	})
```
> If HTTP Client is not provided (nil) on the second parameter, `DefaultClient` from http golang stdlib will be used.

### Read
#### Get data
```go
...
// Get all data
data, err := steinClient.Get("sheet1", gostein.GetParams{})

// Get with offset and limit
data, err := steinClient.Get("sheet1", gostein.GetParams{Offset: 0, Limit: 10})
...
```
The `data` will be in type of `[]map[string]interface{}`. To convert the data to specific struct, I recomended using [maptostructure package](https://github.com/mitchellh/mapstructure).

#### Search data
Look up rows in a sheet by a specific value on column(s).
```go
data, err := steinClient.Get("sheet1", gostein.GetParams{
    Limit: 10,
    Search: map[string]string{
            "column_1": "value_column_1",
            "column_2": "value_column_2",
        }
    })
...
```

### Add
#### Add single data
```go
...
resp, err = steinClient.Add("sheet1",
    map[string]interface{}{
        "column_1": "value_1",
        "column_2": "value_2",
    })
// handle err and do something with resp
...
```

#### Add bulk/multiple data
```go
resp, err := s.Add("gostein", 
    map[string]interface{}{
        "column_1": "value_1-a",
        "column_2": "value_2-a",
    }, 
    map[string]interface{}{
        "column_1": "value_1-b",
        "column_2": "value_2-b",
    })

// with better code readability, utilize the variadic function definition
rows := []map[string]interface{}{
    {"column_1": "value_1-a", "column_2": "value_2-a"},
    {"column_1": "value_1-b", "column_2": "value_2-b"},
}

resp, err = steinClient.Add("sheet1", rows...)
// handle err then do something with resp
...
```
### Update
The update operation return string that indicate ranges of the sheet that was updated. Example: `Sheet1!A3:B3`.
#### Update single row
`Limit 1` indicate to update only the first row that match the `Condition`.
```go
...
resp, err := sc.Update("Sheet1", UpdateParams{
			Condition: map[string]string{
				"column_1": "if_has_this_value",
			},
			Set: map[string]string{
				"column_2": "then_update_this_colum_value",
			},
			Limit: 1,
		})
// handle err and do something with resp
...
```
#### Update with multiple conditions
All `Condition` will be translated to `AND` condition.
```go
...
resp, err := sc.Update("Sheet1", UpdateParams{
			Condition: map[string]string{
				"column_1": "if_has_this_value",
				"column_3": "and_if_has_this_value",
			},
			Set: map[string]string{
				"column_2": "then_update_this_colum_value",
			},
		})
// handle err and do something with resp
...
```
> :warning: **If `Limit` is not set, all rows those match the `Condition` will be updated.**
### Delete
The delete operation return int64 that indicate the number of rows that was deleted.
#### Delete single row
`Limit 1` indicate to delete only the first row that match the `Condition`.
```go
...
resp, err := sc.Delete("Sheet1", DeleteParams{
			Condition: map[string]string{
				"column_1": "if_has_this_value",
			},
			Limit: 1,
		})
// handle err and do something with resp
...
```
#### Delete with multiple conditions
All `Condition` will be translated to `AND` condition.
```go
...
resp, err := sc.Delete("Sheet1", DeleteParams{
			Condition: map[string]string{
				"column_1": "if_has_this_value",
				"column_3": "and_if_has_this_value",
			},
		})
// handle err and do something with resp
...
```
> :warning: **If `Limit` is not set, all rows those match the `Condition` will be deleted.**

## TODO
- [x] Read data (https://docs.steinhq.com/read-data)
- [x] Read data with conditions (https://docs.steinhq.com/search-data)
- [x] Add data (https://docs.steinhq.com/add-rows)
- [x] Update data (https://docs.steinhq.com/update-rows)
- [x] Delete data (https://docs.steinhq.com/delete-rows)
- [x] Authentication (https://docs.steinhq.com/authentication)
