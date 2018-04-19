# las
A go library for reading LAS lidar files

Currently only supporting LAS version 1.1, but soon to include all versions (and hopefully a LAS writer).

## Usage
Open a LAS file and read the contents
```golang
package main

import (
        "fmt"

        "github.com/andygarfield/las"
)

func main() {
        l, _ := las.Open("./example.las")

        for i := 0; i < 300; i++ {
                p := l.ReadPoint()

                fmt.Printf("%f %f %f\n", p.X, p.Y, p.Z)
        }
}
```
