package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type Plan struct {
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type ResourceChange struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Name    string `json:"name"`
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

type ResourceNames []string

type ResultResourceTypeToNames struct {
	Delete  map[string]ResourceNames
	Replace map[string]ResourceNames
}

type ResultResourceAddress struct {
	Create  []string
	Delete  []string
	Update  []string
	Replace []string
}

type Rule struct {
	Resource string `yaml:"resource"`
	Severity string `yaml:"severity"`
}

type ResourceProtectPolicy struct {
	ProtectRules []Rule `yaml:"rules"`
}

type ResourceProtectPolicies struct {
	DeleteProtectPolicy  ResourceProtectPolicy
	ReplaceProtectPolicy ResourceProtectPolicy
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: go run analyze_plan.go plan.json policy_path")
		os.Exit(1)
	}
	planFilePath := os.Args[1]
	policyFilePath := os.Args[2]
	data, err := os.ReadFile(planFilePath)
	if err != nil {
		panic(err)
	}
	var plan Plan
	err = json.Unmarshal(data, &plan)
	if err != nil {
		panic(err)
	}
	resourcePolicy, err := loadPolicies(policyFilePath)
	if err != nil {
		panic(err)
	}
	countResult, addressResult, resourceTypeResult := summarizeResource(plan)
	outPutSummary(countResult, addressResult)
	replaceDetected(addressResult)
	resourcePolicyViolationDetected(resourceTypeResult, resourcePolicy)
}

func summarizeResource(plan Plan) (ResultResourceCount, ResultResourceAddress, ResultResourceTypeToNames) {
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
		}, ResultResourceTypeToNames{
			Delete:  deleteResourceTypeToNames,
			Replace: replaceResourceTypeToName,
		}
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

func replaceDetected(addresses ResultResourceAddress) {
	var isReplaceWarning = len(addresses.Replace) > 0
	fmt.Println("\nReplace Detected\n----------------")
	// replace warn
	if isReplaceWarning {
		for _, address := range addresses.Replace {
			fmt.Println("+/- ", address)
		}
	}
}

func loadPolicies(dir string) (*ResourceProtectPolicies, error) {
	const (
		deletePolicyYamlName  = "delete.yaml"
		replacePolicyYamlName = "replace.yaml"
	)

	var resourceProtectPolicies ResourceProtectPolicies
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := filepath.Join(dir, file.Name())
		switch file.Name() {
		case deletePolicyYamlName:
			protectPolicy, er := loadProtectPolicy(path)
			if er != nil {
				return nil, er
			}
			resourceProtectPolicies.DeleteProtectPolicy = *protectPolicy
		case replacePolicyYamlName:
			protectPolicy, er := loadProtectPolicy(path)
			if er != nil {
				return nil, er
			}
			resourceProtectPolicies.ReplaceProtectPolicy = *protectPolicy
		default:
			continue
		}
	}
	return &resourceProtectPolicies, nil
}

func loadProtectPolicy(path string) (*ResourceProtectPolicy, error) {
	data, er := os.ReadFile(path)
	if er != nil {
		return nil, er
	}
	var policy ResourceProtectPolicy
	err := yaml.Unmarshal(data, &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func resourcePolicyViolationDetected(resourceTypeToNames ResultResourceTypeToNames, resourcePolicies *ResourceProtectPolicies) {
	var deleteProtectRules = resourcePolicies.DeleteProtectPolicy
	var replaceProtectRules = resourcePolicies.ReplaceProtectPolicy
	fmt.Println("\nPolicy Violation\n----------------")
	for _, deleteProtectRule := range deleteProtectRules.ProtectRules {
		deleteResourceNames := resourceTypeToNames.Delete[deleteProtectRule.Resource]
		if deleteResourceNames != nil {
			for _, name := range deleteResourceNames {
				outputPolicyViolation(deleteProtectRule.Severity, fmt.Sprintf("%s %s", "-", deleteProtectRule.Resource), name)
			}
		}
	}

	for _, replaceProtectRule := range replaceProtectRules.ProtectRules {
		replaceResourceNames := resourceTypeToNames.Replace[replaceProtectRule.Resource]
		if replaceResourceNames != nil {
			for _, name := range replaceResourceNames {
				outputPolicyViolation(replaceProtectRule.Severity, fmt.Sprintf("%s %s", "+/-", replaceProtectRule.Resource), name)
			}
		}
	}
}

func outputPolicyViolation(severity string, resourceType, resourceName string) {
	switch strings.ToLower(severity) {
	case "critical":
		fmt.Printf("🚨 %s.%s \n", resourceType, resourceName)
		break
	case "warning", "warn":
		fmt.Printf("⚠️ %s.%s \n", resourceType, resourceName)
		break
	}
}
