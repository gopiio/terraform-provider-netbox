package netbox

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	log "github.com/sirupsen/logrus"
)

func TestAccNetboxPowerPortTemplate_basic(t *testing.T) {
	testSlug := "power_port_template"
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

resource "netbox_power_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "iec-60320-c14"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "type", "iec-60320-c14"),
					resource.TestCheckResourceAttrPair("netbox_power_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_power_port_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxPowerPortTemplate_opts(t *testing.T) {
	testSlug := "power_port_template"
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

resource "netbox_power_port_template" "test" {
	name = "%[1]s"
	description = "%[1]s description"
	label = "%[1]s label"
	device_type_id = netbox_device_type.test.id
	type = "iec-60320-c14"
	maximum_draw = 100
	allocated_draw = 50
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "description", fmt.Sprintf("%s description", testName)),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "label", fmt.Sprintf("%s label", testName)),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "type", "iec-60320-c14"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "maximum_draw", "100"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "allocated_draw", "50"),
					resource.TestCheckResourceAttrPair("netbox_power_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_power_port_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxPowerPortTemplate_moduleType(t *testing.T) {
	testSlug := "power_port_template"
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

resource "netbox_power_port_template" "test" {
	name = "%[1]s"
	module_type_id = netbox_module_type.test.id
	type = "nema-5-15p"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "type", "nema-5-15p"),
					resource.TestCheckResourceAttrPair("netbox_power_port_template.test", "module_type_id", "netbox_module_type.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxPowerPortTemplate_usb(t *testing.T) {
	testSlug := "power_port_template"
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

resource "netbox_power_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "usb-c"
	maximum_draw = 25
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "type", "usb-c"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "maximum_draw", "25"),
					resource.TestCheckResourceAttrPair("netbox_power_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
		},
	})
}

func init() {
	resource.AddTestSweepers("netbox_power_port_template", &resource.Sweeper{
		Name:         "netbox_power_port_template",
		Dependencies: []string{},
		F: func(r string) error {
			return sweepPowerPortTemplates(r)
		},
	})
}

func sweepPowerPortTemplates(region string) error {
	log.Infof("[INFO] sweeping netbox power port templates in region %s", region)
	api := testAccProvider.Meta().(*providerState).NetBoxAPI

	res, err := api.Dcim.DcimPowerPortTemplatesList(dcim.NewDcimPowerPortTemplatesListParams(), nil)
	if err != nil {
		return err
	}

	for _, item := range res.GetPayload().Results {
		if strings.HasPrefix(*item.Name, testPrefix) {
			log.Infof("[INFO] deleting power port template %d", item.ID)
			if err := deletePowerPortTemplate(item.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func deletePowerPortTemplate(id int64) error {
	api := testAccProvider.Meta().(*providerState).NetBoxAPI
	_, err := api.Dcim.DcimPowerPortTemplatesDelete(
		dcim.NewDcimPowerPortTemplatesDeleteParams().WithID(id), nil)
	return err
}
