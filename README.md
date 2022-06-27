![Coverage](https://img.shields.io/badge/Coverage-87.5%25-brightgreen)
# gostein

[Stein](https://steinhq.com/) API wrapper for Go.

## Installation

```go
go get github.com/pararang/gostein
```

## Usage
```go
    ...
    steinClient = gostein.New("http://yourstein.host/v1/storage/your-api-id", nil)

    data, err := steinClient.Get("sheet1", gostein.SearchParams{})
    if err != nil {
        // handle error
    }

    // do something with data
    ...
```
> If HTTP Client is not provided (nil) on the second parameter, `DefaultClient` from http golang stdlib will be used.

The data will be in type of `[]map[string]interface{}`. To convert the data to specific struct, I recomended using [maptostructure package](https://github.com/mitchellh/mapstructure).

## TODO
- [x] Read data (https://docs.steinhq.com/read-data)
- [x] Read data with conditions (https://docs.steinhq.com/search-data)
- [ ] Add data (https://docs.steinhq.com/add-rows)
- [ ] Update data (https://docs.steinhq.com/update-rows)
- [ ] Delete data (https://docs.steinhq.com/delete-rows)
- [ ] Authentication (https://docs.steinhq.com/authentication)
