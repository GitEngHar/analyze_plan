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

type ResultResourceCount struct {
	Create  int
	Delete  int
	Update  int
	Replace int
}

type ResultResourceAddress struct {
	Create  []string
	Delete  []string
	Update  []string
	Replace []string
}

func main() {
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
	countResult, addressResult := summarizeResource(plan)
	outPutSummary(countResult, addressResult)
}

func summarizeResource(plan Plan) (ResultResourceCount, ResultResourceAddress) {
	var (
		createCount, updateCount, delCount, replaceCount            int
		createAddress, updateAddress, deleteAddress, replaceAddress []string
	)
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
			createCount++
			createAddress = append(createAddress, rc.Address)
		case "update":
			updateCount++
			updateAddress = append(updateAddress, rc.Address)
		case "delete":
			delCount++
			deleteAddress = append(deleteAddress, rc.Address)
		case "replace":
			replaceCount++
			replaceAddress = append(replaceAddress, rc.Address)
		}
		if action != "no-op" {
			fmt.Printf("%s -> %s\n", rc.Address, action)
		}
	}
	return ResultResourceCount{
			Create:  createCount,
			Update:  updateCount,
			Delete:  delCount,
			Replace: replaceCount,
		}, ResultResourceAddress{
			Create:  createAddress,
			Update:  updateAddress,
			Delete:  deleteAddress,
			Replace: replaceAddress,
		}
}

func outPutSummary(countResult ResultResourceCount, addressResult ResultResourceAddress) {
	fmt.Println("\nSummary\n---")
	fmt.Println("create:", countResult.Create)
	for _, address := range addressResult.Create {
		fmt.Println("+ ", address)
	}
	fmt.Println("update:", countResult.Update)
	for _, address := range addressResult.Update {
		fmt.Println("~ ", address)
	}
	fmt.Println("delete:", countResult.Delete)
	for _, address := range addressResult.Delete {
		fmt.Println("- ", address)
	}
	fmt.Println("replace:", countResult.Replace)
	for _, address := range addressResult.Replace {
		fmt.Println("+/- ", address)
	}
}
