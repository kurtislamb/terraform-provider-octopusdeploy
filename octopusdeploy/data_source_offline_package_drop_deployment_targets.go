package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOfflinePackageDropDeploymentTargets() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing offline package drop deployment targets.",
		ReadContext: dataSourceOfflinePackageDropDeploymentTargetsRead,
		Schema:      getOfflinePackageDropDeploymentTargetDataSchema(),
	}
}

func dataSourceOfflinePackageDropDeploymentTargetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := machines.MachinesQuery{
		CommunicationStyles: []string{"OfflineDrop"},
		DeploymentID:        d.Get("deployment_id").(string),
		EnvironmentIDs:      expandArray(d.Get("environments").([]interface{})),
		HealthStatuses:      expandArray(d.Get("health_statuses").([]interface{})),
		IDs:                 expandArray(d.Get("ids").([]interface{})),
		IsDisabled:          d.Get("is_disabled").(bool),
		Name:                d.Get("name").(string),
		PartialName:         d.Get("partial_name").(string),
		Roles:               expandArray(d.Get("roles").([]interface{})),
		ShellNames:          expandArray(d.Get("shell_names").([]interface{})),
		Skip:                d.Get("skip").(int),
		Take:                d.Get("take").(int),
		TenantIDs:           expandArray(d.Get("tenants").([]interface{})),
		TenantTags:          expandArray(d.Get("tenant_tags").([]interface{})),
		Thumbprint:          d.Get("thumbprint").(string),
	}

	client := m.(*client.Client)
	existingDeploymentTargets, err := machines.Get(client, d.Get("space_id").(string), query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedOfflinePackageDropDeploymentTargets := []interface{}{}
	for _, deploymentTarget := range existingDeploymentTargets.Items {
		flattenedOfflinePackageDropDeploymentTargets = append(flattenedOfflinePackageDropDeploymentTargets, flattenOfflinePackageDropDeploymentTarget(deploymentTarget))
	}

	d.Set("offline_package_drop_deployment_targets", flattenedOfflinePackageDropDeploymentTargets)
	d.SetId("OfflinePackageDropDeploymentTargets " + time.Now().UTC().String())

	return nil
}
