package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain0(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "*"))
	require.NoError(t, err)

	for _, filename := range files {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			defer func(orig *os.File) {
				os.Stdout = orig
			}(os.Stdout)

			r, w, err := os.Pipe()
			require.NoError(t, err)

			os.Stdout = w

			err = main0([]string{"tfreveal", "--no-color", filepath.Join(filename, "plan.json")})
			require.NoError(t, err)
			err = w.Close()
			require.NoError(t, err)

			out, err := io.ReadAll(r)
			require.NoError(t, err)

			want, err := os.ReadFile(filepath.Join(filename, "output.golden"))
			require.NoError(t, err)

			require.Equal(t, string(want), string(out))
		})
	}
}
