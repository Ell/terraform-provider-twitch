package provider

import (
	"context"
	"fmt"

	"github.com/ell/terraform-provider-twitch/internal/helix"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &channelRewardResource{}
	_ resource.ResourceWithConfigure   = &channelRewardResource{}
	_ resource.ResourceWithImportState = &channelRewardResource{}
)

type channelRewardResource struct {
	TwitchClient *helix.Client
}

func (c *channelRewardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	twitchClient, ok := req.ProviderData.(*helix.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *helix.Client, got %T", req.ProviderData))
	}

	c.TwitchClient = twitchClient
}

func NewChannelRewardResource() resource.Resource {
	return &channelRewardResource{}
}

func (c *channelRewardResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"broadcaster_id": schema.StringAttribute{
				Required: true,
			},
			"title": schema.StringAttribute{
				Required: true,
			},
			"prompt": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Computed: true,
			},
			"cost": schema.Int32Attribute{
				Required: true,
			},
			"background_color": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("#FFBF00"),
				Computed: true,
			},
			"is_global_cooldown_enabled": schema.BoolAttribute{
				Optional: true,
				Default:  booldefault.StaticBool(false),
				Computed: true,
			},
			"global_cooldown_seconds": schema.Int32Attribute{
				Optional: true,
				Default:  int32default.StaticInt32(0),
				Computed: true,
			},
			"is_enabled": schema.BoolAttribute{
				Required: true,
			},
		},
	}
}

type channelRewardResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	BroadcasterId           types.String `tfsdk:"broadcaster_id"`
	Title                   types.String `tfsdk:"title"`
	Prompt                  types.String `tfsdk:"prompt"`
	Cost                    types.Int32  `tfsdk:"cost"`
	BackgroundColor         types.String `tfsdk:"background_color"`
	IsGlobalCooldownEnabled types.Bool   `tfsdk:"is_global_cooldown_enabled"`
	GlobalCooldownSeconds   types.Int32  `tfsdk:"global_cooldown_seconds"`
	IsEnabled               types.Bool   `tfsdk:"is_enabled"`
}

func (c *channelRewardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel_reward"
}

func (c *channelRewardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan channelRewardResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reward, err := c.TwitchClient.CreateChannelReward(plan.BroadcasterId.ValueString(), &helix.CreateChannelRewardRequest{
		Title:                   plan.Title.ValueString(),
		Prompt:                  plan.Prompt.ValueString(),
		Cost:                    int(plan.Cost.ValueInt32()),
		BackgroundColor:         plan.BackgroundColor.ValueString(),
		IsGlobalCooldownEnabled: plan.IsGlobalCooldownEnabled.ValueBool(),
		GlobalCooldownSeconds:   int(plan.GlobalCooldownSeconds.ValueInt32()),
		IsEnabled:               plan.IsEnabled.ValueBool(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Failed to create channel reward", err.Error())
		return
	}

	plan.ID = types.StringValue(reward.ID)

	diags = resp.State.Set(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *channelRewardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state channelRewardResourceModel

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reward, err := c.TwitchClient.GetChannelRewardByID(state.BroadcasterId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get channel reward", err.Error())
		return
	}

	state.ID = types.StringValue(reward.ID)
	state.BroadcasterId = types.StringValue(reward.BroadcasterID)
	state.Title = types.StringValue(reward.Title)
	state.Prompt = types.StringValue(reward.Prompt)
	state.Cost = types.Int32Value(int32(reward.Cost))
	state.BackgroundColor = types.StringValue(reward.BackgroundColor)
	state.IsGlobalCooldownEnabled = types.BoolValue(reward.GlobalCooldownSetting.IsEnabled)
	state.GlobalCooldownSeconds = types.Int32Value(int32(reward.GlobalCooldownSetting.GlobalCooldownSeconds))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *channelRewardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan channelRewardResourceModel
	var state channelRewardResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rewardID := plan.ID.ValueString()
	if plan.ID.IsUnknown() {
		reward, err := c.TwitchClient.GetChannelRewardByName(state.BroadcasterId.ValueString(), state.Title.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to get channel reward by name", err.Error())
			return
		}

		rewardID = reward.ID
	}

	updatedReward, err := c.TwitchClient.UpdateChannelReward(plan.BroadcasterId.ValueString(), rewardID, &helix.UpdateChannelRewardRequest{
		Title:                   plan.Title.ValueString(),
		Prompt:                  plan.Prompt.ValueString(),
		Cost:                    int(plan.Cost.ValueInt32()),
		BackgroundColor:         plan.BackgroundColor.ValueString(),
		IsGlobalCooldownEnabled: plan.IsGlobalCooldownEnabled.ValueBool(),
		GlobalCooldownSeconds:   int(plan.GlobalCooldownSeconds.ValueInt32()),
		IsEnabled:               plan.IsEnabled.ValueBool(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Failed to update channel reward id "+plan.ID.ValueString(), err.Error())
		return
	}

	plan.ID = types.StringValue(updatedReward.ID)
	plan.BackgroundColor = types.StringValue(updatedReward.BackgroundColor)
	plan.IsGlobalCooldownEnabled = types.BoolValue(updatedReward.GlobalCooldownSetting.IsEnabled)
	plan.GlobalCooldownSeconds = types.Int32Value(int32(updatedReward.GlobalCooldownSetting.GlobalCooldownSeconds))
	plan.IsEnabled = types.BoolValue(updatedReward.IsEnabled)
	plan.Cost = types.Int32Value(int32(updatedReward.Cost))
	plan.Prompt = types.StringValue(updatedReward.Prompt)
	plan.Title = types.StringValue(updatedReward.Title)

	diags = resp.State.Set(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *channelRewardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state channelRewardResourceModel

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := c.TwitchClient.DeleteChannelReward(state.BroadcasterId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete channel reward", err.Error())
		return
	}
}

func (c *channelRewardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
