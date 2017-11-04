package encoder_test

import (
	"bytes"
	"io"
	"sort"
	"testing"

	"github.com/konjoot/kaboom/encoder"
)

type rule struct {
	fieldNumber     uint8
	fieldName       string
	fieldType       uint8
	fieldOriginType string
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

func (r *rule) OriginType() string {
	return r.fieldOriginType
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

func TestEncode(t *testing.T) {
	var (
		bts []byte
		err error
	)

	for _, tc := range []struct {
		name     string
		json     io.Reader
		rules    []encoder.Rule
		expBytes []byte
		expErr   error
	}{
		{
			name: "Uint32",
			json: bytes.NewReader([]byte(`{"Uint32": 150}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "Uint32",
					fieldNumber:     1,
					fieldType:       encoder.Varint,
					fieldOriginType: encoder.Uint32,
				},
			},
			expBytes: []byte{
				0x08, 0x96, 0x01,
			},
			expErr: nil,
		},
		{
			name: "Uint64",
			json: bytes.NewReader([]byte(`{"Uint64": 150}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "Uint64",
					fieldNumber:     2,
					fieldType:       encoder.Varint,
					fieldOriginType: encoder.Uint64,
				},
			},
			expBytes: []byte{
				0x10, 0x96, 0x01,
			},
			expErr: nil,
		},
		{
			name: "String",
			json: bytes.NewReader([]byte(`{"String": "String"}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "String",
					fieldNumber:     3,
					fieldType:       encoder.LengthDelimited,
					fieldOriginType: encoder.String,
				},
			},
			expBytes: []byte{
				0x1a, 0x06, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
			},
			expErr: nil,
		},
		{
			name: "Int32Positive",
			json: bytes.NewReader([]byte(`{"Int32": 150}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "Int32",
					fieldNumber:     4,
					fieldType:       encoder.Varint,
					fieldOriginType: encoder.Int32,
				},
			},
			expBytes: []byte{
				0x20, 0x96, 0x01,
			},
			expErr: nil,
		},
		{
			name: "Int32Negative",
			json: bytes.NewReader([]byte(`{"Int32": -150}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "Int32",
					fieldNumber:     4,
					fieldType:       encoder.Varint,
					fieldOriginType: encoder.Int32,
				},
			},
			expBytes: []byte{
				0x20, 0xea, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01,
			},
			expErr: nil,
		},
		{
			name: "Int64Positive",
			json: bytes.NewReader([]byte(`{"Int64": 150}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "Int64",
					fieldNumber:     5,
					fieldType:       encoder.Varint,
					fieldOriginType: encoder.Int64,
				},
			},
			expBytes: []byte{
				0x28, 0x96, 0x01,
			},
			expErr: nil,
		},
		{
			name: "Int64Negative",
			json: bytes.NewReader([]byte(`{"Int64": -150}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "Int64",
					fieldNumber:     5,
					fieldType:       encoder.Varint,
					fieldOriginType: encoder.Int64,
				},
			},
			expBytes: []byte{
				0x28, 0xea, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01,
			},
			expErr: nil,
		},
		{
			name: "Sint32Positive",
			json: bytes.NewReader([]byte(`{"Sint32": 150}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "Sint32",
					fieldNumber:     6,
					fieldType:       encoder.Varint,
					fieldOriginType: encoder.Sint32,
				},
			},
			expBytes: []byte{
				0x30, 0xac, 0x02,
			},
			expErr: nil,
		},
		{
			name: "Sint32Negative",
			json: bytes.NewReader([]byte(`{"Sint32": -150}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "Sint32",
					fieldNumber:     6,
					fieldType:       encoder.Varint,
					fieldOriginType: encoder.Sint32,
				},
			},
			expBytes: []byte{
				0x30, 0xab, 0x02,
			},
			expErr: nil,
		},
		{
			name: "Sint64Positive",
			json: bytes.NewReader([]byte(`{"Sint64": 150}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "Sint64",
					fieldNumber:     7,
					fieldType:       encoder.Varint,
					fieldOriginType: encoder.Sint64,
				},
			},
			expBytes: []byte{
				0x38, 0xac, 0x02,
			},
			expErr: nil,
		},
		{
			name: "Sint64Negative",
			json: bytes.NewReader([]byte(`{"Sint64": -150}`)),
			rules: []encoder.Rule{
				&rule{
					fieldName:       "Sint64",
					fieldNumber:     7,
					fieldType:       encoder.Varint,
					fieldOriginType: encoder.Sint64,
				},
			},
			expBytes: []byte{
				0x38, 0xab, 0x02,
			},
			expErr: nil,
		},
	} {

		t.Run(tc.name, func(t *testing.T) {

			bts, err = encoder.Encode(tc.json, tc.rules)
			t.Log("err =>", err)
			if err != tc.expErr {
				t.Error("Expected =>", tc.expErr)
			}

			t.Logf("bts => % #x", bts)
			if string(bts) != string(tc.expBytes) {
				t.Errorf("Expected => % #x", tc.expBytes)
			}
		})
	}
}
