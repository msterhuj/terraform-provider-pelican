// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	client "terraform-provider-pelican/internal/pelican"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &pelicanProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &pelicanProvider{
			version: version,
		}
	}
}

// pelicanProvider is the provider implementation.
type pelicanProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *pelicanProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pelican"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *pelicanProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// pelicanProviderModel maps provider schema data to a Go type.
type pelicanProviderModel struct {
	Server types.String `tfsdk:"server"`
	Token  types.String `tfsdk:"token"`
}

// Configure prepares a pelican API client for data sources and resources.
func (p *pelicanProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config pelicanProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Server.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Provider Configuration",
			"The provider configuration for `server` is unknown. Please set a value for the `server` attribute in the provider configuration block.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Provider Configuration",
			"The provider configuration for `token` is unknown. Please set a value for the `token` attribute in the provider configuration block.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	server := os.Getenv("PELICAN_SERVER")
	token := os.Getenv("PELICAN_TOKEN")

	if !config.Server.IsNull() {
		server = config.Server.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if server == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server"),
			"Missing Provider Configuration",
			"The provider configuration for `server` is missing. Please set a value for the `server` attribute in the provider configuration block or set the pelican_SERVER environment variable.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Provider Configuration",
			"The provider configuration for `token` is missing. Please set a value for the `token` attribute in the provider configuration block or set the pelican_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	pelicantClient, err := client.NewClient(server, token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Pelican client",
			"An unexpected error occurred while creating the Pelican client: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = pelicantClient
	resp.ResourceData = pelicantClient
}

// DataSources defines the data sources implemented in the provider.
func (p *pelicanProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUsersDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *pelicanProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
