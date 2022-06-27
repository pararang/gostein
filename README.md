![Coverage](https://img.shields.io/badge/Coverage-87.5%25-brightgreen) [![GoReportCard](https://goreportcard.com/badge/github.com/pararang/gostein)](https://goreportcard.com/report/github.com/pararang/gostein) [![Go version of a Go module](https://img.shields.io/github/go-mod/go-version/pararang/gostein.svg)](https://github.com/pararang/gostein) [![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/pararang/gostein/blob/main/LICENSE.md) [![Open Source](https://badges.frapsoft.com/os/v2/open-source.svg?v=103)](#)
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
steinClient = gostein.New("http://yourstein.host/v1/storage/your-api-id", nil)
```
> If HTTP Client is not provided (nil) on the second parameter, `DefaultClient` from http golang stdlib will be used.

### Get/Search data
```go
...
data, err := steinClient.Get("sheet1", gostein.SearchParams{})
// handle error and do something with data
...
```
The `data` will be in type of `[]map[string]interface{}`. To convert the data to specific struct, I recomended using [maptostructure package](https://github.com/mitchellh/mapstructure).

### Add data
#### Add single data
```go
...
resp, err = steinClient.Add("sheet1",
    map[string]interface{}{
        "column_1": "value_1",
        "column_2":  "value_2",
    })
// handle err adn do something with resp
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

## TODO
- [x] Read data (https://docs.steinhq.com/read-data)
- [x] Read data with conditions (https://docs.steinhq.com/search-data)
- [x] Add data (https://docs.steinhq.com/add-rows)
- [ ] Update data (https://docs.steinhq.com/update-rows)
- [ ] Delete data (https://docs.steinhq.com/delete-rows)
- [ ] Authentication (https://docs.steinhq.com/authentication)
