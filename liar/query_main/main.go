package main

import (
	"encoding/json"
	"fmt"
	"liar/query"
)

func main() {
	cfg := query.QueryConfig{
		ServicePort:  9011,
		DomainID:     "562d704a-c442-499a-aff3-223f580bf6b3",
		ChannelID:    "b0ec13df-9ff0-48b9-9cb6-b3be072e7c99",
		Offset:       0,
		Limit:        10,
		ClientSecret: "4fa8890c-5888-48ba-93b8-3e7db1165b65",
		Timeout:      10,
	}
	result := query.FetchMessages(cfg)
	data, _ := json.MarshalIndent(result, "", "    ")
	fmt.Println(string(data))
}
