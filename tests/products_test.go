// +build integration

package tests

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/teddy-schmitz/go-lazada/lazada"
)

var AppKey = os.Getenv("APP_KEY")
var AppSecret = os.Getenv("APP_SECRET")

func TestProductBrands(t *testing.T) {
	c := lazada.NewClient(AppKey, AppSecret, lazada.Singapore)
	br, err := c.Products.Brands(context.Background(), nil)
	require.NoError(t, err)

	assert.Len(t, br, 100)
}

func TestListOptions(t *testing.T) {
	c := lazada.NewClient(AppKey, AppSecret, lazada.Singapore)
	br, err := c.Products.Brands(context.Background(), &lazada.ListOptions{Limit: 50})
	require.NoError(t, err)
	assert.Len(t, br, 50)

	first := br[0]

	br, err = c.Products.Brands(context.Background(), &lazada.ListOptions{Limit: 120, Offset: 50})
	require.NoError(t, err)
	assert.Len(t, br, 120)
	assert.NotEqual(t, first, br[0])
}

func TestCategoryTree(t *testing.T) {
	c := lazada.NewClient(AppKey, AppSecret, lazada.Singapore)
	br, err := c.Products.CategoryTree(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, br)
}

func TestCode(t *testing.T) {
	t.SkipNow()
	c := lazada.NewClient(AppKey, AppSecret, lazada.Singapore)
	_, err := c.Auth.Exchange(context.Background(), "")
	require.NoError(t, err)
}

func TestRefresh(t *testing.T) {
	t.SkipNow()
	c := lazada.NewClient(AppKey, AppSecret, lazada.Singapore)
	_, err := c.Auth.Refresh(context.Background(), "50001500f08BWVnreeai1ffde6e2Tekwfrxk6eFpBqrzHKdtH1izvoEZQDQrd")
	require.NoError(t, err)
}

func TestCategoryAttributes(t *testing.T) {
	c := lazada.NewClient(AppKey, AppSecret, lazada.Singapore)
	br, err := c.Products.CategoryAttributes(context.Background(), 10001996)
	require.NoError(t, err)
	assert.NotEmpty(t, br)
}

func TestImageMigrate(t *testing.T) {
	t.SkipNow()
	cl := lazada.NewClient(AppKey, AppSecret, lazada.Singapore)
	c := cl.NewTokenClient("")

	resp, err := c.Products.MigrateImage(context.Background(), "https://images.unsplash.com/photo-1539784257995-d70d6089a5c8?ixlib=rb-0.3.5&ixid=eyJhcHBfaWQiOjEyMDd9&s=898c8ca83eebb1022781d179e39bb44a&auto=format&fit=crop&w=2134&q=80")
	require.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestCreateProduct(t *testing.T) {
	t.SkipNow()
	cl := lazada.NewClient(AppKey, AppSecret, lazada.Malaysia)
	c := cl.NewTokenClient("")

	product := &lazada.Product{
		PrimaryCategory: "10001958",
		Attributes: &lazada.Attributes{
			Attrs: lazada.StringMap{
				"name":               "test product creation",
				"short_description":  "test product highlights",
				"description":        "test product description",
				"brand":              "Kid Basix",
				"model":              "test model",
				"recommended_gender": "Men",
				"material":           "Cotton",
				"waterproof":         "waterproof",
				"warranty_type":      "Local (Singapore) manufacturer warranty",
				"warranty":           "1 month",
				"Hazmat":             "Battery, Flammable",
			}},
		Skus: []*lazada.Sku{
			{
				SellerSku: "test-product-creation-for-api",
				Images:    &lazada.Images{Image: []string{"https://sg-live.slatic.net/original/b731a8098df7d606ab2e56efc650afcb.jpg"}},
				SkuAttrs: lazada.StringMap{
					"quantity":        "1",
					"color_family":    "Black",
					"special_price":   "0.0",
					"price":           "23.0",
					"package_length":  "1",
					"package_weight":  "1",
					"package_content": "test whats in the box",
					"package_width":   "1",
					"package_height":  "1",
				},
			},
		},
	}

	br, err := c.Products.Create(context.Background(), product)
	require.NoError(t, err)
	assert.NotEmpty(t, br)
}

func TestGetProduct(t *testing.T) {
	t.SkipNow()
	cl := lazada.NewClient(AppKey, AppSecret, lazada.Malaysia)
	c := cl.NewTokenClient("")

	out := lazada.SliceString([]string{"16016131915889"})

	resp, err := c.Products.Get(context.Background(), &lazada.SearchOptions{Filter: "live", Limit: 100, SKUSellerList: &out})
	require.NoError(t, err)
	assert.NotEmpty(t, resp)
}
