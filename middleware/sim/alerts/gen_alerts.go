package alerts

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine_v2/alerts"
	"github.com/google/uuid"
)

func genRandomParams() (alerts.ParamInput, int) {
	structs := []alerts.ParamInput{
		&alerts.Memory{},
		&alerts.CPU{},
		&alerts.Disk{},
		&alerts.Network{},
		&alerts.Power{},
		&alerts.Applications{},
		&alerts.Security{},
	}
	randomChoice := rand.Intn(len(structs))
	randomStruct := structs[randomChoice]

	return randomStruct, randomChoice
}

func genRandomAlert(nodes []string) alerts.AlertInput {
	params, choice := genRandomParams()
	category := fmt.Sprintf("%v", params)
	node := nodes[rand.Intn(len(nodes))]
	source := "Hardware"
	if choice > 4 {
		source = "Software"
	}
	return alerts.AlertInput{
		ID:        uuid.NewString(),
		Category:  category,
		Source:    source,
		Origin:    node,
		Params:    params,
		CreatedAt: time.Now().Format(time.DateTime),
		Handled:   false,
	}
}
