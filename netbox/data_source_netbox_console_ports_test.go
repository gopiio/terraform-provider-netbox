package netbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetboxConsolePortsDataSource_basic(t *testing.T) {
	testSlug := "console_ports_ds_basic"
	testName := testAccGetTestName(testSlug)
	dependencies := testAccNetboxConsolePortsDataSourceDependencies(testName)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dependencies,
			},
			{
				Config: dependencies + fmt.Sprintf(`
data "netbox_console_ports" "by_name" {
  filter {
    name = "name"
    value = "%[1]s"
  }
}

data "netbox_console_ports" "by_tag" {
  filter {
    name = "tag"
    value = "%[1]s"
  }
}

data "netbox_console_ports" "by_device_id" {
  filter {
    name = "device_id"
    value = netbox_device.test.id
  }
}
`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_console_ports.by_name", "console_ports.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_console_ports.by_name", "console_ports.0.name", testName),
					resource.TestCheckResourceAttrPair("data.netbox_console_ports.by_name", "console_ports.0.id", "netbox_device_console_port.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_console_ports.by_name", "console_ports.0.device_id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_console_ports.by_name", "console_ports.0.tag_ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_console_ports.by_tag", "console_ports.#", "3"),
					resource.TestCheckResourceAttr("data.netbox_console_ports.by_device_id", "console_ports.#", "3"),
					resource.TestCheckResourceAttrPair("data.netbox_console_ports.by_device_id", "console_ports.0.device_id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_console_ports.by_device_id", "console_ports.1.device_id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_console_ports.by_device_id", "console_ports.2.device_id", "netbox_device.test", "id"),
				),
			},
			{
				Config: dependencies + testAccNetboxConsolePortsDataSourceNameRegex(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_console_ports.regex_test", "console_ports.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_console_ports.regex_test", "console_ports.0.name", testName+"_regex"),
					resource.TestCheckResourceAttrPair("data.netbox_console_ports.regex_test", "console_ports.0.device_id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_console_ports.regex_test", "console_ports.0.tag_ids.#", "0"),
				),
			},
			{
				Config: dependencies + testAccNetboxConsolePortsDataSourceWithLimit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_console_ports.limited", "console_ports.#", "1"),
				),
			},
			{
				Config: dependencies + testAccNetboxConsolePortsDataSourceDeviceFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_console_ports.device2_only", "console_ports.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_console_ports.device2_only", "console_ports.0.name", testName+"_device2"),
					resource.TestCheckResourceAttrPair("data.netbox_console_ports.device2_only", "console_ports.0.device_id", "netbox_device.test2", "id"),
				),
			},
		},
	})
}

func testAccNetboxConsolePortsDataSourceDependencies(testName string) string {
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

resource "netbox_device_console_port" "test" {
  name = "%[1]s"
  device_id = netbox_device.test.id
  type = "rj-45"
  tags = ["%[1]s"]
}

resource "netbox_device_console_port" "test2" {
  name = "%[1]s_two"
  device_id = netbox_device.test.id
  type = "usb-a"
  tags = ["%[1]s"]
}

resource "netbox_device_console_port" "test_regex" {
  name = "%[1]s_regex"
  device_id = netbox_device.test.id
  type = "de-9"
}

resource "netbox_device_console_port" "test_device2" {
  name = "%[1]s_device2"
  device_id = netbox_device.test2.id
  type = "usb-c"
  tags = ["%[1]s"]
}
`, testName)
}

func testAccNetboxConsolePortsDataSourceNameRegex(testName string) string {
	return fmt.Sprintf(`
data "netbox_console_ports" "regex_test" {
  name_regex = "%s_regex"
}
`, testName)
}

const testAccNetboxConsolePortsDataSourceWithLimit = `
data "netbox_console_ports" "limited" {
  filter {
    name = "tag"
    value = netbox_tag.test.name
  }
  limit = 1
}
`

const testAccNetboxConsolePortsDataSourceDeviceFilter = `
data "netbox_console_ports" "device2_only" {
  filter {
    name = "device_id"
    value = netbox_device.test2.id
  }
}
`
