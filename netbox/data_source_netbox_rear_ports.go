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

func dataSourceNetboxRearPorts() *schema.Resource {
	return &schema.Resource{
		Read:        dataSourceNetboxRearPortRead,
		Description: `:meta:subcategory:DCIM:`,
		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the device to get rear ports for.",
			},
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
			"rear_ports": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"positions": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"color_hex": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mark_connected": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tag_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNetboxRearPortRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)

	params := dcim.NewDcimRearPortsListParams()

	params.Limit = getOptionalInt(d, "limit")

	// Get device_id from required parameter
	deviceID := int64(d.Get("device_id").(int))
	deviceIDString := fmt.Sprintf("%d", deviceID)
	params.DeviceID = &deviceIDString

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
				params.Tag = []string{vString}
			case "type":
				params.Type = &vString
			default:
				return fmt.Errorf("'%s' is not a supported filter parameter", k)
			}
		}
	}

	res, err := api.Dcim.DcimRearPortsList(params, nil)
	if err != nil {
		return err
	}

	if *res.GetPayload().Count == int64(0) {
		return errors.New("no result")
	}

	var filteredRearPorts []*models.RearPort
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, rearPort := range res.GetPayload().Results {
			if r.MatchString(*rearPort.Name) {
				filteredRearPorts = append(filteredRearPorts, rearPort)
			}
		}
	} else {
		filteredRearPorts = res.GetPayload().Results
	}

	var s []map[string]interface{}
	for _, v := range filteredRearPorts {
		var mapping = make(map[string]interface{})
		mapping["id"] = v.ID
		if v.Description != "" {
			mapping["description"] = v.Description
		}
		if v.Name != nil {
			mapping["name"] = *v.Name
		}
		if v.Type != nil {
			mapping["type"] = *v.Type.Value
		}
		mapping["positions"] = v.Positions
		if v.Label != "" {
			mapping["label"] = v.Label
		}
		if v.Color != "" {
			mapping["color_hex"] = v.Color
		}
		mapping["mark_connected"] = v.MarkConnected
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
	return d.Set("rear_ports", s)
}
