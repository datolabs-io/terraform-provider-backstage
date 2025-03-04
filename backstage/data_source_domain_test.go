package backstage

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + testAccDataSourceDomainConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.backstage_domain.test", "api_version", "backstage.io/v1alpha1"),
					resource.TestCheckResourceAttr("data.backstage_domain.test", "kind", "Domain"),
					resource.TestCheckResourceAttr("data.backstage_domain.test", "metadata.annotations.backstage.io/source-location",
						"url:https://github.com/backstage/backstage/tree/master/packages/catalog-model/examples/domains/"),
					resource.TestCheckResourceAttr("data.backstage_domain.test", "metadata.description", "Everything related to artists"),
					resource.TestCheckResourceAttr("data.backstage_domain.test", "relations.0.target_ref", "system:default/artist-engagement-portal"),
					resource.TestCheckResourceAttr("data.backstage_domain.test", "spec.owner", "team-a"),
				),
			},
		},
	})
}

const testAccDataSourceDomainConfig = `
data "backstage_domain" "test" {
  name = "artists"
}
`

func TestAccDataSourceDomain_WithFallback(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "backstage_domain" "test" {
						name = "non_existent_domain_a9ab8"
						namespace = "default"
						fallback = {
							id = "123456"
							name = "fallback_domain"
							kind = "Domain"
							namespace = "fallback_default"
							spec = {
								owner = "team-a"
							}
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.backstage_domain.test", "kind", "Domain"),
					resource.TestCheckResourceAttr("data.backstage_domain.test", "spec.owner", "team-a"),
					resource.TestCheckResourceAttr("data.backstage_domain.test", "name", "fallback_domain"),
					resource.TestCheckResourceAttr("data.backstage_domain.test", "namespace", "fallback_default"),
					resource.TestCheckNoResourceAttr("data.backstage_domain.test", "metadata"),
					resource.TestCheckNoResourceAttr("data.backstage_domain.test", "relations"),
				),
			},
		},
	})
}
