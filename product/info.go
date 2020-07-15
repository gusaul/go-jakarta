package product

import (
	"fmt"
	"log"

	"github.com/gusaul/go-jakarta/product/model"
	"github.com/gusaul/go-jakarta/resource"
	"github.com/jmoiron/sqlx"
)

type Registrar interface {
	New(int64) model.Resource
}

type Options int64

const (
	BasicOpt Options = 1 << iota
	PictureOpt
	StatsOpt
)

type Getter struct {
	reg      Registrar
	prop     model.Resource
	infoType Options
	cacheKey string
}

type QueryGroup struct {
	ids   []int64
	props map[int64]model.Resource
}

type ProductData struct {
	Basic   model.Basic
	Stats   model.Stats
	Picture model.Pictures
}

func GetProductData(productIDs []int64, o Options) (result []ProductData, err error) {

	infoGetter := make(map[int64][]Getter)
	cacheKeys := make(map[string][]string)

	// populate getter object base on registrar
	// generate cache key to append field into same parent key
	for _, pid := range productIDs {
		for _, g := range getSelectedInfo(o) {
			g.prop = g.reg.New(pid)
			g.cacheKey = g.prop.GetCacheKey()
			infoGetter[pid] = append(infoGetter[pid], g)
			cacheKeys[g.cacheKey] = append(cacheKeys[g.cacheKey], g.prop.GetCacheFields()...)
		}
	}

	cacheResult, err := getFromCache(cacheKeys)
	if err != nil {
		log.Println(err)
	}

	queryGroups := make(map[string]*QueryGroup)

	// apply cache result data to every struct fields
	// if any failed process, fallback to database
	for _, pid := range productIDs {
		for i := range infoGetter[pid] {
			getter := infoGetter[pid][i]
			key := getter.cacheKey
			var isCompleted bool
			if res, ok := cacheResult[key]; ok {
				isCompleted = getter.prop.ApplyCache(getter.prop, res)
			}

			if !isCompleted {
				// register incomplete redis result to queryGroup
				query := getter.prop.GetQuery()
				id := getter.prop.GetIdentifier()
				if _, exist := queryGroups[query]; exist {
					queryGroups[query].ids = append(queryGroups[query].ids, id)
					queryGroups[query].props[id] = getter.prop
				} else {
					queryGroups[query] = &QueryGroup{
						ids: []int64{id},
						props: map[int64]model.Resource{
							id: getter.prop,
						},
					}
				}
			}
		}
	}

	if len(queryGroups) > 0 {
		err = getFromDatabase(queryGroups)
		if err != nil {
			log.Println(err)
			return result, err
		}

		go setCache(queryGroups)
	}

	return castAttributeType(infoGetter), nil
}

func getSelectedInfo(o Options) (registry []Getter) {

	if o&BasicOpt > 0 {
		registry = append(registry, Getter{
			reg:      new(model.BasicReg),
			infoType: BasicOpt,
		})
	}

	if o&PictureOpt > 0 {
		registry = append(registry, Getter{
			reg:      new(model.PictureReg),
			infoType: PictureOpt,
		})
	}

	if o&StatsOpt > 0 {
		registry = append(registry, Getter{
			reg:      new(model.StatsReg),
			infoType: StatsOpt,
		})
	}

	return
}

func getFromCache(mapKeys map[string][]string) (map[string]map[string]string, error) {

	result, err := resource.RedisConn.MultiHashGetPipeline(mapKeys)
	if err != nil {
		log.Println(err)
	}

	return result, err
}

func getFromDatabase(queryGroup map[string]*QueryGroup) error {
	fmt.Println("GET FROM DATABASE...")

	for query, mapper := range queryGroup {
		q, args, err := sqlx.In(query, mapper.ids)
		if err != nil {
			log.Println(err)
			return err
		}
		rows, err := resource.DatabaseConn.Queryx(resource.DatabaseConn.Rebind(q), args...)
		if err == nil && rows != nil {
			defer rows.Close()
			for rows.Next() {
				res, err := rows.SliceScan()
				if err != nil || len(res) < 1 {
					log.Println(err)
					return err
				}
				if id, ok := res[0].(int64); ok {
					if _, exist := mapper.props[id]; exist {
						target := mapper.props[id]
						err := rows.StructScan(target)
						if err != nil {
							log.Println("err scan", err)
							return err
						}

						err = target.PostQueryProcess()
						if err != nil {
							log.Println(err)
							return err
						}
					}
				}
			}
		} else if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func setCache(queryGroup map[string]*QueryGroup) {
	cacheData := make(map[string][]string)
	for _, v := range queryGroup {
		for _, prop := range v.props {
			cacheData[prop.GetCacheKey()] = append(cacheData[prop.GetCacheKey()], prop.GetCacheMap()...)
		}
	}

	errs := resource.RedisConn.MultiHashSetPipeline(cacheData)
	if len(errs) > 0 {
		log.Println("Errors", errs)
	}
}

func castAttributeType(infoGetter map[int64][]Getter) []ProductData {
	result := make([]ProductData, len(infoGetter))
	i := 0
	for _, val := range infoGetter {
		info := ProductData{}
		for _, data := range val {
			switch data.infoType {
			case BasicOpt:
				if val, ok := data.prop.(*model.Basic); ok && val != nil {
					info.Basic = *val
				}
			case StatsOpt:
				if val, ok := data.prop.(*model.Stats); ok && val != nil {
					info.Stats = *val
				}
			case PictureOpt:
				if val, ok := data.prop.(*model.Pictures); ok && val != nil {
					info.Picture = *val
				}
			}
		}

		result[i] = info
		i++
	}
	return result
}
