package netbox

import (
	"context"
	"strconv"

	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/fbreckle/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceConsoleServerPortTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConsoleServerPortTemplateCreate,
		ReadContext:   resourceConsoleServerPortTemplateRead,
		UpdateContext: resourceConsoleServerPortTemplateUpdate,
		DeleteContext: resourceConsoleServerPortTemplateDelete,

		Description: `:meta:subcategory:Data Center Inventory Management (DCIM):From the [official documentation](https://netboxlabs.com/docs/netbox/models/dcim/consoleserverporttemplate/):

> A template for a console port that will be created on all instantiations of the parent device type. See the console port documentation for more detail.`,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "One of [de-9, db-25, rj-11, rj-12, rj-45, mini-din-8, usb-a, usb-b, usb-c, usb-mini-a, usb-mini-b, usb-micro-a, usb-micro-b, usb-micro-ab, other]",
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

func resourceConsoleServerPortTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*providerState)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	label := d.Get("label").(string)
	consolePortType := d.Get("type").(string)

	data := models.WritableConsoleServerPortTemplate{
		Name:        &name,
		Description: description,
		Label:       label,
		Type:        consolePortType,
	}

	if deviceTypeID, ok := d.Get("device_type_id").(int); ok && deviceTypeID != 0 {
		data.DeviceType = int64ToPtr(int64(deviceTypeID))
	}
	if moduleTypeID, ok := d.Get("module_type_id").(int); ok && moduleTypeID != 0 {
		data.ModuleType = int64ToPtr(int64(moduleTypeID))
	}
	params := dcim.NewDcimConsoleServerPortTemplatesCreateParams().WithData(&data)
	res, err := api.Dcim.DcimConsoleServerPortTemplatesCreate(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return diags
}

func resourceConsoleServerPortTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*providerState)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)

	var diags diag.Diagnostics

	params := dcim.NewDcimConsoleServerPortTemplatesReadParams().WithID(id)

	res, err := api.Dcim.DcimConsoleServerPortTemplatesRead(params, nil)
	if err != nil {
		if errresp, ok := err.(*dcim.DcimConsoleServerPortTemplatesReadDefault); ok {
			errorcode := errresp.Code()
			if errorcode == 404 {
				// If the ID is updated to blank, this tells Terraform the resource no longer exists (maybe it was destroyed out of band). Just like the destroy callback, the Read function should gracefully handle this case. https://www.terraform.io/docs/extend/writing-custom-providers.html
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	tmpl := res.GetPayload()

	d.Set("name", tmpl.Name)
	d.Set("description", tmpl.Description)
	d.Set("label", tmpl.Label)

	if tmpl.Type.Value != nil {
		d.Set("type", tmpl.Type.Value)
	} else {
		d.Set("type", nil)
	}
	if tmpl.DeviceType != nil {
		d.Set("device_type_id", tmpl.DeviceType.ID)
	}
	if tmpl.ModuleType != nil {
		d.Set("module_type_id", tmpl.ModuleType.ID)
	}

	return diags
}

func resourceConsoleServerPortTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*providerState)

	var diags diag.Diagnostics

	id, _ := strconv.ParseInt(d.Id(), 10, 64)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	label := d.Get("label").(string)
	consolePortType := d.Get("type").(string)

	data := models.WritableConsoleServerPortTemplate{
		Name:        &name,
		Description: description,
		Label:       label,
		Type:        consolePortType,
	}

	if d.HasChange("device_type_id") {
		deviceTypeID := int64(d.Get("device_type_id").(int))
		data.DeviceType = &deviceTypeID
	}

	if d.HasChange("module_type_id") {
		moduleTypeID := int64(d.Get("module_type_id").(int))
		data.ModuleType = &moduleTypeID
	}

	params := dcim.NewDcimConsoleServerPortTemplatesPartialUpdateParams().WithID(id).WithData(&data)
	_, err := api.Dcim.DcimConsoleServerPortTemplatesPartialUpdate(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceConsoleServerPortTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*providerState)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := dcim.NewDcimConsoleServerPortTemplatesDeleteParams().WithID(id)

	_, err := api.Dcim.DcimConsoleServerPortTemplatesDelete(params, nil)
	if err != nil {
		if errresp, ok := err.(*dcim.DcimConsoleServerPortTemplatesDeleteDefault); ok {
			if errresp.Code() == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}
	return nil
}
