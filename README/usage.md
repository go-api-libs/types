To install the library, use the following command:

```shell
go get github.com/go-api-libs/types
```


When using any of the API libraries from the [go-api-libs](https://github.com/go-api-libs/) organization, you can trust that types in successful API responses are valid.

Beyond that, here's how you can use some types from this library:

```go
package main

import (
	"fmt"

	"github.com/go-api-libs/types"
)

func main() {
	// Using Email type
	email := types.Email("user@example.com")
	if err := email.Validate(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Valid email address")
	}
}
```
