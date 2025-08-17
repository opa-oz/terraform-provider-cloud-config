package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NullWhen
// Weird workaround for older versions of terraform
// Example of failing tests in versions >1.0.* <1.3.*:
//
//	https://github.com/opa-oz/terraform-provider-cloud-config/actions/runs/17017998097/job/48243130914
//
// @see https://github.com/hashicorp/terraform-plugin-framework/issues/603
func NullWhen(path path.Path) planmodifier.Object {
	return nullWhen{
		path: path,
	}
}

// nullWhen implements the plan modifier.
type nullWhen struct {
	path path.Path
}

// Description returns a human-readable description of the plan modifier.
func (m nullWhen) Description(_ context.Context) string {
	return "Block is null if it was deleted (doh)"
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m nullWhen) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyObject implements the plan modification logic.
func (m nullWhen) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if !req.PlanValue.IsUnknown() && (req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull()) {

		other := &types.Object{}
		ds := req.Plan.GetAttribute(ctx, m.path, other)
		resp.Diagnostics.Append(ds...)

		if ds.HasError() {
			return
		}

		if other.IsNull() || other.IsUnknown() {
			return
		}

		resp.PlanValue = types.ObjectNull(other.AttributeTypes(ctx))
	}
}
