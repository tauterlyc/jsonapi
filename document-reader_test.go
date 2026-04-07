package jsonapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/tauterlyc/jsonapi"
	jmodels "github.com/tauterlyc/jsonapi/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadResource(t *testing.T) {
	body := `{
  "data": {
    "id": "test",
    "type": "items",
    "attributes": { "description": "description", "enabled": true, "title": "title" },
    "relationships": {
      "initiators": {
        "data": [
          { "id": "e1", "type": "events" },
          { "id": "e2", "type": "events" }
        ]
      },
      "launches": {
        "data": [
          { "id": "a1", "type": "actions" },
          { "id": "a2", "type": "actions" }
        ]
      }
    }
  }
}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	reader, err := Read(req)

	require.NoError(t, err)

	assert.NotEmpty(t, reader)

	resource := jmodels.NewResource("", "")
	err = reader.Resource(resource)

	require.NoError(t, err)

	assert.Equal(t, "test", resource.ID())
	assert.Equal(t, "items", resource.Type())

	xA := jmodels.NewAttributes()
	xA.Set("title", "title")
	xA.Set("description", "description")
	xA.Set("enabled", true)
	assert.Equal(t, xA, resource.Attributes())

	assert.Equal(t, resource.Relationships().Len(), 2)
	assert.Equal(t, []jmodels.ResourceIdentifier{
		{ID: "e1", Type: "events"},
		{ID: "e2", Type: "events"},
	}, resource.Relationships().Get("initiators").Data())
	assert.Equal(t, []jmodels.ResourceIdentifier{
		{ID: "a1", Type: "actions"},
		{ID: "a2", Type: "actions"},
	}, resource.Relationships().Get("launches").Data())

}

func TestReadResourceFail(t *testing.T) {
	body := `{"data": [{ "id": "x", "type": "routines"}]}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	reader, err := Read(req)

	require.NoError(t, err)

	resource := jmodels.NewResource("", "")
	err = reader.Resource(resource)

	assert.Error(t, err)

	assert.Implements(t, (*APIError)(nil), err)

}

func TestReadResources(t *testing.T) {
	body := `{
  "data": [
		{
			"id": "x",
			"type": "routines",
			"attributes": { "description": "17", "enabled": true, "title": "title" },
			"relationships": {
				"initiators": {
					"data": [
						{ "id": "e1", "type": "events" },
						{ "id": "e2", "type": "events" }
					]
				},
				"launches": {
					"data": [
						{ "id": "a1", "type": "actions" },
						{ "id": "a2", "type": "actions" }
					]
				}
			}
		},
		{
			"id": "y",
			"type": "routines",
			"attributes": {"description": "description","enabled": false,"title": "Label"},
			"relationships": {
				"initiators": {
					"data": [
						{ "id": "e3", "type": "events" },
						{ "id": "e4", "type": "events" }
					]
				},
				"launches": {
					"data": [
						{ "id": "a3", "type": "actions" },
						{ "id": "a4", "type": "actions" }
					]
				}
			}
		}
	]
}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	reader, err := Read(req)

	require.NoError(t, err)

	assert.NotEmpty(t, reader)

	resources := make([]jmodels.Resource, 0)
	err = reader.Resources(&resources)

	require.NoError(t, err)

	assert.Len(t, resources, 2)

	assert.Equal(t, "x", resources[0].ID())
	assert.Equal(t, "routines", resources[0].Type())

	xA := jmodels.NewAttributes()
	xA.Set("title", "title")
	xA.Set("description", "17")
	xA.Set("enabled", true)
	assert.Equal(t, xA, resources[0].Attributes())

	assert.Equal(t, resources[0].Relationships().Len(), 2)
	assert.Equal(t, []jmodels.ResourceIdentifier{{ID: "e1", Type: "events"}, {ID: "e2", Type: "events"}}, resources[0].Relationships().Get("initiators").Data())
	assert.Equal(t, []jmodels.ResourceIdentifier{{ID: "a1", Type: "actions"}, {ID: "a2", Type: "actions"}}, resources[0].Relationships().Get("launches").Data())

	assert.Equal(t, "y", resources[1].ID())
	assert.Equal(t, "routines", resources[1].Type())

	xA = jmodels.NewAttributes()
	xA.Set("title", "Label")
	xA.Set("description", "description")
	xA.Set("enabled", false)
	assert.Equal(t, xA, resources[1].Attributes())

	assert.Equal(t, resources[1].Relationships().Len(), 2)
	assert.Equal(t, []jmodels.ResourceIdentifier{{ID: "e3", Type: "events"}, {ID: "e4", Type: "events"}}, resources[1].Relationships().Get("initiators").Data())
	assert.Equal(t, []jmodels.ResourceIdentifier{{ID: "a3", Type: "actions"}, {ID: "a4", Type: "actions"}}, resources[1].Relationships().Get("launches").Data())

}

func TestReadResourcesFail(t *testing.T) {
	body := `{"data": { "id": "x", "type": "routines"}}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	reader, err := Read(req)

	require.NoError(t, err)

	resources := make([]jmodels.Resource, 0)
	err = reader.Resources(&resources)

	assert.Error(t, err)
	assert.Implements(t, (*APIError)(nil), err)
}
