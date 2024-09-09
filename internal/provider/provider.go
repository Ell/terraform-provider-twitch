package provider

import (
	"context"

	"github.com/ell/terraform-provider-twitch/internal/helix"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure TwitchProvider satisfies various provider interfaces.
var _ provider.Provider = &TwitchProvider{}
var _ provider.ProviderWithFunctions = &TwitchProvider{}

// TwitchProvider defines the provider implementation.
type TwitchProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TwitchProviderModel describes the provider data model.
type TwitchProviderModel struct {
	ClientId    types.String `tfsdk:"client_id"`
	AccessToken types.String `tfsdk:"access_token"`
}

func (p *TwitchProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "twitch"
	resp.Version = p.version
}

func (p *TwitchProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Twitch client id",
				Required:            true,
				Sensitive:           true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Twitch access token",
				Sensitive:           true,
				Required:            true,
			},
		},
	}
}

func (p *TwitchProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config TwitchProviderModel

	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ClientId.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("client_id"), "Unknown client_id", "client_id is required")
	}

	if config.AccessToken.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("access_token"), "Unknown access_token", "access_token is required")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	twitchClient := helix.NewHelixClient(config.ClientId.ValueString(), config.AccessToken.ValueString())

	resp.DataSourceData = twitchClient
	resp.ResourceData = twitchClient
}

func (p *TwitchProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewChannelResource,
		NewChannelRewardResource,
	}
}

func (p *TwitchProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewGameDataSource,
	}
}

func (p *TwitchProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TwitchProvider{
			version: version,
		}
	}
}
