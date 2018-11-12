package lazada

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Signature(t *testing.T) {
	c := NewClient("123456", "testsecretnotarealsecret", Singapore)
	req, err := http.NewRequest("GET",
		"https://api.lazada.sg/rest/brands/get?offset=0&limit=100&app_key=123456&sign_method=sha256&timestamp=1537324254708", nil)

	require.NoError(t, err)
	sig := c.Signature(apiNames["GetBrands"], req.URL.Query(), nil)

	assert.Equal(t, "1A4D99631F4059D6C5F565C529916DC4D43E141412700EA371B2F396310943CA", sig)
}

func TestClient_SignatureWithPayload(t *testing.T) {
	c := NewClient("123456", "testsecretnotarealsecret", Singapore)
	c = c.NewTokenClient("faketoken")
	req, err := http.NewRequest("POST",
		"https://api.lazada.sg/rest/product/create?payload=%3C%3Fxml+version%3D%221.0%22+encoding%3D%22UTF-8%22%3F%3E%0A%3CRequest%3E%0A++++%3CProduct%3E%0A++++++++%3CAttributes%3E%0A++++++++++++%3Cname%3Etest+product+creation%3C%2Fname%3E%0A++++++++++++%3Cbrand%3EKid+Basix%3C%2Fbrand%3E%0A++++++++++++%3Cmaterial%3ECotton%3C%2Fmaterial%3E%0A++++++++++++%3Cwaterproof%3Ewaterproof%3C%2Fwaterproof%3E%0A++++++++++++%3Cwarranty_type%3ELocal+%28Singapore%29+manufacturer+warranty%3C%2Fwarranty_type%3E%0A++++++++++++%3Cwarranty%3E1+month%3C%2Fwarranty%3E%0A++++++++++++%3Cshort_description%3Etest+product+highlights%3C%2Fshort_description%3E%0A++++++++++++%3Cdescription%3Etest+product+description%3C%2Fdescription%3E%0A++++++++++++%3Cmodel%3Etest+model%3C%2Fmodel%3E%0A++++++++++++%3Crecommended_gender%3EMen%3C%2Frecommended_gender%3E%0A++++++++++++%3CHazmat%3EBattery%2C+Flammable%3C%2FHazmat%3E%0A++++++++%3C%2FAttributes%3E%0A++++++++%3CPrimaryCategory%3E10001958%3C%2FPrimaryCategory%3E%0A++++++++%3CSkus%3E%0A++++++++++++%3CSku%3E%0A++++++++++++++++%3CImages%3E%0A++++++++++++++++++++%3CImage%3Ehttps%3A%2F%2Fsg-live.slatic.net%2Foriginal%2Fb731a8098df7d606ab2e56efc650afcb.jpg%3C%2FImage%3E%0A++++++++++++++++%3C%2FImages%3E%0A++++++++++++++++%3CSellerSku%3Etest-product-creation-for-api%3C%2FSellerSku%3E%0A++++++++++++++++%3Cquantity%3E1%3C%2Fquantity%3E%0A++++++++++++++++%3Cpackage_length%3E1%3C%2Fpackage_length%3E%0A++++++++++++++++%3Cpackage_content%3Etest+whats+in+the+box%3C%2Fpackage_content%3E%0A++++++++++++++++%3Cpackage_width%3E1%3C%2Fpackage_width%3E%0A++++++++++++++++%3Cpackage_height%3E1%3C%2Fpackage_height%3E%0A++++++++++++++++%3Ccolor_family%3EBlack%3C%2Fcolor_family%3E%0A++++++++++++++++%3Cspecial_price%3E0.0%3C%2Fspecial_price%3E%0A++++++++++++++++%3Cprice%3E23.0%3C%2Fprice%3E%0A++++++++++++++++%3Cpackage_weight%3E1%3C%2Fpackage_weight%3E%0A++++++++++++%3C%2FSku%3E%0A++++++++%3C%2FSkus%3E%0A++++%3C%2FProduct%3E%0A%3C%2FRequest%3E&app_key=123456&sign_method=sha256&timestamp=1539870185083&access_token=faketoken", nil)

	require.NoError(t, err)
	sig := c.Signature(apiNames["CreateProduct"], req.URL.Query(), nil)

	assert.Equal(t, "4F912A7D7FF2B433CE5141291BA3A6B1DB2C069453927B271E6C67E414DAE1F4", sig)
}

func TestSliceString(t *testing.T) {
	out := SliceString([]string{"test"})
	assert.Equal(t, `["test"]`, out)
}

func ExampleProductService_Create() {
	client := NewClient("12345", "example", Singapore)

	userClient := client.NewTokenClient("usertoken") // Set the a token obtained through oauth
	userClient.SetRegion(Malaysia)                   // Change the region to Malaysia

	product := &Product{
		PrimaryCategory: "10001958",
		Attributes: &Attributes{
			Attrs: StringMap{
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
		Skus: []*Sku{
			{
				SellerSku: "test-product-creation-for-api",
				Images:    &Images{Image: []string{"https://sg-live.slatic.net/original/b731a8098df7d606ab2e56efc650afcb.jpg"}},
				SkuAttrs: StringMap{
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

	userClient.Products.Create(context.Background(), product)
}
