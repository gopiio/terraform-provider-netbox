package netbox

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	log "github.com/sirupsen/logrus"
)

func TestAccNetboxPowerOutletTemplate_basic(t *testing.T) {
	testSlug := "power_outlet_template"
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

resource "netbox_power_outlet_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "iec-60320-c13"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "type", "iec-60320-c13"),
					resource.TestCheckResourceAttrPair("netbox_power_outlet_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_power_outlet_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxPowerOutletTemplate_opts(t *testing.T) {
	testSlug := "power_outlet_template"
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
	name = "%[1]s_port"
	device_type_id = netbox_device_type.test.id
	type = "iec-60320-c14"
}

resource "netbox_power_outlet_template" "test" {
	name = "%[1]s"
	description = "%[1]s description"
	label = "%[1]s label"
	device_type_id = netbox_device_type.test.id
	type = "iec-60320-c13"
	power_port_id = netbox_power_port_template.test.id
	feed_leg = "A"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "description", fmt.Sprintf("%s description", testName)),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "label", fmt.Sprintf("%s label", testName)),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "type", "iec-60320-c13"),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "feed_leg", "A"),
					resource.TestCheckResourceAttrPair("netbox_power_outlet_template.test", "device_type_id", "netbox_device_type.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_power_outlet_template.test", "power_port_id", "netbox_power_port_template.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_power_outlet_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxPowerOutletTemplate_moduleType(t *testing.T) {
	testSlug := "power_outlet_template"
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

resource "netbox_power_outlet_template" "test" {
	name = "%[1]s"
	module_type_id = netbox_module_type.test.id
	type = "nema-5-15r"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "type", "nema-5-15r"),
					resource.TestCheckResourceAttrPair("netbox_power_outlet_template.test", "module_type_id", "netbox_module_type.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxPowerOutletTemplate_nema(t *testing.T) {
	testSlug := "power_outlet_template"
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

resource "netbox_power_outlet_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "nema-5-20r"
	feed_leg = "B"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "type", "nema-5-20r"),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "feed_leg", "B"),
					resource.TestCheckResourceAttrPair("netbox_power_outlet_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxPowerOutletTemplate_usb(t *testing.T) {
	testSlug := "power_outlet_template"
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

resource "netbox_power_outlet_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "usb-c"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "type", "usb-c"),
					resource.TestCheckResourceAttrPair("netbox_power_outlet_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxPowerOutletTemplate_feedLegOptions(t *testing.T) {
	testSlug := "power_outlet_template"
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

resource "netbox_power_outlet_template" "test_a" {
	name = "%[1]s_a"
	device_type_id = netbox_device_type.test.id
	type = "iec-60320-c13"
	feed_leg = "A"
}

resource "netbox_power_outlet_template" "test_b" {
	name = "%[1]s_b"
	device_type_id = netbox_device_type.test.id
	type = "iec-60320-c13"
	feed_leg = "B"
}

resource "netbox_power_outlet_template" "test_c" {
	name = "%[1]s_c"
	device_type_id = netbox_device_type.test.id
	type = "iec-60320-c13"
	feed_leg = "C"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test_a", "feed_leg", "A"),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test_b", "feed_leg", "B"),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test_c", "feed_leg", "C"),
				),
			},
		},
	})
}

func init() {
	resource.AddTestSweepers("netbox_power_outlet_template", &resource.Sweeper{
		Name:         "netbox_power_outlet_template",
		Dependencies: []string{},
		F: func(r string) error {
			return sweepPowerOutletTemplates(r)
		},
	})
}

func sweepPowerOutletTemplates(region string) error {
	log.Infof("[INFO] sweeping netbox power outlet templates in region %s", region)
	api := testAccProvider.Meta().(*providerState).NetBoxAPI

	res, err := api.Dcim.DcimPowerOutletTemplatesList(dcim.NewDcimPowerOutletTemplatesListParams(), nil)
	if err != nil {
		return err
	}

	for _, item := range res.GetPayload().Results {
		if strings.HasPrefix(*item.Name, testPrefix) {
			log.Infof("[INFO] deleting power outlet template %d", item.ID)
			if err := deletePowerOutletTemplate(item.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func deletePowerOutletTemplate(id int64) error {
	api := testAccProvider.Meta().(*providerState).NetBoxAPI
	_, err := api.Dcim.DcimPowerOutletTemplatesDelete(
		dcim.NewDcimPowerOutletTemplatesDeleteParams().WithID(id), nil)
	return err
}
