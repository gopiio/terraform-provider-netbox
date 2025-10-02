package netbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccNetboxCablesDataSourceDependencies(testName string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = "%[1]s"
}

resource "netbox_site" "test" {
  name = "%[1]s"
  status = "active"
}

resource "netbox_site" "test2" {
  name = "%[1]s_2"
  status = "active"
}

resource "netbox_tag" "test" {
  name = "%[1]sa"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s"
}

resource "netbox_device_type" "test" {
  model = "%[1]s"
  manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = "%[1]s"
  color_hex = "123456"
}

resource "netbox_rack_role" "test" {
  name = "%[1]s"
  color_hex = "123456"
}

resource "netbox_rack" "test" {
  name = "%[1]s"
  site_id = netbox_site.test.id
  role_id = netbox_rack_role.test.id
  status = "active"
  width = 19
  u_height = 42
}

resource "netbox_rack" "test2" {
  name = "%[1]s_2"
  site_id = netbox_site.test2.id
  role_id = netbox_rack_role.test.id
  status = "active"
  width = 19
  u_height = 42
}

resource "netbox_device" "test" {
  name = "%[1]s"
  device_type_id = netbox_device_type.test.id
  tenant_id = netbox_tenant.test.id
  role_id = netbox_device_role.test.id
  site_id = netbox_site.test.id
  rack_id = netbox_rack.test.id
  rack_position = 1
  rack_face = "front"
}

resource "netbox_device" "test2" {
  name = "%[1]s_2"
  device_type_id = netbox_device_type.test.id
  tenant_id = netbox_tenant.test.id
  role_id = netbox_device_role.test.id
  site_id = netbox_site.test2.id
  rack_id = netbox_rack.test2.id
  rack_position = 1
  rack_face = "front"
}

resource "netbox_device_console_port" "test1" {
  device_id = netbox_device.test.id
  name = "%[1]s1"
}

resource "netbox_device_console_port" "test2" {
  device_id = netbox_device.test2.id
  name = "%[1]s2"
}

resource "netbox_device_console_server_port" "test1" {
  device_id = netbox_device.test.id
  name = "%[1]s1"
}

resource "netbox_device_console_server_port" "test2" {
  device_id = netbox_device.test2.id
  name = "%[1]s2"
}

resource "netbox_cable" "test1" {
  a_termination {
    object_type = "dcim.consoleserverport"
    object_id = netbox_device_console_server_port.test1.id
  }

  b_termination {
    object_type = "dcim.consoleport"
    object_id = netbox_device_console_port.test1.id
  }

  status = "connected"
  label = "%[1]s_cable1"
  type = "cat6"
  tenant_id = netbox_tenant.test.id
  color_hex = "ff0000"
  length = 5
  length_unit = "m"
  description = "%[1]s_description1"
  comments = "%[1]s_comments1"
  tags = ["%[1]sa"]
}

