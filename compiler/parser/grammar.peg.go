package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	identifier     = regexp.MustCompile("^[A-Za-z]+[A-Za-z0-9]")
	prefixVariable = regexp.MustCompile("{\\w*}")
	defaultPrefix  = &ScopePrefix{String: "", Variables: make([]string, 0)}
)

type statementWrapper struct {
	comment   []string
	statement interface{}
}

type exception *Struct

type union *Struct

func newScopePrefix(prefix string) (*ScopePrefix, error) {
	variables := []string{}
	for _, variable := range prefixVariable.FindAllString(prefix, -1) {
		variable = variable[1 : len(variable)-1]
		if len(variable) == 0 || !identifier.MatchString(variable) {
			return nil, fmt.Errorf("parser: invalid prefix variable '%s'", variable)
		}
		variables = append(variables, variable)
	}
	return &ScopePrefix{String: prefix, Variables: variables}, nil
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func ifaceSliceToString(v interface{}) string {
	ifs := toIfaceSlice(v)
	b := make([]byte, len(ifs))
	for i, v := range ifs {
		b[i] = v.([]uint8)[0]
	}
	return string(b)
}

func rawCommentToDocStr(raw string) []string {
	rawLines := strings.Split(raw, "\n")
	comment := make([]string, len(rawLines))
	for i, line := range rawLines {
		comment[i] = strings.TrimLeft(line, "* ")
	}
	return comment
}

// toStruct converts a union to a struct with all fields optional.
func unionToStruct(u union) *Struct {
	st := (*Struct)(u)
	for _, f := range st.Fields {
		f.Modifier = Optional
	}
	return st
}

// toAnnotations converts an interface{} to an Annotation slice.
func toAnnotations(v interface{}) []*Annotation {
	if v == nil {
		return nil
	}
	return v.([]*Annotation)
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 85, col: 1, offset: 2382},
			expr: &actionExpr{
				pos: position{line: 85, col: 12, offset: 2393},
				run: (*parser).callonGrammar1,
				expr: &seqExpr{
					pos: position{line: 85, col: 12, offset: 2393},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 85, col: 12, offset: 2393},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 85, col: 15, offset: 2396},
							label: "statements",
							expr: &zeroOrMoreExpr{
								pos: position{line: 85, col: 26, offset: 2407},
								expr: &seqExpr{
									pos: position{line: 85, col: 28, offset: 2409},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 85, col: 28, offset: 2409},
											name: "Statement",
										},
										&ruleRefExpr{
											pos:  position{line: 85, col: 38, offset: 2419},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 85, col: 45, offset: 2426},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 85, col: 45, offset: 2426},
									name: "EOF",
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 51, offset: 2432},
									name: "SyntaxError",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SyntaxError",
			pos:  position{line: 152, col: 1, offset: 4824},
			expr: &actionExpr{
				pos: position{line: 152, col: 16, offset: 4839},
				run: (*parser).callonSyntaxError1,
				expr: &anyMatcher{
					line: 152, col: 16, offset: 4839,
				},
			},
		},
		{
			name: "Statement",
			pos:  position{line: 156, col: 1, offset: 4897},
			expr: &actionExpr{
				pos: position{line: 156, col: 14, offset: 4910},
				run: (*parser).callonStatement1,
				expr: &seqExpr{
					pos: position{line: 156, col: 14, offset: 4910},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 156, col: 14, offset: 4910},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 156, col: 21, offset: 4917},
								expr: &seqExpr{
									pos: position{line: 156, col: 22, offset: 4918},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 156, col: 22, offset: 4918},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 156, col: 32, offset: 4928},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 156, col: 37, offset: 4933},
							label: "statement",
							expr: &choiceExpr{
								pos: position{line: 156, col: 48, offset: 4944},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 156, col: 48, offset: 4944},
										name: "ThriftStatement",
									},
									&ruleRefExpr{
										pos:  position{line: 156, col: 66, offset: 4962},
										name: "FrugalStatement",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ThriftStatement",
			pos:  position{line: 169, col: 1, offset: 5433},
			expr: &choiceExpr{
				pos: position{line: 169, col: 20, offset: 5452},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 169, col: 20, offset: 5452},
						name: "Include",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 30, offset: 5462},
						name: "Namespace",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 42, offset: 5474},
						name: "Const",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 50, offset: 5482},
						name: "Enum",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 57, offset: 5489},
						name: "TypeDef",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 67, offset: 5499},
						name: "Struct",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 76, offset: 5508},
						name: "Exception",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 88, offset: 5520},
						name: "Union",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 96, offset: 5528},
						name: "Service",
					},
				},
			},
		},
		{
			name: "Include",
			pos:  position{line: 171, col: 1, offset: 5537},
			expr: &actionExpr{
				pos: position{line: 171, col: 12, offset: 5548},
				run: (*parser).callonInclude1,
				expr: &seqExpr{
					pos: position{line: 171, col: 12, offset: 5548},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 171, col: 12, offset: 5548},
							val:        "include",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 171, col: 22, offset: 5558},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 171, col: 24, offset: 5560},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 171, col: 29, offset: 5565},
								name: "Literal",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 171, col: 37, offset: 5573},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Namespace",
			pos:  position{line: 179, col: 1, offset: 5750},
			expr: &actionExpr{
				pos: position{line: 179, col: 14, offset: 5763},
				run: (*parser).callonNamespace1,
				expr: &seqExpr{
					pos: position{line: 179, col: 14, offset: 5763},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 179, col: 14, offset: 5763},
							val:        "namespace",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 179, col: 26, offset: 5775},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 179, col: 28, offset: 5777},
							label: "scope",
							expr: &oneOrMoreExpr{
								pos: position{line: 179, col: 34, offset: 5783},
								expr: &charClassMatcher{
									pos:        position{line: 179, col: 34, offset: 5783},
									val:        "[*a-z.-]",
									chars:      []rune{'*', '.', '-'},
									ranges:     []rune{'a', 'z'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 179, col: 44, offset: 5793},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 179, col: 46, offset: 5795},
							label: "ns",
							expr: &ruleRefExpr{
								pos:  position{line: 179, col: 49, offset: 5798},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 179, col: 60, offset: 5809},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Const",
			pos:  position{line: 186, col: 1, offset: 5934},
			expr: &actionExpr{
				pos: position{line: 186, col: 10, offset: 5943},
				run: (*parser).callonConst1,
				expr: &seqExpr{
					pos: position{line: 186, col: 10, offset: 5943},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 186, col: 10, offset: 5943},
							val:        "const",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 18, offset: 5951},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 186, col: 20, offset: 5953},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 186, col: 24, offset: 5957},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 34, offset: 5967},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 186, col: 36, offset: 5969},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 186, col: 41, offset: 5974},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 52, offset: 5985},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 186, col: 54, offset: 5987},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 58, offset: 5991},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 186, col: 60, offset: 5993},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 186, col: 66, offset: 5999},
								name: "ConstValue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 186, col: 77, offset: 6010},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Enum",
			pos:  position{line: 194, col: 1, offset: 6142},
			expr: &actionExpr{
				pos: position{line: 194, col: 9, offset: 6150},
				run: (*parser).callonEnum1,
				expr: &seqExpr{
					pos: position{line: 194, col: 9, offset: 6150},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 194, col: 9, offset: 6150},
							val:        "enum",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 16, offset: 6157},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 194, col: 18, offset: 6159},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 194, col: 23, offset: 6164},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 34, offset: 6175},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 194, col: 37, offset: 6178},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 41, offset: 6182},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 194, col: 44, offset: 6185},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 194, col: 51, offset: 6192},
								expr: &seqExpr{
									pos: position{line: 194, col: 52, offset: 6193},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 194, col: 52, offset: 6193},
											name: "EnumValue",
										},
										&ruleRefExpr{
											pos:  position{line: 194, col: 62, offset: 6203},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 194, col: 67, offset: 6208},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 71, offset: 6212},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 194, col: 73, offset: 6214},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 194, col: 85, offset: 6226},
								expr: &ruleRefExpr{
									pos:  position{line: 194, col: 85, offset: 6226},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 102, offset: 6243},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EnumValue",
			pos:  position{line: 218, col: 1, offset: 6905},
			expr: &actionExpr{
				pos: position{line: 218, col: 14, offset: 6918},
				run: (*parser).callonEnumValue1,
				expr: &seqExpr{
					pos: position{line: 218, col: 14, offset: 6918},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 218, col: 14, offset: 6918},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 218, col: 21, offset: 6925},
								expr: &seqExpr{
									pos: position{line: 218, col: 22, offset: 6926},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 218, col: 22, offset: 6926},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 218, col: 32, offset: 6936},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 218, col: 37, offset: 6941},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 218, col: 42, offset: 6946},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 218, col: 53, offset: 6957},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 218, col: 55, offset: 6959},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 218, col: 61, offset: 6965},
								expr: &seqExpr{
									pos: position{line: 218, col: 62, offset: 6966},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 218, col: 62, offset: 6966},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 218, col: 66, offset: 6970},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 218, col: 68, offset: 6972},
											name: "IntConstant",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 218, col: 82, offset: 6986},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 218, col: 84, offset: 6988},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 218, col: 96, offset: 7000},
								expr: &ruleRefExpr{
									pos:  position{line: 218, col: 96, offset: 7000},
									name: "TypeAnnotations",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 218, col: 113, offset: 7017},
							expr: &ruleRefExpr{
								pos:  position{line: 218, col: 113, offset: 7017},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "TypeDef",
			pos:  position{line: 234, col: 1, offset: 7415},
			expr: &actionExpr{
				pos: position{line: 234, col: 12, offset: 7426},
				run: (*parser).callonTypeDef1,
				expr: &seqExpr{
					pos: position{line: 234, col: 12, offset: 7426},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 234, col: 12, offset: 7426},
							val:        "typedef",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 22, offset: 7436},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 234, col: 24, offset: 7438},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 234, col: 28, offset: 7442},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 38, offset: 7452},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 234, col: 40, offset: 7454},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 234, col: 45, offset: 7459},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 56, offset: 7470},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 234, col: 58, offset: 7472},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 234, col: 70, offset: 7484},
								expr: &ruleRefExpr{
									pos:  position{line: 234, col: 70, offset: 7484},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 87, offset: 7501},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Struct",
			pos:  position{line: 242, col: 1, offset: 7673},
			expr: &actionExpr{
				pos: position{line: 242, col: 11, offset: 7683},
				run: (*parser).callonStruct1,
				expr: &seqExpr{
					pos: position{line: 242, col: 11, offset: 7683},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 242, col: 11, offset: 7683},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 242, col: 20, offset: 7692},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 242, col: 22, offset: 7694},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 242, col: 25, offset: 7697},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "Exception",
			pos:  position{line: 243, col: 1, offset: 7737},
			expr: &actionExpr{
				pos: position{line: 243, col: 14, offset: 7750},
				run: (*parser).callonException1,
				expr: &seqExpr{
					pos: position{line: 243, col: 14, offset: 7750},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 243, col: 14, offset: 7750},
							val:        "exception",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 243, col: 26, offset: 7762},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 243, col: 28, offset: 7764},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 243, col: 31, offset: 7767},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "Union",
			pos:  position{line: 244, col: 1, offset: 7818},
			expr: &actionExpr{
				pos: position{line: 244, col: 10, offset: 7827},
				run: (*parser).callonUnion1,
				expr: &seqExpr{
					pos: position{line: 244, col: 10, offset: 7827},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 244, col: 10, offset: 7827},
							val:        "union",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 244, col: 18, offset: 7835},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 244, col: 20, offset: 7837},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 244, col: 23, offset: 7840},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "StructLike",
			pos:  position{line: 245, col: 1, offset: 7887},
			expr: &actionExpr{
				pos: position{line: 245, col: 15, offset: 7901},
				run: (*parser).callonStructLike1,
				expr: &seqExpr{
					pos: position{line: 245, col: 15, offset: 7901},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 245, col: 15, offset: 7901},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 245, col: 20, offset: 7906},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 245, col: 31, offset: 7917},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 245, col: 34, offset: 7920},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 245, col: 38, offset: 7924},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 245, col: 41, offset: 7927},
							label: "fields",
							expr: &ruleRefExpr{
								pos:  position{line: 245, col: 48, offset: 7934},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 245, col: 58, offset: 7944},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 245, col: 62, offset: 7948},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 245, col: 64, offset: 7950},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 245, col: 76, offset: 7962},
								expr: &ruleRefExpr{
									pos:  position{line: 245, col: 76, offset: 7962},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 245, col: 93, offset: 7979},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "FieldList",
			pos:  position{line: 256, col: 1, offset: 8196},
			expr: &actionExpr{
				pos: position{line: 256, col: 14, offset: 8209},
				run: (*parser).callonFieldList1,
				expr: &labeledExpr{
					pos:   position{line: 256, col: 14, offset: 8209},
					label: "fields",
					expr: &zeroOrMoreExpr{
						pos: position{line: 256, col: 21, offset: 8216},
						expr: &seqExpr{
							pos: position{line: 256, col: 22, offset: 8217},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 256, col: 22, offset: 8217},
									name: "Field",
								},
								&ruleRefExpr{
									pos:  position{line: 256, col: 28, offset: 8223},
									name: "__",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 265, col: 1, offset: 8404},
			expr: &actionExpr{
				pos: position{line: 265, col: 10, offset: 8413},
				run: (*parser).callonField1,
				expr: &seqExpr{
					pos: position{line: 265, col: 10, offset: 8413},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 265, col: 10, offset: 8413},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 265, col: 17, offset: 8420},
								expr: &seqExpr{
									pos: position{line: 265, col: 18, offset: 8421},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 265, col: 18, offset: 8421},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 265, col: 28, offset: 8431},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 265, col: 33, offset: 8436},
							label: "id",
							expr: &ruleRefExpr{
								pos:  position{line: 265, col: 36, offset: 8439},
								name: "IntConstant",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 265, col: 48, offset: 8451},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 265, col: 50, offset: 8453},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 265, col: 54, offset: 8457},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 265, col: 56, offset: 8459},
							label: "mod",
							expr: &zeroOrOneExpr{
								pos: position{line: 265, col: 60, offset: 8463},
								expr: &ruleRefExpr{
									pos:  position{line: 265, col: 60, offset: 8463},
									name: "FieldModifier",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 265, col: 75, offset: 8478},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 265, col: 77, offset: 8480},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 265, col: 81, offset: 8484},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 265, col: 91, offset: 8494},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 265, col: 93, offset: 8496},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 265, col: 98, offset: 8501},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 265, col: 109, offset: 8512},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 265, col: 112, offset: 8515},
							label: "def",
							expr: &zeroOrOneExpr{
								pos: position{line: 265, col: 116, offset: 8519},
								expr: &seqExpr{
									pos: position{line: 265, col: 117, offset: 8520},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 265, col: 117, offset: 8520},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 265, col: 121, offset: 8524},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 265, col: 123, offset: 8526},
											name: "ConstValue",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 265, col: 136, offset: 8539},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 265, col: 138, offset: 8541},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 265, col: 150, offset: 8553},
								expr: &ruleRefExpr{
									pos:  position{line: 265, col: 150, offset: 8553},
									name: "TypeAnnotations",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 265, col: 167, offset: 8570},
							expr: &ruleRefExpr{
								pos:  position{line: 265, col: 167, offset: 8570},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FieldModifier",
			pos:  position{line: 288, col: 1, offset: 9102},
			expr: &actionExpr{
				pos: position{line: 288, col: 18, offset: 9119},
				run: (*parser).callonFieldModifier1,
				expr: &choiceExpr{
					pos: position{line: 288, col: 19, offset: 9120},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 288, col: 19, offset: 9120},
							val:        "required",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 288, col: 32, offset: 9133},
							val:        "optional",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Service",
			pos:  position{line: 296, col: 1, offset: 9276},
			expr: &actionExpr{
				pos: position{line: 296, col: 12, offset: 9287},
				run: (*parser).callonService1,
				expr: &seqExpr{
					pos: position{line: 296, col: 12, offset: 9287},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 296, col: 12, offset: 9287},
							val:        "service",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 296, col: 22, offset: 9297},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 296, col: 24, offset: 9299},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 296, col: 29, offset: 9304},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 296, col: 40, offset: 9315},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 296, col: 42, offset: 9317},
							label: "extends",
							expr: &zeroOrOneExpr{
								pos: position{line: 296, col: 50, offset: 9325},
								expr: &seqExpr{
									pos: position{line: 296, col: 51, offset: 9326},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 296, col: 51, offset: 9326},
											val:        "extends",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 296, col: 61, offset: 9336},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 296, col: 64, offset: 9339},
											name: "Identifier",
										},
										&ruleRefExpr{
											pos:  position{line: 296, col: 75, offset: 9350},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 296, col: 80, offset: 9355},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 296, col: 83, offset: 9358},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 296, col: 87, offset: 9362},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 296, col: 90, offset: 9365},
							label: "methods",
							expr: &zeroOrMoreExpr{
								pos: position{line: 296, col: 98, offset: 9373},
								expr: &seqExpr{
									pos: position{line: 296, col: 99, offset: 9374},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 296, col: 99, offset: 9374},
											name: "Function",
										},
										&ruleRefExpr{
											pos:  position{line: 296, col: 108, offset: 9383},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 296, col: 114, offset: 9389},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 296, col: 114, offset: 9389},
									val:        "}",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 296, col: 120, offset: 9395},
									name: "EndOfServiceError",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 296, col: 139, offset: 9414},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 296, col: 141, offset: 9416},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 296, col: 153, offset: 9428},
								expr: &ruleRefExpr{
									pos:  position{line: 296, col: 153, offset: 9428},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 296, col: 170, offset: 9445},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EndOfServiceError",
			pos:  position{line: 313, col: 1, offset: 9886},
			expr: &actionExpr{
				pos: position{line: 313, col: 22, offset: 9907},
				run: (*parser).callonEndOfServiceError1,
				expr: &anyMatcher{
					line: 313, col: 22, offset: 9907,
				},
			},
		},
		{
			name: "Function",
			pos:  position{line: 317, col: 1, offset: 9976},
			expr: &actionExpr{
				pos: position{line: 317, col: 13, offset: 9988},
				run: (*parser).callonFunction1,
				expr: &seqExpr{
					pos: position{line: 317, col: 13, offset: 9988},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 317, col: 13, offset: 9988},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 317, col: 20, offset: 9995},
								expr: &seqExpr{
									pos: position{line: 317, col: 21, offset: 9996},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 317, col: 21, offset: 9996},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 317, col: 31, offset: 10006},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 317, col: 36, offset: 10011},
							label: "oneway",
							expr: &zeroOrOneExpr{
								pos: position{line: 317, col: 43, offset: 10018},
								expr: &seqExpr{
									pos: position{line: 317, col: 44, offset: 10019},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 317, col: 44, offset: 10019},
											val:        "oneway",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 317, col: 53, offset: 10028},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 317, col: 58, offset: 10033},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 317, col: 62, offset: 10037},
								name: "FunctionType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 75, offset: 10050},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 317, col: 78, offset: 10053},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 317, col: 83, offset: 10058},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 94, offset: 10069},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 317, col: 96, offset: 10071},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 100, offset: 10075},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 317, col: 103, offset: 10078},
							label: "arguments",
							expr: &ruleRefExpr{
								pos:  position{line: 317, col: 113, offset: 10088},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 317, col: 123, offset: 10098},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 127, offset: 10102},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 317, col: 130, offset: 10105},
							label: "exceptions",
							expr: &zeroOrOneExpr{
								pos: position{line: 317, col: 141, offset: 10116},
								expr: &ruleRefExpr{
									pos:  position{line: 317, col: 141, offset: 10116},
									name: "Throws",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 149, offset: 10124},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 317, col: 151, offset: 10126},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 317, col: 163, offset: 10138},
								expr: &ruleRefExpr{
									pos:  position{line: 317, col: 163, offset: 10138},
									name: "TypeAnnotations",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 317, col: 180, offset: 10155},
							expr: &ruleRefExpr{
								pos:  position{line: 317, col: 180, offset: 10155},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FunctionType",
			pos:  position{line: 345, col: 1, offset: 10806},
			expr: &actionExpr{
				pos: position{line: 345, col: 17, offset: 10822},
				run: (*parser).callonFunctionType1,
				expr: &labeledExpr{
					pos:   position{line: 345, col: 17, offset: 10822},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 345, col: 22, offset: 10827},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 345, col: 22, offset: 10827},
								val:        "void",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 345, col: 31, offset: 10836},
								name: "FieldType",
							},
						},
					},
				},
			},
		},
		{
			name: "Throws",
			pos:  position{line: 352, col: 1, offset: 10958},
			expr: &actionExpr{
				pos: position{line: 352, col: 11, offset: 10968},
				run: (*parser).callonThrows1,
				expr: &seqExpr{
					pos: position{line: 352, col: 11, offset: 10968},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 352, col: 11, offset: 10968},
							val:        "throws",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 352, col: 20, offset: 10977},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 352, col: 23, offset: 10980},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 352, col: 27, offset: 10984},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 352, col: 30, offset: 10987},
							label: "exceptions",
							expr: &ruleRefExpr{
								pos:  position{line: 352, col: 41, offset: 10998},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 352, col: 51, offset: 11008},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "FieldType",
			pos:  position{line: 356, col: 1, offset: 11044},
			expr: &actionExpr{
				pos: position{line: 356, col: 14, offset: 11057},
				run: (*parser).callonFieldType1,
				expr: &labeledExpr{
					pos:   position{line: 356, col: 14, offset: 11057},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 356, col: 19, offset: 11062},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 356, col: 19, offset: 11062},
								name: "BaseType",
							},
							&ruleRefExpr{
								pos:  position{line: 356, col: 30, offset: 11073},
								name: "ContainerType",
							},
							&ruleRefExpr{
								pos:  position{line: 356, col: 46, offset: 11089},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "BaseType",
			pos:  position{line: 363, col: 1, offset: 11214},
			expr: &actionExpr{
				pos: position{line: 363, col: 13, offset: 11226},
				run: (*parser).callonBaseType1,
				expr: &seqExpr{
					pos: position{line: 363, col: 13, offset: 11226},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 363, col: 13, offset: 11226},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 363, col: 18, offset: 11231},
								name: "BaseTypeName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 363, col: 31, offset: 11244},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 363, col: 33, offset: 11246},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 363, col: 45, offset: 11258},
								expr: &ruleRefExpr{
									pos:  position{line: 363, col: 45, offset: 11258},
									name: "TypeAnnotations",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "BaseTypeName",
			pos:  position{line: 370, col: 1, offset: 11394},
			expr: &actionExpr{
				pos: position{line: 370, col: 17, offset: 11410},
				run: (*parser).callonBaseTypeName1,
				expr: &choiceExpr{
					pos: position{line: 370, col: 18, offset: 11411},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 370, col: 18, offset: 11411},
							val:        "bool",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 370, col: 27, offset: 11420},
							val:        "byte",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 370, col: 36, offset: 11429},
							val:        "i16",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 370, col: 44, offset: 11437},
							val:        "i32",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 370, col: 52, offset: 11445},
							val:        "i64",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 370, col: 60, offset: 11453},
							val:        "double",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 370, col: 71, offset: 11464},
							val:        "string",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 370, col: 82, offset: 11475},
							val:        "binary",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ContainerType",
			pos:  position{line: 374, col: 1, offset: 11522},
			expr: &actionExpr{
				pos: position{line: 374, col: 18, offset: 11539},
				run: (*parser).callonContainerType1,
				expr: &labeledExpr{
					pos:   position{line: 374, col: 18, offset: 11539},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 374, col: 23, offset: 11544},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 374, col: 23, offset: 11544},
								name: "MapType",
							},
							&ruleRefExpr{
								pos:  position{line: 374, col: 33, offset: 11554},
								name: "SetType",
							},
							&ruleRefExpr{
								pos:  position{line: 374, col: 43, offset: 11564},
								name: "ListType",
							},
						},
					},
				},
			},
		},
		{
			name: "MapType",
			pos:  position{line: 378, col: 1, offset: 11599},
			expr: &actionExpr{
				pos: position{line: 378, col: 12, offset: 11610},
				run: (*parser).callonMapType1,
				expr: &seqExpr{
					pos: position{line: 378, col: 12, offset: 11610},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 378, col: 12, offset: 11610},
							expr: &ruleRefExpr{
								pos:  position{line: 378, col: 12, offset: 11610},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 378, col: 21, offset: 11619},
							val:        "map<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 378, col: 28, offset: 11626},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 378, col: 31, offset: 11629},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 378, col: 35, offset: 11633},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 378, col: 45, offset: 11643},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 378, col: 48, offset: 11646},
							val:        ",",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 378, col: 52, offset: 11650},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 378, col: 55, offset: 11653},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 378, col: 61, offset: 11659},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 378, col: 71, offset: 11669},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 378, col: 74, offset: 11672},
							val:        ">",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 378, col: 78, offset: 11676},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 378, col: 80, offset: 11678},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 378, col: 92, offset: 11690},
								expr: &ruleRefExpr{
									pos:  position{line: 378, col: 92, offset: 11690},
									name: "TypeAnnotations",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SetType",
			pos:  position{line: 387, col: 1, offset: 11888},
			expr: &actionExpr{
				pos: position{line: 387, col: 12, offset: 11899},
				run: (*parser).callonSetType1,
				expr: &seqExpr{
					pos: position{line: 387, col: 12, offset: 11899},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 387, col: 12, offset: 11899},
							expr: &ruleRefExpr{
								pos:  position{line: 387, col: 12, offset: 11899},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 387, col: 21, offset: 11908},
							val:        "set<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 387, col: 28, offset: 11915},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 387, col: 31, offset: 11918},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 387, col: 35, offset: 11922},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 387, col: 45, offset: 11932},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 387, col: 48, offset: 11935},
							val:        ">",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 387, col: 52, offset: 11939},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 387, col: 54, offset: 11941},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 387, col: 66, offset: 11953},
								expr: &ruleRefExpr{
									pos:  position{line: 387, col: 66, offset: 11953},
									name: "TypeAnnotations",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ListType",
			pos:  position{line: 395, col: 1, offset: 12115},
			expr: &actionExpr{
				pos: position{line: 395, col: 13, offset: 12127},
				run: (*parser).callonListType1,
				expr: &seqExpr{
					pos: position{line: 395, col: 13, offset: 12127},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 395, col: 13, offset: 12127},
							val:        "list<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 395, col: 21, offset: 12135},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 395, col: 24, offset: 12138},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 395, col: 28, offset: 12142},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 395, col: 38, offset: 12152},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 395, col: 41, offset: 12155},
							val:        ">",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 395, col: 45, offset: 12159},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 395, col: 47, offset: 12161},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 395, col: 59, offset: 12173},
								expr: &ruleRefExpr{
									pos:  position{line: 395, col: 59, offset: 12173},
									name: "TypeAnnotations",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CppType",
			pos:  position{line: 403, col: 1, offset: 12336},
			expr: &actionExpr{
				pos: position{line: 403, col: 12, offset: 12347},
				run: (*parser).callonCppType1,
				expr: &seqExpr{
					pos: position{line: 403, col: 12, offset: 12347},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 403, col: 12, offset: 12347},
							val:        "cpp_type",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 403, col: 23, offset: 12358},
							label: "cppType",
							expr: &ruleRefExpr{
								pos:  position{line: 403, col: 31, offset: 12366},
								name: "Literal",
							},
						},
					},
				},
			},
		},
		{
			name: "ConstValue",
			pos:  position{line: 407, col: 1, offset: 12403},
			expr: &choiceExpr{
				pos: position{line: 407, col: 15, offset: 12417},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 407, col: 15, offset: 12417},
						name: "Literal",
					},
					&ruleRefExpr{
						pos:  position{line: 407, col: 25, offset: 12427},
						name: "BoolConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 407, col: 40, offset: 12442},
						name: "DoubleConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 407, col: 57, offset: 12459},
						name: "IntConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 407, col: 71, offset: 12473},
						name: "ConstMap",
					},
					&ruleRefExpr{
						pos:  position{line: 407, col: 82, offset: 12484},
						name: "ConstList",
					},
					&ruleRefExpr{
						pos:  position{line: 407, col: 94, offset: 12496},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "TypeAnnotations",
			pos:  position{line: 409, col: 1, offset: 12508},
			expr: &actionExpr{
				pos: position{line: 409, col: 20, offset: 12527},
				run: (*parser).callonTypeAnnotations1,
				expr: &seqExpr{
					pos: position{line: 409, col: 20, offset: 12527},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 409, col: 20, offset: 12527},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 409, col: 24, offset: 12531},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 409, col: 27, offset: 12534},
							label: "annotations",
							expr: &zeroOrMoreExpr{
								pos: position{line: 409, col: 39, offset: 12546},
								expr: &ruleRefExpr{
									pos:  position{line: 409, col: 39, offset: 12546},
									name: "TypeAnnotation",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 409, col: 55, offset: 12562},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "TypeAnnotation",
			pos:  position{line: 417, col: 1, offset: 12726},
			expr: &actionExpr{
				pos: position{line: 417, col: 19, offset: 12744},
				run: (*parser).callonTypeAnnotation1,
				expr: &seqExpr{
					pos: position{line: 417, col: 19, offset: 12744},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 417, col: 19, offset: 12744},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 417, col: 24, offset: 12749},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 417, col: 35, offset: 12760},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 417, col: 37, offset: 12762},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 417, col: 43, offset: 12768},
								expr: &actionExpr{
									pos: position{line: 417, col: 44, offset: 12769},
									run: (*parser).callonTypeAnnotation8,
									expr: &seqExpr{
										pos: position{line: 417, col: 44, offset: 12769},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 417, col: 44, offset: 12769},
												val:        "=",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 417, col: 48, offset: 12773},
												name: "__",
											},
											&labeledExpr{
												pos:   position{line: 417, col: 51, offset: 12776},
												label: "value",
												expr: &ruleRefExpr{
													pos:  position{line: 417, col: 57, offset: 12782},
													name: "Literal",
												},
											},
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 417, col: 89, offset: 12814},
							expr: &ruleRefExpr{
								pos:  position{line: 417, col: 89, offset: 12814},
								name: "ListSeparator",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 417, col: 104, offset: 12829},
							name: "__",
						},
					},
				},
			},
		},
		{
			name: "BoolConstant",
			pos:  position{line: 428, col: 1, offset: 13025},
			expr: &actionExpr{
				pos: position{line: 428, col: 17, offset: 13041},
				run: (*parser).callonBoolConstant1,
				expr: &choiceExpr{
					pos: position{line: 428, col: 18, offset: 13042},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 428, col: 18, offset: 13042},
							val:        "true",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 428, col: 27, offset: 13051},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "IntConstant",
			pos:  position{line: 432, col: 1, offset: 13106},
			expr: &actionExpr{
				pos: position{line: 432, col: 16, offset: 13121},
				run: (*parser).callonIntConstant1,
				expr: &seqExpr{
					pos: position{line: 432, col: 16, offset: 13121},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 432, col: 16, offset: 13121},
							expr: &charClassMatcher{
								pos:        position{line: 432, col: 16, offset: 13121},
								val:        "[-+]",
								chars:      []rune{'-', '+'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 432, col: 22, offset: 13127},
							expr: &ruleRefExpr{
								pos:  position{line: 432, col: 22, offset: 13127},
								name: "Digit",
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleConstant",
			pos:  position{line: 436, col: 1, offset: 13191},
			expr: &actionExpr{
				pos: position{line: 436, col: 19, offset: 13209},
				run: (*parser).callonDoubleConstant1,
				expr: &seqExpr{
					pos: position{line: 436, col: 19, offset: 13209},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 436, col: 19, offset: 13209},
							expr: &charClassMatcher{
								pos:        position{line: 436, col: 19, offset: 13209},
								val:        "[+-]",
								chars:      []rune{'+', '-'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 436, col: 25, offset: 13215},
							expr: &ruleRefExpr{
								pos:  position{line: 436, col: 25, offset: 13215},
								name: "Digit",
							},
						},
						&litMatcher{
							pos:        position{line: 436, col: 32, offset: 13222},
							val:        ".",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 436, col: 36, offset: 13226},
							expr: &ruleRefExpr{
								pos:  position{line: 436, col: 36, offset: 13226},
								name: "Digit",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 436, col: 43, offset: 13233},
							expr: &seqExpr{
								pos: position{line: 436, col: 45, offset: 13235},
								exprs: []interface{}{
									&charClassMatcher{
										pos:        position{line: 436, col: 45, offset: 13235},
										val:        "['Ee']",
										chars:      []rune{'\'', 'E', 'e', '\''},
										ignoreCase: false,
										inverted:   false,
									},
									&ruleRefExpr{
										pos:  position{line: 436, col: 52, offset: 13242},
										name: "IntConstant",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ConstList",
			pos:  position{line: 440, col: 1, offset: 13312},
			expr: &actionExpr{
				pos: position{line: 440, col: 14, offset: 13325},
				run: (*parser).callonConstList1,
				expr: &seqExpr{
					pos: position{line: 440, col: 14, offset: 13325},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 440, col: 14, offset: 13325},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 440, col: 18, offset: 13329},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 440, col: 21, offset: 13332},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 440, col: 28, offset: 13339},
								expr: &seqExpr{
									pos: position{line: 440, col: 29, offset: 13340},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 440, col: 29, offset: 13340},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 440, col: 40, offset: 13351},
											name: "__",
										},
										&zeroOrOneExpr{
											pos: position{line: 440, col: 43, offset: 13354},
											expr: &ruleRefExpr{
												pos:  position{line: 440, col: 43, offset: 13354},
												name: "ListSeparator",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 440, col: 58, offset: 13369},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 440, col: 63, offset: 13374},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 440, col: 66, offset: 13377},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ConstMap",
			pos:  position{line: 449, col: 1, offset: 13571},
			expr: &actionExpr{
				pos: position{line: 449, col: 13, offset: 13583},
				run: (*parser).callonConstMap1,
				expr: &seqExpr{
					pos: position{line: 449, col: 13, offset: 13583},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 449, col: 13, offset: 13583},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 449, col: 17, offset: 13587},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 449, col: 20, offset: 13590},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 449, col: 27, offset: 13597},
								expr: &seqExpr{
									pos: position{line: 449, col: 28, offset: 13598},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 449, col: 28, offset: 13598},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 449, col: 39, offset: 13609},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 449, col: 42, offset: 13612},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 449, col: 46, offset: 13616},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 449, col: 49, offset: 13619},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 449, col: 60, offset: 13630},
											name: "__",
										},
										&choiceExpr{
											pos: position{line: 449, col: 64, offset: 13634},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 449, col: 64, offset: 13634},
													val:        ",",
													ignoreCase: false,
												},
												&andExpr{
													pos: position{line: 449, col: 70, offset: 13640},
													expr: &litMatcher{
														pos:        position{line: 449, col: 71, offset: 13641},
														val:        "}",
														ignoreCase: false,
													},
												},
											},
										},
										&ruleRefExpr{
											pos:  position{line: 449, col: 76, offset: 13646},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 449, col: 81, offset: 13651},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "FrugalStatement",
			pos:  position{line: 469, col: 1, offset: 14201},
			expr: &ruleRefExpr{
				pos:  position{line: 469, col: 20, offset: 14220},
				name: "Scope",
			},
		},
		{
			name: "Scope",
			pos:  position{line: 471, col: 1, offset: 14227},
			expr: &actionExpr{
				pos: position{line: 471, col: 10, offset: 14236},
				run: (*parser).callonScope1,
				expr: &seqExpr{
					pos: position{line: 471, col: 10, offset: 14236},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 471, col: 10, offset: 14236},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 471, col: 17, offset: 14243},
								expr: &seqExpr{
									pos: position{line: 471, col: 18, offset: 14244},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 471, col: 18, offset: 14244},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 471, col: 28, offset: 14254},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 471, col: 33, offset: 14259},
							val:        "scope",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 471, col: 41, offset: 14267},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 471, col: 44, offset: 14270},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 471, col: 49, offset: 14275},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 471, col: 60, offset: 14286},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 471, col: 63, offset: 14289},
							label: "prefix",
							expr: &zeroOrOneExpr{
								pos: position{line: 471, col: 70, offset: 14296},
								expr: &ruleRefExpr{
									pos:  position{line: 471, col: 70, offset: 14296},
									name: "Prefix",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 471, col: 78, offset: 14304},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 471, col: 81, offset: 14307},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 471, col: 85, offset: 14311},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 471, col: 88, offset: 14314},
							label: "operations",
							expr: &zeroOrMoreExpr{
								pos: position{line: 471, col: 99, offset: 14325},
								expr: &seqExpr{
									pos: position{line: 471, col: 100, offset: 14326},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 471, col: 100, offset: 14326},
											name: "Operation",
										},
										&ruleRefExpr{
											pos:  position{line: 471, col: 110, offset: 14336},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 471, col: 116, offset: 14342},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 471, col: 116, offset: 14342},
									val:        "}",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 471, col: 122, offset: 14348},
									name: "EndOfScopeError",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 471, col: 139, offset: 14365},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 471, col: 141, offset: 14367},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 471, col: 153, offset: 14379},
								expr: &ruleRefExpr{
									pos:  position{line: 471, col: 153, offset: 14379},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 471, col: 170, offset: 14396},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EndOfScopeError",
			pos:  position{line: 493, col: 1, offset: 14993},
			expr: &actionExpr{
				pos: position{line: 493, col: 20, offset: 15012},
				run: (*parser).callonEndOfScopeError1,
				expr: &anyMatcher{
					line: 493, col: 20, offset: 15012,
				},
			},
		},
		{
			name: "Prefix",
			pos:  position{line: 497, col: 1, offset: 15079},
			expr: &actionExpr{
				pos: position{line: 497, col: 11, offset: 15089},
				run: (*parser).callonPrefix1,
				expr: &seqExpr{
					pos: position{line: 497, col: 11, offset: 15089},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 497, col: 11, offset: 15089},
							val:        "prefix",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 497, col: 20, offset: 15098},
							name: "__",
						},
						&ruleRefExpr{
							pos:  position{line: 497, col: 23, offset: 15101},
							name: "PrefixToken",
						},
						&zeroOrMoreExpr{
							pos: position{line: 497, col: 35, offset: 15113},
							expr: &seqExpr{
								pos: position{line: 497, col: 36, offset: 15114},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 497, col: 36, offset: 15114},
										val:        ".",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 497, col: 40, offset: 15118},
										name: "PrefixToken",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PrefixToken",
			pos:  position{line: 502, col: 1, offset: 15249},
			expr: &choiceExpr{
				pos: position{line: 502, col: 16, offset: 15264},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 502, col: 17, offset: 15265},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 502, col: 17, offset: 15265},
								val:        "{",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 502, col: 21, offset: 15269},
								name: "PrefixWord",
							},
							&litMatcher{
								pos:        position{line: 502, col: 32, offset: 15280},
								val:        "}",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 502, col: 39, offset: 15287},
						name: "PrefixWord",
					},
				},
			},
		},
		{
			name: "PrefixWord",
			pos:  position{line: 504, col: 1, offset: 15299},
			expr: &oneOrMoreExpr{
				pos: position{line: 504, col: 15, offset: 15313},
				expr: &charClassMatcher{
					pos:        position{line: 504, col: 15, offset: 15313},
					val:        "[^\\r\\n\\t\\f .{}]",
					chars:      []rune{'\r', '\n', '\t', '\f', ' ', '.', '{', '}'},
					ignoreCase: false,
					inverted:   true,
				},
			},
		},
		{
			name: "Operation",
			pos:  position{line: 506, col: 1, offset: 15331},
			expr: &actionExpr{
				pos: position{line: 506, col: 14, offset: 15344},
				run: (*parser).callonOperation1,
				expr: &seqExpr{
					pos: position{line: 506, col: 14, offset: 15344},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 506, col: 14, offset: 15344},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 506, col: 21, offset: 15351},
								expr: &seqExpr{
									pos: position{line: 506, col: 22, offset: 15352},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 506, col: 22, offset: 15352},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 506, col: 32, offset: 15362},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 506, col: 37, offset: 15367},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 506, col: 42, offset: 15372},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 506, col: 53, offset: 15383},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 506, col: 55, offset: 15385},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 506, col: 59, offset: 15389},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 506, col: 62, offset: 15392},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 506, col: 66, offset: 15396},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 506, col: 77, offset: 15407},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 506, col: 79, offset: 15409},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 506, col: 91, offset: 15421},
								expr: &ruleRefExpr{
									pos:  position{line: 506, col: 91, offset: 15421},
									name: "TypeAnnotations",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 506, col: 108, offset: 15438},
							expr: &ruleRefExpr{
								pos:  position{line: 506, col: 108, offset: 15438},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "Literal",
			pos:  position{line: 523, col: 1, offset: 16024},
			expr: &actionExpr{
				pos: position{line: 523, col: 12, offset: 16035},
				run: (*parser).callonLiteral1,
				expr: &choiceExpr{
					pos: position{line: 523, col: 13, offset: 16036},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 523, col: 14, offset: 16037},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 523, col: 14, offset: 16037},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 523, col: 18, offset: 16041},
									expr: &choiceExpr{
										pos: position{line: 523, col: 19, offset: 16042},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 523, col: 19, offset: 16042},
												val:        "\\\"",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 523, col: 26, offset: 16049},
												val:        "[^\"]",
												chars:      []rune{'"'},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 523, col: 33, offset: 16056},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 523, col: 41, offset: 16064},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 523, col: 41, offset: 16064},
									val:        "'",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 523, col: 46, offset: 16069},
									expr: &choiceExpr{
										pos: position{line: 523, col: 47, offset: 16070},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 523, col: 47, offset: 16070},
												val:        "\\'",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 523, col: 54, offset: 16077},
												val:        "[^']",
												chars:      []rune{'\''},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 523, col: 61, offset: 16084},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 532, col: 1, offset: 16370},
			expr: &actionExpr{
				pos: position{line: 532, col: 15, offset: 16384},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 532, col: 15, offset: 16384},
					exprs: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 532, col: 15, offset: 16384},
							expr: &choiceExpr{
								pos: position{line: 532, col: 16, offset: 16385},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 532, col: 16, offset: 16385},
										name: "Letter",
									},
									&litMatcher{
										pos:        position{line: 532, col: 25, offset: 16394},
										val:        "_",
										ignoreCase: false,
									},
								},
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 532, col: 31, offset: 16400},
							expr: &choiceExpr{
								pos: position{line: 532, col: 32, offset: 16401},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 532, col: 32, offset: 16401},
										name: "Letter",
									},
									&ruleRefExpr{
										pos:  position{line: 532, col: 41, offset: 16410},
										name: "Digit",
									},
									&charClassMatcher{
										pos:        position{line: 532, col: 49, offset: 16418},
										val:        "[._]",
										chars:      []rune{'.', '_'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ListSeparator",
			pos:  position{line: 536, col: 1, offset: 16473},
			expr: &charClassMatcher{
				pos:        position{line: 536, col: 18, offset: 16490},
				val:        "[,;]",
				chars:      []rune{',', ';'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Letter",
			pos:  position{line: 537, col: 1, offset: 16495},
			expr: &charClassMatcher{
				pos:        position{line: 537, col: 11, offset: 16505},
				val:        "[A-Za-z]",
				ranges:     []rune{'A', 'Z', 'a', 'z'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Digit",
			pos:  position{line: 538, col: 1, offset: 16514},
			expr: &charClassMatcher{
				pos:        position{line: 538, col: 10, offset: 16523},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 540, col: 1, offset: 16530},
			expr: &anyMatcher{
				line: 540, col: 15, offset: 16544,
			},
		},
		{
			name: "DocString",
			pos:  position{line: 541, col: 1, offset: 16546},
			expr: &actionExpr{
				pos: position{line: 541, col: 14, offset: 16559},
				run: (*parser).callonDocString1,
				expr: &seqExpr{
					pos: position{line: 541, col: 14, offset: 16559},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 541, col: 14, offset: 16559},
							val:        "/**@",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 541, col: 21, offset: 16566},
							expr: &seqExpr{
								pos: position{line: 541, col: 23, offset: 16568},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 541, col: 23, offset: 16568},
										expr: &litMatcher{
											pos:        position{line: 541, col: 24, offset: 16569},
											val:        "*/",
											ignoreCase: false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 541, col: 29, offset: 16574},
										name: "SourceChar",
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 541, col: 43, offset: 16588},
							val:        "*/",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 547, col: 1, offset: 16768},
			expr: &choiceExpr{
				pos: position{line: 547, col: 12, offset: 16779},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 547, col: 12, offset: 16779},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 547, col: 31, offset: 16798},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 548, col: 1, offset: 16816},
			expr: &seqExpr{
				pos: position{line: 548, col: 21, offset: 16836},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 548, col: 21, offset: 16836},
						expr: &ruleRefExpr{
							pos:  position{line: 548, col: 22, offset: 16837},
							name: "DocString",
						},
					},
					&litMatcher{
						pos:        position{line: 548, col: 32, offset: 16847},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 548, col: 37, offset: 16852},
						expr: &seqExpr{
							pos: position{line: 548, col: 39, offset: 16854},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 548, col: 39, offset: 16854},
									expr: &litMatcher{
										pos:        position{line: 548, col: 40, offset: 16855},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 548, col: 45, offset: 16860},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 548, col: 59, offset: 16874},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 549, col: 1, offset: 16879},
			expr: &seqExpr{
				pos: position{line: 549, col: 37, offset: 16915},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 549, col: 37, offset: 16915},
						expr: &ruleRefExpr{
							pos:  position{line: 549, col: 38, offset: 16916},
							name: "DocString",
						},
					},
					&litMatcher{
						pos:        position{line: 549, col: 48, offset: 16926},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 549, col: 53, offset: 16931},
						expr: &seqExpr{
							pos: position{line: 549, col: 55, offset: 16933},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 549, col: 55, offset: 16933},
									expr: &choiceExpr{
										pos: position{line: 549, col: 58, offset: 16936},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 549, col: 58, offset: 16936},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 549, col: 65, offset: 16943},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 549, col: 71, offset: 16949},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 549, col: 85, offset: 16963},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 550, col: 1, offset: 16968},
			expr: &choiceExpr{
				pos: position{line: 550, col: 22, offset: 16989},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 550, col: 23, offset: 16990},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 550, col: 23, offset: 16990},
								val:        "//",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 550, col: 28, offset: 16995},
								expr: &seqExpr{
									pos: position{line: 550, col: 30, offset: 16997},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 550, col: 30, offset: 16997},
											expr: &ruleRefExpr{
												pos:  position{line: 550, col: 31, offset: 16998},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 550, col: 35, offset: 17002},
											name: "SourceChar",
										},
									},
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 550, col: 53, offset: 17020},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 550, col: 53, offset: 17020},
								val:        "#",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 550, col: 57, offset: 17024},
								expr: &seqExpr{
									pos: position{line: 550, col: 59, offset: 17026},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 550, col: 59, offset: 17026},
											expr: &ruleRefExpr{
												pos:  position{line: 550, col: 60, offset: 17027},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 550, col: 64, offset: 17031},
											name: "SourceChar",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "__",
			pos:  position{line: 552, col: 1, offset: 17047},
			expr: &zeroOrMoreExpr{
				pos: position{line: 552, col: 7, offset: 17053},
				expr: &choiceExpr{
					pos: position{line: 552, col: 9, offset: 17055},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 552, col: 9, offset: 17055},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 552, col: 22, offset: 17068},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 552, col: 28, offset: 17074},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 553, col: 1, offset: 17085},
			expr: &zeroOrMoreExpr{
				pos: position{line: 553, col: 6, offset: 17090},
				expr: &choiceExpr{
					pos: position{line: 553, col: 8, offset: 17092},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 553, col: 8, offset: 17092},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 553, col: 21, offset: 17105},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "WS",
			pos:  position{line: 554, col: 1, offset: 17141},
			expr: &zeroOrMoreExpr{
				pos: position{line: 554, col: 7, offset: 17147},
				expr: &ruleRefExpr{
					pos:  position{line: 554, col: 7, offset: 17147},
					name: "Whitespace",
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 556, col: 1, offset: 17160},
			expr: &charClassMatcher{
				pos:        position{line: 556, col: 15, offset: 17174},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 557, col: 1, offset: 17182},
			expr: &litMatcher{
				pos:        position{line: 557, col: 8, offset: 17189},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 558, col: 1, offset: 17194},
			expr: &choiceExpr{
				pos: position{line: 558, col: 8, offset: 17201},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 558, col: 8, offset: 17201},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 558, col: 8, offset: 17201},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 558, col: 11, offset: 17204},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 558, col: 17, offset: 17210},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 558, col: 17, offset: 17210},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 558, col: 19, offset: 17212},
								expr: &ruleRefExpr{
									pos:  position{line: 558, col: 19, offset: 17212},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 558, col: 38, offset: 17231},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 558, col: 44, offset: 17237},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 558, col: 44, offset: 17237},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 558, col: 47, offset: 17240},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 560, col: 1, offset: 17245},
			expr: &notExpr{
				pos: position{line: 560, col: 8, offset: 17252},
				expr: &anyMatcher{
					line: 560, col: 9, offset: 17253,
				},
			},
		},
	},
}

