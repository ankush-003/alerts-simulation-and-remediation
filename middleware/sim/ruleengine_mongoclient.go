package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	//"rule_engine_demo/mongo"
	"asmr/mailserver"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AlertInput struct {
	ID        string
	Category  string
	Source    string
	Origin    string
	Params    ParamInput
	CreatedAt time.Time
	Handled   bool
}

type AlertContext struct {
	AlertInput  *AlertInput
	AlertOutput *AlertOutput
	AlertParam  *ParamInput
}

func (alertContext *AlertContext) RuleInput() RuleInput {
	return alertContext.AlertInput
}

func (alertContext *AlertContext) RuleOutput() RuleOutput {
	return alertContext.AlertOutput
}

func (alertContext *AlertContext) ParamInput() ParamInput {
	return *alertContext.AlertParam
}

func (alert *AlertInput) DataKey() string {
	return "AlertInput"
}

type AlertOutput struct {
	Severity string
	Remedy   string
}

func (alert *AlertOutput) DataKey() string {
	return "AlertOutput"
}

type Memory struct {
	Usage      uint
	PageFaults uint
	SwapUsge   uint
}

func (mem *Memory) DataKey() string {
	return "MemInput"
}

type CPU struct {
	Utilization uint
	Temperature uint
}

func (cpu *CPU) DataKey() string {
	return "CpuInput"
}

type Disk struct {
	Usage       uint
	IOPs        uint
	ThroughtPut uint
}

func (disk *Disk) DataKey() string {
	return "DiskInput"
}

type Network struct {
	Traffic    uint
	PacketLoss uint
	Latency    uint
}

func (net *Network) DataKey() string {
	return "NetworkInput"
}

type Power struct {
	BatteryLevel uint
	Consumption  uint
	Efficiency   uint
}

func (pow *Power) DataKey() string {
	return "PowerInput"
}

type Applications struct {
	Processes   uint
	MaxCPUusage uint
	MaxMemUsage uint
}

func (app *Applications) DataKey() string {
	return "ApplicationsInput"
}

type Security struct {
	LoginAttempts  uint
	FailedLogins   uint
	SuspectedFiles uint
	IDSEvents      uint
}

func (sec *Security) DataKey() string {
	return "SecurityInput"
}

var knowledgeLibrary = *ast.NewKnowledgeLibrary()

// Rule input object
type RuleInput interface {
	DataKey() string
}

// Rule output object
type RuleOutput interface {
	DataKey() string
}

type ParamInput interface {
	DataKey() string
}

// configs associated with each input
type RuleConfig interface {
	RuleInput() RuleInput
	RuleOutput() RuleOutput
	ParamInput() ParamInput
}

type RuleEngineSvc struct {
}

func NewRuleEngineSvc() *RuleEngineSvc {
	buildRuleEngine()
	return &RuleEngineSvc{}
}

func buildRuleEngine() {
	ruleBuilder := builder.NewRuleBuilder(&knowledgeLibrary)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("rules_sample1")
	collection := db.Collection("intelligent")
	curr, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(bson.D{{Key: "_id", Value: 0}}))
	if err != nil {
		panic(err)
	}
	var results []bson.M
	curr.All(ctx, &results)
	rules, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}

	ruleset, _ := pkg.ParseJSONRuleset(rules)
	fmt.Println(ruleset)
	ruleFile := pkg.NewBytesResource([]byte(ruleset))
	err = ruleBuilder.BuildRuleFromResource("Rules", "0.0.1", ruleFile)

	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}

}

func (svc *RuleEngineSvc) Execute(ruleConf RuleConfig) error {
	// get KnowledgeBase instance to execute particular rule
	knowledgeBase, _ := knowledgeLibrary.NewKnowledgeBaseInstance("Rules", "0.0.1")

	dataCtx := ast.NewDataContext()
	// add input data context
	err := dataCtx.Add(ruleConf.RuleInput().DataKey(), ruleConf.RuleInput())
	if err != nil {
		return err
	}

	// add output data context
	err = dataCtx.Add(ruleConf.RuleOutput().DataKey(), ruleConf.RuleOutput())
	if err != nil {
		return err
	}

	err = dataCtx.Add(ruleConf.ParamInput().DataKey(), ruleConf.ParamInput())
	if err != nil {
		return err
	}
	//fmt.Println("HEYY HI")
	// create rule engine and execute on provided data and knowledge base
	ruleEngine := engine.NewGruleEngine()
	err = ruleEngine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		return err
	}
	//fmt.Println("Done")
	return nil
}

func NewAlert(alertInput *AlertInput, ruleEngineSvc *RuleEngineSvc) {

	defer wg.Done()
	alertContext := AlertContext{
		alertInput,
		&AlertOutput{Remedy: "Too be decided soon...", Severity: "NIL"},
		&alertInput.Params,
	}

	err := ruleEngineSvc.Execute(&alertContext)
	if err != nil {
		panic(err)
	}
	//fmt.Println("Here")
	fmt.Println("Alert -> ", alertInput)
	fmt.Println("Severity -> ", alertContext.AlertOutput.Severity)
	fmt.Println("Remedy -> ", alertContext.AlertOutput.Remedy)

	mailserver.SendEmail(alertInput.ID, alertInput.Category, alertInput.CreatedAt, alertInput.Handled, alertInput.Source, alertInput.Origin, alertContext.AlertOutput.Severity, alertContext.AlertOutput.Remedy, nil)
	
}

var wg sync.WaitGroup

func main() {
	ruleEngineSvc := NewRuleEngineSvc()
	//memoryUsage := 95.0
	userA := AlertInput{
		ID:        "Client1",
		Category:  "Memory",
		Source:    "Hardware",
		Origin:    "NodeA",
		Params:    &Memory{Usage: 10, PageFaults: 30, SwapUsge: 2},
		CreatedAt: time.Now(),
		Handled:   false,
	}

	userB := AlertInput{
		ID:        "Client2",
		Category:  "CPU",
		Source:    "CPU",
		Origin:    "NodeB",
		Params:    &CPU{Utilization: 71, Temperature: 30},
		CreatedAt: time.Now(),
		Handled:   false,
	}
	wg.Add(1)
	wg.Add(1)
	go NewAlert(&userA, ruleEngineSvc)
	go NewAlert(&userB, ruleEngineSvc)
	wg.Wait()
}
