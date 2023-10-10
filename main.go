package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func loadDotEnv() error {
	data, err := os.ReadFile(".env")
	
	if err != nil {
		return errors.New("there was an error opening .env file")
	}

	fileData := string(data)
	pairVals := strings.Split(fileData, "\n")

	for _, v := range pairVals {
		key, val := strings.Split(v, "=")[0], strings.Split(v, "=")[1]
		os.Setenv(key, val)
	}

	return nil
}

func main() {
	loadDotEnv()
	
	port := os.Getenv("PORT")

	if port == "" {
		err := errors.New("error: PORT not found")
		fmt.Println(err)
	}

	fmt.Println("Listening on port", port)
}
