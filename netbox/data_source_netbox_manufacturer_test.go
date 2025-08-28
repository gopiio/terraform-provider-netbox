package netbox

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetboxManufacturerDataSource_basic(t *testing.T) {
	testSlug := "mfg_ds_basic"
	testName := testAccGetTestName(testSlug)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[1]s"
}

data "netbox_manufacturer" "by_name" {
  depends_on = [netbox_manufacturer.test]
  name = "%[1]s"
}
`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.by_name", "id", "netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.by_name", "name", "netbox_manufacturer.test", "name"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.by_name", "name", testName),
				),
			},
		},
	})
}

func TestAccNetboxManufacturerDataSource_withSlug(t *testing.T) {
	testSlug := "mfg_ds_slug"
	testName := testAccGetTestName(testSlug)
	customSlug := testAccGetTestName("custom-slug")
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

data "netbox_manufacturer" "by_name" {
  depends_on = [netbox_manufacturer.test]
  name = "%[1]s"
}
`, testName, customSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.by_name", "id", "netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.by_name", "name", "netbox_manufacturer.test", "name"),
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.by_name", "slug", "netbox_manufacturer.test", "slug"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.by_name", "name", testName),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.by_name", "slug", customSlug),
				),
			},
		},
	})
}

func TestAccNetboxManufacturerDataSource_generatedSlug(t *testing.T) {
	testSlug := "mfg_ds_generated"
	testName := testAccGetTestName(testSlug)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[1]s"
  # slug will be auto-generated
}

data "netbox_manufacturer" "by_name" {
  depends_on = [netbox_manufacturer.test]
  name = "%[1]s"
}
`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.by_name", "id", "netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.by_name", "name", "netbox_manufacturer.test", "name"),
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.by_name", "slug", "netbox_manufacturer.test", "slug"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.by_name", "name", testName),
					resource.TestCheckResourceAttrSet("data.netbox_manufacturer.by_name", "slug"),
				),
			},
		},
	})
}

func TestAccNetboxManufacturerDataSource_exactMatch(t *testing.T) {
	testSlug := "mfg_ds_exact"
	testName := testAccGetTestName(testSlug)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test_1" {
  name = "%[1]s_cisco"
}

resource "netbox_manufacturer" "test_2" {
  name = "%[1]s_cisco_systems"
}

# This should match exactly the first one
data "netbox_manufacturer" "exact" {
  depends_on = [netbox_manufacturer.test_1, netbox_manufacturer.test_2]
  name = "%[1]s_cisco"
}
`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.exact", "id", "netbox_manufacturer.test_1", "id"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.exact", "name", fmt.Sprintf("%s_cisco", testName)),
				),
			},
		},
	})
}

func TestAccNetboxManufacturerDataSource_notFound(t *testing.T) {
	testSlug := "mfg_ds_notfound"
	testName := testAccGetTestName(testSlug)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "netbox_manufacturer" "not_found" {
  name = "%[1]s_nonexistent"
}
`, testName),
				ExpectError: regexp.MustCompile("no manufacturer found matching filter"),
			},
		},
	})
}

func TestAccNetboxManufacturerDataSource_withSpecificName(t *testing.T) {
	testSlug := "mfg_ds_specific"
	testName := testAccGetTestName(testSlug)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[1]s"
}

# Query by name to get exactly the manufacturer we created
data "netbox_manufacturer" "specific" {
  depends_on = [netbox_manufacturer.test]
  name = "%[1]s"
}
`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_manufacturer.specific", "id"),
					resource.TestCheckResourceAttrSet("data.netbox_manufacturer.specific", "name"),
					resource.TestCheckResourceAttrSet("data.netbox_manufacturer.specific", "slug"),
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.specific", "id", "netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.specific", "name", testName),
				),
			},
		},
	})
}

func TestAccNetboxManufacturerDataSource_multipleError(t *testing.T) {
	testSlug := "mfg_ds_multi_err"
	testName := testAccGetTestName(testSlug)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test_1" {
  name = "%[1]s_cisco"
}

resource "netbox_manufacturer" "test_2" {
  name = "%[1]s_dell"
}

# This should fail as no filters are provided and multiple manufacturers exist
data "netbox_manufacturer" "multiple" {
  depends_on = [netbox_manufacturer.test_1, netbox_manufacturer.test_2]
}
`, testName),
				ExpectError: regexp.MustCompile("more than one manufacturer returned"),
			},
		},
	})
}
