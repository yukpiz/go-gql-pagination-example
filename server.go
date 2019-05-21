package main

import (
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

var ShopsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ShopsType",
	Fields: graphql.Fields{
		"nodes": &graphql.Field{Type: graphql.NewList(ShopType)},
	},
})

var ShopType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ShopType",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
	},
})

// Relay-Style cursor pagination
func GetShops() *graphql.Field {
	return &graphql.Field{
		Type: ShopsType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var ss []Shop
			for i := 0; i < 100; i++ {
				ss = append(ss, Shop{
					ID:   i + 1,
					Name: fmt.Sprintf("shop%d", i+1),
				})
			}
			return struct {
				Nodes []Shop `json:"nodes"`
			}{
				Nodes: ss,
			}, nil
		},
	}
}

// Offset-based pagination
func GetUsers() *graphql.Field {
	return &graphql.Field{
		Type: UsersType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
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
				Nodes: us,
			}, nil
		},
	}
}

type Shop struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
