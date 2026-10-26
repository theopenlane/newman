package render

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteTextTemplateBlocksEnvironmentAccess(t *testing.T) {
	t.Setenv("NEWMAN_TEST_SECRET", "super-secret")

	for _, tmpl := range []string{
		`{{ env "NEWMAN_TEST_SECRET" }}`,
		`{{ expandenv "$NEWMAN_TEST_SECRET" }}`,
		`{{ getHostByName "localhost" }}`,
	} {
		t.Run(tmpl, func(t *testing.T) {
			out, err := ExecuteTextTemplate("test", tmpl, nil)
			require.ErrorIs(t, err, ErrTemplateParsingFailed)
			assert.NotContains(t, out, "super-secret")
		})
	}
}

func TestExecuteTextTemplateKeepsDateFuncs(t *testing.T) {
	out, err := ExecuteTextTemplate("test", `{{ now | date "2006" }} {{ .name | upper }}`, map[string]any{"name": "openlane"})
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(time.Now().Year())+" OPENLANE", out)
}

func TestRendererBlocksEnvironmentAccess(t *testing.T) {
	t.Setenv("NEWMAN_TEST_SECRET", "super-secret")

	r := NewRenderer(WithTheme(&Theme{HTML: `{{ env "NEWMAN_TEST_SECRET" }}`, Text: `{{ env "NEWMAN_TEST_SECRET" }}`}))

	_, err := r.GenerateHTML(EmailContent{})
	require.ErrorIs(t, err, ErrTemplateParsingFailed)

	_, err = r.GeneratePlainText(EmailContent{})
	require.ErrorIs(t, err, ErrTemplateParsingFailed)
}
