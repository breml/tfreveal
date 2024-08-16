# tfreveal

[![Test Status](https://github.com/breml/tfreveal/workflows/Main/badge.svg)](https://github.com/breml/tfreveal/actions?query=workflow%3AMain)
 [![Go Report Card](https://goreportcard.com/badge/github.com/breml/tfreveal)](https://goreportcard.com/report/github.com/breml/tfreveal) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

tfreveal is an open-source tool designed to enhance the visibility of Terraform
plan files by displaying all differences in resources and outputs, including
sensitive values. Unlike Terraform, which hides sensitive data, tfreveal reveals
these values to ensure complete transparency in your infrastructure changes.

[![asciicast](https://asciinema.org/a/672302.svg)](https://asciinema.org/a/672302)

## Motivation

Terraform does mask sensitive values in the output (e.g. from `terraform plan`)
in order to protect them from being revealed to unauthorized 3rd parties.

But sometimes it is neccessary to see the exact changes, Terraform will perform
to the infrastructure including all the changes to sensitive values. In
particular, if one observes drift between the Terraform state and the actual
state of the infrastructure, this becomes inevitable. So far, Terraform does not
provide a feature to forcefully unmask the sensitive values in the
[concise diff plan outputs](https://www.hashicorp.com/blog/terraform-0-14-adds-a-new-concise-diff-format-to-terraform-plans).

The general advice given by the Terraform maintainers is to use the JSON output
in such cases. While the JSON output does provide all the necessary information,
it is not perticularely easy to read for humans and to spot small differences.
It gets even more complicated, if the changes are contained in larger JSON
encoded values, that are marked as sensitive.

There exists instructions using for example `jq`, but the process stays manual,
cumbersome and error prone.

`tfreveal` is here to fix this and provide an easy way to show the concise diff
plan outputs with all sensitive values revealed.

## Installation

Download the latest release from the [releases page](https://github.com/breml/tfreveal/releases).

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

## Development

The task to update the test data and the golden files is provided in the
`Taskfile.yml` and can be executed by running `task gen-all`. This requires the
`task` tool to be installed. Please refer to the
[official documentation](https://taskfile.dev/installation/).

Additionally the `terraform` command needs to be present in the `PATH`. Follow
the [official installation instructions](https://developer.hashicorp.com/terraform/install).

## Author

Copyright 2024 by Lucas Bremgartner ([breml](https://github.com/breml))

## License

[MIT License](LICENSE)

## Trademarks

All other trademarks referenced herein are the property of their respective owners.
