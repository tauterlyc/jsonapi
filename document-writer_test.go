package jsonapi_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/tauterlyc/jsonapi"
	"github.com/tauterlyc/jsonapi/models"

	"github.com/stretchr/testify/assert"
)

func setup(query string) (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest("", "/"+query, nil)
}

type k struct{ errors []error }

func (k) Error() string    { return "" }
func (e k) Split() []error { return e.errors }

func TestCreateDocument(t *testing.T) {
	w, r := setup("")

	rA := models.NewResource("0922f90a", "items")
	rA.Attributes().Set("A", "ultra_vibrant_test_99")
	rA.Attributes().Set("B", []string{"alpha_node", "beta_cluster", "", "12345_special_char_!@#"})
	rA.Attributes().Set("C", map[string]interface{}{
		"C1": 9001,
		"C2": 3.1415926535,
		"C3": true,
		"C4": nil,
		"C5": map[string]interface{}{
			"deep_key":  "inception_level_data",
			"timestamp": 1741114870,
		},
	})
	rA.Attributes().Set("D", []interface{}{1, "mixed_type_array", false})
	rA.Attributes().Set("E", 0)

	relA1 := models.NewRelationship()
	relA1.SetData([]models.ResourceIdentifier{{ID: "G1", Type: "groups"}, {ID: "G2", Type: "groups"}})
	rA.Relationships().Add("groups", relA1)

	rB := models.NewResource("1d9c5bbc", "items")
	rB.Attributes().Set("A", "vibrant_99")
	rB.Attributes().Set("B", []string{"alpha", "beta", "12345"})
	rB.Attributes().Set("C", map[string]interface{}{
		"C1": 9001,
		"C2": 3.14,
		"C3": true,
		"C5": map[string]interface{}{"key": "data"},
	})
	rB.Attributes().Set("D", []interface{}{1, "mixed", false})
	rB.Attributes().Set("E", 0)

	relB1 := models.NewRelationship()
	relB1.SetData([]models.ResourceIdentifier{{ID: "G3", Type: "groups"}, {ID: "G2", Type: "groups"}})
	rB.Relationships().Add("groups", relB1)

	New(r).
		SetMeta("build", "vTest-g5907ba4a").
		SetMeta("tag", "vTest").
		LinkTo("self", "http://localhost").
		AddResource(rA).
		Include(rB).
		Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"meta": {
			"build": "vTest-g5907ba4a",
			"tag": "vTest"
		},
		"links": {
			"self":"http://localhost"
		},
		"data": {
			"id": "0922f90a",
			"type": "items",
			"attributes": {
				"A": "ultra_vibrant_test_99",
				"B": [
					"alpha_node",
					"beta_cluster",
					"",
					"12345_special_char_!@#"
				],
				"C": {
					"C1": 9001,
					"C2": 3.1415926535,
					"C3": true,
					"C4": null,
					"C5": {
						"deep_key": "inception_level_data",
						"timestamp": 1741114870
					}
				},
				"D": [
					1,
					"mixed_type_array",
					false
				],
				"E": 0
			},
			"relationships": {
				"groups": {
					"data": [
						{"id": "G1", "type": "groups"},
						{"id": "G2", "type": "groups"}
					]
				}
			}
		},
		"included": [
			{
				"id": "1d9c5bbc",
				"type": "items",
				"attributes": {
					"A": "vibrant_99",
					"B": [
						"alpha",
						"beta",
						"12345"
					],
					"C": {
						"C1": 9001,
						"C2": 3.14,
						"C3": true,
						"C5": {
							"key": "data"
						}
					},
					"D": [
						1,
						"mixed",
						false
					],
					"E": 0
				},
				"relationships": {
					"groups": {
						"data": [
							{"id": "G3", "type": "groups"},
							{"id": "G2", "type": "groups"}
						]
					}
				}
			}
		]
	}`, w.Body.String())
}

func TestSetMeta(t *testing.T) {
	w, r := setup("")

	New(r).
		SetMeta("a", "red").
		SetMeta("b", "green").
		SetMeta("c", "blue").
		Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"meta": {
			"a":"red",
			"b":"green",
			"c":"blue"
		}
	}`, w.Body.String())
}

