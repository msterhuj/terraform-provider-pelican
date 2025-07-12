// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	pelican "terraform-provider-pelican/internal/pelican"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

// NewUsersDataSource is a helper function to simplify the provider implementation.
func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// usersDataSource is the data source implementation.
type usersDataSource struct {
	client *pelican.Pelican
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pelican.Pelican)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *pelican.Pelican, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// add multiple attributes i whan to set multiple filters
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"external_id": schema.StringAttribute{
							Computed: true,
						},
						"uuid": schema.StringAttribute{
							Computed: true,
						},
						"username": schema.StringAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
						"first_name": schema.StringAttribute{
							Computed: true,
						},
						"last_name": schema.StringAttribute{
							Computed: true,
						},
						"language": schema.StringAttribute{
							Computed: true,
						},
						"root_admin": schema.BoolAttribute{
							Computed: true,
						},
						"two_fa": schema.BoolAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// coffeesDataSourceModel maps the data source schema data.
type usersDataSourceModel struct {
	Users []userModel `tfsdk:"users"`
}

// coffeesModel maps coffees schema data.
type userModel struct {
	ID         types.Int64  `tfsdk:"id"`
	ExternalID types.String `tfsdk:"external_id"`
	UUID       types.String `tfsdk:"uuid"`
	Username   types.String `tfsdk:"username"`
	Email      types.String `tfsdk:"email"`
	FirstName  types.String `tfsdk:"first_name"`
	LastName   types.String `tfsdk:"last_name"`
	Language   types.String `tfsdk:"language"`
	RootAdmin  types.Bool   `tfsdk:"root_admin"`
	TwoFA      types.Bool   `tfsdk:"two_fa"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

// Read refreshes the Terraform state with the latest data.
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state usersDataSourceModel

	pelicanUsers, err := d.client.GetUsers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Pelican Users Data Source Error",
			fmt.Sprintf("Unable to retrieve users: %s", err),
		)
		return
	}

	for _, user := range *pelicanUsers {
		userState := userModel{
			ID:         types.Int64Value(int64(user.ID)),
			ExternalID: types.StringValue(user.ExternalID),
			UUID:       types.StringValue(user.UUID),
			Username:   types.StringValue(user.Username),
			Email:      types.StringValue(user.Email),
			FirstName:  types.StringValue(user.FirstName),
			LastName:   types.StringValue(user.LastName),
			Language:   types.StringValue(user.Language),
			RootAdmin:  types.BoolValue(user.RootAdmin),
			TwoFA:      types.BoolValue(user.TwoFA),
			CreatedAt:  types.StringValue(user.CreatedAt),
			UpdatedAt:  types.StringValue(user.UpdatedAt),
		}

		state.Users = append(state.Users, userState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
