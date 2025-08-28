package netbox

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	log "github.com/sirupsen/logrus"
)

func TestAccNetboxConsoleServerPortTemplate_basic(t *testing.T) {
	testSlug := "console_server_port_template"
	testName := testAccGetTestName(testSlug)
	randomSlug := testAccGetTestName(testSlug)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
	name = "%[1]s"
}

resource "netbox_device_type" "test" {
	model = "%[1]s"
	slug = "%[2]s"
	part_number = "%[2]s"
	manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_console_server_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "rj-45"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "type", "rj-45"),
					resource.TestCheckResourceAttrPair("netbox_console_server_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_console_server_port_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxConsoleServerPortTemplate_opts(t *testing.T) {
	testSlug := "console_server_port_template"
	testName := testAccGetTestName(testSlug)
	randomSlug := testAccGetTestName(testSlug)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
	name = "%[1]s"
}

resource "netbox_device_type" "test" {
	model = "%[1]s"
	slug = "%[2]s"
	part_number = "%[2]s"
	manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_console_server_port_template" "test" {
	name = "%[1]s"
	description = "%[1]s description"
	label = "%[1]s label"
	device_type_id = netbox_device_type.test.id
	type = "de-9"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "description", fmt.Sprintf("%s description", testName)),
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "label", fmt.Sprintf("%s label", testName)),
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "type", "de-9"),
					resource.TestCheckResourceAttrPair("netbox_console_server_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_console_server_port_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxConsoleServerPortTemplate_moduleType(t *testing.T) {
	testSlug := "console_server_port_template"
	testName := testAccGetTestName(testSlug)
	randomSlug := testAccGetTestName(testSlug)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
	name = "%[1]s"
}

resource "netbox_module_type" "test" {
	model = "%[1]s"
	manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_console_server_port_template" "test" {
	name = "%[1]s"
	module_type_id = netbox_module_type.test.id
	type = "usb-a"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "type", "usb-a"),
					resource.TestCheckResourceAttrPair("netbox_console_server_port_template.test", "module_type_id", "netbox_module_type.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxConsoleServerPortTemplate_serialTypes(t *testing.T) {
	testSlug := "console_server_port_template"
	testName := testAccGetTestName(testSlug)
	randomSlug := testAccGetTestName(testSlug)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
	name = "%[1]s"
}

resource "netbox_device_type" "test" {
	model = "%[1]s"
	slug = "%[2]s"
	part_number = "%[2]s"
	manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_console_server_port_template" "test_de9" {
	name = "%[1]s_de9"
	device_type_id = netbox_device_type.test.id
	type = "de-9"
}

resource "netbox_console_server_port_template" "test_db25" {
	name = "%[1]s_db25"
	device_type_id = netbox_device_type.test.id
	type = "db-25"
}

resource "netbox_console_server_port_template" "test_rj11" {
	name = "%[1]s_rj11"
	device_type_id = netbox_device_type.test.id
	type = "rj-11"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test_de9", "type", "de-9"),
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test_db25", "type", "db-25"),
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test_rj11", "type", "rj-11"),
				),
			},
		},
	})
}

func TestAccNetboxConsoleServerPortTemplate_usbTypes(t *testing.T) {
	testSlug := "console_server_port_template"
	testName := testAccGetTestName(testSlug)
	randomSlug := testAccGetTestName(testSlug)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
	name = "%[1]s"
}

resource "netbox_device_type" "test" {
	model = "%[1]s"
	slug = "%[2]s"
	part_number = "%[2]s"
	manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_console_server_port_template" "test_usb_a" {
	name = "%[1]s_usb_a"
	device_type_id = netbox_device_type.test.id
	type = "usb-a"
}

resource "netbox_console_server_port_template" "test_usb_c" {
	name = "%[1]s_usb_c"
	device_type_id = netbox_device_type.test.id
	type = "usb-c"
}

resource "netbox_console_server_port_template" "test_usb_micro_b" {
	name = "%[1]s_usb_micro_b"
	device_type_id = netbox_device_type.test.id
	type = "usb-micro-b"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test_usb_a", "type", "usb-a"),
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test_usb_c", "type", "usb-c"),
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test_usb_micro_b", "type", "usb-micro-b"),
				),
			},
		},
	})
}

func TestAccNetboxConsoleServerPortTemplate_withoutType(t *testing.T) {
	testSlug := "console_server_port_template"
	testName := testAccGetTestName(testSlug)
	randomSlug := testAccGetTestName(testSlug)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
	name = "%[1]s"
}

resource "netbox_device_type" "test" {
	model = "%[1]s"
	slug = "%[2]s"
	part_number = "%[2]s"
	manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_console_server_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	# No type specified - should work
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", testName),
					resource.TestCheckResourceAttrPair("netbox_console_server_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
		},
	})
}

func init() {
	resource.AddTestSweepers("netbox_console_server_port_template", &resource.Sweeper{
		Name:         "netbox_console_server_port_template",
		Dependencies: []string{},
		F: func(r string) error {
			return sweepConsoleServerPortTemplates(r)
		},
	})
}

func sweepConsoleServerPortTemplates(region string) error {
	log.Infof("[INFO] sweeping netbox console server port templates in region %s", region)
	api := testAccProvider.Meta().(*providerState).NetBoxAPI

	res, err := api.Dcim.DcimConsoleServerPortTemplatesList(dcim.NewDcimConsoleServerPortTemplatesListParams(), nil)
	if err != nil {
		return err
	}

	for _, item := range res.GetPayload().Results {
		if strings.HasPrefix(*item.Name, testPrefix) {
			log.Infof("[INFO] deleting console server port template %d", item.ID)
			if err := deleteConsoleServerPortTemplate(item.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func deleteConsoleServerPortTemplate(id int64) error {
	api := testAccProvider.Meta().(*providerState).NetBoxAPI
	_, err := api.Dcim.DcimConsoleServerPortTemplatesDelete(
		dcim.NewDcimConsoleServerPortTemplatesDeleteParams().WithID(id), nil)
	return err
}
