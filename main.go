package main

import (
	"fmt"
)

func main() {
	fmt.Println(Config["nginx_bin"])
	fmt.Println(findFreePort())
	fmt.Println(Deploy("llama7"))
	//testOpts := map[string]string{"test": "true", "cool": "true"}
	//addLocation("cooltest", testOpts)
}
