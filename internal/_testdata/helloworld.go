package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	str := "Hello World"
	dat, _ := json.Marshal(str)
	fmt.Println(dat)
	os.Stdout.Write(dat)
}
