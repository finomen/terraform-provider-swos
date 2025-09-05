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

type PortVlanConfigModel struct {
	Port           types.Int32  `tfsdk:"port"`
	Mode           types.String `tfsdk:"mode"`
	Receive        types.String `tfsdk:"receive"`
	DefaultlVlanId types.Int32  `tfsdk:"default_vlan_id"`
	ForceVlanId    types.Bool   `tfsdk:"force_vlan_id"`
	Header         types.String `tfsdk:"header"`
}

var _ resource.Resource = &SwOsResource[PortVlanConfigModel, swos_client.PortForward]{}
var _ resource.ResourceWithImportState = &SwOsResource[PortVlanConfigModel, swos_client.Link]{}

func getPortForward(client *swos_client.SwOsClient, model *PortVlanConfigModel) (*swos_client.PortForward, error) {
	pid := int(model.Port.ValueInt32() - 1)
	if pid >= len(client.Links.Links) {
		return nil, fmt.Errorf("invalid port id %v, valid ids are [1,%v]", model.Port.ValueInt32(), len(client.Links.Links))
	}
	return &client.Fwd.PortForward[pid], nil
}

func NewPortVlanConfig() resource.Resource {
	vlanModes := map[string]swos_client.VlanMode{
		"disabled": swos_client.VlanModeDisabled,
		"optional": swos_client.VlanModeOptional,
		"enabled":  swos_client.VlanModeEnabled,
		"strict":   swos_client.VlanModeStrict,
	}
	vlanReceive := map[string]swos_client.VlanReceive{
		"any":      swos_client.VlanReceiveAny,
		"tagged":   swos_client.VlanReceiveTagged,
		"untagged": swos_client.VlanReceiveUntagged,
	}
	vlanHeader := map[string]swos_client.VlanHeader{
		"leave_as_is":    swos_client.VlanHeaderLeaveAsIs,
		"strip":          swos_client.VlanHeaderStrip,
		"add_If_missing": swos_client.VlanHeaderAddMissing,
	}

	return &SwOsResource[PortVlanConfigModel, swos_client.PortForward]{
		name:        "port_vlan",
		description: "Port VLAN configuration",
		fields: []syncedField[PortVlanConfigModel, swos_client.PortForward]{
			&syncedFieldImpl[int, swos_client.PortForward, PortVlanConfigModel, types.Int32]{
				modelGet: func(model *PortVlanConfigModel) *types.Int32 {
					return &model.Port
				},
				name: "port",
				attribute: schema.Int32Attribute{
					MarkdownDescription: "Port Id",
					Required:            true,
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.RequiresReplace(),
					},
				},
			},
			&syncedFieldImpl[swos_client.VlanMode, swos_client.PortForward, PortVlanConfigModel, types.String]{
				backendGet: func(fwd *swos_client.PortForward) *swos_client.VlanMode {
					return &fwd.VlanMode
				},
				modelGet: func(model *PortVlanConfigModel) *types.String {
					return &model.Mode
				},
				toModel:   mapEnumConverterToModel(vlanModes),
				fromModel: mapEnumConverterFromModel(vlanModes),
				name:      "mode",
				attribute: schema.StringAttribute{
					MarkdownDescription: "VLAN Mode",
					Optional:            true,
					Computed:            true,
				},
			},
			&syncedFieldImpl[swos_client.VlanReceive, swos_client.PortForward, PortVlanConfigModel, types.String]{
				backendGet: func(fwd *swos_client.PortForward) *swos_client.VlanReceive {
					return &fwd.VlanReceive
				},
				modelGet: func(model *PortVlanConfigModel) *types.String {
					return &model.Mode
				},
				toModel:   mapEnumConverterToModel(vlanReceive),
				fromModel: mapEnumConverterFromModel(vlanReceive),
				name:      "vlan_receive",
				attribute: schema.StringAttribute{
					MarkdownDescription: "VLAN Receive",
					Optional:            true,
					Computed:            true,
				},
			},
			&syncedFieldImpl[int, swos_client.PortForward, PortVlanConfigModel, types.Int32]{
				backendGet: func(fwd *swos_client.PortForward) *int {
					return &fwd.DefaultVlanId
				},
				modelGet: func(model *PortVlanConfigModel) *types.Int32 {
					return &model.DefaultlVlanId
				},
				toModel:   intToInt32Value,
				fromModel: int32ValueToInt,
				name:      "default_vlan_id",
				attribute: schema.Int32Attribute{
					MarkdownDescription: "Default VLAN Id",
					Optional:            true,
					Computed:            true,
				},
			},
			&syncedFieldImpl[bool, swos_client.PortForward, PortVlanConfigModel, types.Bool]{
				backendGet: func(fwd *swos_client.PortForward) *bool {
					return &fwd.ForceVlanId
				},
				modelGet: func(model *PortVlanConfigModel) *types.Bool {
					return &model.ForceVlanId
				},
				toModel:   types.BoolValue,
				fromModel: boolValueToBool,
				name:      "force_vlan_id",
				attribute: schema.BoolAttribute{
					MarkdownDescription: "Force VLAN Id",
					Optional:            true,
					Computed:            true,
				},
			},
			&syncedFieldImpl[swos_client.VlanHeader, swos_client.PortForward, PortVlanConfigModel, types.String]{
				backendGet: func(fwd *swos_client.PortForward) *swos_client.VlanHeader {
					return &fwd.VlanHeader
				},
				modelGet: func(model *PortVlanConfigModel) *types.String {
					return &model.Header
				},
				toModel:   mapEnumConverterToModel(vlanHeader),
				fromModel: mapEnumConverterFromModel(vlanHeader),
				name:      "header",
				attribute: schema.StringAttribute{
					MarkdownDescription: "VLAN Header",
					Optional:            true,
					Computed:            true,
				},
			},
		},
		delete: func(client *swos_client.SwOsClient, model *PortVlanConfigModel) error {
			return nil
		},
		create: getPortForward,
		get:    getPortForward,
	}
}
