package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mjagyo/jenkins-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &secretResource{}
	_ resource.ResourceWithConfigure = &secretResource{}
)

// secretResourceModel maps the resource schema data.
type secretResourceModel struct {
	SecretType  types.String          `tfsdk:"secret_type"`
	ID          types.String          `tfsdk:"id"`
	Credential  secretCredentialModel `tfsdk:"credential"`
	LastUpdated types.String          `tfsdk:"last_updated"`
}

// secretCredentialModel maps order item data.
type secretCredentialModel struct {
	Scope       types.String `tfsdk:"scope"`
	ID          types.String `tfsdk:"id"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Description types.String `tfsdk:"description"`
	AppId       types.String `tfsdk:"app_id"`
	PrivateKey  types.String `tfsdk:"private_key"`
	Secret      types.String `tfsdk:"secret"`
	Class       types.String `tfsdk:"class"`
}

// NewSecretResource is a helper function to simplify the provider implementation.
func NewSecretResource() resource.Resource {
	return &secretResource{}
}

// secretResource is the resource implementation.
type secretResource struct {
	client *jenkins.Client
}

// Metadata returns the resource type name.
func (r *secretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

// Schema defines the schema for the resource.
func (r *secretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"credential": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"scope": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"id": schema.StringAttribute{
						Optional: true,
					},
					"username": schema.StringAttribute{
						Optional: true,
					},
					"password": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
					"app_id": schema.StringAttribute{
						Optional: true,
					},
					"private_key": schema.StringAttribute{
						Optional: true,
					},
					"secret": schema.StringAttribute{
						Optional: true,
					},
					"class": schema.StringAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

// Create a new resource.
func (r *secretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan secretResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	secretType := strings.Trim(plan.SecretType.String(), `" `)
	fmt.Println("Secret type: ", secretType)

	var payload jenkins.CredentialRequest
	if secretType == "auth_pair" {
		payload = jenkins.CredentialRequest{
			Credentials: jenkins.Credential{
				Scope:       plan.Credential.Scope.ValueString(),
				ID:          plan.Credential.ID.ValueString(),
				Username:    plan.Credential.Username.ValueString(),
				Password:    plan.Credential.Password.ValueString(),
				Description: plan.Credential.Description.ValueString(),
				Class:       "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl",
			},
		}
		plan.Credential.Class = types.StringValue("com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl")
	} else if secretType == "github" {
		payload = jenkins.CredentialRequest{
			Credentials: jenkins.Credential{
				Scope:       plan.Credential.Scope.ValueString(),
				ID:          plan.Credential.ID.ValueString(),
				AppID:       plan.Credential.AppId.ValueString(),
				PrivateKey:  plan.Credential.PrivateKey.ValueString(),
				Description: plan.Credential.Description.ValueString(),
				Class:       "org.jenkinsci.plugins.github_branch_source.GitHubAppCredentials",
			},
		}
		plan.Credential.Class = types.StringValue("org.jenkinsci.plugins.github_branch_source.GitHubAppCredentials")
	} else {
		fmt.Println("Unknown secret type")
		return
	}

	tflog.Info(ctx, "Checking crendentials request", map[string]any{"success": payload})
	err := r.client.CreateSecret(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(payload.Credentials.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	if plan.SecretType.String() == "auth_pair" {
		plan.Credential = secretCredentialModel{
			Scope:       types.StringValue(payload.Credentials.Scope),
			ID:          types.StringValue(payload.Credentials.ID),
			Username:    types.StringValue(payload.Credentials.Username),
			Password:    types.StringValue(payload.Credentials.Password),
			Description: types.StringValue(payload.Credentials.Description),
			Class:       types.StringValue("com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"),
		}
	} else if plan.SecretType.String() == "github" {
		plan.Credential = secretCredentialModel{
			Scope:       types.StringValue(payload.Credentials.Scope),
			ID:          types.StringValue(payload.Credentials.ID),
			AppId:       types.StringValue(payload.Credentials.AppID),
			PrivateKey:  types.StringValue(payload.Credentials.PrivateKey),
			Description: types.StringValue(payload.Credentials.Description),
			Class:       types.StringValue("org.jenkinsci.plugins.github_branch_source.GitHubAppCredentials"),
		}
	}
	tflog.Info(ctx, "Checking plan ", map[string]any{"plan": plan})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *secretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state secretResourceModel
	diags := req.State.Get(ctx, &state)
	tflog.Info(ctx, "State retrieved", map[string]any{"state": state})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.GetSecret(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Jenkins Secret",
			"Could not read Jenkins Secret ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, "Checking secret on read", map[string]any{"secret": secret})

	splitDisplayname := strings.Split(secret.DisplayName, "/")

	// var credentialClass string
	if secret.TypeName == "Username with password" {
		state.Credential = secretCredentialModel{
			ID:          types.StringValue(secret.ID),
			Username:    types.StringValue(splitDisplayname[0]),
			Description: types.StringValue(secret.Description),
			Class:       types.StringValue(state.Credential.Class.ValueString()),
			Password:    state.Credential.Password,
			Scope:       state.Credential.Scope,
		}
	} else if secret.TypeName == "GitHub App" {
		state.Credential = secretCredentialModel{
			ID:          types.StringValue(secret.ID),
			Description: types.StringValue(secret.Description),
			Class:       types.StringValue(state.Credential.Class.ValueString()),
			PrivateKey:  state.Credential.PrivateKey,
			AppId:       state.Credential.AppId,
			Password:    state.Credential.Password,
			Scope:       state.Credential.Scope,
		}
	} else {
		fmt.Println("Unknown secret type")
		return
	}

	tflog.Info(ctx, "Checking state on read", map[string]any{"state": state})

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *secretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan secretResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Checking the plan", map[string]any{"plan": plan})

	if plan.SecretType.IsNull() || plan.SecretType.IsUnknown() {
		resp.Diagnostics.AddError(
			"Error Updating Jenkins Secret",
			"Secret type is not set in the plan.",
		)
		return
	}

	var payload jenkins.Credential
	tflog.Info(ctx, "Checking secret type", map[string]any{"type": plan.SecretType.ValueString()})

	if plan.SecretType.ValueString() == "auth_pair" {
		payload = jenkins.Credential{
			StaperClass: "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl",
			Class:       "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl",
			ID:          plan.Credential.ID.ValueString(),
			Description: plan.Credential.Description.ValueString(),
			Username:    plan.Credential.Username.ValueString(),
			Scope:       plan.Credential.Scope.ValueString(),
		}
	} else if plan.SecretType.ValueString() == "github" {
		payload = jenkins.Credential{
			StaperClass: "org.jenkinsci.plugins.github_branch_source.GitHubAppCredentials",
			Class:       "org.jenkinsci.plugins.github_branch_source.GitHubAppCredentials",
			ID:          plan.Credential.ID.ValueString(),
			Description: plan.Credential.Description.ValueString(),
			AppID:       plan.Credential.AppId.ValueString(),
			PrivateKey:  plan.Credential.PrivateKey.ValueString(),
			Scope:       plan.Credential.Scope.ValueString(),
		}
	}

	err := r.client.UpdateSecret(payload)
	tflog.Info(ctx, "Checking payload", map[string]any{"payload": payload})
	tflog.Info(ctx, "Debugging logs", map[string]any{"error": err})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Jenkins Secret",
			"Could not update secret, unexpected error: "+err.Error(),
		)
		return
	}

	secret, err := r.client.GetSecret(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Jenkins Secret",
			"Could not read Jenkins Secret ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	var credentialClass string
	if secret.TypeName == "Username with password" {
		credentialClass = "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"
	} else if secret.TypeName == "GitHub App" {
		credentialClass = "org.jenkinsci.plugins.github_branch_source.GitHubAppCredentials"
	} else {
		resp.Diagnostics.AddError(
			"Error Reading Jenkins Secret",
			"Unknown secret type: "+secret.TypeName,
		)
		return
	}

	splitDisplayname := strings.Split(secret.DisplayName, "/")

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	if plan.SecretType.ValueString() == "auth_pair" {
		plan.Credential = secretCredentialModel{
			Scope:       types.StringValue(payload.Scope),
			ID:          types.StringValue(payload.ID),
			Username:    types.StringValue(splitDisplayname[0]),
			Password:    plan.Credential.Password,
			Description: types.StringValue(payload.Description),
			Class:       types.StringValue(credentialClass),
		}
	} else if plan.SecretType.ValueString() == "github" {
		plan.Credential = secretCredentialModel{
			Scope:       types.StringValue(payload.Scope),
			ID:          types.StringValue(payload.ID),
			AppId:       plan.Credential.AppId,
			PrivateKey:  plan.Credential.PrivateKey,
			Description: types.StringValue(payload.Description),
			Class:       types.StringValue(credentialClass),
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *secretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state secretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Deleting Jenkins Secret", map[string]any{"id": state.ID.ValueString()})

	// Delete existing order
	err := r.client.DeleteSecret(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Order",
			"Could not delete order, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Jenkins Secret deleted", map[string]any{"id": state.ID.ValueString()})

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Configure adds the provider configured client to the resource.
func (r *secretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jenkins.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
