# analyze_plan

<p align="center">
  <img src="analyze_plan.png" width="600">
</p>

<p align="center">
  <b>🚨 Detect risky Terraform changes before they break production</b><br>
  <em>Analyze Terraform plans with risk detection, policy guard, and GitHub Actions integration</em>
</p>

---

## ❗ Why

Terraform plans are hard to review.

- Did you miss a destructive change?
- Will this replacement cause downtime?
- Is it really safe to merge?

👉 One mistake can cause serious production issues.

**analyze_plan helps you detect risks instantly.**

---

## 🆚 Before vs After

### ❌ Terraform plan
- Hard to read
- Risk is unclear
- Easy to miss critical changes

### ✨ Example Output
```bash
module.alb.aws_lb.ecs_alb -> create
module.ecs.aws_ecs_service.app_service -> update
module.db.aws_ecs_service.db_service -> delete
module.ecs.aws_ecs_task_definition.app_task_def -> replace

Terraform Plan Summary
----------------------
create: 1
update: 1
delete: 1
replace: 1
+  module.alb.aws_lb.ecs_alb
~  module.ecs.aws_ecs_service.app_service
-  module.db.aws_ecs_service.db_service
+/-  module.ecs.aws_ecs_task_definition.app_task_def

Replace Detected
----------------
+/-  module.ecs.aws_ecs_task_definition.app_task_def

Policy Violation
----------------
🚨 - aws_ecs_service.db_service 
⚠️ +/- aws_ecs_task_definition.app_task_def 
```

### 🚀Quick Start
```bash
$ git clone git@github.com:GitEngHar/analyze_plan.git
$ cd analyze_plan

$ go run ./analyze_plan.go ./test_data/plan.json test_policy
```

### ⚡ GitHub Actions
Automatically analyze Terraform plans in PRs and detect risky changes before merge.
```bash
- name: summarize terraform plan
  uses: GitEngHar/analyze_plan@v1.2.3
  with:
    plan-path: plan.json
```

#### 📦 Full Example
```yaml
name: terraform

on:
  workflow_dispatch:
    inputs:
      command:
        description: "Terraform command"
        required: true
        type: choice
        options:
          - plan

permissions:
  id-token: write
  contents: read

jobs:
  terraform:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: arn:aws:iam::YOUR_AWS_ACCOUNT_ID:role/YOUR_ROLE_NAME
          aws-region: ap-northeast-1

      - name: Terraform Plan
        if: github.event.inputs.command == 'plan'
        run: | 
          terraform init
          terraform plan -out tfplan
          terraform show -json tfplan > plan.json

      - name: analyze terraform plan
        uses: GitEngHar/analyze_plan@v1.2.3
        with:
          plan-path: plan.json
```

### ⚙️ Features
- 📃 Plan summary 
  - Summarizes create / update / delete / replace
- 🔍 Replace detection 
  - Detects resources that will be replaced
- 🛡 Policy guard 
  - Detects risky changes (delete / replace) based on custom rules


### 🧩 Use Cases
- Prevent destructive infrastructure changes
- Enforce infrastructure policies
- Improve SRE / DevOps workflows

### ⭐ Support
If this project helps you, please ⭐️ the repo!

### 📜 License
MIT License