package gojsondiff_test

import (
	"testing"

	"github.com/breml/tfreveal/formatter"
	"github.com/breml/tfreveal/gojsondiff"
	origformatter "github.com/breml/tfreveal/gojsondiff/formatter"
	"github.com/ghetzel/testify/require"
)

func Test(t *testing.T) {
	tt := []struct {
		name    string
		disable bool

		before map[string]interface{}
		after  map[string]interface{}

		want string
	}{
		{
			name:    "array - equal",
			disable: true,

			before: map[string]interface{}{
				"a": []interface{}{
					"foo",
					"bar",
				},
			},
			after: map[string]interface{}{
				"a": []interface{}{
					"foo",
					"bar",
				},
			},

			want: `{}
`,
		},
		{
			name: "array - add item middle",

			before: map[string]interface{}{
				"a": []interface{}{
					"foo",
					"bar",
				},
			},
			after: map[string]interface{}{
				"a": []interface{}{
					"foo",
					"more",
					"more2",
					"more3",
					"bar",
				},
			},

			want: `{
  "a": {
    "1": [
      "more"
    ],
    "_t": "a"
  }
}
`,
		},
		{
			name:    "array - add item start",
			disable: true,

			before: map[string]interface{}{
				"a": []interface{}{
					"foo",
					"bar",
				},
			},
			after: map[string]interface{}{
				"a": []interface{}{
					"more",
					"foo",
					"bar",
				},
			},

			want: `{
  "a": [
    [
      "foo",
      "bar"
    ],
    [
      "more",
      "foo",
      "bar"
    ]
  ]
}
`,
		},
		{
			name:    "array - add item end",
			disable: true,

			before: map[string]interface{}{
				"a": []interface{}{
					"foo",
					"bar",
				},
			},
			after: map[string]interface{}{
				"a": []interface{}{
					"foo",
					"bar",
					"more",
				},
			},

			want: `{
  "a": [
    [
      "foo",
      "bar"
    ],
    [
      "foo",
      "bar",
      "more"
    ]
  ]
}
`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.disable {
				t.Skip()
			}

			diff := gojsondiff.New().CompareObjects(tc.before, tc.after)

			dd := formatter.NewTerraformFormatter(tc.before, formatter.TerraformFormatterConfig{Coloring: true})
			gotOutput, err := dd.Format(diff)
			require.NoError(t, err)
			t.Log(gotOutput)

			d := origformatter.NewDeltaFormatter()
			got, err := d.Format(diff)
			require.NoError(t, err)

			require.Equal(t, tc.want, got)
		})
	}
}
