package netbox

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetboxRearPortTemplate_basic(t *testing.T) {
	testSlug := "rear_port_template"
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

resource "netbox_rear_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "8p8c"
	positions = 1
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", "8p8c"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", "1"),
					resource.TestCheckResourceAttrPair("netbox_rear_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_rear_port_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxRearPortTemplate_opts(t *testing.T) {
	testSlug := "rear_port_template_opts"
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

resource "netbox_rear_port_template" "test" {
	name = "%[1]s"
	description = "%[1]s description"
	label = "%[1]s label"
	device_type_id = netbox_device_type.test.id
	type = "lc"
	positions = 4
	color_hex = "ff0000"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", fmt.Sprintf("%s description", testName)),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", fmt.Sprintf("%s label", testName)),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", "lc"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", "4"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "color_hex", "ff0000"),
					resource.TestCheckResourceAttrPair("netbox_rear_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_rear_port_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxRearPortTemplate_moduleType(t *testing.T) {
	testSlug := "rear_port_template_module"
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

resource "netbox_rear_port_template" "test" {
	name = "%[1]s"
	module_type_id = netbox_module_type.test.id
	type = "sc"
	positions = 2
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", "sc"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", "2"),
					resource.TestCheckResourceAttrPair("netbox_rear_port_template.test", "module_type_id", "netbox_module_type.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_rear_port_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxRearPortTemplate_fiberTypes(t *testing.T) {
	testSlug := "rear_port_template_fiber"
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

resource "netbox_rear_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "mpo"
	positions = 12
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", "mpo"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", "12"),
					resource.TestCheckResourceAttrPair("netbox_rear_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxRearPortTemplate_copperTypes(t *testing.T) {
	testSlug := "rear_port_template_copper"
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

resource "netbox_rear_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "110-punch"
	positions = 8
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", "110-punch"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", "8"),
					resource.TestCheckResourceAttrPair("netbox_rear_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxRearPortTemplate_withModuleId(t *testing.T) {
	testSlug := "rear_port_template_module_id"
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

resource "netbox_module_type" "test" {
	model = "%[1]s_module"
	manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_rear_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "fc"
	positions = 2
	module_id = netbox_module_type.test.id
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", "fc"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", "2"),
					resource.TestCheckResourceAttrPair("netbox_rear_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_rear_port_template.test", "module_id", "netbox_module_type.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxRearPortTemplate_update(t *testing.T) {
	testSlug := "rear_port_template_update"
	testName := testAccGetTestName(testSlug)
	testNameUpdated := testAccGetTestName(testSlug + "_updated")
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

resource "netbox_rear_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "8p8c"
	positions = 1
	description = "Initial description"
	label = "Initial label"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", "8p8c"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", "1"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", "Initial label"),
				),
			},
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

resource "netbox_rear_port_template" "test" {
	name = "%[3]s"
	device_type_id = netbox_device_type.test.id
	type = "lc"
	positions = 2
	description = "Updated description"
	label = "Updated label"
	color_hex = "00ff00"
}`, testName, randomSlug, testNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", testNameUpdated),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", "lc"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", "2"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", "Updated label"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "color_hex", "00ff00"),
				),
			},
		},
	})
}

func init() {
	resource.AddTestSweepers("netbox_rear_port_template", &resource.Sweeper{
		Name:         "netbox_rear_port_template",
		Dependencies: []string{},
		F: func(region string) error {
			m, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("Error getting client: %s", err)
			}
			api := m.(*providerState)
			params := dcim.NewDcimRearPortTemplatesListParams()
			res, err := api.Dcim.DcimRearPortTemplatesList(params, nil)
			if err != nil {
				return err
			}
			for _, rearPortTemplate := range res.GetPayload().Results {
				if strings.HasPrefix(*rearPortTemplate.Name, testPrefix) {
					deleteParams := dcim.NewDcimRearPortTemplatesDeleteParams().WithID(rearPortTemplate.ID)
					_, err := api.Dcim.DcimRearPortTemplatesDelete(deleteParams, nil)
					if err != nil {
						return err
					}
					log.Print("[DEBUG] Deleted a rear port template")
				}
			}
			return nil
		},
	})
}
