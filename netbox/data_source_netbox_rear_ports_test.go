package netbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestNetboxRearPortsDataSource(t *testing.T) {
	dataSource := dataSourceNetboxRearPorts()

	// Test that the data source is properly configured
	if dataSource.Read == nil {
		t.Fatal("Expected Read function to be set")
	}

	// Test required schema fields
	deviceIDSchema, ok := dataSource.Schema["device_id"]
	if !ok {
		t.Fatal("Expected device_id to be in schema")
	}
	if !deviceIDSchema.Required {
		t.Error("Expected device_id to be required")
	}
	if deviceIDSchema.Type != schema.TypeInt {
		t.Error("Expected device_id to be TypeInt")
	}

	// Test optional schema fields
	filterSchema, ok := dataSource.Schema["filter"]
	if !ok {
		t.Fatal("Expected filter to be in schema")
	}
	if filterSchema.Required {
		t.Error("Expected filter to be optional")
	}

	limitSchema, ok := dataSource.Schema["limit"]
	if !ok {
		t.Fatal("Expected limit to be in schema")
	}
	if limitSchema.Required {
		t.Error("Expected limit to be optional")
	}

	// Test computed schema fields
	rearPortsSchema, ok := dataSource.Schema["rear_ports"]
	if !ok {
		t.Fatal("Expected rear_ports to be in schema")
	}
	if !rearPortsSchema.Computed {
		t.Error("Expected rear_ports to be computed")
	}
	if rearPortsSchema.Type != schema.TypeList {
		t.Error("Expected rear_ports to be TypeList")
	}
}

func TestAccNetboxRearPortsDataSource_basic(t *testing.T) {
	testSlug := "rear_ports_ds_basic"
	testName := testAccGetTestName(testSlug)
	dependencies := testAccNetboxRearPortsDataSourceDependencies(testName)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dependencies,
			},
			{
				Config: dependencies + fmt.Sprintf(`
data "netbox_rear_ports" "by_device_id" {
  device_id = netbox_device.test.id
}

data "netbox_rear_ports" "by_name_filter" {
  device_id = netbox_device.test.id
  filter {
    name = "name"
    value = "%[1]s"
  }
}

data "netbox_rear_ports" "by_type_filter" {
  device_id = netbox_device.test.id
  filter {
    name = "type"
    value = "8p8c"
  }
}

data "netbox_rear_ports" "by_tag_filter" {
  device_id = netbox_device.test.id
  filter {
    name = "tag"
    value = "%[1]s"
  }
}
`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_device_id", "rear_ports.#", "2"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_filter", "rear_ports.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_filter", "rear_ports.0.name", testName),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_filter", "rear_ports.0.type", "8p8c"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_filter", "rear_ports.0.description", testName+"_description"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_filter", "rear_ports.0.positions", "1"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_filter", "rear_ports.0.label", testName+"_label"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_filter", "rear_ports.0.color_hex", "123456"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_filter", "rear_ports.0.mark_connected", "true"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_filter", "rear_ports.0.tag_ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_type_filter", "rear_ports.#", "2"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_tag_filter", "rear_ports.#", "2"),
				),
			},
		},
	})
}

func TestAccNetboxRearPortsDataSource_nameRegex(t *testing.T) {
	testSlug := "rear_ports_ds_regex"
	testName := testAccGetTestName(testSlug)
	dependencies := testAccNetboxRearPortsDataSourceDependencies(testName)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dependencies,
			},
			{
				Config: dependencies + fmt.Sprintf(`
data "netbox_rear_ports" "by_name_regex" {
  device_id = netbox_device.test.id
  name_regex = "%[1]s.*"
}

data "netbox_rear_ports" "by_name_regex_specific" {
  device_id = netbox_device.test.id
  name_regex = "%[1]s$"
}
`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_regex", "rear_ports.#", "2"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_regex_specific", "rear_ports.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_rear_ports.by_name_regex_specific", "rear_ports.0.name", testName),
				),
			},
		},
	})
}

func TestAccNetboxRearPortsDataSource_limit(t *testing.T) {
	testSlug := "rear_ports_ds_limit"
	testName := testAccGetTestName(testSlug)
	dependencies := testAccNetboxRearPortsDataSourceDependencies(testName)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dependencies,
			},
			{
				Config: dependencies + `
data "netbox_rear_ports" "limited" {
  device_id = netbox_device.test.id
  limit = 1
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rear_ports.limited", "rear_ports.#", "1"),
				),
			},
		},
	})
}

func testAccNetboxRearPortsDataSourceDependencies(testName string) string {
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

resource "netbox_device_rear_port" "test" {
  device_id = netbox_device.test.id
  name = "%[1]s"
  type = "8p8c"
  positions = 1
  description = "%[1]s_description"
  label = "%[1]s_label"
  color_hex = "123456"
  mark_connected = true
  tags = ["%[1]s"]
}

resource "netbox_device_rear_port" "test2" {
  device_id = netbox_device.test.id
  name = "%[1]s_two"
  type = "8p8c"
  positions = 1
  description = "%[1]s_two_description"
  label = "%[1]s_two_label"
  color_hex = "654321"
  mark_connected = false
  tags = ["%[1]s"]
}
`, testName)
}
