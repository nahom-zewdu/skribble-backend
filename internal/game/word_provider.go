// internal/game/word_provider.go
// This file defines the WordProvider interface and its implementation for generating word choices for the game.
// The WordProvider is responsible for providing a list of words that players can choose from during their turn.
// internal/game/word_provider.go
package game

import (
	"math/rand"
	"time"
)

type Category string
type Difficulty int

const (
	CategoryFood       Category = "food"
	CategoryAnimal     Category = "animal"
	CategoryObject     Category = "object"
	CategoryAction     Category = "action"
	CategoryFantasy    Category = "fantasy"
	CategoryPlace      Category = "place"
	CategoryProfession Category = "profession"
	CategoryEmotion    Category = "emotion"
	CategorySports     Category = "sports"
	CategoryTechnology Category = "technology"
)

const (
	Easy Difficulty = iota + 1
	Medium
	Hard
)

type Word struct {
	Text       string
	Category   Category
	Difficulty Difficulty
	Enabled    bool
}

type WordProvider interface {
	GenerateChoices(n int) []string
}

type StaticWordProvider struct {
	words []Word
	rng   *rand.Rand
}

func NewStaticWordProvider() *StaticWordProvider {
	return &StaticWordProvider{
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
		words: loadWords(),
	}
}

// GenerateChoices returns randomized words with category diversity.
func (p *StaticWordProvider) GenerateChoices(n int) []string {
	if n <= 0 {
		return []string{}
	}

	if len(p.words) == 0 {
		return []string{}
	}

	if len(p.words) <= n {
		return p.extractTexts(p.words)
	}

	// Group words by category
	categoryBuckets := p.groupByCategory(p.words)

	selected := make([]Word, 0, n)
	usedWords := make(map[string]bool)

	// Extract category list
	categories := make([]Category, 0, len(categoryBuckets))
	for category := range categoryBuckets {
		categories = append(categories, category)
	}

	// Shuffle categories
	p.rng.Shuffle(len(categories), func(i, j int) {
		categories[i], categories[j] = categories[j], categories[i]
	})

	// PASS 1:
	// Try selecting one word per category
	for _, category := range categories {
		if len(selected) >= n {
			break
		}

		words := categoryBuckets[category]

		word := p.pickRandomUnusedWord(words, usedWords)
		if word == nil {
			continue
		}

		selected = append(selected, *word)
		usedWords[word.Text] = true
	}

	// PASS 2:
	// Fill remaining slots from any category
	if len(selected) < n {
		allWords := p.words

		for len(selected) < n {
			word := p.pickRandomUnusedWord(allWords, usedWords)

			if word == nil {
				break
			}

			selected = append(selected, *word)
			usedWords[word.Text] = true
		}
	}

	// Final shuffle
	p.rng.Shuffle(len(selected), func(i, j int) {
		selected[i], selected[j] = selected[j], selected[i]
	})

	return p.extractTexts(selected)
}

func (p *StaticWordProvider) pickRandomUnusedWord(
	words []Word,
	usedWords map[string]bool,
) *Word {
	candidates := make([]Word, 0)

	for _, word := range words {
		if usedWords[word.Text] {
			continue
		}

		candidates = append(candidates, word)
	}

	if len(candidates) == 0 {
		return nil
	}

	chosen := candidates[p.rng.Intn(len(candidates))]
	return &chosen
}

func (p *StaticWordProvider) groupByCategory(words []Word) map[Category][]Word {
	result := make(map[Category][]Word)

	for _, word := range words {
		result[word.Category] = append(result[word.Category], word)
	}

	return result
}

func (p *StaticWordProvider) extractTexts(words []Word) []string {
	result := make([]string, 0, len(words))

	for _, word := range words {
		result = append(result, word.Text)
	}

	return result
}
