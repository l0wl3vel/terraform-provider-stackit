package sqlserverflex

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

func TestMapDataSourceFields(t *testing.T) {
	tests := []struct {
		description string
		input       *sqlserverflex.GetUserResponse
		expected    DataSourceModel
		isValid     bool
	}{
		{
			"default_values",
			&sqlserverflex.GetUserResponse{
				Item: &sqlserverflex.InstanceResponseUser{},
			},
			DataSourceModel{
				Id:         types.StringValue("pid,iid,uid"),
				UserId:     types.StringValue("uid"),
				InstanceId: types.StringValue("iid"),
				ProjectId:  types.StringValue("pid"),
				Username:   types.StringNull(),
				Roles:      types.SetNull(types.StringType),
				Host:       types.StringNull(),
				Port:       types.Int64Null(),
			},
			true,
		},
		{
			"simple_values",
			&sqlserverflex.GetUserResponse{
				Item: &sqlserverflex.InstanceResponseUser{
					Roles: &[]string{
						"role_1",
						"role_2",
						"",
					},
					Username: utils.Ptr("username"),
					Host:     utils.Ptr("host"),
					Port:     utils.Ptr(int64(1234)),
				},
			},
			DataSourceModel{
				Id:         types.StringValue("pid,iid,uid"),
				UserId:     types.StringValue("uid"),
				InstanceId: types.StringValue("iid"),
				ProjectId:  types.StringValue("pid"),
				Username:   types.StringValue("username"),
				Roles: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("role_1"),
					types.StringValue("role_2"),
					types.StringValue(""),
				}),
				Host: types.StringValue("host"),
				Port: types.Int64Value(1234),
			},
			true,
		},
		{
			"null_fields_and_int_conversions",
			&sqlserverflex.GetUserResponse{
				Item: &sqlserverflex.InstanceResponseUser{
					Id:       utils.Ptr("uid"),
					Roles:    &[]string{},
					Username: nil,
					Host:     nil,
					Port:     utils.Ptr(int64(2123456789)),
				},
			},
			DataSourceModel{
				Id:         types.StringValue("pid,iid,uid"),
				UserId:     types.StringValue("uid"),
				InstanceId: types.StringValue("iid"),
				ProjectId:  types.StringValue("pid"),
				Username:   types.StringNull(),
				Roles:      types.SetValueMust(types.StringType, []attr.Value{}),
				Host:       types.StringNull(),
				Port:       types.Int64Value(2123456789),
			},
			true,
		},
		{
			"nil_response",
			nil,
			DataSourceModel{},
			false,
		},
		{
			"nil_response_2",
			&sqlserverflex.GetUserResponse{},
			DataSourceModel{},
			false,
		},
		{
			"no_resource_id",
			&sqlserverflex.GetUserResponse{
				Item: &sqlserverflex.InstanceResponseUser{},
			},
			DataSourceModel{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			state := &DataSourceModel{
				ProjectId:  tt.expected.ProjectId,
				InstanceId: tt.expected.InstanceId,
				UserId:     tt.expected.UserId,
			}
			err := mapDataSourceFields(tt.input, state)
			if !tt.isValid && err == nil {
				t.Fatalf("Should have failed")
			}
			if tt.isValid && err != nil {
				t.Fatalf("Should not have failed: %v", err)
			}
			if tt.isValid {
				diff := cmp.Diff(state, &tt.expected)
				if diff != "" {
					t.Fatalf("Data does not match: %s", diff)
				}
			}
		})
	}
}
