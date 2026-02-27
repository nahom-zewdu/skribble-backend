// internal/game/word_provider.go
// This file defines the WordProvider interface and its implementation for generating word choices for the game.
// The WordProvider is responsible for providing a list of words that players can choose from during their turn.
package game

type WordProvider interface {
	GenerateChoices(n int) []string
}
