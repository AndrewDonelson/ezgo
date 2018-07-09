package main

import (
	"fmt"

	cfg "github.com/NlaakStudios/ezgo/config"
)

func main() {
	myConfig, err := cfg.LoadConfig("", "/")
	if err != nil {
		fmt.Println("Error loading config.")
		myConfig = cfg.DefaultConfig()
		if myConfig != nil {
			fmt.Println("Default configuration set.")
		}
	}
	fmt.Printf("Log loaded %+v", myConfig)
}
