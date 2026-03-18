
<p align="center">
  <img src="analyze_plan.png" width="600">
</p>

<p align="center">
    <em>Terraform plan more comfortably and safety with GitHub Actions</em>
</p>

```bash
start diff & replace detect & protect policy
planFilePath: ./test_data/plan.json 
,policyFilePath: test_policy
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

<hr>

## 🐝 Debug 
go run start

1. git clone
```bash
git clone git@github.com:GitEngHar/analyze_plan.git 
```

2. run debug
```bash
go run ./analyze_plan.go ./test_data/plan.json test_policy
```

## 🚀 Quick Start
use analyze_plan with GithubActions
```yaml
  - name: summarize terraform plan
    uses: GitEngHar/analyze_plan@v1.2.3
    with:
      plan-path: plan.json
```

sample yaml
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
      - name: summarize terraform plan
        uses: GitEngHar/analyze_plan@v1.2.3
        with:
          plan-path: plan.json
```

## ⚙️ Features
- **📃 Plan summary**: Summarizes Terraform plan changes (create, update, delete, replace)
- **🔍 Replace detected**: Detects resource replacements
- **🛡 Policy alert**: Alerts on risky changes based on predefined policies


## 📜 License
This project is distributed under the **MIT** license.



