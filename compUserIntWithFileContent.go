// compate user input string with string from file
// 
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	// set var for scan
	var userInt string

	// greeting message
	fmt.Println("Please input your string: ")

	// take input from user
	fmt.Scan(&userInt)

	// read the file (output is bytes)
	dataFromFile, err := os.ReadFile("/Users/user/file.txt")
	if err != nil {
		log.Printf("File %s does not exist or can not be read...", dataFromFile)
	}
	// convert bytes to string
	conToStr := string(dataFromFile)

	// convert to slice ( tab and new line ignore )
	conToSlice := (strings.Fields(conToStr))
	endOfFile := len(conToSlice)
	//fmt.Println(endOfFile)
	// loop for slice
	for i, el := range conToSlice {
		//sl[len(sl)-1] read the last element for slice

		if el == userInt {
			fmt.Printf("Host %s will be reboot ...\n", el)
		} else if i == endOfFile-1 {
			fmt.Printf("I'm sorry I don't know host %s\n", userInt)
		} else {
		}
	}
}
