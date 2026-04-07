package models_test

import (
	"encoding/json"
	"testing"

	. "github.com/tauterlyc/jsonapi/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceUnmarshal(t *testing.T) {
	body := []byte(`{
	"id": "test",
	"type": "items",
	"links": {"self": "http://localhost"},
	"meta": {"build": "vTest"},
	"attributes": {
		"a": "value"
	},
	"relationships": {
		"initiators": {
			"data": [
				{ "id": "e1", "type": "events" },
				{ "id": "e2", "type": "events" }
			]
		}
	}
}`)

	resource := NewResource("", "")
	err := json.Unmarshal(body, resource)

	require.NoError(t, err)

	expectedAttrs := NewAttributes()
	expectedAttrs.Set("a", "value")

	assert.Equal(t, "test", resource.ID())
	assert.Equal(t, "items", resource.Type())
	assert.Equal(t, expectedAttrs, resource.Attributes())
	assert.Equal(t, []ResourceIdentifier{{ID: "e1", Type: "events"}, {ID: "e2", Type: "events"}}, resource.Relationships().Get("initiators").Data())
	assert.Equal(t, Meta{"build": "vTest"}, *resource.Meta())
	assert.Equal(t, Links{"self": "http://localhost"}, *resource.Links())
}

func TestResourceUnmarshalWithNils(t *testing.T) {
	body := []byte(`{"id": "test","type": "items"}`)

	resource := NewResource("", "")
	err := json.Unmarshal(body, resource)

	require.NoError(t, err)

	assert.Equal(t, "test", resource.ID())
	assert.Equal(t, "items", resource.Type())

	assert.Empty(t, resource.Attributes())
	assert.Empty(t, resource.Relationships())
	assert.Empty(t, *resource.Meta())
	assert.Empty(t, *resource.Links())

	assert.Empty(t, resource.Relationships().Get(""))
	assert.Empty(t, resource.Meta().Get(""))
}

func TestResourceMarshal(t *testing.T) {
	body := `{
	"id": "test",
	"type": "items",
	"links": {"self": "http://localhost"},
	"meta": {"build": "vTest"},
	"attributes": {
		"a": "value"
	},
	"relationships": {
		"initiators": {
			"data": [
				{ "id": "e1", "type": "events" },
				{ "id": "e2", "type": "events" }
			]
		}
	}
}`

	resource := NewResource("test", "items")
	resource.Links().Add("self", "http://localhost")
	resource.Meta().Add("build", "vTest")
	resource.Attributes().Set("a", "value")
	resource.Relationships().Add("initiators", NewRelationship())
	resource.Relationships().Get("initiators").SetData([]ResourceIdentifier{{ID: "e1", Type: "events"}, {ID: "e2", Type: "events"}})

	actual, err := json.Marshal(resource)

	require.NoError(t, err)

	assert.JSONEq(t, body, string(actual))
}