func (c *current) onGrammar1(statements interface{}) (interface{}, error) {
	thrift := &Thrift{
		Includes:       []*Include{},
		Namespaces:     []*Namespace{},
		Typedefs:       []*TypeDef{},
		Constants:      []*Constant{},
		Enums:          []*Enum{},
		Structs:        []*Struct{},
		Exceptions:     []*Struct{},
		Unions:         []*Struct{},
		Services:       []*Service{},
		typedefIndex:   make(map[string]*TypeDef),
		namespaceIndex: make(map[string]*Namespace),
	}
	stmts := toIfaceSlice(statements)
	frugal := &Frugal{
		Thrift:         thrift,
		Scopes:         []*Scope{},
		ParsedIncludes: make(map[string]*Frugal),
	}

	for _, st := range stmts {
		wrapper := st.([]interface{})[0].(*statementWrapper)
		switch v := wrapper.statement.(type) {
		case *Namespace:
			thrift.Namespaces = append(thrift.Namespaces, v)
			thrift.namespaceIndex[v.Scope] = v
		case *Constant:
			v.Comment = wrapper.comment
			thrift.Constants = append(thrift.Constants, v)
		case *Enum:
			v.Comment = wrapper.comment
			thrift.Enums = append(thrift.Enums, v)
		case *TypeDef:
			v.Comment = wrapper.comment
			thrift.Typedefs = append(thrift.Typedefs, v)
			thrift.typedefIndex[v.Name] = v
		case *Struct:
			v.Type = StructTypeStruct
			v.Comment = wrapper.comment
			thrift.Structs = append(thrift.Structs, v)
		case exception:
			strct := (*Struct)(v)
			strct.Type = StructTypeException
			strct.Comment = wrapper.comment
			thrift.Exceptions = append(thrift.Exceptions, strct)
		case union:
			strct := unionToStruct(v)
			strct.Type = StructTypeUnion
			strct.Comment = wrapper.comment
			thrift.Unions = append(thrift.Unions, strct)
		case *Service:
			v.Comment = wrapper.comment
			thrift.Services = append(thrift.Services, v)
		case *Include:
			thrift.Includes = append(thrift.Includes, v)
		case *Scope:
			v.Comment = wrapper.comment
			v.Frugal = frugal
			frugal.Scopes = append(frugal.Scopes, v)
		default:
			return nil, fmt.Errorf("parser: unknown value %#v", v)
		}
	}
	return frugal, nil
}

