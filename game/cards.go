package game

import uuid "github.com/satori/go.uuid"

// BlackCard is a phrase with <form> placeholders.
// Each placeholder suppose to be filled with corresponding form of WhiteCard
type BlackCard struct {
	id       string
	question string
}

// NewBlackCard creates new BlackCard with unique ID
func NewBlackCard(question string) *BlackCard {
	return &BlackCard{uuid.NewV4().String(), question}
}

// WhiteCard is a set of form:value values filling BlackCards <form> placeholders
type WhiteCard struct {
	id      string
	answers map[string]string
}

// NewWhiteCard create new WhiteCard with unique ID
func NewWhiteCard(answers map[string]string) *WhiteCard {
	return &WhiteCard{uuid.NewV4().String(), answers}
}
