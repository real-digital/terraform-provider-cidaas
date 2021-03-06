package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/real-digital/terraform-provider-cidaas/internal/client"
)

type resourcePasswordPolicyType struct{}
type resourcePasswordPolicy struct {
	p provider
}

func (r resourcePasswordPolicyType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "`cidaas_password_policy` controls the password policies in the tenant",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Description: "Unique identifier of the policy",
			},
			"policy_name": {
				Type:        types.StringType,
				Required:    true,
				Description: "Display name of the policy",
			},
			"lower_and_upper_case": {
				Type:        types.BoolType,
				Required:    true,
				Description: "Indicates if passwords are required to have lower and upper case letters",
			},
			"minimum_length": {
				Type:        types.Int64Type,
				Required:    true,
				Description: "Minimum length of the passwords",
			},
			"no_of_digits": {
				Type:        types.Int64Type,
				Required:    true,
				Description: "Number of digits that need to be included in the password",
			},
			"no_of_special_chars": {
				Type:        types.Int64Type,
				Required:    true,
				Description: "Number of special chars that need to be included in the password",
			},
		},
	}, nil
}

func (r resourcePasswordPolicyType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourcePasswordPolicy{
		p: *(p.(*provider)),
	}, nil
}

func (r resourcePasswordPolicy) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan PasswordPolicy

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	plannedPolicy := client.PasswordPolicy{
		PolicyName:        plan.PolicyName.Value,
		MinimumLength:     plan.MinimumLength.Value,
		NoOfDigits:        plan.NoOfDigits.Value,
		LowerAndUpperCase: plan.LowerAndUpperCase.Value,
		NoOfSpecialChars:  plan.NoOfSpecialChars.Value,
	}

	policy, err := r.p.client.UpdatePasswordPolicy(plannedPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating password policy",
			"Could not create policy, unexpected error: "+err.Error(),
		)
		return
	}

	result := PasswordPolicy{
		ID:                types.String{Value: policy.ID},
		PolicyName:        types.String{Value: policy.PolicyName},
		MinimumLength:     types.Int64{Value: policy.MinimumLength},
		NoOfDigits:        types.Int64{Value: policy.NoOfDigits},
		LowerAndUpperCase: types.Bool{Value: policy.LowerAndUpperCase},
		NoOfSpecialChars:  types.Int64{Value: policy.NoOfSpecialChars},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourcePasswordPolicy) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state PasswordPolicy
	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := state.ID.Value

	policy, err := r.p.client.GetPasswordPolicy(policyID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading hook",
			"Could not read hookID "+policyID+": "+err.Error(),
		)
		return
	}

	state.ID.Value = policy.ID
	state.PolicyName.Value = policy.PolicyName
	state.LowerAndUpperCase.Value = policy.LowerAndUpperCase
	state.MinimumLength.Value = policy.MinimumLength
	state.NoOfDigits.Value = policy.NoOfDigits
	state.NoOfSpecialChars.Value = policy.NoOfSpecialChars

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourcePasswordPolicy) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan PasswordPolicy
	var state PasswordPolicy

	req.State.Get(ctx, &state)
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	plannedPolicy := client.PasswordPolicy{
		ID:                state.ID.Value,
		PolicyName:        plan.PolicyName.Value,
		MinimumLength:     plan.MinimumLength.Value,
		NoOfDigits:        plan.NoOfDigits.Value,
		LowerAndUpperCase: plan.LowerAndUpperCase.Value,
		NoOfSpecialChars:  plan.NoOfSpecialChars.Value,
	}

	policy, err := r.p.client.UpdatePasswordPolicy(plannedPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating password policy",
			"Could not update policy, unexpected error: "+err.Error(),
		)
		return
	}

	result := PasswordPolicy{
		ID:                state.ID,
		PolicyName:        types.String{Value: policy.PolicyName},
		MinimumLength:     types.Int64{Value: policy.MinimumLength},
		NoOfDigits:        types.Int64{Value: policy.NoOfDigits},
		LowerAndUpperCase: types.Bool{Value: policy.LowerAndUpperCase},
		NoOfSpecialChars:  types.Int64{Value: policy.NoOfSpecialChars},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourcePasswordPolicy) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var state PasswordPolicy

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.p.client.DeletePasswordPolicy(state.ID.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting password policy",
			"Could not delete policy, unexpected error: "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}