func (p *parser) callonGrammar1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammar1(stack["statements"])
}

func (c *current) onSyntaxError1() (interface{}, error) {
	return nil, errors.New("parser: syntax error")
}

func (p *parser) callonSyntaxError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSyntaxError1()
}

func (c *current) onStatement1(docstr, statement interface{}) (interface{}, error) {
	wrapper := &statementWrapper{statement: statement}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		wrapper.comment = rawCommentToDocStr(raw)
	}
	return wrapper, nil
}

func (p *parser) callonStatement1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStatement1(stack["docstr"], stack["statement"])
}

func (c *current) onInclude1(file interface{}) (interface{}, error) {
	name := file.(string)
	if ix := strings.LastIndex(name, "."); ix > 0 {
		name = name[:ix]
	}
	return &Include{Name: name, Value: file.(string)}, nil
}

func (p *parser) callonInclude1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInclude1(stack["file"])
}

func (c *current) onNamespace1(scope, ns interface{}) (interface{}, error) {
	return &Namespace{
		Scope: ifaceSliceToString(scope),
		Value: string(ns.(Identifier)),
	}, nil
}

func (p *parser) callonNamespace1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNamespace1(stack["scope"], stack["ns"])
}

func (c *current) onConst1(typ, name, value interface{}) (interface{}, error) {
	return &Constant{
		Name:  string(name.(Identifier)),
		Type:  typ.(*Type),
		Value: value,
	}, nil
}

