package schemas

import (
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var displaySettingsControlTypeNames = struct {
	CheckBox       string
	MultiLineText  string
	Select         string
	SingleLineText string
}{
	"Checkbox",
	"MultiLineText",
	"Select",
	"SingleLineText",
}

var displaySettingsControlTypes = []string{
	displaySettingsControlTypeNames.SingleLineText,
	displaySettingsControlTypeNames.MultiLineText,
	displaySettingsControlTypeNames.CheckBox,
	displaySettingsControlTypeNames.Select,
}

func VariablePromptOptionsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		SchemaAttributeNames.Description: types.StringType,
		VariableSchemaAttributeNames.DisplaySettings: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: VariableDisplaySettingsObjectType(),
			},
		},
		VariableSchemaAttributeNames.IsRequired: types.BoolType,
		VariableSchemaAttributeNames.Label:      types.StringType,
	}
}

func VariableDisplaySettingsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		VariableSchemaAttributeNames.ControlType: types.StringType,
		VariableSchemaAttributeNames.SelectOption: types.ListType{
			ElemType: types.ObjectType{AttrTypes: VariableSelectOptionsObjectType()},
		},
	}
}

func VariableSelectOptionsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		VariableSchemaAttributeNames.Value:       types.StringType,
		VariableSchemaAttributeNames.DisplayName: types.StringType,
	}
}

func MapFromVariablePromptOptions(variablePromptOptions *variables.VariablePromptOptions) attr.Value {
	if variablePromptOptions == nil {
		return nil
	}

	attrs := map[string]attr.Value{
		SchemaAttributeNames.Description:             types.StringValue(variablePromptOptions.Description),
		VariableSchemaAttributeNames.IsRequired:      types.BoolValue(variablePromptOptions.IsRequired),
		VariableSchemaAttributeNames.Label:           types.StringValue(variablePromptOptions.Label),
		VariableSchemaAttributeNames.DisplaySettings: types.ListNull(types.ObjectType{AttrTypes: VariableDisplaySettingsObjectType()}),
	}
	if variablePromptOptions.DisplaySettings != nil {
		attrs[VariableSchemaAttributeNames.DisplaySettings] = types.ListValueMust(
			types.ObjectType{
				AttrTypes: VariableDisplaySettingsObjectType(),
			},
			[]attr.Value{
				MapFromDisplaySettings(variablePromptOptions.DisplaySettings),
			},
		)
	}

	return types.ObjectValueMust(VariablePromptOptionsObjectType(), attrs)
}

func MapFromDisplaySettings(displaySettings *resources.DisplaySettings) attr.Value {
	if displaySettings == nil {
		return nil
	}

	attrs := map[string]attr.Value{
		VariableSchemaAttributeNames.ControlType: types.StringValue(string(displaySettings.ControlType)),
	}
	if displaySettings.ControlType == resources.ControlTypeSelect {
		if len(displaySettings.SelectOptions) > 0 {
			attrs[VariableSchemaAttributeNames.SelectOption] = types.ListValueMust(
				types.ObjectType{AttrTypes: VariableSelectOptionsObjectType()},
				MapFromSelectOptions(displaySettings.SelectOptions),
			)
		}
	} else {
		attrs[VariableSchemaAttributeNames.SelectOption] = types.ListNull(types.ObjectType{AttrTypes: VariableSelectOptionsObjectType()})
	}

	return types.ObjectValueMust(
		VariableDisplaySettingsObjectType(),
		attrs,
	)
}

func MapFromSelectOptions(selectOptions []*resources.SelectOption) []attr.Value {
	options := make([]attr.Value, 0, len(selectOptions))
	for _, option := range selectOptions {
		options = append(options, types.ObjectValueMust(
			VariableSelectOptionsObjectType(),
			map[string]attr.Value{
				VariableSchemaAttributeNames.Value:       types.StringValue(option.Value),
				VariableSchemaAttributeNames.DisplayName: types.StringValue(option.DisplayName),
			},
		))
	}
	return options
}

func MapToVariablePromptOptions(flattenedVariablePromptOptions types.List) *variables.VariablePromptOptions {
	if flattenedVariablePromptOptions.IsNull() || len(flattenedVariablePromptOptions.Elements()) == 0 {
		return nil
	}

	obj := flattenedVariablePromptOptions.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	var promptOptions variables.VariablePromptOptions
	if description, ok := attrs[SchemaAttributeNames.Description].(types.String); ok && !description.IsNull() {
		promptOptions.Description = description.ValueString()
	}

	if isRequired, ok := attrs[VariableSchemaAttributeNames.IsRequired].(types.Bool); ok && !isRequired.IsNull() {
		promptOptions.IsRequired = isRequired.ValueBool()
	}

	if label, ok := attrs[VariableSchemaAttributeNames.Label].(types.String); ok && !label.IsNull() {
		promptOptions.Label = label.ValueString()
	}

	if displaySettings, ok := attrs[VariableSchemaAttributeNames.DisplaySettings].(types.List); ok && !displaySettings.IsNull() {
		promptOptions.DisplaySettings = MapToDisplaySettings(displaySettings)
	}

	return &promptOptions
}

