package netbox

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/fbreckle/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetboxConsolePorts() *schema.Resource {
	return &schema.Resource{
		Read:        dataSourceNetboxConsolePortRead,
		Description: `:meta:subcategory:dcim:`,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"limit": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				Default:          0,
				Description:      "The limit of objects to return from the API lookup.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"console_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"device_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNetboxConsolePortRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)

	//params := virtualization.NewVirtualizationInterfacesListParams()
	params := dcim.NewDcimConsolePortsListParams()
	params.Limit = getOptionalInt(d, "limit")

	if filter, ok := d.GetOk("filter"); ok {
		var filterParams = filter.(*schema.Set)
		for _, f := range filterParams.List() {
			k := f.(map[string]interface{})["name"]
			v := f.(map[string]interface{})["value"]
			vString := v.(string)
			switch k {
			case "name":
				params.Name = &vString
			case "tag":
				params.Tag = []string{vString} //TODO: switch schema to list?
			case "device_id":
				params.DeviceID = &vString
			default:
				return fmt.Errorf("'%s' is not a supported filter parameter", k)
			}
		}
	}

	res, err := api.Dcim.DcimConsolePortsList(params, nil)
	if err != nil {
		return err
	}

	if *res.GetPayload().Count == int64(0) {
		return errors.New("no result")
	}

	var filteredConsolePorts []*models.ConsolePort
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, consolePort := range res.GetPayload().Results {
			if r.MatchString(*consolePort.Name) {
				filteredConsolePorts = append(filteredConsolePorts, consolePort)
			}
		}
	} else {
		filteredConsolePorts = res.GetPayload().Results
	}

	var s []map[string]interface{}
	for _, v := range filteredConsolePorts {
		var mapping = make(map[string]interface{})
		mapping["id"] = v.ID
		mapping["device_id"] = v.Device.ID
		if v.Name != nil {
			mapping["name"] = *v.Name
		}
		if v.Description != "" {
			mapping["description"] = v.Description
		}
		if v.Tags != nil {
			var tags []int64
			for _, t := range v.Tags {
				tags = append(tags, t.ID)
			}
			mapping["tag_ids"] = tags
		}

		s = append(s, mapping)
	}
	d.SetId(id.UniqueId())
	return d.Set("console_ports", s)
}
