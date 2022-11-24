<h1 align="center">fastreq</h1>
<p align="center">Fast, convenient and simple HTTP client based on fasthttp for Go (inspired by Fiber and fasthttp)</p>

## Features

- Extreme performance
- Low memory footprint
- Middleware support

## Installation

```go
go get github.com/wnanbei/fastreq

import "github.com/wnanbei/fastreq"
```

## Usage

```go
import "github.com/wnanbei/fastreq"

resp, err := fastreq.Get("https://hello-world", fastreq.NewArgs())
if err != nil {
    panic(err)
}
fmt.Println(resp.BodyString())
```

## BenchMark