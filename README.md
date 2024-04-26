# tfreveal

Terraform does mask sensitive values in the output (e.g. from `terraform plan`)
in order to protect them from being revealed to 3rd parties.

Sometimes it is neccessary to see the exact changes, Terraform will perform to the
infrastructure including all the changes to sensitive values. So far, Terraform
does not provide a feature to forcefully unmask the sensitive values in the
[concise diff plan outputs](https://www.hashicorp.com/blog/terraform-0-14-adds-a-new-concise-diff-format-to-terraform-plans).
The general advice given by the Terraform maintainers is to use the JSON output
in such cases. While the JSON output does provide all the necessary information,
it is not perticularely easy to read for humans and to spot small differences.
It gets even more complicated, if the sensitive values contain larger JSON
encoded values.

There exists instructions using for example `jq`, but the process stays manual,
cumbersome and error prone.

`tfreveal` is here to fix this and provide an easy way to show the concise diff
plan outputs with all sensitive values revealed.

## Usage

The plan file generated from Terraform can be directly piped to `tfreveal`:

```bash
$ terraform plan -out plan.out
$ terraform show -json plan.out | tfreveal
```

Alternatively, the plan file can also be passed as argument:

```bash
$ terraform plan -out plan.out
$ terraform show -json plan.out > plan.json
$ tfreveal plan.json
```

## Trademarks

All other trademarks referenced herein are the property of their respective owners.
