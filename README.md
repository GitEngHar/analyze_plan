
<p align="center">
  <img src="analyze_plan.png" width="600">
</p>

<p align="center">
    <em>Terraform plan more comfortably and safety with GitHub Actions</em>
</p>

<hr>

<!-- explain><!-->

<!-- invite to oss><!-->

## Sample Yaml
use analyze_plan
```yaml
GitEngHar/analyze_plan@v1.2.1
```

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

## Features
<!-- features desc><!-->
- **Plan summary**
- **Replace detected**
- **Policy alert**



## 🛠 Usage


## 📜 License
This project is distributed under the **MIT** license.



