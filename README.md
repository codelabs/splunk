# splunk
Splunk Libraries in golang

* [![Build Status](https://travis-ci.org/codelabs/splunk.svg?branch=master)](https://travis-ci.org/codelabs/splunk)
* [![Coverage Status](https://coveralls.io/repos/github/codelabs/splunk/badge.svg?branch=master)](https://coveralls.io/github/codelabs/splunk?branch=master)

```{go}
package main

import (
    "fmt"
    "github.com/codelabs/splunk"
)

func main() {

    var user = &User{
        username: "admin",
        password: "changeme",
    }

    session, err := splunk.Connect(user, "localhost", 5500)
    if err != nil {
        fmt.Println("Error " + err)
    }
}
```
