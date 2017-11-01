package encoder_test

import (
	"sort"
	"testing"

	"github.com/konjoot/kaboom/encoder"
)


func TestParseRules(t *testing.T) {
	type rule struct {
		Number, Type uint8
		Name string
	}
	var (
		ruleString string
		rules      []*rule
		err        error
	)



	for _, tc := range []struct{
		name string
		ruleString string
		expRules []*Rule
		expErr error
	}{
		{
			name: "success",
			ruleString: "one:string;second:int"
			expRules: []*Rule{
				{
					
				}
			}
		}
	}

	rules, err = encoder.ParseRules(ruleString)
	if err != tc.err {
		t.Log("Actual =>", err)
		t.Error("Expected =>", tc.expErr)
	}

	if len(rules) != len(tc.expRules) {
		t.Log("Actual =>", len(rules))
		t.Error("Expected =>", len(tc.expRules))
	}

	sort.Sort(encoder.RuleSorter(rules))
	sort.Sort(encoder.RuleSorter(tc.expRules))
	for i, rule := range rules {
		if rule.fieldName != tc.expRules[i].fieldName {
			t.Log("Actual =>", rule.fieldName)
			t.Error("Exptected =>", tc.expRules[i].fieldName)
		}
		if rule.fieldNumber != tc.expRules[i].fieldNumber {
			t.Log("Actual =>", rule.fieldNumber)
			t.Error("Exptected =>", tc.expRules[i].fieldNumber)
		}
		if rule.fieldType != tc.expRules[i].fieldType {
			t.Log("Actual =>", rule.fieldType)
			t.Error("Exptected =>", tc.expRules[i].fieldType)
		}
	}
}
