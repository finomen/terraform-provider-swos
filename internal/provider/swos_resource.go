package provider

import (
	"context"
	"fmt"

	swos_client "github.com/finomen/swos-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type SwOsResource[M any, B any] struct {
	client      *swos_client.SwOsClient
	name        string
	description string
	fields      []syncedField[M, B]

	delete func(client *swos_client.SwOsClient, model *M) error
	create func(client *swos_client.SwOsClient, model *M) (*B, error)
	get    func(client *swos_client.SwOsClient, model *M) (*B, error)
}

func (s *SwOsResource[M, B]) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = fmt.Sprintf("%s_%s", request.ProviderTypeName, s.name)
}

func (s *SwOsResource[M, B]) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: s.description,
		Attributes:          map[string]schema.Attribute{},
	}
	for _, f := range s.fields {
		response.Schema.Attributes[f.Name()] = f.Attribute()
	}
}

func (r *SwOsResource[M, B]) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *SwOsResource[M, B]) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data M
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	res, err := s.create(s.client, &data)

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Unable to create %s", s.name), err.Error())
		return
	}
	for _, field := range s.fields {
		field.Sync(res, &data)
	}

	err = s.client.Save()

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Unable to save config for %s", s.name), err.Error())
		return
	}

	tflog.Trace(ctx, "created a resource")

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (s *SwOsResource[M, B]) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data M

	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	res, err := s.get(s.client, &data)

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Unable to get %s", s.name), err.Error())
		return
	}

	for _, field := range s.fields {
		field.Read(res, &data)
	}
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (s *SwOsResource[M, B]) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data M
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	res, err := s.get(s.client, &data)

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Unable to get %s", s.name), err.Error())
		return
	}

	for _, field := range s.fields {
		field.Sync(res, &data)
	}

	err = s.client.Save()

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Unable to save config for %s", s.name), err.Error())
		return
	}

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (s *SwOsResource[M, B]) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data M
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	err := s.delete(s.client, &data)

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Unable to delete %s", s.name), err.Error())
		return
	}

	err = s.client.Save()

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Unable to save config for %s", s.name), err.Error())
		return
	}

	if response.Diagnostics.HasError() {
		return
	}
}

func (s *SwOsResource[M, B]) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}
