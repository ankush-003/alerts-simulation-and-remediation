package alerts

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/google/uuid"
)

func genRandomParams(r *rand.Rand) (ParamInput, int) {
	structs := []ParamInput{
		&Memory{},
		&CPU{},
		&Disk{},
		&Network{},
		&Power{},
		&Applications{},
		&Security{},
	}
	randomChoice := r.Intn(len(structs))
	randomStruct := structs[randomChoice]

	return randomStruct, randomChoice
}

func GenRandomAlert(nodes []string) AlertInput {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	params, choice := genRandomParams(r)
	category := reflect.TypeOf(params).Elem().Name()
	node := nodes[r.Intn(len(nodes))]
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
