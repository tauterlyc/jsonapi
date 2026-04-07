package models_test

import (
	"encoding/json"
	"testing"

	. "github.com/tauterlyc/jsonapi/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRelationshipsUnmarshal(t *testing.T) {

	body := `{
	"groups": {"data": [
		{ "id": "G1", "type": "groups" },
		{ "id": "G2", "type": "groups" }
	]},
	"users": {"data": [
		{ "id": "P1", "type": "people" },
		{ "id": "P2", "type": "people" }
	]}
}`

	relationships := NewRelationships()
	err := json.Unmarshal([]byte(body), relationships)

	require.NoError(t, err)

	assert.IsType(t, NewRelationship(), relationships.Get("groups"))
	assert.IsType(t, NewRelationship(), relationships.Get("users"))

	assert.Equal(t, []ResourceIdentifier{{ID: "G1", Type: "groups"}, {ID: "G2", Type: "groups"}}, relationships.Get("groups").Data())
	assert.Equal(t, []ResourceIdentifier{{ID: "P1", Type: "people"}, {ID: "P2", Type: "people"}}, relationships.Get("users").Data())

	assert.Nil(t, relationships.Get("does not exist"))

}

func TestRelationshipsAdd(t *testing.T) {

	relationships := NewRelationships()

	groups := NewRelationship()
	groups.SetData([]ResourceIdentifier{
		{ID: "G1", Type: "groups"},
		{ID: "G2", Type: "groups"},
	})

	people := NewRelationship()
	people.SetData([]ResourceIdentifier{
		{ID: "P1", Type: "people"},
		{ID: "P2", Type: "people"},
	})

	relationships.Add("groups", groups)
	relationships.Add("users", people)

	assert.Equal(t, []ResourceIdentifier{{ID: "G1", Type: "groups"}, {ID: "G2", Type: "groups"}}, relationships.Get("groups").Data())
	assert.Equal(t, []ResourceIdentifier{{ID: "P1", Type: "people"}, {ID: "P2", Type: "people"}}, relationships.Get("users").Data())

}

func TestRelationshipsEach(t *testing.T) {

	relationships := NewRelationships()

	groups := NewRelationship()
	groups.SetData([]ResourceIdentifier{
		{ID: "G1", Type: "groups"},
		{ID: "G2", Type: "groups"},
	})

	people := NewRelationship()
	people.SetData([]ResourceIdentifier{
		{ID: "P1", Type: "people"},
		{ID: "P2", Type: "people"},
	})

	relationships.Add("groups", groups)
	relationships.Add("users", people)

	i := 0
	expectedName := []string{"groups", "users"}
	expectedType := []string{"groups", "people"}
	pf := []string{"G", "P"}

	relationships.Each(func(name string, relationship Relationship) {

		relationship.Links().Add("self", name)

		assert.Equal(t, expectedName[i], name)
		assert.Equal(t, []ResourceIdentifier{
			{ID: pf[i] + "1", Type: expectedType[i]},
			{ID: pf[i] + "2", Type: expectedType[i]},
		}, relationship.Data())

		i++

	})

	assert.Equal(t, Links{"self": "groups"}, *relationships.Get("groups").Links())
	assert.Equal(t, Links{"self": "users"}, *relationships.Get("users").Links())

}

// func TestRelationshipsStubAdd(t *testing.T) {

// 	relationships := NewRelationshipsStub()

// 	groups := NewRelationship()
// 	groups.SetData([]ResourceIdentifier{
// 		{ID: "G1", Type: "groups"},
// 		{ID: "G2", Type: "groups"},
// 	})

// 	people := NewRelationship()
// 	people.SetData([]ResourceIdentifier{
// 		{ID: "P1", Type: "people"},
// 		{ID: "P2", Type: "people"},
// 	})

// 	relationships.Add("groups", groups)
// 	relationships.Add("users", people)

// 	assert.Equal(t, NewRelationship(), relationships.Get("groups"))
// 	assert.Equal(t, NewRelationship(), relationships.Get("people"))
// 	assert.Equal(t, NewRelationship(), relationships.Get("does not exist"))

// 	assert.Nil(t, relationships.Get("users").Data())
// 	assert.Nil(t, relationships.Get("people").Data())
// 	assert.Nil(t, relationships.Get("does not exist").Data())

// }
