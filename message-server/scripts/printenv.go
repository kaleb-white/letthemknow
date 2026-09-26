package main 

import (
	"fmt"
	"os"
)

func main() {
	for _, str := range os.Environ() {
		fmt.Println(str)
	}
}