func (p *parser) callonConst1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConst1(stack["typ"], stack["name"], stack["value"])
}

func (c *current) onEnum1(name, values, annotations interface{}) (interface{}, error) {
	vs := toIfaceSlice(values)
	en := &Enum{
		Name:        string(name.(Identifier)),
		Values:      make([]*EnumValue, len(vs)),
		Annotations: toAnnotations(annotations),
	}
	// Assigns numbers in order. This will behave badly if some values are
	// defined and other are not, but I think that's ok since that's a silly
	// thing to do.
	next := 0
	for idx, v := range vs {
		ev := v.([]interface{})[0].(*EnumValue)
		if ev.Value < 0 {
			ev.Value = next
		}
		if ev.Value >= next {
			next = ev.Value + 1
		}
		en.Values[idx] = ev
	}
	return en, nil
}

func (p *parser) callonEnum1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnum1(stack["name"], stack["values"], stack["annotations"])
}

func (c *current) onEnumValue1(docstr, name, value, annotations interface{}) (interface{}, error) {
	ev := &EnumValue{
		Name:        string(name.(Identifier)),
		Value:       -1,
		Annotations: toAnnotations(annotations),
	}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		ev.Comment = rawCommentToDocStr(raw)
	}
	if value != nil {
		ev.Value = int(value.([]interface{})[2].(int64))
	}
	return ev, nil
}

