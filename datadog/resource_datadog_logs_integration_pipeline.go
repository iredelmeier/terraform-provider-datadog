package datadog

import (
	"fmt"

	datadogV1 "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDatadogLogsIntegrationPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatadogLogsIntegrationPipelineCreate,
		Update: resourceDatadogLogsIntegrationPipelineUpdate,
		Read:   resourceDatadogLogsIntegrationPipelineRead,
		Delete: resourceDatadogLogsIntegrationPipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"is_enabled": {Type: schema.TypeBool, Optional: true},
		},
	}
}

func resourceDatadogLogsIntegrationPipelineCreate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("cannot create an integration pipeline, please import it first to make changes")
}

func resourceDatadogLogsIntegrationPipelineRead(d *schema.ResourceData, meta interface{}) error {
	providerConf := meta.(*ProviderConfiguration)
	datadogClientV1 := providerConf.DatadogClientV1
	authV1 := providerConf.AuthV1
	ddPipeline, httpresp, err := datadogClientV1.LogsPipelinesApi.GetLogsPipeline(authV1, d.Id()).Execute()
	if err != nil {
		if httpresp != nil && httpresp.StatusCode == 400 {
			d.SetId("")
			return nil
		}
		return translateClientError(err, "error getting logs integration pipeline")
	}
	if !ddPipeline.GetIsReadOnly() {
		d.SetId("")
		return nil
	}
	if err := d.Set("is_enabled", ddPipeline.GetIsEnabled()); err != nil {
		return err
	}
	return nil
}

func resourceDatadogLogsIntegrationPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	var ddPipeline datadogV1.LogsPipeline
	ddPipeline.SetIsEnabled(d.Get("is_enabled").(bool))
	providerConf := meta.(*ProviderConfiguration)
	datadogClientV1 := providerConf.DatadogClientV1
	authV1 := providerConf.AuthV1
	updatedPipeline, _, err := datadogClientV1.LogsPipelinesApi.UpdateLogsPipeline(authV1, d.Id()).Body(ddPipeline).Execute()
	if err != nil {
		return translateClientError(err, "error updating logs integration pipeline")
	}
	d.SetId(*updatedPipeline.Id)
	return resourceDatadogLogsIntegrationPipelineRead(d, meta)
}

func resourceDatadogLogsIntegrationPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
