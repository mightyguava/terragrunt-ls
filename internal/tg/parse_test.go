package tg

import (
	"os"
	"path/filepath"
	"terragrunt-ls/internal/testutils"
	"testing"

	"github.com/gruntwork-io/terragrunt/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTerragruntBuffer(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, tmpDir string) string // returns the path to the file to parse
		content  string
		wantCfg  *config.TerragruntConfig
		wantDiag bool
	}{
		{
			name: "basic terragrunt config",
			setup: func(t *testing.T, tmpDir string) string {
				t.Helper()

				path := filepath.Join(tmpDir, "terragrunt.hcl")
				return path
			},
			content: `
terraform {
  source = "./modules/example"
}
`,
			wantDiag: false,
		},
		{
			name: "dependency with missing outputs should not show diagnostic",
			setup: func(t *testing.T, tmpDir string) string {
				t.Helper()

				// Create base module directory and its terragrunt.hcl
				baseDir := filepath.Join(tmpDir, "base")
				require.NoError(t, os.MkdirAll(baseDir, 0755))

				baseTg := filepath.Join(baseDir, "terragrunt.hcl")
				require.NoError(t, os.WriteFile(baseTg, []byte(`
terraform {
  source = "./module"
}
`), 0644))

				// Create unit directory and return path to its terragrunt.hcl
				unitDir := filepath.Join(tmpDir, "unit")
				require.NoError(t, os.MkdirAll(unitDir, 0755))
				return filepath.Join(unitDir, "terragrunt.hcl")
			},
			content: `
dependency "base" {
  config_path = "../base"
}

terraform {
  source = "./modules/example"
}

inputs = {
  vpc_id = dependency.base.outputs.vpc_id
}
`,
			wantDiag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test
			tmpDir := t.TempDir()

			// Run setup and get file to parse
			filename := tt.setup(t, tmpDir)

			l := testutils.NewTestLogger(t)
			cfg, diags := ParseTerragruntBuffer(l, filename, tt.content)

			if tt.wantDiag {
				assert.NotEmpty(t, diags, "expected diagnostics but got none")
			} else {
				assert.Empty(t, diags, "expected no diagnostics but got: %v", diags)
			}

			assert.NotNil(t, cfg, "expected config to not be nil")
		})
	}
}
