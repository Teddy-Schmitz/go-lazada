package lazada

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

// The Product Service deals with any methods under the "Product" category of the open platform
type ProductService service

// A brand object returned from the open platform
type Brand struct {
	Name             string `json:"name"`
	BrandID          int    `json:"brand_id"`
	GlobalIdentifier string `json:"global_identifier"`
}

// Brands returns a list of brands in the region set
// If opts is nil then the default options are used
func (p *ProductService) Brands(ctx context.Context, opts *ListOptions) ([]*Brand, error) {
	if opts == nil {
		opts = &DefaultListOptions
	}

	u, err := addOptions(apiNames["GetBrands"], opts)
	if err != nil {
		return nil, err
	}

	req, err := p.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	brands := []*Brand{}
	_, err = p.client.Do(ctx, req, &brands)
	if err != nil {
		return nil, err
	}

	return brands, nil
}

type CategoryTree struct {
	CategoryID int             `json:"category_id"`
	Children   []*CategoryTree `json:"children"`
	Var        bool            `json:"var"`
	Name       string          `json:"name"`
	Leaf       bool            `json:"leaf"`
}

// CategoryTree returns all the categories available in the region set
func (p *ProductService) CategoryTree(ctx context.Context) ([]*CategoryTree, error) {
	req, err := p.client.NewRequest("GET", apiNames["CategoryTree"], nil)
	if err != nil {
		return nil, err
	}

	tree := []*CategoryTree{}
	_, err = p.client.Do(ctx, req, &tree)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

type ImageReq struct {
	XMLName xml.Name `xml:"Request"`
	URL     string   `xml:"Image>Url"`
}

type ImageResponse struct {
	Image struct {
		HashCode string `json:"hash_code"`
		URL      string `json:"url"`
	}
}

// MigrateImage lets you move any publicly accessible image into the Lazada platform
// Requires a client access token
func (p *ProductService) MigrateImage(ctx context.Context, imgURL string) (*ImageResponse, error) {
	if p.client.accessToken == "" {
		return nil, errors.New("an access token is required for this api call")
	}

	req, err := p.client.NewRequest("POST", apiNames["ImageMigrate"], ImageReq{URL: imgURL})
	if err != nil {
		return nil, err
	}

	img := &ImageResponse{}
	_, err = p.client.Do(ctx, req, img)
	if err != nil {
		return nil, err
	}

	return img, nil
}

type CategoryAttributes struct {
	Label         string             `json:"label"`
	Name          string             `json:"name"`
	IsMandatory   int                `json:"is_mandatory"`
	AttributeType string             `json:"attribute_type"`
	InputType     string             `json:"input_type"`
	Options       []*CategoryOptions `json:"options"`
	IsSale        int                `json:"is_sale_prop"`
}

type CategoryOptions struct {
	Name string
}

// CategoryAttributes returns all the attributes related to the category id provided
func (p *ProductService) CategoryAttributes(ctx context.Context, id int) ([]*CategoryAttributes, error) {
	req, err := p.client.NewRequest("GET", fmt.Sprintf("%s?primary_category_id=%d", apiNames["CategoryAttributes"], id), nil)
	if err != nil {
		return nil, err
	}

	attr := []*CategoryAttributes{}
	_, err = p.client.Do(ctx, req, &attr)
	if err != nil {
		return nil, err
	}

	return attr, nil
}

type CreateProductResponse struct {
	ItemID  int64     `json:"item_id"`
	SKUList []SKUItem `json:"sku_list"`
}

type SKUItem struct {
	ShopSKU   string `json:"shop_sku"`
	SellerSKU string `json:"seller_sku"`
	SKUID     string `json:"sku_id"`
}

// Create lets you create a new product on the open platform.
//
//
//
// Requires a client access token
func (p *ProductService) Create(ctx context.Context, pReq *Product) (*CreateProductResponse, error) {
	if p.client.accessToken == "" {
		return nil, errors.New("an access token is required for this api call")
	}

	req, err := p.client.NewRequest("POST", apiNames["CreateProduct"], &ProductRequest{Product: pReq})
	if err != nil {
		return nil, err
	}

	resp := &CreateProductResponse{}
	_, err = p.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update lets you update an existing product on the open platform
// Requires a client access token
func (p *ProductService) Update(ctx context.Context, pReq *Product) error {
	if p.client.accessToken == "" {
		return errors.New("an access token is required for this api call")
	}

	req, err := p.client.NewRequest("POST", apiNames["UpdateProduct"], &ProductRequest{Product: pReq})
	if err != nil {
		return err
	}

	_, err = p.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

type SearchOptions struct {
	// Filter the product status
	Filter string `url:"filter"`

	// Search for products with this name or seller sku
	Search *string `url:"search,omitempty"`

	// Set to "1" to get more stock information
	Options *string `url:"options,omitempty"`

	// Return only products that have this seller sku
	SKUSellerList *string `url:"sku_seller_list,omitempty"`

	// Offset the results by
	Offset int `url:"offset"`

	// Limit the amount of returned results
	Limit int `url:"limit"`
}

type ProductSKU struct {
	Status          string          `json:"Status"`
	Quantity        int             `json:"quantity"`
	ProductWeight   string          `json:"product_weight"`
	Images          []string        `json:"Images"`
	SellerSKU       string          `json:"SellerSku"`
	ShopSKU         string          `json:"ShopSku"`
	URL             string          `json:"Url"`
	PackageWidth    string          `json:"package_width"`
	SpecialToTime   string          `json:"special_to_time"`
	SpecialFromTime string          `json:"special_from_time"`
	PackageHeight   string          `json:"package_height"`
	SpecialPrice    decimal.Decimal `json:"special_price"`
	Price           decimal.Decimal `json:"price"`
	PackageLength   string          `json:"package_length"`
	PackageWeight   string          `json:"package_weight"`
	Available       int             `json:"Available"`
	SkuID           int             `json:"SkuId"`
	SpecialToDate   string          `json:"special_to_date"`
}

type GetProduct struct {
	ItemID          int               `json:"item_id"`
	PrimaryCategory int               `json:"primary_category"`
	Attributes      map[string]string `json:"attributes"`
	SKUs            []*ProductSKU     `json:"skus"`
}

type GetProductResponse struct {
	TotalProducts int           `json:"total_products"`
	Products      []*GetProduct `json:"products"`
}

// Get lets you retrieve all products in a specific region
func (p *ProductService) Get(ctx context.Context, opts *SearchOptions) (*GetProductResponse, error) {
	if opts == nil {
		opts = &SearchOptions{
			Offset: DefaultListOptions.Offset,
			Limit:  DefaultListOptions.Limit,
		}
	}

	if opts.Filter == "" {
		opts.Filter = "live"
	}

	u, err := addOptions(apiNames["GetProducts"], opts)
	if err != nil {
		return nil, err
	}

	req, err := p.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp := &GetProductResponse{}
	_, err = p.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
