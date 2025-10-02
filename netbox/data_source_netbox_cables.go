package netbox

import (
	"errors"
	"regexp"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/fbreckle/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetboxCables() *schema.Resource {
	return &schema.Resource{
		Read:        dataSourceNetboxCablesRead,
		Description: `:meta:subcategory:Data Center Inventory Management (DCIM):`,
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
			"cables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"color_hex": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"length": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"length_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeInt,
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
						"a_terminations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"object_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"b_terminations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"object_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"tag_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"custom_fields": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNetboxCablesRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)

	params := dcim.NewDcimCablesListParams()

	if limit, ok := d.GetOk("limit"); ok {
		limitInt := int64(limit.(int))
		params.Limit = &limitInt
	}

	if filter, ok := d.GetOk("filter"); ok {
		var filterParams = filter.(*schema.Set)
		for _, f := range filterParams.List() {
			k := f.(map[string]interface{})["name"]
			v := f.(map[string]interface{})["value"]
			vString := v.(string)
			switch k {
			case "status":
				params.Status = &vString
			case "type":
				params.Type = &vString
			case "label":
				params.Label = &vString
			case "color":
				params.Color = &vString
			case "tenant_id":
				params.TenantID = &vString
			case "tenant":
				params.Tenant = &vString
			case "site_id":
				params.SiteID = &vString
			case "site":
				params.Site = &vString
			case "rack_id":
				params.RackID = &vString
			case "rack":
				params.Rack = &vString
			case "device_id":
				params.DeviceID = &vString
			case "device":
				params.Device = &vString
			default:
				return errors.New("'" + k.(string) + "' is not a supported filter parameter")
			}
		}
	}

	res, err := api.Dcim.DcimCablesList(params, nil)
	if err != nil {
		return err
	}

	if *res.GetPayload().Count == int64(0) {
		return errors.New("no result")
	}

	var filteredCables []*models.Cable
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, cable := range res.GetPayload().Results {
			if cable.Label != "" && r.MatchString(cable.Label) {
				filteredCables = append(filteredCables, cable)
			}
		}
	} else {
		filteredCables = res.GetPayload().Results
	}

	var s []map[string]interface{}
	for _, v := range filteredCables {
		var mapping = make(map[string]interface{})
		mapping["id"] = v.ID
		mapping["label"] = v.Label

		if v.Status != nil {
			mapping["status"] = v.Status.Value
		}

		mapping["type"] = v.Type
		mapping["color_hex"] = v.Color
		mapping["length"] = v.Length

		if v.LengthUnit != nil {
			mapping["length_unit"] = v.LengthUnit.Value
		}

		if v.Tenant != nil {
			mapping["tenant_id"] = v.Tenant.ID
		}

		mapping["description"] = v.Description
		mapping["comments"] = v.Comments

		// Handle A terminations
		if v.ATerminations != nil {
			var aTerms []map[string]interface{}
			for _, term := range v.ATerminations {
				termMap := make(map[string]interface{})
				termMap["object_type"] = term.ObjectType
				termMap["object_id"] = term.ObjectID
				aTerms = append(aTerms, termMap)
			}
			mapping["a_terminations"] = aTerms
		}

		// Handle B terminations
		if v.BTerminations != nil {
			var bTerms []map[string]interface{}
			for _, term := range v.BTerminations {
				termMap := make(map[string]interface{})
				termMap["object_type"] = term.ObjectType
				termMap["object_id"] = term.ObjectID
				bTerms = append(bTerms, termMap)
			}
			mapping["b_terminations"] = bTerms
		}

		// Handle tags
		if v.Tags != nil {
			var tags []int64
			for _, t := range v.Tags {
				tags = append(tags, t.ID)
			}
			mapping["tag_ids"] = tags
		}

		// Handle custom fields
		if v.CustomFields != nil {
			cf := getCustomFields(v.CustomFields)
			if cf != nil {
				mapping["custom_fields"] = cf
			}
		}

		s = append(s, mapping)
	}

	d.SetId(id.UniqueId())
	return d.Set("cables", s)
}