func (p *parser) callonEnumValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumValue1(stack["docstr"], stack["name"], stack["value"], stack["annotations"])
}

func (c *current) onTypeDef1(typ, name, annotations interface{}) (interface{}, error) {
	return &TypeDef{
		Name:        string(name.(Identifier)),
		Type:        typ.(*Type),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonTypeDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeDef1(stack["typ"], stack["name"], stack["annotations"])
}

func (c *current) onStruct1(st interface{}) (interface{}, error) {
	return st.(*Struct), nil
}

func (p *parser) callonStruct1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStruct1(stack["st"])
}

func (c *current) onException1(st interface{}) (interface{}, error) {
	return exception(st.(*Struct)), nil
}

func (p *parser) callonException1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onException1(stack["st"])
}

func (c *current) onUnion1(st interface{}) (interface{}, error) {
	return union(st.(*Struct)), nil
}

func (p *parser) callonUnion1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnion1(stack["st"])
}

func (c *current) onStructLike1(name, fields, annotations interface{}) (interface{}, error) {
	st := &Struct{
		Name:        string(name.(Identifier)),
		Annotations: toAnnotations(annotations),
	}
	if fields != nil {
		st.Fields = fields.([]*Field)
	}
	return st, nil
}

func (p *parser) callonStructLike1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructLike1(stack["name"], stack["fields"], stack["annotations"])
}

