package model

import (
	"fmt"
	"strconv"
)

// StatsReg - standard type for attribute registry
type StatsReg struct{}

func (r *StatsReg) New(productID int64) Resource {
	s := new(Stats)
	s.Key = fmt.Sprintf(statsCacheKey, productID)
	s.Fields = []string{"product_id", "view", "transaction", "review", "talk"}

	s.Query = `
		SELECT product_id, view, transactions, review, talk
		FROM gj_stats
		WHERE product_id IN (?)
	`
	s.Identifier = productID
	return s
}

type Stats struct {
	ResourceGetter

	ProductID        int64 `db:"product_id" cache:"product_id" setter:"SetProductID"`
	ViewCount        int64 `db:"view" cache:"view" setter:"SetViewCount"`
	TransactionCount int64 `db:"transactions" cache:"transaction" setter:"SetTransactionCount"`
	ReviewCount      int64 `db:"review" cache:"review" setter:"SetReviewCount"`
	TalkCount        int64 `db:"talk" cache:"talk" setter:"SetTalkCount"`
}

func (s *Stats) SetProductID(value string) (err error) {
	s.ProductID, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (s *Stats) SetViewCount(value string) (err error) {
	s.ViewCount, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (s *Stats) SetTransactionCount(value string) (err error) {
	s.TransactionCount, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (s *Stats) SetReviewCount(value string) (err error) {
	s.ReviewCount, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (s *Stats) SetTalkCount(value string) (err error) {
	s.TalkCount, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (s *Stats) GetCacheMap() []string {
	return []string{
		"product_id", strconv.FormatInt(s.ProductID, 10),
		"view", strconv.FormatInt(s.ViewCount, 10),
		"transaction", strconv.FormatInt(s.TransactionCount, 10),
		"review", strconv.FormatInt(s.ReviewCount, 10),
		"talk", strconv.FormatInt(s.TalkCount, 10),
	}
}
