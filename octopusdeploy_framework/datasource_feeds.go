package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
)

type feedsDataSource struct {
	*Config
}

func NewFeedsDataSource() datasource.DataSource {
	return &feedsDataSource{}
}

func (*feedsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("feeds")
}

func (e *feedsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	e.Config = DataSourceConfiguration(req, resp)
}

func (*feedsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema.Schema{
		Description: "Provides information about existing feeds.",
		Attributes:  schemas.GetFeedsDataSourceSchema(),
	}
}

func (e *feedsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data schemas.FeedsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := feeds.FeedsQuery{
		Name:        data.Name.ValueString(),
		IDs:         util.GetIds(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        util.GetNumber(data.Skip),
		Take:        util.GetNumber(data.Take),
	}

	util.DatasourceReading(ctx, "feeds", query)

	existingFeeds, err := feeds.Get(e.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load feeds", err.Error())
		return
	}

	util.DatasourceResultCount(ctx, "feeds", len(existingFeeds.Items))

	flattenedFeeds := []interface{}{}
	for _, feed := range existingFeeds.Items {
		feedResource, err := feeds.ToFeedResource(feed)
		if err != nil {
			resp.Diagnostics.AddError("Unable to map to feeds: %s", err.Error())
			return
		}

		flattenedFeeds = append(flattenedFeeds, schemas.FlattenFeed(feedResource))
	}

	data.Feeds, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.FeedObjectType()}, flattenedFeeds)
	data.ID = types.StringValue("Feeds " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