func (c *current) onFieldList1(fields interface{}) (interface{}, error) {
	fs := fields.([]interface{})
	flds := make([]*Field, len(fs))
	for i, f := range fs {
		flds[i] = f.([]interface{})[0].(*Field)
	}
	return flds, nil
}

func (p *parser) callonFieldList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldList1(stack["fields"])
}

func (c *current) onField1(docstr, id, mod, typ, name, def, annotations interface{}) (interface{}, error) {
	f := &Field{
		ID:          int(id.(int64)),
		Name:        string(name.(Identifier)),
		Type:        typ.(*Type),
		Annotations: toAnnotations(annotations),
	}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		f.Comment = rawCommentToDocStr(raw)
	}
	if mod != nil {
		f.Modifier = mod.(FieldModifier)
	} else {
		f.Modifier = Default
	}

	if def != nil {
		f.Default = def.([]interface{})[2]
	}
	return f, nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["docstr"], stack["id"], stack["mod"], stack["typ"], stack["name"], stack["def"], stack["annotations"])
}

func (c *current) onFieldModifier1() (interface{}, error) {
	if bytes.Equal(c.text, []byte("required")) {
		return Required, nil
	} else {
		return Optional, nil
	}
}

func (p *parser) callonFieldModifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldModifier1()
}

