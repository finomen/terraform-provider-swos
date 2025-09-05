package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

func int32ValueToInt(v types.Int32) int {
	return int(v.ValueInt32())
}

func intToInt32Value(v int) types.Int32 {
	return types.Int32Value(int32(v))
}

func boolValueToBool(v types.Bool) bool {
	return v.ValueBool()
}

func stringValueToString(v types.String) string {
	return v.ValueString()
}

func mapEnumConverterFromModel[T any](m map[string]T) func(v types.String) T {
	return func(v types.String) T {
		return m[v.ValueString()]
	}
}

func mapEnumConverterToModel[T comparable](m map[string]T) func(v T) types.String {
	return func(value T) types.String {
		for k, v := range m {
			if v == value {
				return types.StringValue(k)
			}
		}
		panic("Unknown value")
	}
}
