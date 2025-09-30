package netbox

import (
	"strconv"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/fbreckle/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetboxFrontPortTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxFrontPortTemplateCreate,
		Read:   resourceNetboxFrontPortTemplateRead,
		Update: resourceNetboxFrontPortTemplateUpdate,
		Delete: resourceNetboxFrontPortTemplateDelete,

		Description: `:meta:subcategory:Data Center Inventory Management (DCIM):From the [official documentation](https://docs.netbox.dev/en/stable/models/dcim/frontporttemplate/):

> Front ports are pass-through ports which represent physical cable connections that comprise part of a longer path. For example, the ports on the front face of a UTP patch panel would be modeled in NetBox as front ports. Each port is assigned a physical type, and must be mapped to a specific rear port on the same device. A single rear port may be mapped to multiple front ports, using numeric positions to annotate the specific alignment of each.`,

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
			"rear_port_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rear_port_position": {
				Type:     schema.TypeInt,
				Required: true,
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

func resourceNetboxFrontPortTemplateCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)

	data := models.WritableFrontPortTemplate{
		Name:             strToPtr(d.Get("name").(string)),
		Type:             strToPtr(d.Get("type").(string)),
		RearPort:         int64ToPtr(int64(d.Get("rear_port_id").(int))),
		RearPortPosition: int64(d.Get("rear_port_position").(int)),
		Label:            getOptionalStr(d, "label", false),
		Color:            getOptionalStr(d, "color_hex", false),
		Description:      getOptionalStr(d, "description", false),
	}

	if deviceTypeID, ok := d.Get("device_type_id").(int); ok && deviceTypeID != 0 {
		data.DeviceType = int64ToPtr(int64(deviceTypeID))
	}
	if moduleTypeID, ok := d.Get("module_type_id").(int); ok && moduleTypeID != 0 {
		data.ModuleType = int64ToPtr(int64(moduleTypeID))
	}

	params := dcim.NewDcimFrontPortTemplatesCreateParams().WithData(&data)
	var err error
	res, err := api.Dcim.DcimFrontPortTemplatesCreate(params, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxFrontPortTemplateRead(d, m)
}

func resourceNetboxFrontPortTemplateRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := dcim.NewDcimFrontPortTemplatesReadParams().WithID(id)

	res, err := api.Dcim.DcimFrontPortTemplatesRead(params, nil)

	if err != nil {
		errorcode := err.(*dcim.DcimFrontPortTemplatesReadDefault).Code()
		if errorcode == 404 {
			// If the ID is updated to blank, this tells Terraform the resource no longer exists (maybe it was destroyed out of band). Just like the destroy callback, the Read function should gracefully handle this case. https://www.terraform.io/docs/extend/writing-custom-providers.html
			d.SetId("")
			return nil
		}
		return err
	}

	frontPort := res.GetPayload()

	if frontPort.DeviceType != nil {
		d.Set("device_type_id", frontPort.DeviceType.ID)
	} else {
		d.Set("device_type_id", nil)
	}

	d.Set("name", frontPort.Name)

	if frontPort.Type != nil {
		d.Set("type", frontPort.Type.Value)
	} else {
		d.Set("type", nil)
	}

	if frontPort.RearPort != nil {
		d.Set("rear_port_id", frontPort.RearPort.ID)
	} else {
		d.Set("rear_port_id", nil)
	}

	d.Set("rear_port_position", frontPort.RearPortPosition)

	if frontPort.ModuleType != nil {
		d.Set("module_type_id", frontPort.ModuleType.ID)
	} else {
		d.Set("module_type_id", nil)
	}

	d.Set("label", frontPort.Label)
	d.Set("color_hex", frontPort.Color)
	d.Set("description", frontPort.Description)

	return nil
}

func resourceNetboxFrontPortTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)

	data := models.WritableFrontPortTemplate{
		Name:             strToPtr(d.Get("name").(string)),
		Type:             strToPtr(d.Get("type").(string)),
		RearPort:         int64ToPtr(int64(d.Get("rear_port_id").(int))),
		RearPortPosition: int64(d.Get("rear_port_position").(int)),
		Label:            getOptionalStr(d, "label", false),
		Color:            getOptionalStr(d, "color_hex", false),
		Description:      getOptionalStr(d, "description", false),
	}

	if deviceTypeID, ok := d.Get("device_type_id").(int); ok && deviceTypeID != 0 {
		data.DeviceType = int64ToPtr(int64(deviceTypeID))
	}
	if moduleTypeID, ok := d.Get("module_type_id").(int); ok && moduleTypeID != 0 {
		data.ModuleType = int64ToPtr(int64(moduleTypeID))
	}

	var err error
	params := dcim.NewDcimFrontPortTemplatesPartialUpdateParams().WithID(id).WithData(&data)

	_, err = api.Dcim.DcimFrontPortTemplatesPartialUpdate(params, nil)
	if err != nil {
		return err
	}

	return resourceNetboxFrontPortTemplateRead(d, m)
}

func resourceNetboxFrontPortTemplateDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := dcim.NewDcimFrontPortTemplatesDeleteParams().WithID(id)

	_, err := api.Dcim.DcimFrontPortTemplatesDelete(params, nil)
	if err != nil {
		return err
	}
	return nil
}
