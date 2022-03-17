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

func safeListGet(d *schema.ResourceData, key string) []string {
	v := d.Get(key).([]interface{})
	var s []string
	for i := range v {
		s = append(s, v[i].(string))
	}
	return s
}

func safeMapGet(d *schema.ResourceData, key string) map[string]interface{} {
	v := d.Get(key).(map[string]interface{})
	return v
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
		Name:            safeGet(d, "name"),
		CommitSha:       safeGet(d, "commit_sha"),
		Version:         safeGet(d, "version"),
		Cloud:           safeGet(d, "cloud"),
		CloudProjectId:  safeGet(d, "cloud_project_id"),
		Subsystem:       safeGet(d, "subsystem"),
		Type:            safeGet(d, "type"),
		MonitoringLink:  safeGet(d, "monitoring_link"),
		DeploymentLink:  safeGet(d, "deployment_link"),
		MonitoringLinks: safeListGet(d, "monitoring_links"),
		LoggingLinks:    safeListGet(d, "logging_links"),
		DeploymentLinks: safeListGet(d, "deployment_links"),
		Parameters:      safeMapGet(d, "parameters"),
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
	must(s.Set("deployment_links", d.DeploymentLinks))
	must(s.Set("monitoring_link", d.MonitoringLink))
	must(s.Set("monitoring_links", d.MonitoringLinks))
	must(s.Set("logging_links", d.LoggingLinks))
	must(s.Set("parameters", d.Parameters))
	return s
}
