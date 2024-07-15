package octopusdeploy_framework

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExpandLifecycleWithNil(t *testing.T) {
	lifecycle := expandLifecycle(nil)
	require.Nil(t, lifecycle)
}

func TestExpandLifecycle(t *testing.T) {
	description := "test-description"
	name := "test-name"
	spaceID := "test-space-id"
	Id := "test-id"
	releaseRetention := core.NewRetentionPeriod(0, "Days", true)
	tentacleRetention := core.NewRetentionPeriod(2, "Items", false)

	data := &lifecycleTypeResourceModel{
		ID:          types.StringValue(Id),
		Description: types.StringValue(description),
		Name:        types.StringValue(name),
		SpaceID:     types.StringValue(spaceID),
		ReleaseRetentionPolicy: types.ListValueMust(
			types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
			[]attr.Value{
				types.ObjectValueMust(
					getRetentionPeriodAttrTypes(),
					map[string]attr.Value{
						"quantity_to_keep":    types.Int64Value(int64(releaseRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(releaseRetention.ShouldKeepForever),
						"unit":                types.StringValue(releaseRetention.Unit),
					},
				),
			},
		),
		TentacleRetentionPolicy: types.ListValueMust(
			types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
			[]attr.Value{
				types.ObjectValueMust(
					getRetentionPeriodAttrTypes(),
					map[string]attr.Value{
						"quantity_to_keep":    types.Int64Value(int64(tentacleRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(tentacleRetention.ShouldKeepForever),
						"unit":                types.StringValue(tentacleRetention.Unit),
					},
				),
			},
		),
	}

	lifecycle := expandLifecycle(data)

	require.Equal(t, description, lifecycle.Description)
	require.NotEmpty(t, lifecycle.ID)
	require.NotNil(t, lifecycle.Links)
	require.Empty(t, lifecycle.Links)
	require.Equal(t, name, lifecycle.Name)
	require.Empty(t, lifecycle.Phases)
	require.Equal(t, releaseRetention, lifecycle.ReleaseRetentionPolicy)
	require.Equal(t, tentacleRetention, lifecycle.TentacleRetentionPolicy)
	require.Equal(t, spaceID, lifecycle.SpaceID)
}

func TestExpandPhasesWithEmptyInput(t *testing.T) {
	emptyList := types.ListValueMust(types.ObjectType{AttrTypes: getPhaseAttrTypes()}, []attr.Value{})
	phases := expandPhases(emptyList)
	require.Nil(t, phases)
}

func TestExpandPhasesWithNullInput(t *testing.T) {
	nullList := types.ListNull(types.ObjectType{AttrTypes: getPhaseAttrTypes()})
	phases := expandPhases(nullList)
	require.Nil(t, phases)
}

func TestExpandPhasesWithUnknownInput(t *testing.T) {
	unknownList := types.ListUnknown(types.ObjectType{AttrTypes: getPhaseAttrTypes()})
	phases := expandPhases(unknownList)
	require.Nil(t, phases)
}

func TestExpandAndFlattenPhasesWithSensibleDefaults(t *testing.T) {
	phase := createTestPhase("TestPhase", []string{"AutoTarget1", "AutoTarget2"}, true, 5)

	flattenedPhases := flattenPhases([]*lifecycles.Phase{phase})
	require.NotNil(t, flattenedPhases)
	require.Equal(t, 1, len(flattenedPhases.Elements()))

	expandedPhases := expandPhases(flattenedPhases)
	require.NotNil(t, expandedPhases)
	require.Len(t, expandedPhases, 1)

	expandedPhase := expandedPhases[0]
	require.NotEmpty(t, expandedPhase.ID)
	require.Equal(t, phase.AutomaticDeploymentTargets, expandedPhase.AutomaticDeploymentTargets)
	require.Equal(t, phase.IsOptionalPhase, expandedPhase.IsOptionalPhase)
	require.EqualValues(t, phase.MinimumEnvironmentsBeforePromotion, expandedPhase.MinimumEnvironmentsBeforePromotion)
	require.Equal(t, phase.Name, expandedPhase.Name)
	require.Equal(t, phase.ReleaseRetentionPolicy, expandedPhase.ReleaseRetentionPolicy)
	require.Equal(t, phase.TentacleRetentionPolicy, expandedPhase.TentacleRetentionPolicy)
}

func TestExpandAndFlattenMultiplePhasesWithSensibleDefaults(t *testing.T) {
	phase1 := createTestPhase("Phase1", []string{"AutoTarget1", "AutoTarget2"}, true, 5)
	phase2 := createTestPhase("Phase2", []string{"AutoTarget3", "AutoTarget4"}, false, 3)

	flattenedPhases := flattenPhases([]*lifecycles.Phase{phase1, phase2})
	require.NotNil(t, flattenedPhases)
	require.Equal(t, 2, len(flattenedPhases.Elements()))

	expandedPhases := expandPhases(flattenedPhases)
	require.NotNil(t, expandedPhases)
	require.Len(t, expandedPhases, 2)

	require.NotEmpty(t, expandedPhases[0].ID)
	require.Equal(t, phase1.AutomaticDeploymentTargets, expandedPhases[0].AutomaticDeploymentTargets)
	require.Equal(t, phase1.IsOptionalPhase, expandedPhases[0].IsOptionalPhase)
	require.EqualValues(t, phase1.MinimumEnvironmentsBeforePromotion, expandedPhases[0].MinimumEnvironmentsBeforePromotion)
	require.Equal(t, phase1.Name, expandedPhases[0].Name)
	require.Equal(t, phase1.ReleaseRetentionPolicy, expandedPhases[0].ReleaseRetentionPolicy)
	require.Equal(t, phase1.TentacleRetentionPolicy, expandedPhases[0].TentacleRetentionPolicy)

	require.NotEmpty(t, expandedPhases[1].ID)
	require.Equal(t, phase2.AutomaticDeploymentTargets, expandedPhases[1].AutomaticDeploymentTargets)
	require.Equal(t, phase2.IsOptionalPhase, expandedPhases[1].IsOptionalPhase)
	require.EqualValues(t, phase2.MinimumEnvironmentsBeforePromotion, expandedPhases[1].MinimumEnvironmentsBeforePromotion)
	require.Equal(t, phase2.Name, expandedPhases[1].Name)
	require.Equal(t, phase2.ReleaseRetentionPolicy, expandedPhases[1].ReleaseRetentionPolicy)
	require.Equal(t, phase2.TentacleRetentionPolicy, expandedPhases[1].TentacleRetentionPolicy)
}

func createTestPhase(name string, autoTargets []string, isOptional bool, minEnvs int32) *lifecycles.Phase {
	phase := lifecycles.NewPhase(name)
	phase.AutomaticDeploymentTargets = autoTargets
	phase.IsOptionalPhase = isOptional
	phase.MinimumEnvironmentsBeforePromotion = minEnvs
	phase.ReleaseRetentionPolicy = core.NewRetentionPeriod(15, "Items", false)
	phase.TentacleRetentionPolicy = core.NewRetentionPeriod(0, "Days", true)
	phase.ID = name + "-Id"
	return phase
}