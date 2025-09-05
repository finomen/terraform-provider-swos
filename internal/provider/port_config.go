package provider

import (
	"fmt"

	swos_client "github.com/finomen/swos-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PortConfigModel struct {
	Id          types.Int32  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	FlowControl types.Bool   `tfsdk:"flow_control"`
	PoeOut      types.String `tfsdk:"poe_out"`
	PoePriority types.Int32  `tfsdk:"poe_priority"`
}

var _ resource.Resource = &SwOsResource[PortConfigModel, swos_client.Link]{}
var _ resource.ResourceWithImportState = &SwOsResource[PortConfigModel, swos_client.Link]{}

func getPort(client *swos_client.SwOsClient, model *PortConfigModel) (*swos_client.Link, error) {
	pid := int(model.Id.ValueInt32() - 1)
	if pid >= len(client.Links.Links) {
		return nil, fmt.Errorf("invalid port id %v, valid ids are [1,%v]", model.Id.ValueInt32(), len(client.Links.Links))
	}
	return client.Links.Links[pid], nil
}

func NewPortConfig() resource.Resource {
	poeModes := map[string]swos_client.PoeMode{
		"off":   swos_client.Off,
		"auto":  swos_client.Auto,
		"on":    swos_client.On,
		"calib": swos_client.Calib,
	}

	return &SwOsResource[PortConfigModel, swos_client.Link]{
		name:        "port",
		description: "Port configuration",
		fields: []syncedField[PortConfigModel, swos_client.Link]{
			&syncedFieldImpl[int, swos_client.Link, PortConfigModel, types.Int32]{
				modelGet: func(model *PortConfigModel) *types.Int32 {
					return &model.Id
				},
				name: "id",
				attribute: schema.Int32Attribute{
					MarkdownDescription: "Port Id",
					Required:            true,
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.RequiresReplace(),
					},
				},
			},
			&syncedFieldImpl[string, swos_client.Link, PortConfigModel, types.String]{
				backendGet: func(link *swos_client.Link) *string {
					return &link.Name
				},
				modelGet: func(model *PortConfigModel) *types.String {
					return &model.Name
				},
				toModel:   types.StringValue,
				fromModel: stringValueToString,
				name:      "name",
				attribute: schema.StringAttribute{
					MarkdownDescription: "Port Name",
					Optional:            true,
					Computed:            true,
				},
			},
			&syncedFieldImpl[bool, swos_client.Link, PortConfigModel, types.Bool]{
				backendGet: func(link *swos_client.Link) *bool {
					return &link.Enabled
				},
				modelGet: func(model *PortConfigModel) *types.Bool {
					return &model.Enabled
				},
				toModel:   types.BoolValue,
				fromModel: boolValueToBool,
				name:      "enabled",
				attribute: schema.BoolAttribute{
					MarkdownDescription: "Port enabled",
					Optional:            true,
					Computed:            true,
				},
			},
			&syncedFieldImpl[bool, swos_client.Link, PortConfigModel, types.Bool]{
				backendGet: func(link *swos_client.Link) *bool {
					return &link.FlowControl
				},
				modelGet: func(model *PortConfigModel) *types.Bool {
					return &model.FlowControl
				},
				toModel:   types.BoolValue,
				fromModel: boolValueToBool,
				name:      "flow_control",
				attribute: schema.BoolAttribute{
					MarkdownDescription: "Port enabled",
					Optional:            true,
					Computed:            true,
				},
			},
			&syncedFieldImpl[swos_client.PoeMode, swos_client.Link, PortConfigModel, types.String]{
				backendGet: func(link *swos_client.Link) *swos_client.PoeMode {
					return &link.PoeMode
				},
				modelGet: func(model *PortConfigModel) *types.String {
					return &model.PoeOut
				},
				toModel:   mapEnumConverterToModel(poeModes),
				fromModel: mapEnumConverterFromModel(poeModes),
				name:      "poe_out",
				attribute: schema.StringAttribute{
					MarkdownDescription: "PoE Out",
					Optional:            true,
					Computed:            true,
				},
			},
			&syncedFieldImpl[int, swos_client.Link, PortConfigModel, types.Int32]{
				backendGet: func(link *swos_client.Link) *int {
					return &link.PoePrio
				},
				modelGet: func(model *PortConfigModel) *types.Int32 {
					return &model.PoePriority
				},
				toModel:   intToInt32Value,
				fromModel: int32ValueToInt,
				name:      "poe_priority",
				attribute: schema.Int32Attribute{
					MarkdownDescription: "PoE priority",
					Optional:            true,
					Computed:            true,
				},
			},
		},
		delete: func(client *swos_client.SwOsClient, model *PortConfigModel) error {
			return nil
		},
		create: getPort,
		get:    getPort,
	}
}
