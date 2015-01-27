# Drip
Custom spin on leaky bucket algorithm to help with rate limited APIs. The go routine was overkill (could have just done some math) but I was learning about them + channels and decided to use both.


## Example usage

```go
package main

import (
    "github.com/joncalhoun/drip"
    "fmt"
)

func main() {
    b := Bucket{
        consumed: 32,
        capacity: 60,
        dripInterval: 3 * time.Second,
        perDrip: 60
    }


    b.Start() // this returns an error if the bucket was already started. generally this doesn't matter.
    defer b.Stop() // this needs to be called to stop the go routine that is created when Start() is called

    count := 0
    for {
        err := b.Consume(2)
        if err != nil {
            fmt.Println("Sleeping a bit.")
            time.Sleep(time.Second)
        } else {
            count++
            fmt.Println("Consuming 2")
        }

        if count > 100 {
            break
        }
    }

}

```
