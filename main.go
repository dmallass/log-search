package main

import (
	"cribl_take_home/loggenerator"
	"cribl_take_home/search"
	"fmt"
	"os"
)

func main() {
	if err := loggenerator.GenerateLog("app.log"); err != nil {
		fmt.Println("error generating log", err)
		os.Exit(1)
	}
	if err := search.RunWebserver(); err != nil {
		fmt.Println("error running webserver", err)
		os.Exit(1)
	}
}
