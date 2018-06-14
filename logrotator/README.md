## example
```go
package main

import (
	"log"
	"logrotator"
	"time"
)

func main() {
	r := logrotator.New("rotate.log", 8)
	log.SetOutput(r)

	count := 0
	for {
		time.Sleep(time.Second * 10)
		count++
		log.Printf("%d\n", count)
	}
}
```

