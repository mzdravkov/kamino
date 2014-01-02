package main

import (
	"fmt"
)

func main() {
	fmt.Println(Config["nginx_bin"])
	fmt.Println(findFreePort())
	fmt.Println(Deploy("llama"))
}
