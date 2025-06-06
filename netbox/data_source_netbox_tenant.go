package netbox

import (
	"errors"
	"strconv"

	"github.com/fbreckle/go-netbox/netbox/client/tenancy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetboxTenant() *schema.Resource {
	return &schema.Resource{
		Read:        dataSourceNetboxTenantRead,
		Description: `:meta:subcategory:Tenancy:`,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				AtLeastOneOf: []string{"name", "slug"},
			},
			"slug": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"name", "slug"},
			},
			"group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNetboxTenantRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*providerState)
	params := tenancy.NewTenancyTenantsListParams()

	if name, ok := d.Get("name").(string); ok && name != "" {
		params.Name = &name
	}

	if slug, ok := d.Get("slug").(string); ok && slug != "" {
		params.Slug = &slug
	}

	limit := int64(2) // Limit of 2 is enough
	params.Limit = &limit

	res, err := api.Tenancy.TenancyTenantsList(params, nil)
	if err != nil {
		return err
	}

	if *res.GetPayload().Count > int64(1) {
		return errors.New("more than one tenant returned, specify a more narrow filter")
	}
	if *res.GetPayload().Count == int64(0) {
		return errors.New("no tenant found matching filter")
	}
	result := res.GetPayload().Results[0]
	d.SetId(strconv.FormatInt(result.ID, 10))
	d.Set("name", result.Name)
	d.Set("slug", result.Slug)
	d.Set("description", result.Description)
	if result.Group != nil {
		d.Set("group_id", result.Group.ID)
	}
	return nil
}
