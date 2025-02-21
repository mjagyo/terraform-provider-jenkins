package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mjagyo/jenkins-client-go"
)

// coffeesDataSourceModel maps the data source schema data.
type jenkinsDataSourceModel struct {
	Jobs []jobModel `tfsdk:"jobs"`
}

type jobModel struct {
	Description types.String `tfsdk:"description" json:"description"`
	Name        types.String `tfsdk:"name" json:"name"`
	URL         types.String `tfsdk:"url" json:"url"`
	Buildable   types.Bool   `tfsdk:"buildable" json:"buildable"`
	Color       types.String `tfsdk:"color" json:"color"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &jenkinsDataSource{}
	_ datasource.DataSourceWithConfigure = &jenkinsDataSource{}
)

// NewJenkinsDataSource is a helper function to simplify the provider implementation.
func NewJenkinsDataSource() datasource.DataSource {
	return &jenkinsDataSource{}
}

// jenkinsDataSource is the data source implementation.
type jenkinsDataSource struct {
	client *jenkins.Client
}

// Metadata returns the data source type name.
func (d *jenkinsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobs"
}

// Schema defines the schema for the data source.
func (d *jenkinsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"jobs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"url": schema.StringAttribute{
							Computed: true,
						},
						"color": schema.StringAttribute{
							Computed: true,
						},
						"buildable": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *jenkinsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state jenkinsDataSourceModel

	jobs, err := d.client.GetJobs(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Jenkins Jobs",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, job := range jobs.Jobs {
		coffeeState := jobModel{
			Name:        types.StringValue(job.Name),
			Description: types.StringValue(job.Description),
			URL:         types.StringValue(job.URL),
			Buildable:   types.BoolValue(job.Buildable),
			Color:       types.StringValue(job.Color),
		}

		state.Jobs = append(state.Jobs, coffeeState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *jenkinsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jenkins.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jenkins.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
