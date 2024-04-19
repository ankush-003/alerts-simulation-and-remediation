package alerts

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func genRandomParams() (ParamInput, int) {
	structs := []ParamInput{
		&Memory{},
		&CPU{},
		&Disk{},
		&Network{},
		&Power{},
		&Applications{},
		&Security{},
	}
	randomChoice := rand.Intn(len(structs))
	randomStruct := structs[randomChoice]

	return randomStruct, randomChoice
}

func genRandomAlert(nodes []string) AlertInput {
	params, choice := genRandomParams()
	category := fmt.Sprintf("%v", params)
	node := nodes[rand.Intn(len(nodes))]
	source := "Hardware"
	if choice > 4 {
		source = "Software"
	}
	return AlertInput{
		ID:        uuid.NewString(),
		Category:  category,
		Source:    source,
		Origin:    node,
		Params:    params,
		CreatedAt: time.Now().Format(time.DateTime),
		Handled:   false,
	}
}
