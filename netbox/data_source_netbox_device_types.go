package netbox

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/fbreckle/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetboxDeviceTypes() *schema.Resource {
	return &schema.Resource{
		Read:        dataSourceNetboxDeviceTypesRead,
		Description: ":meta:subcategory:Data Center Inventory Management (DCIM):",
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
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"device_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_type_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"part_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"manufacturer_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"manufacturer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_full_depth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"u_height": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"custom_fields": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_updated": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchemaRead,
					},
				},
			},
		},
	}
}

func dataSourceNetboxDeviceTypesRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)

	params := dcim.NewDcimDeviceTypesListParams()

	if filter, ok := d.GetOk("filter"); ok {
		var filterParams = filter.(*schema.Set)
		for _, f := range filterParams.List() {
			k := f.(map[string]interface{})["name"]
			v := f.(map[string]interface{})["value"]
			vString := v.(string)
			switch k {
			case "manufacturer":
				params.Manufacturer = &vString
			case "manufacturer_id":
				params.ManufacturerID = &vString
			case "model":
				params.Model = &vString
			case "slug":
				params.Slug = &vString
			case "part_number":
				params.PartNumber = &vString
			case "is_full_depth":
				params.IsFullDepth = &vString
			case "subdevice_role":
				params.SubdeviceRole = &vString
			case "u_height":
				params.UHeight = &vString
			case "tags":
				params.Tag = strings.Split(vString, ",")
			case "q":
				params.Q = &vString
			default:
				return fmt.Errorf("'%s' is not a supported filter parameter", k)
			}
		}
	}

	if limit, ok := d.GetOk("limit"); ok {
		limitInt := int64(limit.(int))
		params.Limit = &limitInt
	}

	res, err := api.Dcim.DcimDeviceTypesList(params, nil)
	if err != nil {
		return err
	}

	var filteredDeviceTypes []*models.DeviceType
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, deviceType := range res.GetPayload().Results {
			if r.MatchString(*deviceType.Model) {
				filteredDeviceTypes = append(filteredDeviceTypes, deviceType)
			}
		}
	} else {
		filteredDeviceTypes = res.GetPayload().Results
	}

	var s []map[string]interface{}
	for _, deviceType := range filteredDeviceTypes {
		var mapping = make(map[string]interface{})
		mapping["device_type_id"] = deviceType.ID
		if deviceType.Model != nil {
			mapping["model"] = *deviceType.Model
		}
		if deviceType.Slug != nil {
			mapping["slug"] = *deviceType.Slug
		}
		if deviceType.PartNumber != "" {
			mapping["part_number"] = deviceType.PartNumber
		}
		if deviceType.Description != "" {
			mapping["description"] = deviceType.Description
		}
		if deviceType.Comments != "" {
			mapping["comments"] = deviceType.Comments
		}
		if deviceType.Manufacturer != nil {
			mapping["manufacturer_id"] = deviceType.Manufacturer.ID
			if deviceType.Manufacturer.Name != nil {
				mapping["manufacturer_name"] = *deviceType.Manufacturer.Name
			}
		}
		mapping["is_full_depth"] = deviceType.IsFullDepth
		if deviceType.UHeight != nil {
			mapping["u_height"] = *deviceType.UHeight
		}
		if deviceType.CustomFields != nil {
			mapping["custom_fields"] = deviceType.CustomFields
		}
		if deviceType.Tags != nil {
			mapping["tags"] = getTagListFromNestedTagList(deviceType.Tags)
		}
		if deviceType.Created != nil {
			mapping["created"] = deviceType.Created.String()
		}
		if deviceType.LastUpdated != nil {
			mapping["last_updated"] = deviceType.LastUpdated.String()
		}
		s = append(s, mapping)
	}

	d.SetId(id.UniqueId())
	return d.Set("device_types", s)
}
