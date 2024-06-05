---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "octopusdeploy_polling_subscription_id Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  A unique polling subscription ID that can be used by polling tentacles.
---

# octopusdeploy_polling_subscription_id (Resource)

A unique polling subscription ID that can be used by polling tentacles.

## Example Usage

```terraform
resource "octopusdeploy_polling_subscription_id" "example" {}

resource "octopusdeploy_polling_subscription_id" "example_with_dependencies" {
  dependencies = {
    "target" = octopusdeploy_kubernetes_agent_deployment_target.example.id
  }
}

# Usage
resource "octopusdeploy_kubernetes_agent_deployment_target" "agent" {
  name         = "agent"
  environments = ["environments-1"]
  roles        = ["role-1", "role-2"]
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  uri          = octopusdeploy_polling_subscription_id.example_with_dependencies.polling_uri
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `dependencies` (Map of String) Optional map of dependencies that when modified will trigger a re-creation of this resource.

### Read-Only

- `id` (String) The generated polling subscription ID.
- `polling_uri` (String) The URI of the polling subscription ID.

