package main

import (
	"fmt"
	"reflect"

	"github.com/gusaul/go-jakarta/product"
	"github.com/gusaul/go-jakarta/resource"
)

func main() {
	resource.InitDatabase()
	resource.InitRedisConn()

	data, err := product.GetProductData([]int64{123, 234, 345, 456, 567}, product.BasicOpt|product.StatsOpt|product.PictureOpt)

	fmt.Println("Error:", err)
	for _, d := range data {
		print(d.Basic)
		print(d.Picture)
		print(d.Stats)
	}
}

// just dummy func for print struct primary field only
func print(v interface{}) {
	objType := reflect.TypeOf(v)
	valType := reflect.ValueOf(v)
	for i := 0; i < objType.NumField(); i++ {
		if objType.Field(i).Tag.Get("cache") == "" {
			continue
		}
		fmt.Print(objType.Field(i).Name, "\t: ")
		fmt.Println(valType.Field(i).Interface())
	}
}
