package main

import (
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func main() {
	e := gin.Default()

	gcfg := graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"GetShops": GetShops(),
			"GetUsers": GetUsers(),
		},
	}

	scfg := graphql.SchemaConfig{
		Query: graphql.NewObject(gcfg),
	}

	schema, err := graphql.NewSchema(scfg)
	if err != nil {
		panic(err)
	}
	gqlh := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	e.POST("/graphql", func(gc *gin.Context) {
		gqlh.ServeHTTP(gc.Writer, gc.Request)
	})

	e.Run(":8888")
}

var UsersType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UsersType",
	Fields: graphql.Fields{
		"nodes": &graphql.Field{Type: graphql.NewList(UserType)},
	},
})

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserType",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
	},
})

var ShopType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ShopType",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
	},
})

type Shop struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ShopEdge struct {
	Node   Shop   `json:"node"`
	Cursor string `json:"cursor"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Relay-Style cursor pagination
func GetShops() *graphql.Field {
	return &graphql.Field{
		Type: RelayStylePaginationType(fmt.Sprintf("Pagination%s", ShopType.PrivateName), ShopType, graphql.Fields{}),
		Args: RelayStylePaginationArgs(graphql.FieldConfigArgument{}),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			first, last, after, before, err := GetRelayStylePaginationParams(p)
			if err != nil {
				return nil, err
			}
			_, _, _, _ = first, last, after, before

			var ss []Shop
			for i := 0; i < 100; i++ {
				ss = append(ss, Shop{
					ID:   i + 1,
					Name: fmt.Sprintf("shop%d", i+1),
				})
			}

			var ses []ShopEdge
			for _, s := range ss {
				ses = append(ses, ShopEdge{
					Node:   s,
					Cursor: base64.StdEncoding.EncodeToString([]byte("cursor")),
				})
			}
			return struct {
				Edges []ShopEdge `json:"edges"`
			}{
				Edges: ses,
			}, nil
		},
	}
}

// Offset-based pagination
func GetUsers() *graphql.Field {
	return &graphql.Field{
		Type: UsersType,
		Args: OffsetBasedPaginationArgs(graphql.FieldConfigArgument{}),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			offset, limit, err := GetOffsetBasedPaginationParams(p)
			if err != nil {
				return nil, err
			}

			var us []User
			for i := 0; i < 100; i++ {
				us = append(us, User{
					ID:   i + 1,
					Name: fmt.Sprintf("yukpiz%d", i+1),
				})
			}
			return struct {
				Nodes []User `json:"nodes"`
			}{
				Nodes: us[limit*(offset-1) : limit*offset],
			}, nil
		},
	}
}

func OffsetBasedPaginationArgs(args graphql.FieldConfigArgument) graphql.FieldConfigArgument {
	args["offset"] = &graphql.ArgumentConfig{Type: graphql.Int}
	args["limit"] = &graphql.ArgumentConfig{Type: graphql.Int}
	return args
}

func RelayStylePaginationArgs(args graphql.FieldConfigArgument) graphql.FieldConfigArgument {
	args["first"] = &graphql.ArgumentConfig{Type: graphql.Int}
	args["after"] = &graphql.ArgumentConfig{Type: graphql.String}
	args["before"] = &graphql.ArgumentConfig{Type: graphql.String}
	args["last"] = &graphql.ArgumentConfig{Type: graphql.Int}
	return args
}

func GetOffsetBasedPaginationParams(p graphql.ResolveParams) (offset, limit int, err error) {
	var ok bool

	offset, ok = p.Args["offset"].(int)
	if !ok {
		err = fmt.Errorf("error args: offset")
		return
	}

	limit, ok = p.Args["limit"].(int)
	if !ok {
		err = fmt.Errorf("error args: limit")
		return
	}
	return
}
func GetRelayStylePaginationParams(p graphql.ResolveParams) (first, last int, after, before string, err error) {
	var ok bool

	first, ok = p.Args["first"].(int)
	if !ok {
		err = fmt.Errorf("error args: first")
		return
	}

	after, _ = p.Args["after"].(string)
	last, _ = p.Args["last"].(int)
	before, _ = p.Args["before"].(string)

	return
}

func RelayStylePaginationType(name string, node graphql.Output, fields graphql.Fields) *graphql.Object {
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
					Name: name + "Edges",
					Fields: graphql.Fields{
						"node":   &graphql.Field{Type: node},
						"cursor": &graphql.Field{Type: graphql.String},
					},
				})),
			},
		},
	})

	for k, v := range fields {
		obj.AddFieldConfig(k, v)
	}
	return obj
}
