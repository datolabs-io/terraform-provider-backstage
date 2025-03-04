package backstage

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceApi(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + testAccDataSourceApiConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.backstage_api.test", "api_version", "backstage.io/v1alpha1"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "kind", "API"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "metadata.annotations.backstage.io/edit-url",
						"https://github.com/backstage/backstage/edit/master/packages/catalog-model/examples/apis/streetlights-api.yaml"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "metadata.description",
						"The Smartylighting Streetlights API allows you to remotely manage the city lights."),
					resource.TestCheckResourceAttr("data.backstage_api.test", "metadata.tags.0", "mqtt"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "relations.0.target_ref", "component:default/petstore"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "spec.lifecycle", "production"),
				),
			},
		},
	})
}

const testAccDataSourceApiConfig = `
data "backstage_api" "test" {
  name = "streetlights"
}
`

func TestAccApiDataSource_WithFallback(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
                    data "backstage_api" "test" {
                        name = "non_existent_api_a9ab8"
                        namespace = "default"
                        fallback = {
							id = "123456"
                            name = "fallback_api"
                            namespace = "default"
							metadata = {
								labels = {
									"key" = "value"
								}
							}
                            spec = {
                                type = "openapi"
                                lifecycle = "production"
                                owner = "team-a"
                                definition = "https://example.com/api-spec"
                                system = "system-x"
                            }
                        }
                    }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.backstage_api.test", "name", "fallback_api"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "namespace", "default"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "spec.type", "openapi"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "spec.lifecycle", "production"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "spec.owner", "team-a"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "spec.definition", "https://example.com/api-spec"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "spec.system", "system-x"),
					resource.TestCheckResourceAttr("data.backstage_api.test", "metadata.labels.key", "value"),
				),
			},
		},
	})
}