func (c *current) onService1(name, extends, methods, annotations interface{}) (interface{}, error) {
	ms := methods.([]interface{})
	svc := &Service{
		Name:        string(name.(Identifier)),
		Methods:     make([]*Method, len(ms)),
		Annotations: toAnnotations(annotations),
	}
	if extends != nil {
		svc.Extends = string(extends.([]interface{})[2].(Identifier))
	}
	for i, m := range ms {
		mt := m.([]interface{})[0].(*Method)
		svc.Methods[i] = mt
	}
	return svc, nil
}

func (p *parser) callonService1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onService1(stack["name"], stack["extends"], stack["methods"], stack["annotations"])
}

func (c *current) onEndOfServiceError1() (interface{}, error) {
	return nil, errors.New("parser: expected end of service")
}

func (p *parser) callonEndOfServiceError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEndOfServiceError1()
}

func (c *current) onFunction1(docstr, oneway, typ, name, arguments, exceptions, annotations interface{}) (interface{}, error) {
	m := &Method{
		Name:        string(name.(Identifier)),
		Annotations: toAnnotations(annotations),
	}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		m.Comment = rawCommentToDocStr(raw)
	}
	t := typ.(*Type)
	if t.Name != "void" {
		m.ReturnType = t
	}
	if oneway != nil {
		m.Oneway = true
	}
	if arguments != nil {
		m.Arguments = arguments.([]*Field)
	}
	if exceptions != nil {
		m.Exceptions = exceptions.([]*Field)
		for _, e := range m.Exceptions {
			e.Modifier = Optional
		}
	}
	return m, nil
}

