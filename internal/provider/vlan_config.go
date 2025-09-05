package provider

import (
	swos_client "github.com/finomen/swos-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &SwOsResource[VlanConfigModel, swos_client.Vlan]{}
var _ resource.ResourceWithImportState = &SwOsResource[VlanConfigModel, swos_client.Vlan]{}

type VlanConfigModel struct {
	Id                    types.Int32 `tfsdk:"id"`
	IndependentVlanLookup types.Bool  `tfsdk:"independent_vlan_lookup"`
	IgmpSnooping          types.Bool  `tfsdk:"igmp_snooping"`
}

func NewVlanConfig() resource.Resource {
	return &SwOsResource[VlanConfigModel, swos_client.Vlan]{
		name:        "vlan",
		description: "VLAN configuration",
		fields: []syncedField[VlanConfigModel, swos_client.Vlan]{
			&syncedFieldImpl[int, swos_client.Vlan, VlanConfigModel, types.Int32]{
				backendGet: func(vlan *swos_client.Vlan) *int {
					return &vlan.Id
				},
				modelGet: func(model *VlanConfigModel) *types.Int32 {
					return &model.Id
				},
				fromModel: int32ValueToInt,
				toModel:   intToInt32Value,
				name:      "id",
				attribute: schema.Int32Attribute{
					MarkdownDescription: "Vlan Id",
					Required:            true,
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.RequiresReplace(),
					},
				},
			},
			&syncedFieldImpl[bool, swos_client.Vlan, VlanConfigModel, types.Bool]{
				backendGet: func(vlan *swos_client.Vlan) *bool {
					return &vlan.IndependentVlanLookup
				},
				modelGet: func(model *VlanConfigModel) *types.Bool {
					return &model.IndependentVlanLookup
				},
				fromModel: boolValueToBool,
				toModel:   types.BoolValue,
				name:      "independent_vlan_lookup",
				attribute: schema.BoolAttribute{
					MarkdownDescription: "Independent Vlan Lookup",
					Optional:            true,
					Computed:            true,
				},
			},
			&syncedFieldImpl[bool, swos_client.Vlan, VlanConfigModel, types.Bool]{
				backendGet: func(vlan *swos_client.Vlan) *bool {
					return &vlan.IgmpSnooping
				},
				modelGet: func(model *VlanConfigModel) *types.Bool {
					return &model.IgmpSnooping
				},
				fromModel: boolValueToBool,
				toModel:   types.BoolValue,
				name:      "igmp_snooping",
				attribute: schema.BoolAttribute{
					MarkdownDescription: "IGMP Snooping",
					Computed:            true,
					Optional:            true,
				},
			},
		},
		delete: func(client *swos_client.SwOsClient, model *VlanConfigModel) error {
			client.Vlan.DeleteVlan(int(model.Id.ValueInt32()))
			return nil
		},
		create: func(client *swos_client.SwOsClient, model *VlanConfigModel) (*swos_client.Vlan, error) {
			return client.Vlan.AddVlan(int(model.Id.ValueInt32()))
		},
		get: func(client *swos_client.SwOsClient, model *VlanConfigModel) (*swos_client.Vlan, error) {
			return client.Vlan.GetVlan(int(model.Id.ValueInt32()))
		},
	}
}
