# bricklinkapi

A Golang API client for the [bricklink.com](https://www.bricklink.com) API.

You need to provide your Consumer key and secret als well as a Token key and secret. Head over to the [Bricklink Store API](http://apidev.bricklink.com/redmine/projects/bricklink-api/wiki) page and register.

## Usage

```go
package main

import (
    "fmt"
    "github.com/brickventory/bricklinkapi"
)

func main() {
    // Load keys and secrets for accessing the Bricklink Api

    // we need a new Bricklink
    bl := bricklinkapi.New(CONSUMER_KEY, CONSUMER_SECRET, TOKEN_VALUE, TOKEN_SECRET)

    // try some simple query
    // query for part #3004, which is the basic 1x4 brick
    fmt.Println(bl.GetItem("part", "3004"))
}
```