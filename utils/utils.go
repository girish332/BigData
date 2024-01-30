package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrints(x interface{}) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("error while marshalling during pretty print")
	}
	fmt.Println(string(b))
}
