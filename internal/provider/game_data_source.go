package provider

import (
	"context"

	"github.com/ell/terraform-provider-twitch/internal/helix"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &gameDataSource{}
	_ datasource.DataSourceWithConfigure = &gameDataSource{}
)

type gameDataSource struct {
	client *helix.Client
}

type gameDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	BoxArtURL types.String `tfsdk:"box_art_url"`
	IGDBID    types.String `tfsdk:"igdb_id"`
}

func (g *gameDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*helix.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected DataSource Configure Type", "Expected *helix.Client")

		return
	}

	g.client = client
}

func (g *gameDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_game"
}

func (g *gameDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state gameDataSourceModel

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	games, err := g.client.GetGameByName(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get game", err.Error())

		return
	}

	if len(games.Data) == 0 {
		resp.Diagnostics.AddError("No game found", "No game found with the provided name")

		return
	}

	game := games.Data[0]

	resp.State.Set(ctx, gameDataSourceModel{
		ID:        types.StringValue(game.ID),
		Name:      types.StringValue(game.Name),
		BoxArtURL: types.StringValue(game.BoxArtURL),
		IGDBID:    types.StringValue(game.IGDBID),
	})
}

func (g *gameDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"box_art_url": schema.StringAttribute{
				Computed: true,
			},
			"igdb_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func NewGameDataSource() datasource.DataSource {
	return &gameDataSource{}
}
