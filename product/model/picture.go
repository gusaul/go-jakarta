package model

import (
	"encoding/json"
	"fmt"
)

// PictureReg - standard type for attribute registry
type PictureReg struct{}

func (r *PictureReg) New(productID int64) Resource {
	p := new(Pictures)
	p.Key = fmt.Sprintf(coreCacheKey, productID)
	p.Fields = []string{"pictures"}

	p.Query = `
		SELECT product_id, picture_id, file_path, file_name
		FROM gj_picture
		WHERE product_id IN (?)
	`
	p.Identifier = productID
	return p
}

type Pictures struct {
	ResourceGetter

	//embed real struct to catch value from db
	Picture

	Data []Picture `cache:"pictures" setter:"SetPictures"`
}

type Picture struct {
	ProductID int64  `db:"product_id" json:"product_id"`
	PictureID int64  `db:"picture_id" json:"picture_id"`
	FilePath  string `db:"file_path" json:"file_path"`
	FileName  string `db:"file_name" json:"file_name"`
}

func (p *Pictures) SetPictures(value string) (err error) {
	return json.Unmarshal([]byte(value), &p.Data)
}

func (p *Pictures) GetCacheMap() []string {
	m, err := json.Marshal(p.Data)
	if err != nil {
		return []string{}
	}
	return []string{
		"pictures", string(m),
	}
}

// PostQueryProcess - processing after query per row
func (p *Pictures) PostQueryProcess() error {
	p.Data = append(p.Data, p.Picture)
	return nil
}