resource "netbox_cable" "test2" {
  a_termination {
    object_type = "dcim.consoleserverport"
    object_id = netbox_device_console_server_port.test2.id
  }

  b_termination {
    object_type = "dcim.consoleport"
    object_id = netbox_device_console_port.test2.id
  }

  status = "planned"
  label = "%[1]s_cable2"
  type = "cat5e"
  tenant_id = netbox_tenant.test.id
  color_hex = "00ff00"
  length = 10
  length_unit = "ft"
  description = "%[1]s_description2"
  comments = "%[1]s_comments2"
  tags = ["%[1]sa"]
}
`, testName)
}

func TestAccNetboxCablesDataSource_basic(t *testing.T) {
	testSlug := "cables_ds_basic"
	testName := testAccGetTestName(testSlug)
	dependencies := testAccNetboxCablesDataSourceDependencies(testName)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dependencies,
			},
			{
				Config: dependencies + testAccNetboxCablesDataSourceFilterStatus(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.status", "connected"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.label", testName+"_cable1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.type", "cat6"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.color_hex", "ff0000"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.length", "5"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.length_unit", "m"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.description", testName+"_description1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.comments", testName+"_comments1"),
					resource.TestCheckResourceAttrPair("data.netbox_cables.by_status", "cables.0.tenant_id", "netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.tag_ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.a_terminations.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.b_terminations.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.a_terminations.0.object_type", "dcim.consoleserverport"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_status", "cables.0.b_terminations.0.object_type", "dcim.consoleport"),
				),
			},
			{
				Config: dependencies + testAccNetboxCablesDataSourceFilterType(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cables.by_type", "cables.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_type", "cables.0.type", "cat5e"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_type", "cables.0.label", testName+"_cable2"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_type", "cables.0.status", "planned"),
				),
			},
			{
				Config: dependencies + testAccNetboxCablesDataSourceFilterSite(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cables.by_site", "cables.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_site", "cables.0.label", testName+"_cable2"),
				),
			},
			{
				Config: dependencies + testAccNetboxCablesDataSourceFilterDevice(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cables.by_device", "cables.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_device", "cables.0.label", testName+"_cable1"),
				),
			},
			{
				Config: dependencies + testAccNetboxCablesDataSourceFilterRack(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cables.by_rack", "cables.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_rack", "cables.0.label", testName+"_cable2"),
				),
			},
			{
				Config: dependencies + testAccNetboxCablesDataSourceFilterTenant(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cables.by_tenant", "cables.#", "2"),
					resource.TestCheckResourceAttrPair("data.netbox_cables.by_tenant", "cables.0.tenant_id", "netbox_tenant.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_cables.by_tenant", "cables.1.tenant_id", "netbox_tenant.test", "id"),
				),
			},
			{
				Config: dependencies + testAccNetboxCablesDataSourceNameRegex(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cables.by_regex", "cables.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cables.by_regex", "cables.0.label", testName+"_cable2"),
				),
			},
			{
				Config: dependencies + testAccNetboxCablesDataSourceLimit(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cables.limited", "cables.#", "1"),
				),
			},
			{
				Config: dependencies + testAccNetboxCablesDataSourceMultipleFilters(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cables.multiple_filters", "cables.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cables.multiple_filters", "cables.0.status", "connected"),
					resource.TestCheckResourceAttr("data.netbox_cables.multiple_filters", "cables.0.type", "cat6"),
				),
			},
		},
	})
}

func testAccNetboxCablesDataSourceFilterStatus(testName string) string {
	return `
data "netbox_cables" "by_status" {
  filter {
    name  = "status"
    value = "connected"
  }
}
`
}

func testAccNetboxCablesDataSourceFilterType(testName string) string {
	return `
data "netbox_cables" "by_type" {
  filter {
    name  = "type"
    value = "cat5e"
  }
}
`
}

func testAccNetboxCablesDataSourceFilterSite(testName string) string {
	return fmt.Sprintf(`
data "netbox_cables" "by_site" {
  filter {
    name  = "site_id"
    value = netbox_site.test2.id
  }
}
`)
}

func testAccNetboxCablesDataSourceFilterDevice(testName string) string {
	return fmt.Sprintf(`
data "netbox_cables" "by_device" {
  filter {
    name  = "device_id"
    value = netbox_device.test.id
  }
}
`)
}

func testAccNetboxCablesDataSourceFilterRack(testName string) string {
	return fmt.Sprintf(`
data "netbox_cables" "by_rack" {
  filter {
    name  = "rack_id"
    value = netbox_rack.test2.id
  }
}
`)
}

func testAccNetboxCablesDataSourceFilterTenant(testName string) string {
	return fmt.Sprintf(`
data "netbox_cables" "by_tenant" {
  filter {
    name  = "tenant_id"
    value = netbox_tenant.test.id
  }
}
`)
}

func testAccNetboxCablesDataSourceNameRegex(testName string) string {
	return fmt.Sprintf(`
data "netbox_cables" "by_regex" {
  name_regex = "%[1]s_cable2"
}
`, testName)
}

func testAccNetboxCablesDataSourceLimit(testName string) string {
	return `
data "netbox_cables" "limited" {
  limit = 1
}
`
}

func testAccNetboxCablesDataSourceMultipleFilters(testName string) string {
	return `
data "netbox_cables" "multiple_filters" {
  filter {
    name  = "status"
    value = "connected"
  }
  filter {
    name  = "type"
    value = "cat6"
  }
}
`
}
