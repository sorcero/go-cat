package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func safeGet(d *schema.ResourceData, key string) string {
	v, ok := d.Get(key).(string)
	if !ok {
		return ""
	}
	return v
}


func NewInfraFromSchemaResourceData(d *schema.ResourceData) *infrastructure.Metadata {
	infra := &infrastructure.Metadata{
		Name:           safeGet(d, "name"),
		CommitSha:      safeGet(d, "commit_sha"),
		Version:        safeGet(d, "version"),
		Cloud:          safeGet(d, "cloud"),
		CloudProjectId: safeGet(d, "cloud_project_id"),
		Subsystem:      safeGet(d, "subsystem"),
		Type:           safeGet(d, "type"),
		MonitoringLink: safeGet(d, "monitoring_link"),
		DeploymentLink: safeGet(d, "deployment_link"),
	}
	return infra
}

func NewSchemaResourceDataFromInfra(d *infrastructure.Metadata, s *schema.ResourceData) *schema.ResourceData {
	must(s.Set("name", d.Name))
	must(s.Set("commit_sha", d.CommitSha))
	must(s.Set("version", d.Version))
	must(s.Set("cloud", d.Cloud))
	must(s.Set("cloud_project_id", d.CloudProjectId))
	must(s.Set("type", d.Type))
	must(s.Set("subsystem", d.Subsystem))
	must(s.Set("deployment_link", d.DeploymentLink))

	return s
}
