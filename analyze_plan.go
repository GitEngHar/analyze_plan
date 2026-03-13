package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Plan struct {
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type ResourceChange struct {
	Address string `json:"address"`
	Change  Change `json:"change"`
}

type Change struct {
	Actions []string `json:"actions"`
}

func main() {
	var (
		create, update, del, noop int
	)
	if len(os.Args) < 2 {
		fmt.Println("usage: go run analyze_plan.go plan.json")
		os.Exit(1)
	}

	file := os.Args[1]

	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var plan Plan
	err = json.Unmarshal(data, &plan)
	if err != nil {
		panic(err)
	}

	for _, rc := range plan.ResourceChanges {
		for _, action := range rc.Change.Actions {

			switch action {
			case "create":
				create++
			case "update":
				update++
			case "delete":
				del++
			}
			if action != "no-op" {
				fmt.Printf("%s -> %s\n", rc.Address, action)
			}

		}
	}

	fmt.Println("\nSummary")
	fmt.Println("create:", create)
	fmt.Println("update:", update)
	fmt.Println("delete:", del)
	fmt.Println("noop:", noop)
}
