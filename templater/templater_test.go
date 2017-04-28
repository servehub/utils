package templater

import (
	"testing"

	"github.com/servehub/utils/gabs"
)

type processorTestCase struct {
	in     string
	expect string
}

func TestUtilsTemplater(t *testing.T) {
	runAllProcessorTests(t, map[string]processorTestCase{
		"simple": {
			in:     `var`,
			expect: `var`,
		},

		"simple resolve with digit": {
			in:     `{{ var1 }}`,
			expect: `var1`,
		},

		"simple resolve with sep": {
			in:     `{{ var-var }}`,
			expect: `var-var`,
		},

		"simple resolve with dot": {
			in:     `{{ var.var }}`,
			expect: `1`,
		},

		"multi resolve": {
			in:     `{{ feature }}-{{ feature-suffix }}`,
			expect: `value-unknown-value-unknown`,
		},

		"replace": {
			in:     `{{ var--v |  replace('\W','_') }}`,
			expect: `var__v`,
		},

		"replace with whitespace": {
			in:     `{{ var--v | replace('\W',  '*') }}`,
			expect: `var**v`,
		},

		"multi resolve and replace": {
			in:     `{{ version | replace('\W',  '*') }}`,
			expect: `value*unknown*value*unknown`,
		},

		"multi resolve and replace with breaks": {
			in:     `{{ version | replace('[a-b]',  '*') }}`,
			expect: `v*lue-unknown-v*lue-unknown`,
		},

		"array value must print first element": {
			in:     `{{ list }}`,
			expect: `1`,
		},

		"lower": {
			in:     `{{ name | lower  }}`,
			expect: `some&name with_simbols`,
		},

		"lower & replace ": {
			in:     `{{ name | lower | replace('\W|_', '-') }}`,
			expect: `some-name-with-simbols`,
		},

		"percent": {
			in:     `{{ memory | percent("50") }}`,
			expect: `32`,
		},
	})
}

func runAllProcessorTests(t *testing.T, cases map[string]processorTestCase) {
	json := `{
		"var1": "var1",
		"var-var": "var-var",
		"var": {"var": "1"},
		"version": "{{ feature }}-{{ feature-suffix }}",
		"feature": "value-unknown",
		"feature-suffix": "{{ feature }}",
		"name": "Some&Name With_simbols",
		"memory": 64,
		"list": [1, 2, 3]
	}`

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			tree, err := gabs.ParseJSON([]byte(json))
			if err != nil {
				t.Fatalf("%v: failed!\n", err)
			}

			if res, err := Template(test.in, tree); err == nil {
				if test.expect != res {
					t.Errorf("%v: %v != %v: failed!\n", name, test.expect, res)
				}
			} else {
				t.Errorf("Error on run `%s`: %v", name, err)
			}
		})
	}
}
