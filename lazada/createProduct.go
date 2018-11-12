package lazada

import "encoding/xml"

type ProductRequest struct {
	XMLName xml.Name `xml:"Request,omitempty"`
	Product *Product `xml:"Product,omitempty"`
}

type Attributes struct {
	XMLName xml.Name `xml:"Attributes,omitempty"`
	Attrs StringMap
}

type Images struct {
	XMLName xml.Name `xml:"Images,omitempty"`
	Image []string `xml:"Image,omitempty"`
}

type Product struct {
	XMLName xml.Name `xml:"Product,omitempty"`
	Attributes *Attributes `xml:"Attributes,omitempty"`
	PrimaryCategory string `xml:"PrimaryCategory,omitempty"`
	Skus []*Sku `xml:"Skus>Sku,omitempty"`
}

type Sku struct {
	XMLName   xml.Name `xml:"Sku,omitempty"`
	Images    *Images `xml:"Images,omitempty"`
	SellerSku string `xml:"SellerSku"`
	SkuAttrs  StringMap
}
