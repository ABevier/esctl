package esctl

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/mitchellh/mapstructure"
)

type ESCtl struct {
	client *elasticsearch.Client
}

type IndexSummary struct {
	Name         string
	UUID         string     `mapstructure:"uuid"`
	PrimaryStats IndexStats `mapstructure:"primaries"`
	TotalStats   IndexStats `mapstructure:"total"`
}

type IndexStats struct {
	IndexStore `mapstructure:"store"`
}

type IndexStore struct {
	Size int `mapstructure:"size_in_bytes"`
}

func NewClient() *ESCtl {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	return &ESCtl{
		client: es,
	}
}

// DoSomething should do something when called :troll:
func (es *ESCtl) DoSomething() {
	fmt.Println(elasticsearch.Version)
	fmt.Println(es.client.Info())

	res, err := es.client.Indices.Stats()
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		fmt.Printf("status code= %v \n", res.StatusCode)
	}

	fmt.Println(res.Body)

	var body map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		panic(err)
	}

	//fmt.Printf("body= %v\n", body["indices"])

	printKeys(body["indices"].(map[string]interface{}))

	stats := decode(body["indices"].(map[string]interface{}))
	fmt.Printf("stats=%v\n", stats)
}

func decode(indicies map[string]interface{}) []IndexSummary {
	result := []IndexSummary{}
	for k, v := range indicies {
		stats := IndexSummary{}
		err := mapstructure.Decode(v.(map[string]interface{}), &stats)
		if err != nil {
			fmt.Printf("Got an err: %v\n", err)
			continue
		}

		stats.Name = k
		result = append(result, stats)
	}
	return result
}

func printKeys(m map[string]interface{}) {
	for k := range m {
		fmt.Println(k)
	}
}
