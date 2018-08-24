package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/crit/inventory/cmd/inventory/cmd"
	"github.com/crit/inventory/internal/storage/models"
	"github.com/crit/inventory/internal/storage/providers"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	store, err := providers.Bolt(usr.HomeDir + "/.inventory")
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	models.Register(store)

	cmd.Execute()
}
