package netbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccNetboxDeviceTypesDataSourceDependencies(testName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[1]s"
}

resource "netbox_manufacturer" "test2" {
  name = "%[1]s_2"
}

resource "netbox_tag" "test" {
  name = "%[1]sa"
}

resource "netbox_device_type" "test1" {
  model           = "%[1]s_model1"
  slug            = "%[1]s-model1"
  manufacturer_id = netbox_manufacturer.test.id
  part_number     = "%[1]s_part1"
  tags            = [netbox_tag.test.name]
}

resource "netbox_device_type" "test2" {
  model           = "%[1]s_model2"
  slug            = "%[1]s-model2"
  manufacturer_id = netbox_manufacturer.test.id
  part_number     = "%[1]s_part2"
}

resource "netbox_device_type" "test3" {
  model           = "%[1]s_model3_regex"
  slug            = "%[1]s-model3-regex"
  manufacturer_id = netbox_manufacturer.test2.id
}
`, testName)
}

func TestAccNetboxDeviceTypesDataSource_basic(t *testing.T) {
	testSlug := "dev_types_ds"
	testName := testAccGetTestName(testSlug)
	dependencies := testAccNetboxDeviceTypesDataSourceDependencies(testName)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dependencies,
			},
			{
				Config: dependencies + testAccNetboxDeviceTypesDataSourceFilterModel(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_types.by_model", "device_types.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_device_types.by_model", "device_types.0.model", testName+"_model1"),
					resource.TestCheckResourceAttr("data.netbox_device_types.by_model", "device_types.0.slug", testName+"-model1"),
					resource.TestCheckResourceAttr("data.netbox_device_types.by_model", "device_types.0.part_number", testName+"_part1"),
					resource.TestCheckResourceAttrPair("data.netbox_device_types.by_model", "device_types.0.manufacturer_id", "netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_types.by_model", "device_types.0.manufacturer_name", testName),
					resource.TestCheckResourceAttr("data.netbox_device_types.by_model", "device_types.0.tags.#", "1"),
				),
			},
			{
				Config: dependencies + testAccNetboxDeviceTypesDataSourceFilterManufacturer(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_types.by_manufacturer", "device_types.#", "2"),
				),
			},
			{
				Config: dependencies + testAccNetboxDeviceTypesDataSourceNameRegex(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_types.by_regex", "device_types.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_device_types.by_regex", "device_types.0.model", testName+"_model3_regex"),
				),
			},
			{
				Config: dependencies + testAccNetboxDeviceTypesDataSourceLimit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_types.limited", "device_types.#", "1"),
				),
			},
			{
				Config: dependencies + testAccNetboxDeviceTypesDataSourceFilterSlug(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_types.by_slug", "device_types.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_device_types.by_slug", "device_types.0.slug", testName+"-model2"),
					resource.TestCheckResourceAttr("data.netbox_device_types.by_slug", "device_types.0.model", testName+"_model2"),
				),
			},
		},
	})
}

func testAccNetboxDeviceTypesDataSourceFilterModel(testName string) string {
	return fmt.Sprintf(`
data "netbox_device_types" "by_model" {
  filter {
    name  = "model"
    value = "%[1]s_model1"
  }
}`, testName)
}

func testAccNetboxDeviceTypesDataSourceFilterManufacturer(testName string) string {
	return `
data "netbox_device_types" "by_manufacturer" {
  filter {
    name  = "manufacturer_id"
    value = netbox_manufacturer.test.id
  }
}`
}

func testAccNetboxDeviceTypesDataSourceNameRegex(testName string) string {
	return fmt.Sprintf(`
data "netbox_device_types" "by_regex" {
  name_regex = "%[1]s.*_regex"
}`, testName)
}

const testAccNetboxDeviceTypesDataSourceLimit = `
data "netbox_device_types" "limited" {
  limit = 1
}`

func testAccNetboxDeviceTypesDataSourceFilterSlug(testName string) string {
	return fmt.Sprintf(`
data "netbox_device_types" "by_slug" {
  filter {
    name  = "slug"
    value = "%[1]s-model2"
  }
}`, testName)
}
