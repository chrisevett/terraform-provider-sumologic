package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/erikvanbrakel/terraform-provider-sumologic/go-sumologic"
	"strconv"
)

func resourceSumologicPollingSource() *schema.Resource {
	return &schema.Resource {
		Create: resourceSumologicPollingSourceCreate,
		Read: resourceSumologicPollingSourceRead,
		Delete: resourceSumologicPollingSourceDelete,

		Schema: map[string]*schema.Schema {
			"name" : {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"category" : {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content_type" : {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scan_interval" : {
				Type: schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"paused" : {
				Type: schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"collector_id" : {
				Type: schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"authentication" : {
				Type: schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource {
					Schema: map[string]*schema.Schema {
						"access_key" : {
							Type: schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"secret_key" : {
							Type: schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"path" : {
				Type: schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type: schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"path_expression": {
							Type: schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceSumologicPollingSourceCreate(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*sumologic.SumologicClient)

	authSettings := sumologic.PollingAuthentication{}

	auths := d.Get("authentication").([]interface{})

	if len(auths) > 0 {
		auth := auths[0].(map[string]interface{})
		authSettings.Type = "S3BucketAuthentication"
		authSettings.AwsId = auth["access_key"].(string)
		authSettings.AwsKey = auth["secret_key"].(string)
	}

	pathSettings := sumologic.PollingPath{}
	paths := d.Get("path").([]interface{})

	if len(paths) > 0 {
		path := paths[0].(map[string]interface{})
		pathSettings.Type = "S3BucketPathExpression"
		pathSettings.BucketName = path["bucket_name"].(string)
		pathSettings.PathExpression = path["path_expression"].(string)
	}

	sourceId, err := c.CreatePollingSource(
		d.Get("name").(string),
		d.Get("content_type").(string),
		d.Get("category").(string),
		d.Get("scan_interval").(int),
		d.Get("paused").(bool),
		d.Get("collector_id").(int),
		authSettings,
		pathSettings,
	)

	if err != nil {
		return err
	}


	id := strconv.Itoa(sourceId)

	d.SetId(id)

	return resourceSumologicPollingSourceRead(d, meta)
}

func resourceSumologicPollingSourceRead(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*sumologic.SumologicClient)

	id, err := strconv.Atoi(d.Id())
	collector_id := d.Get("collector_id").(int)

	if err != nil {
		return err
	}

	_, err = c.GetPollingSource(collector_id, id)


	return nil
}

func resourceSumologicPollingSourceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
