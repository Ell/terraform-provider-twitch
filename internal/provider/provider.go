package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
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
	Endpoint types.String `tfsdk:"endpoint"`
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
			},
			"app_token": schema.StringAttribute{
				MarkdownDescription: "Twitch app token",
				Required:            true,
			},
		},
	}
}

func (p *TwitchProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TwitchProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	// client := http.DefaultClient
	// resp.DataSourceData = client
	// resp.ResourceData = client
}

func (p *TwitchProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *TwitchProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
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
