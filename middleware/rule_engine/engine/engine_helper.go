package rule_engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
)

type Pair struct {
	first  alerts.AlertOutput
	second int
}

type Key struct {
	category string
	node     string
}

type Cache struct {
	Map map[Key]Pair
}

func (cache *Cache) New() Cache {
	cache.Map = make(map[Key]Pair)
	return *cache
}

type AlertContext struct {
	AlertInput  *alerts.AlertInput
	AlertOutput *alerts.AlertOutput
	AlertParam  *alerts.ParamInput
}

func (alertContext AlertContext) RuleInput() RuleInput {
	return alertContext.AlertInput
}

func (alertContext AlertContext) RuleOutput() RuleOutput {
	return alertContext.AlertOutput
}

func (alertContext AlertContext) ParamInput() ParamInput {
	return *alertContext.AlertParam
}

func (alertContext AlertContext) NotifyRestServer() {
	jsonData := map[string]interface{}{
		"ID":        alertContext.AlertInput.ID,
		"Category":  alertContext.AlertInput.Category,
		"Source":    alertContext.AlertInput.Source,
		"Origin":    alertContext.AlertInput.Origin,
		"CreatedAt": alertContext.AlertInput.CreatedAt,
		"Severity":  alertContext.AlertOutput.Severity,
		"Remedy":    alertContext.AlertOutput.Remedy,
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	req, err := http.NewRequest("POST", "http://"+host+":8000/postRemedy", bytes.NewBuffer(jsonBytes))

	if err != nil {
		fmt.Println("Error in creating request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in Connecting to Rest Server: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Convert the response body to a string and print it
	fmt.Println("Response Body:", string(body))
}

func CacheChecker(category, node string, alertOutput *alerts.AlertOutput, cache *Cache) bool {
	key := Key{category: category, node: node}
	pair, exists := cache.Map[key]
	if !exists {
		pair = Pair{
			first:  *alertOutput,
			second: 1,
		}
		cache.Map[key] = pair
		return false
	} else if pair.first == *alertOutput {
		pair.second++
		cache.Map[key] = pair
		return pair.second > 4
	} else {
		pair.first = *alertOutput
		pair.second = 1
		cache.Map[key] = pair
		return false
	}
}

func PrintStruct(alert *alerts.AlertInput, output *alerts.AlertOutput) {
	fmt.Println("------------ALERT--------------------")
	fmt.Println("ID: ", alert.ID)
	fmt.Println("Category: ", alert.Category)
	fmt.Println("Origin: ", alert.Origin)
	fmt.Println("Source: ", alert.Source)
	b, err := json.MarshalIndent(alert.Params, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print("Params: ", string(b))
	if output != nil {
		fmt.Println("--------------OUTPUT-----------------")
		fmt.Println("Severity -> ", output.Severity)
		fmt.Println("Remedy -> ", output.Remedy)
	}
	fmt.Println("-------------------------------------")
}
