// Package graphql implements the GraphQL driving adapter for video operations.
package graphql

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"
	"github.com/yinebebt/hexagonal-architecture/internal/core/port"
)

// personType defines the GraphQL type for Person
var personType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Person",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"firstname": &graphql.Field{Type: graphql.String},
		"lastname":  &graphql.Field{Type: graphql.String},
		"age":       &graphql.Field{Type: graphql.Int},
		"email":     &graphql.Field{Type: graphql.String},
	},
})

// videoType defines the GraphQL type for Video
var videoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Video",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.String},
		"title":       &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"url":         &graphql.Field{Type: graphql.String},
		"author": &graphql.Field{
			Type: personType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if video, ok := p.Source.(entity.Video); ok {
					return video.Director, nil
				}
				return nil, nil
			},
		},
	},
})

// personInput defines the GraphQL input type for Person in mutations
var personInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "PersonInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"firstname": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"lastname":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"age":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
		"email":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})

// NewHandler creates an http.Handler that serves the GraphQL API.
func NewHandler(videoService port.VideoService) (http.Handler, error) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"videos": &graphql.Field{
				Type:        graphql.NewList(videoType),
				Description: "List all videos",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return videoService.FindAll(p.Context)
				},
			},
			"video": &graphql.Field{
				Type:        videoType,
				Description: "Get a video by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := int64(p.Args["id"].(int))
					return videoService.FindByID(p.Context, id)
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createVideo": &graphql.Field{
				Type:        videoType,
				Description: "Create a new video",
				Args: graphql.FieldConfigArgument{
					"title":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"description": &graphql.ArgumentConfig{Type: graphql.String},
					"url":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"author":      &graphql.ArgumentConfig{Type: personInput},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					video := entity.Video{
						Title: p.Args["title"].(string),
						URL:   p.Args["url"].(string),
					}
					if desc, ok := p.Args["description"].(string); ok {
						video.Description = desc
					}
					if author, ok := p.Args["author"].(map[string]interface{}); ok {
						person, err := mapToPerson(author)
						if err != nil {
							return nil, err
						}
						video.Director = person
					}
					return videoService.Save(p.Context, video)
				},
			},
			"updateVideo": &graphql.Field{
				Type:        videoType,
				Description: "Update an existing video",
				Args: graphql.FieldConfigArgument{
					"id":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"title":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"description": &graphql.ArgumentConfig{Type: graphql.String},
					"url":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"author":      &graphql.ArgumentConfig{Type: personInput},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					video := entity.Video{
						ID:    int64(p.Args["id"].(int)),
						Title: p.Args["title"].(string),
						URL:   p.Args["url"].(string),
					}
					if desc, ok := p.Args["description"].(string); ok {
						video.Description = desc
					}
					if author, ok := p.Args["author"].(map[string]interface{}); ok {
						person, err := mapToPerson(author)
						if err != nil {
							return nil, err
						}
						video.Director = person
					}
					return videoService.Update(p.Context, video)
				},
			},
			"deleteVideo": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete a video by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := int64(p.Args["id"].(int))
					if err := videoService.Delete(p.Context, id); err != nil {
						return false, fmt.Errorf("failed to delete video")
					}
					return true, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		return nil, err
	}

	return handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	}), nil
}

// mapToPerson converts a GraphQL input map to a Person entity.
func mapToPerson(m map[string]interface{}) (entity.Person, error) {
	p := entity.Person{}
	if v, ok := m["firstname"].(string); ok {
		p.FirstName = v
	}
	if v, ok := m["lastname"].(string); ok {
		p.LastName = v
	}
	if v, ok := m["age"].(int); ok {
		if v < 10 || v > 127 {
			return entity.Person{}, fmt.Errorf("age must be between 10 and 127")
		}
		p.Age = int8(v)
	}
	if v, ok := m["email"].(string); ok {
		p.Email = v
	}
	return p, nil
}
