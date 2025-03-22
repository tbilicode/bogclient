package translate

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tbilicode/bogclient/pkg/bogapi"
)

func Test_IsGeorgian(t *testing.T) {
	t.Parallel()
	assert.True(t, IsGeorgian("საქართველო"))
	assert.True(t, IsGeorgian("სს \"საქართველოს ბანკი\""))
	assert.True(t, IsGeorgian("სს 12345"))
	assert.True(t, IsGeorgian("სს english"))
	assert.False(t, IsGeorgian("english "))
	assert.False(t, IsGeorgian(""))
}

func Test_Extract(t *testing.T) {
	t.Parallel()

	tr := NewTranslator()
	assert.NotNil(t, tr)
	defer func() {
		assert.NoError(t, tr.Close())
	}()

	data, err := os.ReadFile("../bogapi/testdata/statement_feb.json")
	require.NoError(t, err)

	var doc bogapi.AccountStatements
	err = json.Unmarshal(data, &doc)
	require.NoError(t, err)

	texts, err := tr.Extract(&doc)
	require.NoError(t, err)
	assert.Equal(t, 6, len(texts))
}

func Test_TranslateBulk_Cached(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tr := NewTranslator()
	assert.NotNil(t, tr)
	defer func() {
		assert.NoError(t, tr.Close())
	}()

	data, err := os.ReadFile("../bogapi/testdata/statement_feb.json")
	require.NoError(t, err)

	var doc bogapi.AccountStatements
	err = json.Unmarshal(data, &doc)
	require.NoError(t, err)

	// without cache
	texts, err := tr.Extract(&doc)
	require.NoError(t, err)
	assert.Equal(t, 6, len(texts))

	err = tr.LoadDictionary("testdata/dict.json", false)
	require.NoError(t, err)

	// Nothing should be translated, as cached data is loaded
	texts, err = tr.Extract(&doc)
	require.NoError(t, err)
	assert.Equal(t, 0, len(texts))

	replaced, err := tr.Update(ctx, &doc)
	require.NoError(t, err)
	assert.Equal(t, 6, len(replaced))

	data, err = json.MarshalIndent(&doc, "", "  ")
	require.NoError(t, err)
	eng, err := os.ReadFile("testdata/statement_feb_eng.json")
	assert.Equal(t, string(eng), string(data))
}

func Test_TranslateBulk_Real(t *testing.T) {
	// Uncomment the following line to run the real test
	t.Skip("skipping real test")
	t.Parallel()

	ctx := context.Background()
	tr := NewTranslator()
	assert.NotNil(t, tr)
	defer func() {
		assert.NoError(t, tr.Close())
	}()

	data, err := os.ReadFile("../bogapi/testdata/statement_feb.json")
	require.NoError(t, err)

	var doc bogapi.AccountStatements
	err = json.Unmarshal(data, &doc)
	require.NoError(t, err)

	texts, err := tr.Extract(&doc)
	require.NoError(t, err)
	assert.Equal(t, 6, len(texts))

	err = tr.Translate(ctx, "google", texts)
	require.NoError(t, err)

	assert.Equal(t, 6, len(tr.translated))
	tr.SaveDictionary("testdata/dict.json")

	replaced, err := tr.Update(ctx, &doc)
	require.NoError(t, err)
	assert.Equal(t, 6, len(replaced))

	data, err = json.MarshalIndent(&doc, "", "  ")
	err = os.WriteFile("testdata/statement_feb_eng.json", data, 0644)
	require.NoError(t, err)
}
