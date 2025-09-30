package netbox

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetboxFrontPortTemplate_basic(t *testing.T) {
	testSlug := "front_port_template"
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
	name = "%[1]s_rear"
	device_type_id = netbox_device_type.test.id
	type = "8p8c"
	positions = 1
}

resource "netbox_front_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "8p8c"
	rear_port_id = netbox_rear_port_template.test.id
	rear_port_position = 1
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", "8p8c"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", "1"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test", "rear_port_id", "netbox_rear_port_template.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_front_port_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxFrontPortTemplate_opts(t *testing.T) {
	testSlug := "front_port_template_opts"
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
	name = "%[1]s_rear"
	device_type_id = netbox_device_type.test.id
	type = "lc"
	positions = 4
}

resource "netbox_front_port_template" "test" {
	name = "%[1]s"
	description = "%[1]s description"
	label = "%[1]s label"
	device_type_id = netbox_device_type.test.id
	type = "lc"
	rear_port_id = netbox_rear_port_template.test.id
	rear_port_position = 2
	color_hex = "ff0000"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "description", fmt.Sprintf("%s description", testName)),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "label", fmt.Sprintf("%s label", testName)),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", "lc"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", "2"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "color_hex", "ff0000"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test", "rear_port_id", "netbox_rear_port_template.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_front_port_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxFrontPortTemplate_moduleType(t *testing.T) {
	testSlug := "front_port_template_module"
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
	name = "%[1]s_rear"
	module_type_id = netbox_module_type.test.id
	type = "sc"
	positions = 2
}

resource "netbox_front_port_template" "test" {
	name = "%[1]s"
	module_type_id = netbox_module_type.test.id
	type = "sc"
	rear_port_id = netbox_rear_port_template.test.id
	rear_port_position = 1
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", "sc"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", "1"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test", "module_type_id", "netbox_module_type.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test", "rear_port_id", "netbox_rear_port_template.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_front_port_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxFrontPortTemplate_multiplePositions(t *testing.T) {
	testSlug := "front_port_template_multi"
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
	name = "%[1]s_rear"
	device_type_id = netbox_device_type.test.id
	type = "mpo"
	positions = 12
}

resource "netbox_front_port_template" "test1" {
	name = "%[1]s_1"
	device_type_id = netbox_device_type.test.id
	type = "lc"
	rear_port_id = netbox_rear_port_template.test.id
	rear_port_position = 1
}

resource "netbox_front_port_template" "test2" {
	name = "%[1]s_2"
	device_type_id = netbox_device_type.test.id
	type = "lc"
	rear_port_id = netbox_rear_port_template.test.id
	rear_port_position = 2
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test1", "name", fmt.Sprintf("%s_1", testName)),
					resource.TestCheckResourceAttr("netbox_front_port_template.test1", "type", "lc"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test1", "rear_port_position", "1"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test2", "name", fmt.Sprintf("%s_2", testName)),
					resource.TestCheckResourceAttr("netbox_front_port_template.test2", "type", "lc"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test2", "rear_port_position", "2"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test1", "rear_port_id", "netbox_rear_port_template.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test2", "rear_port_id", "netbox_rear_port_template.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxFrontPortTemplate_fiberTypes(t *testing.T) {
	testSlug := "front_port_template_fiber"
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
	name = "%[1]s_rear"
	device_type_id = netbox_device_type.test.id
	type = "mpo"
	positions = 12
}

resource "netbox_front_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "lc-apc"
	rear_port_id = netbox_rear_port_template.test.id
	rear_port_position = 1
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", "lc-apc"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", "1"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test", "rear_port_id", "netbox_rear_port_template.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxFrontPortTemplate_copperTypes(t *testing.T) {
	testSlug := "front_port_template_copper"
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
	name = "%[1]s_rear"
	device_type_id = netbox_device_type.test.id
	type = "110-punch"
	positions = 8
}

resource "netbox_front_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "gg45"
	rear_port_id = netbox_rear_port_template.test.id
	rear_port_position = 4
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", "gg45"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", "4"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test", "device_type_id", "netbox_device_type.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_front_port_template.test", "rear_port_id", "netbox_rear_port_template.test", "id"),
				),
			},
		},
	})
}

func TestAccNetboxFrontPortTemplate_update(t *testing.T) {
	testSlug := "front_port_template_update"
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
	name = "%[1]s_rear"
	device_type_id = netbox_device_type.test.id
	type = "8p8c"
	positions = 4
}

resource "netbox_front_port_template" "test" {
	name = "%[1]s"
	device_type_id = netbox_device_type.test.id
	type = "8p8c"
	rear_port_id = netbox_rear_port_template.test.id
	rear_port_position = 1
	description = "Initial description"
	label = "Initial label"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", "8p8c"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", "1"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "label", "Initial label"),
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
	name = "%[1]s_rear"
	device_type_id = netbox_device_type.test.id
	type = "8p8c"
	positions = 4
}

resource "netbox_front_port_template" "test" {
	name = "%[3]s"
	device_type_id = netbox_device_type.test.id
	type = "8p8c"
	rear_port_id = netbox_rear_port_template.test.id
	rear_port_position = 3
	description = "Updated description"
	label = "Updated label"
	color_hex = "00ff00"
}`, testName, randomSlug, testNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", testNameUpdated),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", "8p8c"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", "3"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "label", "Updated label"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "color_hex", "00ff00"),
				),
			},
		},
	})
}

func init() {
	resource.AddTestSweepers("netbox_front_port_template", &resource.Sweeper{
		Name:         "netbox_front_port_template",
		Dependencies: []string{},
		F: func(region string) error {
			m, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("Error getting client: %s", err)
			}
			api := m.(*providerState)
			params := dcim.NewDcimFrontPortTemplatesListParams()
			res, err := api.Dcim.DcimFrontPortTemplatesList(params, nil)
			if err != nil {
				return err
			}
			for _, frontPortTemplate := range res.GetPayload().Results {
				if strings.HasPrefix(*frontPortTemplate.Name, testPrefix) {
					deleteParams := dcim.NewDcimFrontPortTemplatesDeleteParams().WithID(frontPortTemplate.ID)
					_, err := api.Dcim.DcimFrontPortTemplatesDelete(deleteParams, nil)
					if err != nil {
						return err
					}
					log.Print("[DEBUG] Deleted a front port template")
				}
			}
			return nil
		},
	})
}
