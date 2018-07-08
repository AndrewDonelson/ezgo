package ezgo

import (
	"fmt"

	cfg "github.com/NlaakStudios/ezgo/config"
)

//LoadConfig either load existing config file or creates and returns a default one
func LoadConfig() {
	myConfig, err := cfg.LoadConfig("", "/")
	if err != nil {
		fmt.Println("Error creating config.")
	}
	fmt.Printf("Log loaded %+v", myConfig)
}
