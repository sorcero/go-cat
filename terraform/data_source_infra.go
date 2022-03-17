package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.com/sorcero/community/go-cat/ops"
)

func dataSourceInfraRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(ops.GoCatContext)
	if !ok {
		return diag.FromErr(ErrorInvalidConfig)
	}

	var diags diag.Diagnostics
	id := d.Get("id").(string)

	metadata, err := ops.CatFromStorage(c.Storage, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(metadata) == 0 {
		return diag.FromErr(fmt.Errorf("no resource with id '%s' exists", id))
	} else if len(metadata) > 1 {
		return diag.FromErr(fmt.Errorf("more than one resource was returned when '%s' id was requested which is not supported yet", id))
	}

	err = d.Set("deployment_link", metadata[0].DeploymentLink)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("deployment_links", metadata[0].DeploymentLinks)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("commit_sha", metadata[0].CommitSha)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("name", metadata[0].Name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return diags
}

func dataSourceInfra() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInfraRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_sha": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
