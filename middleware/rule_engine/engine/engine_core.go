package rule_engine

import (
	"fmt"
	// "os"

	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine/mongo"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
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

// type Rule struct {
// 	Name        string   `json:"name" bson:"name"`
// 	Description string   `json:"desc" bson:"desc"`
// 	Salience    int32    `json:"salience" bson:"salience"`
// 	When        string   `json:"when" bson:"when"`
// 	Then        []string `json:"then" bson:"then"`
// }

func buildRuleEngine() {
	ruleBuilder := builder.NewRuleBuilder(&knowledgeLibrary)

	// Read rule from GRULE file and build rules
	// ruleFile := pkg.NewFileResource("/home/abhayjo/Desktop/HPE_CTY/rule_engine/engine/rules.grl")

	// Read rules from JSON file
	// rules, _ := os.Open("/home/abhayjo/Desktop/HPE_CTY/middleware/rule_engine_v2/engine/rules.json")
	// ruleResource := pkg.NewReaderResource(rules)
	// ruleFile, _ := pkg.NewJSONResourceFromResource(ruleResource)
	// err := ruleBuilder.BuildRuleFromResource("Rules", "0.0.1", ruleFile)

	// Read rules from MongoDB
	rules, err := mongo.GetRules()
	if err != nil {
		panic(err)
	}

	ruleset, _ := pkg.ParseJSONRuleset(rules)
	// fmt.Println(ruleset)
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

	// create rule engine and execute on provided data and knowledge base
	ruleEngine := engine.NewGruleEngine()
	err = ruleEngine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		return err
	}
	return nil
}
