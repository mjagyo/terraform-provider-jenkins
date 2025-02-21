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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mjagyo/jenkins-client-go"
)

// hashicupsProviderModel maps provider schema data to a Go type.
type jenkinsProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Token    types.String `tfsdk:"token"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &jenkinsProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &jenkinsProvider{
			version: version,
		}
	}
}

// jenkinsProvider is the provider implementation.
type jenkinsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *jenkinsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jenkins"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *jenkinsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *jenkinsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring HashiCups client")
	// Retrieve provider data from configuration
	var config jenkinsProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Jenkins API Host",
			"The provider cannot create the Jenkins API client as there is an unknown configuration value for the Jenkins API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the Jenkins_HOST environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Jenkins API Token",
			"The provider cannot create the Jenkins API client as there is an unknown configuration value for the Jenkins API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the Jenkins_TOKEN environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Jenkins API Username",
			"The provider cannot create the Jenkins API client as there is an unknown configuration value for the Jenkins API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the Jenkins_USERNAME environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("JENKINS_HOST")
	token := os.Getenv("JENKINS_TOKEN")
	username := os.Getenv("JENKINS_USERNAME")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing JENKINS API Host",
			"The provider cannot create the JENKINS API client as there is a missing or empty value for the JENKINS API host. "+
				"Set the host value in the configuration or use the JENKINS_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing JENKINS API Token",
			"The provider cannot create the JENKINS API client as there is a missing or empty value for the JENKINS API token. "+
				"Set the token value in the configuration or use the JENKINS_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing JENKINS API Username",
			"The provider cannot create the JENKINS API client as there is a missing or empty value for the JENKINS API username. "+
				"Set the username value in the configuration or use the JENKINS_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "jenkins_host", host)
	ctx = tflog.SetField(ctx, "jenkins_username", username)
	ctx = tflog.SetField(ctx, "jenkins_token", token)

	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "jenkins_token")

	tflog.Debug(ctx, "Creating HashiCups client")

	// Create a new Jenkins client using the configuration values
	client, err := jenkins.NewClient(&host, &username, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Authenticate Jenkins Token Client",
			"An unexpected error occurred when authenticating the Jenkins API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Jenkins Client Error: "+err.Error(),
		)
		return
	}

	// Make the Jenkins client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Jenkins client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *jenkinsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewJenkinsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *jenkinsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSecretResource,
		NewJobResource,
	}
}
