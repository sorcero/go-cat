package terraform

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/meta"
	"gitlab.com/sorcero/community/go-cat/ops"
	"gitlab.com/sorcero/community/go-cat/storage"
	"os"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"git_url": {
				Optional: true,
				Type: schema.TypeString,
				Default: os.Getenv(meta.GitUrlEnvVar),
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					if s == "" {
						return nil, []error{errors.New(fmt.Sprintf("git_url or %s environment variable is not set", meta.GitUrlEnvVar ))}
					}
					return nil, nil
				},
			},
			"git_username": {
				Optional: true,
				Type: schema.TypeString,
				Default: os.Getenv(meta.GitUsernameEnvVar),
			},
			"git_password": {
				Sensitive: true,
				Optional: true,
				Type: schema.TypeString,
				Default: os.Getenv(meta.GitPasswordEnvVar),
			},
		},
		ResourcesMap:   map[string]*schema.Resource{
			"gocat_infra": resourceInfra(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gocat_infra": dataSourceInfra(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	gitUrl := d.Get("git_url").(string)
	gitUsername := d.Get("git_username").(string)
	gitPassword := d.Get("git_password").(string)

	cfg := config.GlobalConfig{
		GitRepository: gitUrl,
		GitUsername:   gitUsername,
		GitPassword:   gitPassword,
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	repo, fs, err := storage.Clone(cfg)
	if err != nil {
		return diag.FromErr(err), nil
	}

	ctx := ops.GoCatContext{
		Repo:    repo,
		Storage: fs,
		Config:  cfg,
	}


	return ctx, diags
}