func TestAddResource(t *testing.T) {
	w, r := setup("")

	rA := models.NewResource("0922f90a", "items")
	rA.Attributes().Set("A", "ultra_vibrant_test_99")
	rA.Attributes().Set("B", []string{"alpha_node", "beta_cluster", "", "12345_special_char_!@#"})
	rA.Attributes().Set("C", map[string]interface{}{
		"C1": 9001,
		"C2": 3.1415926535,
		"C3": true,
		"C4": nil,
		"C5": map[string]interface{}{
			"deep_key":  "inception_level_data",
			"timestamp": 1741114870,
		},
	})
	rA.Attributes().Set("D", []interface{}{1, "mixed_type_array", false})
	rA.Attributes().Set("E", 0)

	New(r).AddResource(rA).Write(w)

	acc := w.Body.String()

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"data": {
			"id": "0922f90a",
			"type": "items",
			"attributes": {
				"A": "ultra_vibrant_test_99",
				"B": ["alpha_node","beta_cluster","","12345_special_char_!@#"],
				"C": {
					"C1": 9001,
					"C2": 3.1415926535,
					"C3": true,
					"C4": null,
					"C5": {"deep_key": "inception_level_data","timestamp": 1741114870}
				},
				"D": [1,"mixed_type_array",false],
				"E": 0
			}
		}
	}`, acc)
}

func TestAddResources(t *testing.T) {
	w, r := setup("")

	rA := models.NewResource("0922f90a", "items")
	rA.Attributes().Set("A", "blue")

	rB := models.NewResource("8ec04634", "items")
	rB.Attributes().Set("B", "green")

	New(r).AddResources([]models.Resource{rA, rB}).Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"data": [{
			"id": "0922f90a",
			"type": "items",
			"attributes": {"A": "blue"}
		},
		{
			"id": "8ec04634",
			"type": "items",
			"attributes": {"B": "green"}
		}]
	}`, w.Body.String())
}

func TestSortResources(t *testing.T) {
	w, r := setup("?sort=a")

	rA := models.NewResource("1", "items")
	rA.Attributes().Set("a", "d")

	rB := models.NewResource("2", "items")
	rB.Attributes().Set("a", "b")

	rC := models.NewResource("3", "items")
	rC.Attributes().Set("a", "c")

	rD := models.NewResource("4", "items")
	rD.Attributes().Set("a", "a")

	New(r).
		Sortable().
		AddResources([]models.Resource{rA, rB, rC, rD}).
		Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"data":[
			{"id": "4", "type":"items", "attributes": {"a":"a"}},
			{"id": "2", "type":"items", "attributes": {"a":"b"}},
			{"id": "3", "type":"items", "attributes": {"a":"c"}},
			{"id": "1", "type":"items", "attributes": {"a":"d"}}
		]
	}`, w.Body.String())
}

func TestSortResourcesReversed(t *testing.T) {
	w, r := setup("?sort=-a")

	rA := models.NewResource("1", "items")
	rA.Attributes().Set("a", "d")

	rB := models.NewResource("2", "items")
	rB.Attributes().Set("a", "b")

	rC := models.NewResource("3", "items")
	rC.Attributes().Set("a", "c")

	rD := models.NewResource("4", "items")
	rD.Attributes().Set("a", "a")

	New(r).
		Sortable().
		AddResources([]models.Resource{rA, rB, rC, rD}).
		Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"data":[
			{"id": "1", "type":"items", "attributes": {"a":"d"}},
			{"id": "3", "type":"items", "attributes": {"a":"c"}},
			{"id": "2", "type":"items", "attributes": {"a":"b"}},
			{"id": "4", "type":"items", "attributes": {"a":"a"}}
		]
	}`, w.Body.String())
}

