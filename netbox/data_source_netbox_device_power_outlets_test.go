package netbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetboxDevicePowerOutletsDataSource_basic(t *testing.T) {
	testSlug := "device_power_outlets_ds_basic"
	testName := testAccGetTestName(testSlug)
	dependencies := testAccNetboxDevicePowerOutletsDataSourceDependencies(testName)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dependencies,
			},
			{
				Config: dependencies + fmt.Sprintf(`
data "netbox_device_power_outlets" "by_name" {
  filter {
    name = "name"
    value = "%[1]s"
  }
}

data "netbox_device_power_outlets" "by_tag" {
  filter {
    name = "tag"
    value = "%[1]s"
  }
}

data "netbox_device_power_outlets" "by_device_id" {
  filter {
    name = "device_id"
    value = netbox_device.test.id
  }
}

data "netbox_device_power_outlets" "by_type" {
  filter {
    name = "type"
    value = "iec-60320-c13"
  }
}
`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_name", "power_outlets.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_name", "power_outlets.0.name", testName),
					resource.TestCheckResourceAttrPair("data.netbox_device_power_outlets.by_name", "power_outlets.0.id", "netbox_device_power_outlet.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_device_power_outlets.by_name", "power_outlets.0.device_id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_name", "power_outlets.0.tag_ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_name", "power_outlets.0.type", "iec-60320-c13"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_name", "power_outlets.0.mark_connected", "false"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_tag", "power_outlets.#", "3"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_device_id", "power_outlets.#", "3"),
					resource.TestCheckResourceAttrPair("data.netbox_device_power_outlets.by_device_id", "power_outlets.0.device_id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_device_id", "power_outlets.0.mark_connected", "false"),
					resource.TestCheckResourceAttrPair("data.netbox_device_power_outlets.by_device_id", "power_outlets.1.device_id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_device_id", "power_outlets.1.mark_connected", "false"),
					resource.TestCheckResourceAttrPair("data.netbox_device_power_outlets.by_device_id", "power_outlets.2.device_id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_device_id", "power_outlets.2.mark_connected", "false"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.by_type", "power_outlets.#", "2"),
				),
			},
			{
				Config: dependencies + testAccNetboxDevicePowerOutletsDataSourceNameRegex(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.regex_test", "power_outlets.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.regex_test", "power_outlets.0.name", testName+"_regex"),
					resource.TestCheckResourceAttrPair("data.netbox_device_power_outlets.regex_test", "power_outlets.0.device_id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.regex_test", "power_outlets.0.tag_ids.#", "0"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.regex_test", "power_outlets.0.mark_connected", "false"),
				),
			},
			{
				Config: dependencies + testAccNetboxDevicePowerOutletsDataSourceWithLimit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.limited", "power_outlets.#", "1"),
				),
			},
			{
				Config: dependencies + testAccNetboxDevicePowerOutletsDataSourceDeviceFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.device2_only", "power_outlets.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.device2_only", "power_outlets.0.name", testName+"_device2"),
					resource.TestCheckResourceAttrPair("data.netbox_device_power_outlets.device2_only", "power_outlets.0.device_id", "netbox_device.test2", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_power_outlets.device2_only", "power_outlets.0.mark_connected", "false"),
				),
			},
		},
	})
}

func testAccNetboxDevicePowerOutletsDataSourceDependencies(testName string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = "%[1]s"
}

resource "netbox_site" "test" {
  name = "%[1]s"
  status = "active"
}

resource "netbox_device_role" "test" {
  name = "%[1]s"
  color_hex = "123456"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s"
}

resource "netbox_device_type" "test" {
  model = "%[1]s"
  manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name = "%[1]s"
  device_type_id = netbox_device_type.test.id
  role_id = netbox_device_role.test.id
  site_id = netbox_site.test.id
}

resource "netbox_device" "test2" {
  name = "%[1]s_2"
  device_type_id = netbox_device_type.test.id
  role_id = netbox_device_role.test.id
  site_id = netbox_site.test.id
}

resource "netbox_device_power_port" "test" {
  name = "%[1]s_port"
  device_id = netbox_device.test.id
  type = "iec-60320-c19"
}

resource "netbox_device_power_outlet" "test" {
  name = "%[1]s"
  device_id = netbox_device.test.id
  type = "iec-60320-c13"
  power_port_id = netbox_device_power_port.test.id
  tags = ["%[1]s"]
}

resource "netbox_device_power_outlet" "test2" {
  name = "%[1]s_two"
  device_id = netbox_device.test.id
  type = "iec-60320-c13"
  tags = ["%[1]s"]
}

resource "netbox_device_power_outlet" "test_regex" {
  name = "%[1]s_regex"
  device_id = netbox_device.test.id
  type = "nema-5-15r"
}

resource "netbox_device_power_outlet" "test_device2" {
  name = "%[1]s_device2"
  device_id = netbox_device.test2.id
  type = "nema-5-20r"
  tags = ["%[1]s"]
}
`, testName)
}

func testAccNetboxDevicePowerOutletsDataSourceNameRegex(testName string) string {
	return fmt.Sprintf(`
data "netbox_device_power_outlets" "regex_test" {
  name_regex = "%s_regex"
}
`, testName)
}

const testAccNetboxDevicePowerOutletsDataSourceWithLimit = `
data "netbox_device_power_outlets" "limited" {
  filter {
    name = "tag"
    value = netbox_tag.test.name
  }
  limit = 1
}
`

const testAccNetboxDevicePowerOutletsDataSourceDeviceFilter = `
data "netbox_device_power_outlets" "device2_only" {
  filter {
    name = "device_id"
    value = netbox_device.test2.id
  }
}
`
