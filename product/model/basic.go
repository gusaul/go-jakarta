package model

import (
	"fmt"
	"strconv"
)

// BasicReg - standard type for attribute registry
type BasicReg struct{}

func (r *BasicReg) New(productID int64) Resource {
	b := new(Basic)
	b.Key = fmt.Sprintf(coreCacheKey, productID)
	b.Fields = []string{"product_id", "shop_id", "name", "desc", "price"}

	b.Query = `
		SELECT product_id, shop_id, product_name, product_desc, product_price
		FROM gj_product
		WHERE product_id IN (?)
	`
	b.Identifier = productID
	return b
}

type Basic struct {
	ResourceGetter

	ProductID int64   `db:"product_id" cache:"product_id" setter:"SetProductID"`
	ShopID    int64   `db:"shop_id" cache:"shop_id" setter:"SetShopID"`
	Name      string  `db:"product_name" cache:"name" setter:"SetName"`
	Desc      string  `db:"product_desc" cache:"desc" setter:"SetDesc"`
	Price     float64 `db:"product_price" cache:"price" setter:"SetPrice"`
}

func (b *Basic) SetProductID(value string) (err error) {
	b.ProductID, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (b *Basic) SetShopID(value string) (err error) {
	b.ShopID, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (b *Basic) SetName(value string) error {
	b.Name = value
	return nil
}

func (b *Basic) SetDesc(value string) error {
	b.Desc = value
	return nil
}

func (b *Basic) SetPrice(value string) (err error) {
	b.Price, err = strconv.ParseFloat(value, 64)
	return err
}

func (b *Basic) GetCacheMap() []string {
	return []string{
		"product_id", strconv.FormatInt(b.ProductID, 10),
		"shop_id", strconv.FormatInt(b.ShopID, 10),
		"name", b.Name,
		"desc", b.Desc,
		"price", strconv.FormatFloat(b.Price, 'f', -1, 64),
	}
}
