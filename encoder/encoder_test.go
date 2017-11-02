package encoder_test

import (
	"sort"
	"testing"

	"github.com/konjoot/kaboom/encoder"
)

type rule struct {
	fieldNumber uint8
	fieldName   string
	fieldType   uint8
}

func (r *rule) Number() uint8 {
	return r.fieldNumber
}

func (r *rule) Name() string {
	return r.fieldName
}

func (r *rule) Type() uint8 {
	return r.fieldType
}

func TestParseRules(t *testing.T) {
	var (
		rules []encoder.Rule
		err   error
	)

	for _, tc := range []struct {
		name       string
		ruleString string
		expRules   []encoder.Rule
		expErr     error
	}{
		{
			name:       "Success",
			ruleString: "one:string;two:int",
			expRules: []encoder.Rule{
				&rule{
					fieldName:   "one",
					fieldNumber: 1,
					fieldType:   encoder.LengthDelimited,
				},
				&rule{
					fieldName:   "two",
					fieldNumber: 2,
					fieldType:   encoder.Varint,
				},
			},
		},
		{
			name:       "EmptyRuleString",
			ruleString: "",
			expRules:   []encoder.Rule{},
		},
		{
			name:       "WrongSepInRuleString",
			ruleString: "one:string&two:int",
			expRules: []encoder.Rule{
				&rule{
					fieldName:   "one",
					fieldNumber: 1,
					fieldType:   encoder.Undefined,
				},
			},
		},
		{
			name:       "RandomString",
			ruleString: "kasjofdjwa[e0j0ifjw[ifjs9a8 !№;%:?*()_ufmw3r",
			expRules: []encoder.Rule{
				&rule{
					fieldName:   "%",
					fieldNumber: 2,
					fieldType:   encoder.Undefined,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {

			rules, err = encoder.ParseRules(tc.ruleString)
			t.Log("err =>", err)
			if err != tc.expErr {
				t.Error("Expected =>", tc.expErr)
			}

			t.Log("len(rules) =>", len(rules))
			if len(rules) != len(tc.expRules) {
				t.Error("Expected =>", len(tc.expRules))
				t.FailNow()
			}

			sort.Sort(encoder.RuleSorter(rules))
			sort.Sort(encoder.RuleSorter(tc.expRules))
			for i, rule := range rules {
				t.Log("rule.Name() =>", rule.Name())
				if rule.Name() != tc.expRules[i].Name() {
					t.Error("Expected =>", tc.expRules[i].Name())
				}
				t.Log("rule.Number() =>", rule.Number())
				if rule.Number() != tc.expRules[i].Number() {
					t.Error("Expected =>", tc.expRules[i].Number())
				}
				t.Log("rule.Type() =>", rule.Type())
				if rule.Type() != tc.expRules[i].Type() {
					t.Error("Expected =>", tc.expRules[i].Type())
				}
			}
		})
	}
}