func (p *parser) callonFunction1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunction1(stack["docstr"], stack["oneway"], stack["typ"], stack["name"], stack["arguments"], stack["exceptions"], stack["annotations"])
}

func (c *current) onFunctionType1(typ interface{}) (interface{}, error) {
	if t, ok := typ.(*Type); ok {
		return t, nil
	}
	return &Type{Name: string(c.text)}, nil
}

func (p *parser) callonFunctionType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunctionType1(stack["typ"])
}

func (c *current) onThrows1(exceptions interface{}) (interface{}, error) {
	return exceptions, nil
}

func (p *parser) callonThrows1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onThrows1(stack["exceptions"])
}

func (c *current) onFieldType1(typ interface{}) (interface{}, error) {
	if t, ok := typ.(Identifier); ok {
		return &Type{Name: string(t)}, nil
	}
	return typ, nil
}

func (p *parser) callonFieldType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldType1(stack["typ"])
}

func (c *current) onBaseType1(name, annotations interface{}) (interface{}, error) {
	return &Type{
		Name:        name.(string),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonBaseType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBaseType1(stack["name"], stack["annotations"])
}

func (c *current) onBaseTypeName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonBaseTypeName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBaseTypeName1()
}

func (c *current) onContainerType1(typ interface{}) (interface{}, error) {
	return typ, nil
}

func (p *parser) callonContainerType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContainerType1(stack["typ"])
}

func (c *current) onMapType1(key, value, annotations interface{}) (interface{}, error) {
	return &Type{
		Name:        "map",
		KeyType:     key.(*Type),
		ValueType:   value.(*Type),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonMapType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMapType1(stack["key"], stack["value"], stack["annotations"])
}

func (c *current) onSetType1(typ, annotations interface{}) (interface{}, error) {
	return &Type{
		Name:        "set",
		ValueType:   typ.(*Type),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonSetType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSetType1(stack["typ"], stack["annotations"])
}

func (c *current) onListType1(typ, annotations interface{}) (interface{}, error) {
	return &Type{
		Name:        "list",
		ValueType:   typ.(*Type),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonListType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListType1(stack["typ"], stack["annotations"])
}

func (c *current) onCppType1(cppType interface{}) (interface{}, error) {
	return cppType, nil
}

func (p *parser) callonCppType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCppType1(stack["cppType"])
}

func (c *current) onTypeAnnotations1(annotations interface{}) (interface{}, error) {
	var anns []*Annotation
	for _, ann := range annotations.([]interface{}) {
		anns = append(anns, ann.(*Annotation))
	}
	return anns, nil
}

func (p *parser) callonTypeAnnotations1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeAnnotations1(stack["annotations"])
}

func (c *current) onTypeAnnotation8(value interface{}) (interface{}, error) {
	return value, nil
}

func (p *parser) callonTypeAnnotation8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeAnnotation8(stack["value"])
}

func (c *current) onTypeAnnotation1(name, value interface{}) (interface{}, error) {
	var optValue string
	if value != nil {
		optValue = value.(string)
	}
	return &Annotation{
		Name:  string(name.(Identifier)),
		Value: optValue,
	}, nil
}

func (p *parser) callonTypeAnnotation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeAnnotation1(stack["name"], stack["value"])
}

func (c *current) onBoolConstant1() (interface{}, error) {
	return string(c.text) == "true", nil
}

func (p *parser) callonBoolConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolConstant1()
}

func (c *current) onIntConstant1() (interface{}, error) {
	return strconv.ParseInt(string(c.text), 10, 64)
}

func (p *parser) callonIntConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntConstant1()
}

func (c *current) onDoubleConstant1() (interface{}, error) {
	return strconv.ParseFloat(string(c.text), 64)
}

func (p *parser) callonDoubleConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleConstant1()
}

func (c *current) onConstList1(values interface{}) (interface{}, error) {
	valueSlice := values.([]interface{})
	vs := make([]interface{}, len(valueSlice))
	for i, v := range valueSlice {
		vs[i] = v.([]interface{})[0]
	}
	return vs, nil
}

func (p *parser) callonConstList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstList1(stack["values"])
}

func (c *current) onConstMap1(values interface{}) (interface{}, error) {
	if values == nil {
		return nil, nil
	}
	vals := values.([]interface{})
	kvs := make([]KeyValue, len(vals))
	for i, kv := range vals {
		v := kv.([]interface{})
		kvs[i] = KeyValue{
			Key:   v[0],
			Value: v[4],
		}
	}
	return kvs, nil
}

func (p *parser) callonConstMap1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstMap1(stack["values"])
}

func (c *current) onScope1(docstr, name, prefix, operations, annotations interface{}) (interface{}, error) {
	ops := operations.([]interface{})
	scope := &Scope{
		Name:        string(name.(Identifier)),
		Operations:  make([]*Operation, len(ops)),
		Prefix:      defaultPrefix,
		Annotations: toAnnotations(annotations),
	}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		scope.Comment = rawCommentToDocStr(raw)
	}
	if prefix != nil {
		scope.Prefix = prefix.(*ScopePrefix)
	}
	for i, o := range ops {
		op := o.([]interface{})[0].(*Operation)
		scope.Operations[i] = op
	}
	return scope, nil
}

func (p *parser) callonScope1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onScope1(stack["docstr"], stack["name"], stack["prefix"], stack["operations"], stack["annotations"])
}

func (c *current) onEndOfScopeError1() (interface{}, error) {
	return nil, errors.New("parser: expected end of scope")
}

func (p *parser) callonEndOfScopeError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEndOfScopeError1()
}

func (c *current) onPrefix1() (interface{}, error) {
	prefix := strings.TrimSpace(strings.TrimPrefix(string(c.text), "prefix"))
	return newScopePrefix(prefix)
}

func (p *parser) callonPrefix1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrefix1()
}

func (c *current) onOperation1(docstr, name, typ, annotations interface{}) (interface{}, error) {
	o := &Operation{
		Name:        string(name.(Identifier)),
		Type:        &Type{Name: string(typ.(Identifier))},
		Annotations: toAnnotations(annotations),
	}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		o.Comment = rawCommentToDocStr(raw)
	}
	return o, nil
}

func (p *parser) callonOperation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOperation1(stack["docstr"], stack["name"], stack["typ"], stack["annotations"])
}

func (c *current) onLiteral1() (interface{}, error) {
	if len(c.text) != 0 && c.text[0] == '\'' {
		intermediate := strings.Replace(string(c.text[1:len(c.text)-1]), `\'`, `'`, -1)
		return strconv.Unquote(`"` + strings.Replace(intermediate, `"`, `\"`, -1) + `"`)
	}

	return strconv.Unquote(string(c.text))
}

func (p *parser) callonLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteral1()
}

func (c *current) onIdentifier1() (interface{}, error) {
	return Identifier(string(c.text)), nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1()
}

func (c *current) onDocString1() (interface{}, error) {
	comment := string(c.text)
	comment = strings.TrimPrefix(comment, "/**@")
	comment = strings.TrimSuffix(comment, "*/")
	return strings.TrimSpace(comment), nil
}

func (p *parser) callonDocString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDocString1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
