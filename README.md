<p align="center" style="font-size: 1.5em;">
    <em>Terraform provider for Mikrotik SwOS switches</em>
</p><br>

[![CodeQL](https://github.com/finomen/terraform-provider-swos/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/finomen/terraform-provider-swos/actions/workflows/github-code-scanning/codeql)
[![Release](https://github.com/finomen/terraform-provider-swos/actions/workflows/release.yml/badge.svg)](https://github.com/finomen/terraform-provider-swos/actions/workflows/release.yml)
---

This is a hobby project to create a Terraform provider for [SwOS](https://help.mikrotik.com/docs/spaces/SWOS/pages/328415/SwOS) switches that are web-managed. It relies on [goquery](https://github.com/finomen/swos-client).

Check the documentation at:

- Terraform: [HRUI Provider](https://registry.terraform.io/providers/finomen/swos)
- OpenToFu: [HRUI Provider](https://search.opentofu.org/provider/finomen/swos)


## Getting Started

1.  Configure the provider in your Terraform configuration:

    ```terraform
    terraform {
      required_providers {
        swos = {
          source  = "finomen/swos"
        }
      }
    }

    provider "swos" {
      url      = "http://192.168.2.1"
      username = "admin"
      password = "XXX"
    }
    ```

## Contributing

Contributions are welcome!

## License

This provider is licensed under the MIT License.