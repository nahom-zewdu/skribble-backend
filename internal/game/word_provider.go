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
	CategoryFood    Category = "food"
	CategoryAnimal  Category = "animal"
	CategoryObject  Category = "object"
	CategoryAction  Category = "action"
	CategoryFantasy Category = "fantasy"
	CategoryPlace   Category = "place"
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
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
		words: []Word{
			{"apple", CategoryFood, Easy, true},
			{"pizza", CategoryFood, Easy, true},
			{"sushi", CategoryFood, Medium, true},

			{"dog", CategoryAnimal, Easy, true},
			{"octopus", CategoryAnimal, Medium, true},
			{"giraffe", CategoryAnimal, Easy, true},

			{"phone", CategoryObject, Easy, true},
			{"guitar", CategoryObject, Medium, true},
			{"microscope", CategoryObject, Hard, true},

			{"dancing", CategoryAction, Medium, true},
			{"juggling", CategoryAction, Hard, true},

			{"dragon", CategoryFantasy, Medium, true},
			{"wizard", CategoryFantasy, Easy, true},

			{"volcano", CategoryPlace, Medium, true},
			{"airport", CategoryPlace, Easy, true},
		},
	}
}

// GenerateChoices returns a balanced and randomized set of words.
func (p *StaticWordProvider) GenerateChoices(n int) []string {
	if n <= 0 {
		return []string{}
	}

	// STEP 1: FILTER INVALID WORDS
	validWords := p.filterValidWords()

	if len(validWords) == 0 {
		return []string{}
	}

	if len(validWords) <= n {
		return p.extractTexts(validWords)
	}

	// STEP 2: GROUP WORDS BY CATEGORY
	categoryBuckets := p.groupByCategory(validWords)

	// STEP 3: BUILD DIFFICULTY TARGETS
	difficultyTargets := p.buildDifficultyTargets(n)

	// STEP 4: SELECT BALANCED WORDS
	selected := p.selectBalancedWords(
		categoryBuckets,
		difficultyTargets,
		n,
	)

	// STEP 5: FINAL RANDOM SHUFFLE
	p.rng.Shuffle(len(selected), func(i, j int) {
		selected[i], selected[j] = selected[j], selected[i]
	})

	return p.extractTexts(selected)
}

func (p *StaticWordProvider) filterValidWords() []Word {
	result := make([]Word, 0)

	seen := make(map[string]bool)

	for _, word := range p.words {
		if !word.Enabled {
			continue
		}

		if word.Text == "" {
			continue
		}

		if seen[word.Text] {
			continue
		}

		seen[word.Text] = true
		result = append(result, word)
	}

	return result
}

func (p *StaticWordProvider) groupByCategory(words []Word) map[Category][]Word {
	result := make(map[Category][]Word)

	for _, word := range words {
		result[word.Category] = append(result[word.Category], word)
	}

	return result
}

func (p *StaticWordProvider) buildDifficultyTargets(n int) []Difficulty {
	targets := make([]Difficulty, 0, n)

	// Balanced distribution:
	// Easy -> Medium -> Hard rotation
	pattern := []Difficulty{
		Easy,
		Medium,
		Hard,
	}

	for i := 0; i < n; i++ {
		targets = append(targets, pattern[i%len(pattern)])
	}

	// Shuffle targets slightly to avoid predictability
	p.rng.Shuffle(len(targets), func(i, j int) {
		targets[i], targets[j] = targets[j], targets[i]
	})

	return targets
}

func (p *StaticWordProvider) selectBalancedWords(
	categoryBuckets map[Category][]Word,
	difficultyTargets []Difficulty,
	n int,
) []Word {
	selected := make([]Word, 0, n)

	usedWords := make(map[string]bool)
	usedCategories := make(map[Category]bool)

	categories := make([]Category, 0, len(categoryBuckets))
	for category := range categoryBuckets {
		categories = append(categories, category)
	}

	// Shuffle categories for randomness
	p.rng.Shuffle(len(categories), func(i, j int) {
		categories[i], categories[j] = categories[j], categories[i]
	})

	for _, targetDifficulty := range difficultyTargets {
		if len(selected) >= n {
			break
		}

		var chosen *Word

		// First pass:
		// prefer unused categories
		for _, category := range categories {
			if usedCategories[category] {
				continue
			}

			word := p.pickWord(
				categoryBuckets[category],
				targetDifficulty,
				usedWords,
			)

			if word != nil {
				chosen = word
				usedCategories[category] = true
				break
			}
		}

		// Second pass:
		// allow reused categories if needed
		if chosen == nil {
			for _, category := range categories {
				word := p.pickWord(
					categoryBuckets[category],
					targetDifficulty,
					usedWords,
				)

				if word != nil {
					chosen = word
					break
				}
			}
		}

		if chosen != nil {
			selected = append(selected, *chosen)
			usedWords[chosen.Text] = true
		}
	}

	return selected
}

func (p *StaticWordProvider) pickWord(
	words []Word,
	targetDifficulty Difficulty,
	usedWords map[string]bool,
) *Word {
	candidates := make([]Word, 0)

	// First priority:
	// exact difficulty match
	for _, word := range words {
		if usedWords[word.Text] {
			continue
		}

		if word.Difficulty == targetDifficulty {
			candidates = append(candidates, word)
		}
	}

	// Fallback:
	// any unused word
	if len(candidates) == 0 {
		for _, word := range words {
			if usedWords[word.Text] {
				continue
			}

			candidates = append(candidates, word)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	chosen := candidates[p.rng.Intn(len(candidates))]
	return &chosen
}

func (p *StaticWordProvider) extractTexts(words []Word) []string {
	result := make([]string, 0, len(words))

	for _, word := range words {
		result = append(result, word.Text)
	}

	return result
}