func TestSortResourcesSecondary(t *testing.T) {
	w, r := setup("?sort=a,b")

	rA := models.NewResource("1", "items")
	rA.Attributes().Set("a", "d")
	rA.Attributes().Set("b", 0)

	rB := models.NewResource("2", "items")
	rB.Attributes().Set("a", "c")
	rB.Attributes().Set("b", 1)

	rC := models.NewResource("3", "items")
	rC.Attributes().Set("a", "c")
	rC.Attributes().Set("b", 0)

	rD := models.NewResource("4", "items")
	rD.Attributes().Set("a", "a")
	rD.Attributes().Set("b", 0)

	New(r).
		Sortable().
		AddResources([]models.Resource{rA, rB, rC, rD}).
		Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"data":[
			{"id": "4", "type":"items", "attributes": {"a":"a","b":0}},
			{"id": "3", "type":"items", "attributes": {"a":"c","b":0}},
			{"id": "2", "type":"items", "attributes": {"a":"c","b":1}},
			{"id": "1", "type":"items", "attributes": {"a":"d","b":0}}
		]
	}`, w.Body.String())
}

func TestSortResourcesReversedSecondary(t *testing.T) {
	w, r := setup("?sort=-a,b")

	rA := models.NewResource("1", "items")
	rA.Attributes().Set("a", "d")
	rA.Attributes().Set("b", 0)

	rB := models.NewResource("2", "items")
	rB.Attributes().Set("a", "c")
	rB.Attributes().Set("b", 1)

	rC := models.NewResource("3", "items")
	rC.Attributes().Set("a", "c")
	rC.Attributes().Set("b", 0)

	rD := models.NewResource("4", "items")
	rD.Attributes().Set("a", "a")
	rD.Attributes().Set("b", 0)

	New(r).
		Sortable().
		AddResources([]models.Resource{rA, rB, rC, rD}).
		Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"data":[
			{"id": "1", "type":"items", "attributes": {"a":"d","b":0}},
			{"id": "3", "type":"items", "attributes": {"a":"c","b":0}},
			{"id": "2", "type":"items", "attributes": {"a":"c","b":1}},
			{"id": "4", "type":"items", "attributes": {"a":"a","b":0}}
		]
	}`, w.Body.String())
}

func TestSortResourcesSecondaryReversed(t *testing.T) {
	w, r := setup("?sort=a,-b")

	rA := models.NewResource("1", "items")
	rA.Attributes().Set("a", "d")
	rA.Attributes().Set("b", 0)

	rB := models.NewResource("2", "items")
	rB.Attributes().Set("a", "c")
	rB.Attributes().Set("b", 1)

	rC := models.NewResource("3", "items")
	rC.Attributes().Set("a", "c")
	rC.Attributes().Set("b", 0)

	rD := models.NewResource("4", "items")
	rD.Attributes().Set("a", "a")
	rD.Attributes().Set("b", 0)

	New(r).
		Sortable().
		AddResources([]models.Resource{rA, rB, rC, rD}).
		Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"data":[
			{"id": "4", "type":"items", "attributes": {"a":"a","b":0}},
			{"id": "2", "type":"items", "attributes": {"a":"c","b":1}},
			{"id": "3", "type":"items", "attributes": {"a":"c","b":0}},
			{"id": "1", "type":"items", "attributes": {"a":"d","b":0}}
		]
	}`, w.Body.String())
}

func TestSortResourcesNotAllowed(t *testing.T) {
	w, r := setup("?sort=a")

	rA := models.NewResource("1", "items")
	rA.Attributes().Set("a", "d")
	rA.Attributes().Set("b", 0)

	rB := models.NewResource("2", "items")
	rB.Attributes().Set("a", "c")
	rB.Attributes().Set("b", 1)

	rC := models.NewResource("3", "items")
	rC.Attributes().Set("a", "c")
	rC.Attributes().Set("b", 0)

	rD := models.NewResource("4", "items")
	rD.Attributes().Set("a", "a")
	rD.Attributes().Set("b", 0)

	New(r).
		AddResources([]models.Resource{rA, rB, rC, rD}).
		Write(w)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"errors":[{
			"title": "query unsupported",
			"status": "400",
			"code": "115",
			"detail": "unsupported query: sort"
		}]
	}`, w.Body.String())
}

