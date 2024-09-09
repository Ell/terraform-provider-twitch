package provider

import (
	"context"
	"fmt"

	"github.com/ell/terraform-provider-twitch/internal/helix"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &channelResource{}
	_ resource.ResourceWithConfigure   = &channelResource{}
	_ resource.ResourceWithImportState = &channelResource{}
)

type channelResource struct {
	TwitchClient *helix.Client
}

func (c *channelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	twitchClient, ok := req.ProviderData.(*helix.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *helix.Client, got %T", req.ProviderData))
	}

	c.TwitchClient = twitchClient
}

func NewChannelResource() resource.Resource {
	return &channelResource{}
}

// Schema implements resource.Resource.
func (c *channelResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"title": schema.StringAttribute{
				Required: true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"game_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

type channelResourceModel struct {
	ID     types.String   `tfsdk:"id"`
	Title  types.String   `tfsdk:"title"`
	Tags   []types.String `tfsdk:"tags"`
	GameID types.String   `tfsdk:"game_id"`
}

func (c *channelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

// Unimplemented
func (c *channelResource) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {}

func (c *channelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state channelResourceModel

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get twitch channel information
	channelInfos, err := c.TwitchClient.GetChannel(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get channel information", err.Error())
		return
	}

	if len(channelInfos.Data) == 0 {
		resp.Diagnostics.AddError("Channel not found", "Channel not found")
		return
	}

	channelInfo := channelInfos.Data[0]

	state.Title = types.StringValue(channelInfo.Title)
	state.GameID = types.StringValue(channelInfo.GameId)

	state.Tags = []types.String{}
	for _, tag := range channelInfo.Tags {
		state.Tags = append(state.Tags, types.StringValue(tag))
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *channelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state channelResourceModel

	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateRequest helix.UpdateChannelRequest

	if !state.GameID.IsNull() {
		updateRequest.GameID = state.GameID.ValueString()
	}

	if !state.Title.IsNull() {
		updateRequest.Title = state.Title.ValueString()
	}

	updateRequest.Tags = []string{}
	if state.Tags != nil {
		updateRequest.Tags = []string{}
		for _, tag := range state.Tags {
			updateRequest.Tags = append(updateRequest.Tags, tag.ValueString())
		}
	}

	err := c.TwitchClient.UpdateChannel(state.ID.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update channel", err.Error())
		return
	}

	channelInfos, err := c.TwitchClient.GetChannel(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get channel information", err.Error())
		return
	}

	if len(channelInfos.Data) == 0 {
		resp.Diagnostics.AddError("Channel not found", "Channel not found")
		return
	}

	channelInfo := channelInfos.Data[0]

	state.Title = types.StringValue(channelInfo.Title)
	state.GameID = types.StringValue(channelInfo.GameId)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Unimplemented
func (c *channelResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {}

func (c *channelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
