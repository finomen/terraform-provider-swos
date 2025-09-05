package provider

import (
	"context"
	"fmt"

	swos_client "github.com/finomen/swos-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SwOsConfig{}
var _ resource.ResourceWithImportState = &SwOsConfig{}

func NewSwOsConfig() resource.Resource {
	return &SwOsConfig{}
}

// SwOsConfig defines the resource implementation.
type SwOsConfig struct {
	client *swos_client.SwOsClient
}

// SwOsConfigModel describes the resource data model.
type SwOsConfigModel struct {
	Identity types.String `tfsdk:"identity"`
}

func (r *SwOsConfig) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (r *SwOsConfig) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "SwOs general configuration",

		Attributes: map[string]schema.Attribute{
			"identity": schema.StringAttribute{
				MarkdownDescription: "Switch identity",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SwOsConfig) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*swos_client.SwOsClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *swos_client.SwOsClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *SwOsConfig) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SwOsConfigModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SwOsConfig) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SwOsConfigModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Identity = types.StringValue(r.client.Sys.Identity)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SwOsConfig) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SwOsConfigModel

	// Read Terraform plan data into the model
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.client.Sys.Identity = data.Identity.ValueString()

	err := r.client.Save()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update configuration",
			fmt.Sprintf("Error: %v", err),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SwOsConfig) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SwOsConfigModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SwOsConfig) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
