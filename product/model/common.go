package model

import (
	"reflect"
)

const (
	coreCacheKey  string = "core:pid:%d"
	statsCacheKey string = "stats:pid:%d"

	cacheTag  string = "cache"
	setterTag string = "setter"
)

// Resource - common interface to get data properties
type Resource interface {
	// redis cache
	GetCacheKey() string
	GetCacheFields() []string
	ApplyCache(interface{}, map[string]string) bool
	GetCacheMap() []string

	// persistent DB
	GetQuery() string
	GetIdentifier() int64
	PostQueryProcess() error
}

// ResourceGetter - generic implementation for Resource Interface
type ResourceGetter struct {
	Key    string
	Fields []string

	Query      string
	Identifier int64
}

func (c *ResourceGetter) GetCacheKey() string {
	return c.Key
}

func (c *ResourceGetter) GetCacheFields() []string {
	return c.Fields
}

// ApplyCache - to fill in map cache to struct
// it will find map key based on struct field tag and use setter func to assign value
func (c *ResourceGetter) ApplyCache(dest interface{}, data map[string]string) (isCompleted bool) {
	var applied int

	// get reflection type and value from struct parent
	objType := reflect.TypeOf(dest)
	objVals := reflect.ValueOf(dest)
	for i := 0; i < objType.Elem().NumField(); i++ {
		// get struct field attribute
		cacheField := objType.Elem().Field(i).Tag.Get(cacheTag)
		setterName := objType.Elem().Field(i).Tag.Get(setterTag)
		// ensue field has cache and setter tag
		if cacheField == "" || setterName == "" {
			continue
		}

		// ensure setter method exist before invoke
		_, setterExist := objType.MethodByName(setterName)
		if !setterExist {
			continue
		}

		//  method setter must has only 1 string arg and 1 return field
		setterMethod := objVals.MethodByName(setterName)
		if setterMethod.Type().NumIn() != 1 || setterMethod.Type().In(0).Kind() != reflect.String || setterMethod.Type().NumOut() != 1 {
			continue
		}

		// get value from cache map
		cacheVal, cacheExist := data[cacheField]
		if !cacheExist {
			continue
		}

		// invoke setter method
		result := setterMethod.Call([]reflect.Value{reflect.ValueOf(cacheVal)})
		if len(result) > 0 && result[0].IsNil() {
			applied++
		}

	}
	return applied == len(c.Fields)
}

// GetCacheMap - abstract method, override this func on child
// use to get map of struct field for set cache
// value should be []string{field1, value1, field2, value2...}
func (c *ResourceGetter) GetCacheMap() []string {
	return []string{}
}

func (c *ResourceGetter) GetQuery() string {
	return c.Query
}

func (c *ResourceGetter) GetIdentifier() int64 {
	return c.Identifier
}

// PostQueryProcess - abstract method, override this func on child
// use to post processing after query db scan
// usually for formatting some fields
func (c *ResourceGetter) PostQueryProcess() error {
	return nil
}
