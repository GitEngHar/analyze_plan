package summary

import (
	"context"
	"fmt"
	"terraform-sammary/internal/domain/repository"
	"time"
)

type Confirm interface {
	Execute(ctx context.Context) (error, *ResultResourceAddress, *ResultResourceTypeToNames)
}

type ConfirmImpl struct {
	repo repository.PlanRepository
}

func NewConfirm(repo repository.PlanRepository) Confirm {
	return ConfirmImpl{
		repo: repo,
	}
}

func (c ConfirmImpl) Execute(ctx context.Context) (error, *ResultResourceAddress, *ResultResourceTypeToNames) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	plan, err := c.repo.ReadPlan()
	if err != nil {
		return err, nil, nil
	}
	var (
		createCount, updateCount, delCount, replaceCount            int
		createAddress, updateAddress, deleteAddress, replaceAddress []string
		deleteResourceNames, replaceResourceNames                   ResourceNames
	)
	deleteResourceTypeToNames := make(map[string]ResourceNames, len(plan.ResourceChanges))
	replaceResourceTypeToName := make(map[string]ResourceNames, len(plan.ResourceChanges))
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
			deleteResourceNames = append(deleteResourceNames, rc.Name)
			deleteResourceTypeToNames[rc.Type] = deleteResourceNames
		case "replace":
			replaceCount++
			replaceAddress = append(replaceAddress, rc.Address)
			replaceResourceNames = append(replaceResourceNames, rc.Name)
			replaceResourceTypeToName[rc.Type] = replaceResourceNames
		}
		if action != "no-op" {
			fmt.Printf("%s -> %s\n", rc.Address, action)
		}
	}
	resourceCount := ResultResourceCount{
		Create:  createCount,
		Update:  updateCount,
		Delete:  delCount,
		Replace: replaceCount,
	}
	resourceAddress := ResultResourceAddress{
		Create:  createAddress,
		Update:  updateAddress,
		Delete:  deleteAddress,
		Replace: replaceAddress,
	}
	resourceTypeToNames := ResultResourceTypeToNames{
		Delete:  deleteResourceTypeToNames,
		Replace: replaceResourceTypeToName,
	}
	outPutSummary(resourceCount, resourceAddress)
	return nil, &resourceAddress, &resourceTypeToNames
}

func outPutSummary(countResult ResultResourceCount, addressResult ResultResourceAddress) {
	fmt.Println("\nTerraform Plan Summary\n----------------------")
	fmt.Println("create:", countResult.Create)
	fmt.Println("update:", countResult.Update)
	fmt.Println("delete:", countResult.Delete)
	fmt.Println("replace:", countResult.Replace)
	for _, address := range addressResult.Create {
		fmt.Println("+ ", address)
	}
	for _, address := range addressResult.Update {
		fmt.Println("~ ", address)
	}
	for _, address := range addressResult.Delete {
		fmt.Println("- ", address)
	}
	for _, address := range addressResult.Replace {
		fmt.Println("+/- ", address)
	}
}
