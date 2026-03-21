package main

import (
	"fmt"
	"os"
	"terraform-sammary/internal/handler"
)

var isUsePolicy bool

func main() {
	h := handler.NewActionsHandler()
	err := h.Handle(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
