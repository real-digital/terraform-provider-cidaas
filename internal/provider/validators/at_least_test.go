package validators_test

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/real-digital/terraform-provider-cidaas/internal/provider/validators"
	"testing"
)

func TestAtLeastValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		min         int64
		expectError bool
	}

	tests := map[string]testCase{
		"not an Int64": {
			val:         types.Bool{Value: true},
			expectError: true,
		},
		"unknown Int64": {
			val: types.Int64{Unknown: true},
			min: 1,
		},
		"null Int64": {
			val: types.Int64{Null: true},
			min: 1,
		},
		"valid integer as Int64": {
			val: types.Int64{Value: 2},
			min: 1,
		},
		"to small integer as Int64": {
			val:         types.Int64{Value: -1},
			min:         1,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			request := tfsdk.ValidateAttributeRequest{
				AttributePath:   path.Root("test"),
				AttributeConfig: test.val,
			}

			response := tfsdk.ValidateAttributeResponse{}
			validators.AtLeast(test.min).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
