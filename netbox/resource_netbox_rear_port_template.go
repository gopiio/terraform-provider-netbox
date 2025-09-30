package netbox

import (
	"strconv"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/fbreckle/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetboxRearPortTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxRearPortTemplateCreate,
		Read:   resourceNetboxRearPortTemplateRead,
		Update: resourceNetboxRearPortTemplateUpdate,
		Delete: resourceNetboxRearPortTemplateDelete,

		Description: `:meta:subcategory:Data Center Inventory Management (DCIM):From the [official documentation](https://docs.netbox.dev/en/stable/models/dcim/rearporttemplate/):

> Like front ports, rear ports are pass-through ports which represent the continuation of a path from one cable to the next. Each rear port is defined with its physical type and a number of positions: Rear ports with more than one position can be mapped to multiple front ports. This can be useful for modeling instances where multiple paths share a common cable (for example, six discrete two-strand fiber connections sharing a 12-strand MPO cable).`,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "One of [8p8c, 8p6c, 8p4c, 8p2c, 6p6c, 6p4c, 6p2c, 4p4c, 4p2c, gg45, tera-4p, tera-2p, tera-1p, 110-punch, bnc, f, n, mrj21, fc, lc, lc-pc, lc-upc, lc-apc, lsh, lsh-pc, lsh-upc, lsh-apc, mpo, mtrj, sc, sc-pc, sc-upc, sc-apc, st, cs, sn, sma-905, sma-906, urm-p2, urm-p4, urm-p8, splice, other]",
			},
			"positions": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"module_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"color_hex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"device_type_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"device_type_id", "module_type_id"},
				ForceNew:     true,
			},
			"module_type_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"device_type_id", "module_type_id"},
				ForceNew:     true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceNetboxRearPortTemplateCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)

	data := models.WritableRearPortTemplate{
		Name:        strToPtr(d.Get("name").(string)),
		Type:        strToPtr(d.Get("type").(string)),
		Positions:   int64(d.Get("positions").(int)),
		Label:       getOptionalStr(d, "label", false),
		Color:       getOptionalStr(d, "color_hex", false),
		Description: getOptionalStr(d, "description", false),
	}

	if deviceTypeID, ok := d.Get("device_type_id").(int); ok && deviceTypeID != 0 {
		data.DeviceType = int64ToPtr(int64(deviceTypeID))
	}
	if moduleTypeID, ok := d.Get("module_type_id").(int); ok && moduleTypeID != 0 {
		data.ModuleType = int64ToPtr(int64(moduleTypeID))
	}

	var err error
	params := dcim.NewDcimRearPortTemplatesCreateParams().WithData(&data)

	res, err := api.Dcim.DcimRearPortTemplatesCreate(params, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxRearPortTemplateRead(d, m)
}

func resourceNetboxRearPortTemplateRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := dcim.NewDcimRearPortTemplatesReadParams().WithID(id)

	res, err := api.Dcim.DcimRearPortTemplatesRead(params, nil)

	if err != nil {
		errorcode := err.(*dcim.DcimRearPortTemplatesReadDefault).Code()
		if errorcode == 404 {
			// If the ID is updated to blank, this tells Terraform the resource no longer exists (maybe it was destroyed out of band). Just like the destroy callback, the Read function should gracefully handle this case. https://www.terraform.io/docs/extend/writing-custom-providers.html
			d.SetId("")
			return nil
		}
		return err
	}

	rearPort := res.GetPayload()

	if rearPort.DeviceType != nil {
		d.Set("device_type_id", rearPort.DeviceType.ID)
	} else {
		d.Set("device_type_id", nil)
	}

	d.Set("name", rearPort.Name)

	if rearPort.Type != nil {
		d.Set("type", rearPort.Type.Value)
	} else {
		d.Set("type", nil)
	}

	d.Set("positions", rearPort.Positions)

	if rearPort.ModuleType != nil {
		d.Set("module_type_id", rearPort.ModuleType.ID)
	} else {
		d.Set("module_type_id", nil)
	}

	d.Set("label", rearPort.Label)
	d.Set("color_hex", rearPort.Color)
	d.Set("description", rearPort.Description)

	return nil
}

func resourceNetboxRearPortTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)

	data := models.WritableRearPortTemplate{
		Name:        strToPtr(d.Get("name").(string)),
		Type:        strToPtr(d.Get("type").(string)),
		Positions:   int64(d.Get("positions").(int)),
		Label:       getOptionalStr(d, "label", true),
		Color:       getOptionalStr(d, "color_hex", false),
		Description: getOptionalStr(d, "description", true),
	}

	if deviceTypeID, ok := d.Get("device_type_id").(int); ok && deviceTypeID != 0 {
		data.DeviceType = int64ToPtr(int64(deviceTypeID))
	}
	if moduleTypeID, ok := d.Get("module_type_id").(int); ok && moduleTypeID != 0 {
		data.ModuleType = int64ToPtr(int64(moduleTypeID))
	}

	var err error
	params := dcim.NewDcimRearPortTemplatesPartialUpdateParams().WithID(id).WithData(&data)

	_, err = api.Dcim.DcimRearPortTemplatesPartialUpdate(params, nil)
	if err != nil {
		return err
	}

	return resourceNetboxRearPortTemplateRead(d, m)
}

func resourceNetboxRearPortTemplateDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := dcim.NewDcimRearPortTemplatesDeleteParams().WithID(id)

	_, err := api.Dcim.DcimRearPortTemplatesDelete(params, nil)
	if err != nil {
		return err
	}
	return nil
}
