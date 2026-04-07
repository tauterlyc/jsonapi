package models_test

import (
	"encoding/json"
	"testing"

	"tauterlyc/jsonapi/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRelationshipUnmarshal(t *testing.T) {
	body := `{
	"links": {"self": "http://test"},
	"data": [
		{ "id": "e1", "type": "events" },
		{ "id": "e2", "type": "events" }
	]
}`

	relationship := models.NewRelationship()
	err := json.Unmarshal([]byte(body), relationship)

	require.NoError(t, err)

	assert.Equal(t, models.Links{"self": "http://test"}, *relationship.Links())
	assert.Equal(t, []models.ResourceIdentifier{
		{ID: "e1", Type: "events"},
		{ID: "e2", Type: "events"},
	}, relationship.Data())
}

func TestRelationshipUnmarshalSlice(t *testing.T) {
	body := `[{
	"links": {"self": "http://test"},
	"data": [
		{ "id": "e1", "type": "events" },
		{ "id": "e2", "type": "events" }
	]
}]`

	relationship := []models.Relationship{models.NewRelationship()}
	err := json.Unmarshal([]byte(body), &relationship)

	require.NoError(t, err)

	assert.Equal(t, models.Links{"self": "http://test"}, *relationship[0].Links())
	assert.Equal(t, []models.ResourceIdentifier{
		{ID: "e1", Type: "events"},
		{ID: "e2", Type: "events"},
	}, relationship[0].Data())
}
