// Be careful it works as grep 
package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

func checkData(x string) bool {

	// read file
	data, err := os.ReadFile("/Users/username/file.txt")
	if err != nil {
		log.Println("File does not exist or can not read it!")
	}
	isExist, err := regexp.Match(x, data)
	if err != nil {
		panic(err)
	}
	return isExist

}
func main() {

	// get input fro user
	fmt.Println("Please input your string:")
	var userInput string
	fmt.Scan(&userInput)
	boolValue := checkData(userInput)

	if boolValue {
		fmt.Printf("File hosts contain %s\n", userInput)
	} else {
		fmt.Printf("I don't know about %s\n", userInput)
	}

}
