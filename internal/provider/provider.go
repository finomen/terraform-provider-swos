package provider

import (
	"context"
	"fmt"

	"github.com/finomen/swos-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type swosProvider struct {
	version string
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &swosProvider{
			version: version,
		}
	}
}

type swosProviderModel struct {
	Url      types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *swosProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "swos"
	resp.Version = p.version
}

func (p *swosProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *swosProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	// Retrieve provider data from configuration
	var config swosProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Url.IsUnknown() || config.Url.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"No Url",
			"Url is required",
		)
	}
	if config.Username.IsUnknown() || config.Username.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"No Username",
			"Username is required",
		)
	}
	if config.Password.IsUnknown() || config.Password.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"No Password",
			"Password is required",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	swClient, err := swos_client.NewSwOsClient(config.Url.ValueString(), config.Username.ValueString(), config.Password.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Connect failed",
			fmt.Sprintf("Failed to connect to switch: %v", err),
		)
		return
	}
	err = swClient.Fetch()

	if err != nil {
		resp.Diagnostics.AddError(
			"State fetch failed",
			fmt.Sprintf("Failed to connect to switch: %v", err),
		)
		return
	}

	resp.ResourceData = swClient
}

func (p *swosProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *swosProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSwOsConfig,
		NewVlanConfig,
		NewPortConfig,
		NewPortVlanConfig,
	}
}
