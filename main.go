package main

import (
	"github.com/leartgjoni/umigrate/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		cmd.LogErr("%s\n", err)
	}
}
