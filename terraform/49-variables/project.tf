data "octopusdeploy_lifecycles" "lifecycle_default_lifecycle" {
  ids          = null
  partial_name = "Default Lifecycle"
  skip         = 0
  take         = 1
}

resource "octopusdeploy_project_group" "project_group_test" {
  name        = "Test"
  description = "Test Description"
}



resource "octopusdeploy_project" "test_project" {
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Test project"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
  name                                 = "Test"
  project_group_id                     = octopusdeploy_project_group.project_group_test.id
  tenanted_deployment_participation    = "Untenanted"
  space_id                             = var.octopus_space_id
  included_library_variable_sets       = []
  versioning_strategy {
    template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.LastPatch}.#{Octopus.Version.NextRevision}"
  }

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
  template {
    name             = "Project Template Variable"
    label            = "Test"
    default_value    = "Test"
    display_settings = { "Octopus.ControlType" = "SingleLineText" }
  }
}

resource "octopusdeploy_deployment_process" "test_deployment_process" {
  project_id = octopusdeploy_project.test_project.id
  step {
    condition           = "Success"
    name                = "Hello world (using PowerShell)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using PowerShell)"
      script_body                        = <<-EOT
          Write-Host 'Hello world, using PowerShell'
          #TODO: Experiment with steps of your own :)
          Write-Host '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
        EOT
      run_on_server                      = true
    }
  }
}

resource "octopusdeploy_channel" "test_channel" {
  depends_on  = [octopusdeploy_deployment_process.test_deployment_process]
  description = "Test Channel"
  name = "Test Channel"
  project_id = octopusdeploy_project.test_project.id

  rule {
    version_range = "1.0"
  }
}