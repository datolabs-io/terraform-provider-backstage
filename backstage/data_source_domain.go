package backstage

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/datolabs-io/go-backstage/v3"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &domainDataSource{}
	_ datasource.DataSourceWithConfigure = &domainDataSource{}
)

// NewDomainDataSource is a helper function to simplify the provider implementation.
func NewDomainDataSource() datasource.DataSource {
	return &domainDataSource{}
}

// domainDataSource is the data source implementation.
type domainDataSource struct {
	client *backstage.Client
}

type domainDataSourceModel struct {
	ID         types.String          `tfsdk:"id"`
	Name       types.String          `tfsdk:"name"`
	Namespace  types.String          `tfsdk:"namespace"`
	ApiVersion types.String          `tfsdk:"api_version"`
	Kind       types.String          `tfsdk:"kind"`
	Metadata   *entityMetadataModel  `tfsdk:"metadata"`
	Relations  []entityRelationModel `tfsdk:"relations"`
	Spec       *domainSpecModel      `tfsdk:"spec"`
	Fallback   *domainFallbackModel  `tfsdk:"fallback"`
}

type domainFallbackModel struct {
	ID         types.String          `tfsdk:"id"`
	Name       types.String          `tfsdk:"name"`
	Namespace  types.String          `tfsdk:"namespace"`
	ApiVersion types.String          `tfsdk:"api_version"`
	Kind       types.String          `tfsdk:"kind"`
	Metadata   *entityMetadataModel  `tfsdk:"metadata"`
	Relations  []entityRelationModel `tfsdk:"relations"`
	Spec       *domainSpecModel      `tfsdk:"spec"`
}

type domainSpecModel struct {
	Owner types.String `tfsdk:"owner"`
}

const (
	descriptionDomainSpecOwner = "An entity reference to the owner of the domain."
	descriptionDomainFallback  = "A complete replica of the `Domain` as it would exist in backstage. Set this to provide a fallback in case the Backstage instance is not functioning, is down, or is unrealiable."
)

