---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "octopusdeploy_project_group Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  
---

# octopusdeploy_project_group (Resource)



## Example Usage

```terraform
resource "octopusdeploy_project_group" "example" {
  description  = "The development project group."
  name         = "Development Project Group (OK to Delete)"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of this resource.

### Optional

- `description` (String) The description of this project group.
- `id` (String) The unique ID for this resource.
- `space_id` (String) The space ID associated with this project group.

## Import

Import is supported using the following syntax:

```shell
terraform import [options] octopusdeploy_project_group.<name> <project_group-id>
```
