[![Build Status](https://travis-ci.org/sbabiv/xml2map.svg?branch=master)](https://travis-ci.org/sbabiv/xml2map)
[![cover.run](https://cover.run/go/github.com/sbabiv/xml2map.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fsbabiv%2Fxml2map)
[![Go Report Card](https://goreportcard.com/badge/github.com/sbabiv/xml2map)](https://goreportcard.com/report/github.com/sbabiv/xml2map)
[![GoDoc](https://godoc.org/github.com/sbabiv/xml2map?status.svg)](https://godoc.org/github.com/sbabiv/xml2map)


# xml2map
XML to MAP converter written Golang

Sometimes there is a need for the representation of previously unknown structures. Such a universal representation is usually a string in the form of JSON, XML, or the structure of data map. similar to the map[string]interface{} or map[interface{}]interface{}.

This is a converter from the old XML format to map[string]interface{} Golang

For example, the map[string]interface{} can be used as a universal type in template generation. Golang "text/template" and etc.

## Getting started

#### 1. install

``` sh
go get -u github.com/sbabiv/xml2map
```

Or, using dep:

``` sh
dep ensure -add github.com/sbabiv/xml2map
```


#### 2. use it

```go

func main() {
	data := `<container uid="FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6" lifetime="2019-10-10T18:00:11">
				<cats>
					<cat>
						<id>CDA035B6-D453-4A17-B090-84295AE2DEC5</id>
						<name>moritz</name>
						<age>7</age> 	
						<items>
							<n>1293</n>
							<n>1255</n>
							<n>1257</n>
						</items>
					</cat>
					<cat>
						<id>1634C644-975F-4302-8336-1EF1366EC6A4</id>
						<name>oliver</name>
						<age>12</age>
					</cat>
					<dog color="gray">hello</dog>
				</cats>
				<color>white</color>
				<city>NY</city>
			</container>`

	decoder := xml2map.NewDecoder(strings.NewReader(data))
	result, err := decoder.Decode()

	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", result)
	}
	
	v := result["container"].
		(map[string]interface{})["cats"].
			(map[string]interface{})["cat"].
				([]map[string]interface{})[0]["items"].
					(map[string]interface{})["n"].([]string)[1]
					
	fmt.Printf("n[1]: %v\n", v)
}

```

## Output

```go
map[container:map[@uid:FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6 @lifetime:2019-10-10T18:00:11 cats:map[cat:[map[id:CDA035B6-D453-4A17-B090-84295AE2DEC5 name:moritz age:7 items:map[n:[1293 1255 1257]]] map[id:1634C644-975F-4302-8336-1EF1366EC6A4 name:oliver age:12]] dog:map[@color:gray #text:hello]] color:white city:NY]]

result: 1255
```

## Benchmark


```go
$ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/sbabiv/xml2map
BenchmarkDecoder-8         50000             29773 ns/op           15032 B/op        261 allocs/op
PASS
ok      github.com/sbabiv/xml2map       1.805s
```


## Test

```go
$ go test -cover -v .
=== RUN   TestStartAttrs
--- PASS: TestStartAttrs (0.00s)
=== RUN   TestPars
--- PASS: TestPars (0.00s)
=== RUN   TestFuzz1000
--- PASS: TestFuzz1000 (0.00s)
=== RUN   TestErrDecoder
--- PASS: TestErrDecoder (0.00s)
    decoder_test.go:101: result: map[] err: XML syntax error on line 10: unexpected EOF
=== RUN   TestEmpty
--- PASS: TestEmpty (0.00s)
=== RUN   TestSpaces
--- PASS: TestSpaces (0.00s)
=== RUN   TestInvalidStartIndex
--- PASS: TestInvalidStartIndex (0.00s)
=== RUN   TestDecode
--- PASS: TestDecode (0.00s)
PASS
coverage: 93.4% of statements
ok      github.com/sbabiv/xml2map       (cached)        coverage: 93.4% of statements
```

## Licence
[MIT](https://opensource.org/licenses/MIT)

## Author 
Babiv Sergey
