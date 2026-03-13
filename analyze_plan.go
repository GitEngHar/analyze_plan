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
		create, update, del, replace int
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
		var action string
		actions := rc.Change.Actions

		switch {
		case len(actions) == 1 && actions[0] == "create":
			action = "create"
		case len(actions) == 1 && actions[0] == "update":
			action = "update"
		case len(actions) == 1 && actions[0] == "delete":
			action = "delete"
		case len(actions) == 2 && actions[0] == "delete" && actions[1] == "create":
			action = "replace"
		default:
			action = "no-op"
		}

		switch action {
		case "create":
			create++
		case "update":
			update++
		case "delete":
			del++
		case "replace":
			replace++
		}
		if action != "no-op" {
			fmt.Printf("%s -> %s\n", rc.Address, action)
		}

	}

	fmt.Println("\nSummary")
	fmt.Println("create:", create)
	fmt.Println("update:", update)
	fmt.Println("delete:", del)
	fmt.Println("replace:", replace)

}
