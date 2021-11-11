package terraform

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.com/sorcero/community/go-cat/ops"
)

func resourceInfra() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cloud": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cloud_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deployment_link": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployment_links": {
				Type:     schema.TypeList,
				Elem: 	&schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"monitoring_link": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subsystem": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: 	&schema.Schema{Type: schema.TypeString},
			},
			"commit_sha": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		CreateContext: resourceInfraCreate,
		ReadContext:   resourceInfraRead,
		DeleteContext: resourceInfraDelete,
		UpdateContext: resourceInfraUpdate,
	}
}

func resourceInfraCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(ops.GoCatContext)
	if !ok {
		return diag.FromErr(InvalidConfigError)
	}

	var diags diag.Diagnostics

	infra := NewInfraFromSchemaResourceData(d)
	err := ops.UpsertFromStorage(c.Config, c.Repo, c.Storage, infra)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infra.GetId())
	return diags
}

func resourceInfraRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(ops.GoCatContext)
	if !ok {
		return diag.FromErr(InvalidConfigError)
	}

	id := d.Id()

	var diags diag.Diagnostics

	infraGroup, err := ops.CatFromStorage(c.Storage, id)
	if len(infraGroup) == 0 {
		return diag.FromErr(errors.New(fmt.Sprintf("no resource with id '%s' exists", id)))
	} else if len(infraGroup) > 1 {
		return diag.FromErr(errors.New(fmt.Sprintf("more than one resource was returned when '%s' id was requested which is not supported yet", id)))
	}

	if err != nil {
		return diag.FromErr(err)
	}
	infra := infraGroup[0]
	NewSchemaResourceDataFromInfra(infra, d)

	d.SetId(infra.GetId())
	return diags
}

func resourceInfraDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(ops.GoCatContext)
	if !ok {
		return diag.FromErr(InvalidConfigError)
	}

	var diags diag.Diagnostics

	infra := NewInfraFromSchemaResourceData(d)
	id := infra.GetId()
	err := ops.RemoveFromStorage(c.Config, c.Repo, c.Storage, infra.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return diags
}

func resourceInfraUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceInfraCreate(ctx, d, m)
}
