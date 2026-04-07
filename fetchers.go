package jsonapi

import "github.com/tauterlyc/jsonapi/models"

var Getters map[string]func(ID string) (models.Resource, error) = map[string]func(ID string) (models.Resource, error){}

func AddResourceGetter(rType string, fn func(ID string) (models.Resource, error)) {
	Getters[rType] = fn
}

func GetResource(rType string, ID string) (models.Resource, error) {
	return Getters[rType](ID)
}

var linkDecorator map[string]func(models.Resource) = map[string]func(models.Resource){}

func AddLinkDecorator(rType string, fn func(models.Resource)) {
	linkDecorator[rType] = fn
}

func GetLinkDecorator(rType string) func(models.Resource) {
	return linkDecorator[rType]
}
