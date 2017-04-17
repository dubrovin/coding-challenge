package main

import (
	"fmt"
	"time"
)

func main() {

	tNow := time.Now()
	tBefore := tNow.Add(-time.Second * 60)

	fmt.Println(tNow.String(), tBefore.String())
}
