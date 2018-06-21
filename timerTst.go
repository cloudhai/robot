package main

import (
	"time"
	"fmt"
)

func main(){
	t := time.NewTimer(time.Second*5)
	fmt.Println("start")
	<- t.C
	fmt.Printf("finsh")
}
