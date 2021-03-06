# gomiddlewares

- logger gin+zap
- cors

## Usage

### Start using it

Download and install it:

```sh
$ go get github.com/judascrow/gomiddlewares

```

Import it in your code:

```go
import "github.com/judascrow/gomiddlewares"
```

## Example

See the [example](example/main.go).

[embedmd]: # "example/main.go go"

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/judascrow/gomiddlewares"
)

func main() {
	r := gin.New()

	r.Use(gomiddlewares.GoLogger())
	r.Use(gomiddlewares.GoCors())

	// Example ping request.
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "API is Online",
		})
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

```
