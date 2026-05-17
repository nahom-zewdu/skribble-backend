// internal/game/word_provider.go
// This file defines the WordProvider interface and its implementation for generating word choices for the game.
// The WordProvider is responsible for providing a list of words that players can choose from during their turn.
package game

import (
	"math/rand"
	"time"
)

type WordProvider interface {
	GenerateChoices(n int) []string
}

type StaticWordProvider struct {
	words []string
}

func NewStaticWordProvider() *StaticWordProvider {
	return &StaticWordProvider{
		words: []string{
			"apple", "car", "house", "dog",
			"guitar", "tree", "pizza", "phone",
			"mountain", "river", "sun", "computer",
		},
	}
}

// GenerateChoices returns a random selection of n words from the provider's list.
func (p *StaticWordProvider) GenerateChoices(n int) []string {
	if len(p.words) < n {
		return p.words
	}

	rand.Seed(time.Now().UnixNano())

	result := make([]string, 0, n)
	used := make(map[int]bool)

	for len(result) < n {
		i := rand.Intn(len(p.words))
		if !used[i] {
			used[i] = true
			result = append(result, p.words[i])
		}
	}

	return result
}
