# go-tomcat

## Getting started

### Prerequisites

- **[Go](https://go.dev/)**: any one of the **five latest major** [releases](https://go.dev/doc/devel/release) (we test it with these).

### Getting Gin

run the following Go command to install the `gin` package:

```sh
$ go get -u github.com/alikWu/go-tomcat
```

### Running Gin

First you need to import Gin package for using Gin, one simplest example likes the follow `example.go`:

```go
package main

import (
  "fmt"

	GoTomcat "github.com/alikWu/go-tomcat"
	"github.com/alikWu/go-tomcat/servlet"
	"github.com/alikWu/go-tomcat/webapps"
	"github.com/alikWu/go-tomcat/webapps/test"
)

func main() {
  bs := GoTomcat.NewBootStrap(":8008")
	bs.SetMaxConnections(200)
	bs.RegisterServlets([]servlet.Servlet{&webapps.HelloServlet{}, &webapps.StaticResourceProcessor{}})
	bs.RegisterFilters([]servlet.Filter{&test.TestFilter{}})
	if err := bs.Start(); err != nil {
		fmt.Println("Start Fail ,err ", err)
	}
}
```

And use the Go command to run the demo:

```
# run example.go and visit http://localhost:8008/servlet/hello on browser
$ go run example.go
```

### Learn more examples

todo