func TestAddResourceMeta(t *testing.T) {
	w, r := setup("")

	rA := models.NewResource("d370ef1f", "items")
	rA.Meta().Add("t", 12456765432)

	New(r).AddResource(rA).Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"data": {
			"id": "d370ef1f",
			"type": "items",
			"meta": {"t": 12456765432}
		}
	}`, w.Body.String())
}

func TestAddResourceLink(t *testing.T) {
	w, r := setup("")

	res := models.NewResource("d370ef1f", "items")
	res.Links().Add("self", "http://localhost/items/d370ef1f")

	New(r).AddResource(res).Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"data": {
			"id": "d370ef1f",
			"type": "items",
			"links": {"self": "http://localhost/items/d370ef1f"}
		}
	}`, w.Body.String())
}

func TestInclude(t *testing.T) {
	w, r := setup("")

	r1 := models.NewResource("0922f90a", "items")
	r1.Attributes().Set("A", "blue")
	r2 := models.NewResource("8ec04634", "items")
	r2.Attributes().Set("B", "green")

	New(r).Include(r1).Include(r2).Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"included": [{
			"id": "0922f90a",
			"type": "items",
			"attributes": {"A": "blue"}
		},
		{
			"id": "8ec04634",
			"type": "items",
			"attributes": {"B": "green"}
		}]
	}`, w.Body.String())
}

// ## AddError

func TestAddError(t *testing.T) {
	w, r := setup("")

	New(r).AddError(errors.New("test error")).Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"errors": [{"detail": "test error"}]
	}`, w.Body.String())
}

func TestAddErrorRemovesData(t *testing.T) {
	w, r := setup("")

	New(r).AddResource(models.NewResource("", "")).AddError(errors.New("test error")).Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"errors": [{"detail": "test error"}]
	}`, w.Body.String())
}

func TestAddErrorRemovesDataSlice(t *testing.T) {
	w, r := setup("")

	New(r).
		AddResources(make([]models.Resource, 2)).
		AddError(errors.New("test error")).
		Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"errors": [{"detail": "test error"}]
	}`, w.Body.String())
}

func TestAddErrorRemovesIncluded(t *testing.T) {
	w, r := setup("")

	New(r).
		Include(models.NewResource("", "")).
		AddError(errors.New("test error")).
		Write(w)

	assert.JSONEq(t, `{
		"jsonapi": {"version":"1.1"},
		"errors": [{"detail": "test error"}]
	}`, w.Body.String())
}

func TestAddErrorSetsDefaultStatus(t *testing.T) {
	w, r := setup("")

	New(r).AddError(errors.New("")).Write(w)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestAddErrorSetsStatusToAPIErrorCode(t *testing.T) {
	w, r := setup("")

	New(r).AddError(Error(DataNotCollection)).Write(w)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Result().StatusCode)
}

func TestAddErrorSetsStatusToAPIErrorCodeMultiple(t *testing.T) {
	w, r := setup("")

	New(r).
		AddError(Error(ReferenceNotExist)).
		AddError(Error(ResourceNotExist)).
		Write(w)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
}

func TestAddErrorSetsStatusToAPIErrorCodeMultiError(t *testing.T) {
	w, r := setup("")

	err := k{errors: []error{Error(ReferenceNotExist), Error(ResourceNotExist)}}

	New(r).AddError(err).Write(w)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
}

func TestAddErrorSetsStatusToAPIErrorCodeMany(t *testing.T) {
	w, r := setup("")

	New(r).
		AddError(Error(ReferenceNotExist)).
		AddError(Error(ResourceNotExist)).
		AddError(k{errors: []error{
			Error(DataNotCollection),
			Error(DataNotResource),
		}}).
		AddError(errors.New("")).
		Write(w)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestAddErrorNestedGroupErrors(t *testing.T) {
	w, r := setup("")

	err := k{errors: []error{
		k{errors: []error{
			Error(DataNotCollection),
			k{errors: []error{
				Error(DataNotResource),
			}},
		}},
		Error(DataNotCollection),
	}}

	New(r).AddError(err).Write(w)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Result().StatusCode)
}
