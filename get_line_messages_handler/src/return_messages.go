package main

type returnMessages struct {
	Messages []botioMessage `json:"messages" bson:"messages"`
}