// Metadata returns the data source type name.
func (d *domainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

// Schema defines the schema for the data source.
func (d *domainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get a specific " +
			"[Domain entity](https://backstage.io/docs/features/software-catalog/descriptor-format#kind-domain) from Backstage Software Catalog.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: descriptionEntityMetadataUID},
			"name": schema.StringAttribute{Required: true, Description: descriptionEntityMetadataName, Validators: []validator.String{
				stringvalidator.LengthBetween(1, 63),
				stringvalidator.RegexMatches(
					regexp.MustCompile(patternEntityName),
					"must follow Backstage format restrictions",
				),
			}},
			"namespace": schema.StringAttribute{Optional: true, Description: descriptionEntityMetadataNamespace, Validators: []validator.String{
				stringvalidator.LengthBetween(1, 63),
				stringvalidator.RegexMatches(
					regexp.MustCompile(patternEntityName),
					"must follow Backstage format restrictions",
				),
			}},
			"api_version": schema.StringAttribute{Computed: true, Description: descriptionEntityApiVersion},
			"kind":        schema.StringAttribute{Computed: true, Description: descriptionEntityKind},
			"metadata": schema.SingleNestedAttribute{Computed: true, Description: descriptionEntityMetadata, Attributes: map[string]schema.Attribute{
				"uid":         schema.StringAttribute{Computed: true, Description: descriptionEntityMetadataUID},
				"etag":        schema.StringAttribute{Computed: true, Description: descriptionEntityMetadataEtag},
				"name":        schema.StringAttribute{Computed: true, Description: descriptionEntityMetadataName},
				"namespace":   schema.StringAttribute{Computed: true, Description: descriptionEntityMetadataNamespace},
				"title":       schema.StringAttribute{Computed: true, Description: descriptionEntityMetadataTitle},
				"description": schema.StringAttribute{Computed: true, Description: descriptionEntityMetadataDescription},
				"labels":      schema.MapAttribute{Computed: true, Description: descriptionEntityMetadataLabels, ElementType: types.StringType},
				"annotations": schema.MapAttribute{Computed: true, Description: descriptionEntityMetadataAnnotations, ElementType: types.StringType},
				"tags":        schema.ListAttribute{Computed: true, Description: descriptionEntityMetadataTags, ElementType: types.StringType},
				"links": schema.ListNestedAttribute{Computed: true, Description: descriptionEntityMetadataLinks, NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"url":   schema.StringAttribute{Computed: true, Description: descriptionEntityLinkURL},
						"title": schema.StringAttribute{Computed: true, Description: descriptionEntityLinkTitle},
						"icon":  schema.StringAttribute{Computed: true, Description: descriptionEntityLinkIco},
						"type":  schema.StringAttribute{Computed: true, Description: descriptionEntityLinkType},
					},
				}},
			}},
			"relations": schema.ListNestedAttribute{Computed: true, Description: descriptionEntityRelations, NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"type":       schema.StringAttribute{Computed: true, Description: descriptionEntityRelationType},
					"target_ref": schema.StringAttribute{Computed: true, Description: descriptionEntityRelationTargetRef},
					"target": schema.SingleNestedAttribute{Computed: true, Description: descriptionEntityRelationTarget,
						Attributes: map[string]schema.Attribute{
							"name":      schema.StringAttribute{Computed: true, Description: descriptionEntityRelationTargetName},
							"kind":      schema.StringAttribute{Computed: true, Description: descriptionEntityRelationTargetKind},
							"namespace": schema.StringAttribute{Computed: true, Description: descriptionEntityRelationTargetNamespace},
						}},
				},
			}},
			"spec": schema.SingleNestedAttribute{Computed: true, Description: descriptionEntitySpec, Attributes: map[string]schema.Attribute{
				"owner": schema.StringAttribute{Computed: true, Description: descriptionDomainSpecOwner},
			}},
			"fallback": schema.SingleNestedAttribute{Optional: true, Description: descriptionDomainFallback, Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{Optional: true, Description: descriptionEntityMetadataUID},
				"name": schema.StringAttribute{Required: true, Description: descriptionEntityMetadataName, Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					stringvalidator.RegexMatches(
						regexp.MustCompile(patternEntityName),
						"must follow Backstage format restrictions",
					),
				}},
				"namespace": schema.StringAttribute{Optional: true, Description: descriptionEntityMetadataNamespace, Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					stringvalidator.RegexMatches(
						regexp.MustCompile(patternEntityName),
						"must follow Backstage format restrictions",
					),
				}},
				"api_version": schema.StringAttribute{Optional: true, Description: descriptionEntityApiVersion},
				"kind":        schema.StringAttribute{Optional: true, Description: descriptionEntityKind},
				"metadata": schema.SingleNestedAttribute{Optional: true, Description: descriptionEntityMetadata, Attributes: map[string]schema.Attribute{
					"uid":         schema.StringAttribute{Optional: true, Description: descriptionEntityMetadataUID},
					"etag":        schema.StringAttribute{Optional: true, Description: descriptionEntityMetadataEtag},
					"name":        schema.StringAttribute{Optional: true, Description: descriptionEntityMetadataName},
					"namespace":   schema.StringAttribute{Optional: true, Description: descriptionEntityMetadataNamespace},
					"title":       schema.StringAttribute{Optional: true, Description: descriptionEntityMetadataTitle},
					"description": schema.StringAttribute{Optional: true, Description: descriptionEntityMetadataDescription},
					"labels":      schema.MapAttribute{Optional: true, Description: descriptionEntityMetadataLabels, ElementType: types.StringType},
					"annotations": schema.MapAttribute{Optional: true, Description: descriptionEntityMetadataAnnotations, ElementType: types.StringType},
					"tags":        schema.ListAttribute{Optional: true, Description: descriptionEntityMetadataTags, ElementType: types.StringType},
					"links": schema.ListNestedAttribute{Optional: true, Description: descriptionEntityMetadataLinks, NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"url":   schema.StringAttribute{Optional: true, Description: descriptionEntityLinkURL},
							"title": schema.StringAttribute{Optional: true, Description: descriptionEntityLinkTitle},
							"icon":  schema.StringAttribute{Optional: true, Description: descriptionEntityLinkIco},
							"type":  schema.StringAttribute{Optional: true, Description: descriptionEntityLinkType},
						},
					}},
				}},
				"relations": schema.ListNestedAttribute{Optional: true, Description: descriptionEntityRelations, NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type":       schema.StringAttribute{Optional: true, Description: descriptionEntityRelationType},
						"target_ref": schema.StringAttribute{Optional: true, Description: descriptionEntityRelationTargetRef},
						"target": schema.SingleNestedAttribute{Optional: true, Description: descriptionEntityRelationTarget,
							Attributes: map[string]schema.Attribute{
								"name":      schema.StringAttribute{Optional: true, Description: descriptionEntityRelationTargetName},
								"kind":      schema.StringAttribute{Optional: true, Description: descriptionEntityRelationTargetKind},
								"namespace": schema.StringAttribute{Optional: true, Description: descriptionEntityRelationTargetNamespace},
							}},
					},
				}},
				"spec": schema.SingleNestedAttribute{Optional: true, Description: descriptionEntitySpec, Attributes: map[string]schema.Attribute{
					"owner": schema.StringAttribute{Optional: true, Description: descriptionDomainSpecOwner},
				}},
			}},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *domainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*backstage.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *domainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state domainDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Namespace.IsNull() {
		state.Namespace = types.StringValue(backstage.DefaultNamespaceName)
	}

	tflog.Debug(ctx, fmt.Sprintf("Getting Domain kind %s/%s from Backstage API", state.Name.ValueString(), state.Namespace.ValueString()))
	domain, response, err := d.client.Catalog.Domains.Get(ctx, state.Name.ValueString(), state.Namespace.ValueString())
	if err != nil {
		const shortErr = "Error reading Backstage Domain kind"
		longErr := fmt.Sprintf("Could not read Backstage Domain kind %s/%s: %s", state.Namespace.ValueString(), state.Name.ValueString(), err.Error())
		if state.Fallback == nil {
			resp.Diagnostics.AddError(shortErr, longErr)
			return
		}
		resp.Diagnostics.AddWarning(shortErr, longErr)
	}

	if response.StatusCode != http.StatusOK {
		const shortErr = "Error reading Backstage Domain kind"
		longErr := fmt.Sprintf("Could not read Backstage Domain kind %s/%s: %s", state.Namespace.ValueString(), state.Name.ValueString(), response.Status)
		if state.Fallback == nil {
			resp.Diagnostics.AddError(shortErr, longErr)
			return
		}
		resp.Diagnostics.AddWarning(shortErr, longErr)
	}

	if (err != nil || response.StatusCode != http.StatusOK) && state.Fallback != nil {
		if state.Fallback.ID.IsNull() {
			state.Fallback.ID = types.StringValue("123456789")
		}
		if state.Fallback.ApiVersion.IsNull() {
			state.Fallback.ApiVersion = types.StringValue("backstage.io/v1alpha1")
		}
		if state.Fallback.Kind.IsNull() {
			state.Fallback.Kind = types.StringValue(backstage.KindDomain)
		}
		state.ID = state.Fallback.ID
		state.Name = state.Fallback.Name
		state.Namespace = state.Fallback.Namespace
		state.ApiVersion = state.Fallback.ApiVersion
		state.Kind = state.Fallback.Kind
		state.Metadata = state.Fallback.Metadata
		state.Relations = state.Fallback.Relations
		state.Spec = state.Fallback.Spec
	}

	if err == nil && response.StatusCode == http.StatusOK {
		state.ID = types.StringValue(domain.Metadata.UID)
		state.ApiVersion = types.StringValue(domain.ApiVersion)
		state.Kind = types.StringValue(domain.Kind)

		for _, i := range domain.Relations {
			state.Relations = append(state.Relations, entityRelationModel{
				Type:      types.StringValue(i.Type),
				TargetRef: types.StringValue(i.TargetRef),
				Target: &entityRelationTargetModel{
					Kind:      types.StringValue(i.Target.Kind),
					Name:      types.StringValue(i.Target.Name),
					Namespace: types.StringValue(i.Target.Namespace)},
			})
		}

		state.Spec = &domainSpecModel{
			Owner: types.StringValue(domain.Spec.Owner),
		}

		state.Metadata = &entityMetadataModel{
			UID:         types.StringValue(domain.Metadata.UID),
			Etag:        types.StringValue(domain.Metadata.Etag),
			Name:        types.StringValue(domain.Metadata.Name),
			Namespace:   types.StringValue(domain.Metadata.Namespace),
			Title:       types.StringValue(domain.Metadata.Title),
			Description: types.StringValue(domain.Metadata.Description),
			Annotations: map[string]string{},
			Labels:      map[string]string{},
		}

		for k, v := range domain.Metadata.Labels {
			state.Metadata.Labels[k] = v
		}

		for k, v := range domain.Metadata.Annotations {
			state.Metadata.Annotations[k] = v
		}

		for _, v := range domain.Metadata.Tags {
			state.Metadata.Tags = append(state.Metadata.Tags, types.StringValue(v))
		}

		for _, v := range domain.Metadata.Links {
			state.Metadata.Links = append(state.Metadata.Links, entityLinkModel{
				URL:   types.StringValue(v.URL),
				Title: types.StringValue(v.Title),
				Icon:  types.StringValue(v.Icon),
				Type:  types.StringValue(v.Type),
			})
		}
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
