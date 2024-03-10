package rule_engine

import (
	"encoding/json"
	"fmt"
	"rule_engine_demo/mongo"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var knowledgeLibrary = *ast.NewKnowledgeLibrary()

// Rule input object
type RuleInput interface {
	DataKey() string
}

// Rule output object
type RuleOutput interface {
	DataKey() string
}

// configs associated with each rule
type RuleConfig interface {
	RuleName() string
	RuleInput() RuleInput
	RuleOutput() RuleOutput
}

type RuleEngineSvc struct {
}

func NewRuleEngineSvc() *RuleEngineSvc {
	// you could add your cloud provider here instead of keeping rule file in your code.
	buildRuleEngine()
	return &RuleEngineSvc{}
}

type Rule struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"desc" bson:"desc"`
	Salience    int32    `json:"salience" bson:"salience"`
	When        string   `json:"when" bson:"when"`
	Then        []string `json:"then" bson:"then"`
}

func buildRuleEngine() {
	ruleBuilder := builder.NewRuleBuilder(&knowledgeLibrary)

	// Read rule from file and build rules
	// ruleFile := pkg.NewFileResource("/home/abhayjo/Desktop/HPE_CTY/rule_engine/engine/rules.grl")
	client, ctx, cancelFunc, err := mongo.Connect("mongodb://127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer mongo.Close(client, ctx, cancelFunc)

	db := client.Database("hpe_cty")
	collection := db.Collection("Rules")
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

	// rules, _ := os.Open("/home/abhayjo/Desktop/HPE_CTY/rule_engine/engine/rules.json")
	// ruleResource := pkg.NewReaderResource(rules)
	// ruleFile, _ := pkg.NewJSONResourceFromResource(ruleResource)
	// err = ruleBuilder.BuildRuleFromResource("Rules", "0.0.1", ruleFile)

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

	// create rule engine and execute on provided data and knowledge base
	ruleEngine := engine.NewGruleEngine()
	err = ruleEngine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		return err
	}
	return nil
}