func MapToDisplaySettings(displaySettings types.List) *resources.DisplaySettings {
	if displaySettings.IsNull() || len(displaySettings.Elements()) == 0 {
		return nil
	}

	obj := displaySettings.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	ct, _ := attrs[VariableSchemaAttributeNames.ControlType].(types.String)
	controlType := resources.ControlType(ct.ValueString())

	var selectOptions []*resources.SelectOption
	if controlType == resources.ControlTypeSelect {
		selectOptions = MapToSelectOptions(attrs[VariableSchemaAttributeNames.SelectOption].(types.List))
	}

	return resources.NewDisplaySettings(controlType, selectOptions)
}

func MapToSelectOptions(selectOptions types.List) []*resources.SelectOption {
	if selectOptions.IsNull() || selectOptions.IsUnknown() || len(selectOptions.Elements()) == 0 {
		return nil
	}

	options := make([]*resources.SelectOption, 0, len(selectOptions.Elements()))
	for _, option := range selectOptions.Elements() {
		attrs := option.(types.Object).Attributes()
		options = append(options, &resources.SelectOption{
			DisplayName: attrs[VariableSchemaAttributeNames.DisplayName].(types.String).ValueString(),
			Value:       attrs[VariableSchemaAttributeNames.Value].(types.String).ValueString(),
		})
	}

	return options
}

func getVariablePromptDatasourceSchema() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				SchemaAttributeNames.Description:             util.GetDescriptionDatasourceSchema("variable prompt option"),
				VariableSchemaAttributeNames.DisplaySettings: getDisplaySettingsDatasourceSchema(),
				VariableSchemaAttributeNames.IsRequired: datasourceSchema.BoolAttribute{
					Computed: true,
				},
				VariableSchemaAttributeNames.Label: datasourceSchema.StringAttribute{
					Computed: true,
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func getDisplaySettingsDatasourceSchema() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				VariableSchemaAttributeNames.ControlType: datasourceSchema.StringAttribute{
					Description: fmt.Sprintf("The type of control for rendering this prompted variable. Valid types are %s.", strings.Join(util.Map(displaySettingsControlTypes, func(item string) string { return fmt.Sprintf("`%s`", item) }), ", ")),
					Computed:    true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							displaySettingsControlTypes...,
						),
					},
				},
				VariableSchemaAttributeNames.SelectOption: datasourceSchema.ListNestedAttribute{
					Description: fmt.Sprintf("If the `%s` is `%s`, then this value defines an option.", VariableSchemaAttributeNames.ControlType, displaySettingsControlTypeNames.Select),
					Computed:    true,
					NestedObject: datasourceSchema.NestedAttributeObject{
						Attributes: map[string]datasourceSchema.Attribute{
							VariableSchemaAttributeNames.Value: datasourceSchema.StringAttribute{
								Description: "The select value",
								Computed:    true,
							},
							VariableSchemaAttributeNames.DisplayName: datasourceSchema.StringAttribute{
								Description: "The display name for the select value",
								Computed:    true,
							},
						},
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func getVariablePromptResourceSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				SchemaAttributeNames.Description: util.GetDescriptionResourceSchema("variable prompt option"),
				VariableSchemaAttributeNames.IsRequired: resourceSchema.BoolAttribute{
					Optional: true,
				},
				VariableSchemaAttributeNames.Label: resourceSchema.StringAttribute{
					Optional: true,
				},
			},
			Blocks: map[string]resourceSchema.Block{
				VariableSchemaAttributeNames.DisplaySettings: getDisplaySettingsResourceSchema(),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func getDisplaySettingsResourceSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				VariableSchemaAttributeNames.ControlType: resourceSchema.StringAttribute{
					Description: fmt.Sprintf("The type of control for rendering this prompted variable. Valid types are %s.", strings.Join(util.Map(displaySettingsControlTypes, func(item string) string { return fmt.Sprintf("`%s`", item) }), ", ")),
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							displaySettingsControlTypes...,
						),
					},
				},
			},
			Blocks: map[string]resourceSchema.Block{
				VariableSchemaAttributeNames.SelectOption: resourceSchema.ListNestedBlock{
					Description: fmt.Sprintf("If the `%s` is `%s`, then this value defines an option.", VariableSchemaAttributeNames.ControlType, displaySettingsControlTypeNames.Select),
					NestedObject: resourceSchema.NestedBlockObject{
						Attributes: map[string]resourceSchema.Attribute{
							VariableSchemaAttributeNames.Value: resourceSchema.StringAttribute{
								Description: "The select value",
								Required:    true,
							},
							VariableSchemaAttributeNames.DisplayName: resourceSchema.StringAttribute{
								Description: "The display name for the select value",
								Required:    true,
							},
						},
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}
