// internal/game/word_pool.go
package game

func loadWords() []Word {
	return []Word{
		// =========================
		// FOOD
		// =========================
		{"apple", CategoryFood, Easy, true},
		{"pizza", CategoryFood, Easy, true},
		{"burger", CategoryFood, Easy, true},
		{"sushi", CategoryFood, Medium, true},
		{"spaghetti", CategoryFood, Medium, true},
		{"donut", CategoryFood, Easy, true},
		{"taco", CategoryFood, Easy, true},
		{"avocado", CategoryFood, Medium, true},
		{"ice cream", CategoryFood, Easy, true},
		{"pancakes", CategoryFood, Easy, true},
		{"watermelon", CategoryFood, Easy, true},
		{"hotdog", CategoryFood, Easy, true},
		{"croissant", CategoryFood, Medium, true},
		{"ramen", CategoryFood, Medium, true},
		{"popcorn", CategoryFood, Easy, true},
		{"cupcake", CategoryFood, Easy, true},
		{"milkshake", CategoryFood, Medium, true},
		{"cheesecake", CategoryFood, Medium, true},
		{"fried chicken", CategoryFood, Medium, true},
		{"sandwich", CategoryFood, Easy, true},

		// =========================
		// ANIMALS
		// =========================
		{"dog", CategoryAnimal, Easy, true},
		{"cat", CategoryAnimal, Easy, true},
		{"giraffe", CategoryAnimal, Easy, true},
		{"octopus", CategoryAnimal, Medium, true},
		{"penguin", CategoryAnimal, Easy, true},
		{"kangaroo", CategoryAnimal, Medium, true},
		{"elephant", CategoryAnimal, Easy, true},
		{"shark", CategoryAnimal, Easy, true},
		{"hedgehog", CategoryAnimal, Medium, true},
		{"dolphin", CategoryAnimal, Easy, true},
		{"flamingo", CategoryAnimal, Medium, true},
		{"crocodile", CategoryAnimal, Medium, true},
		{"gorilla", CategoryAnimal, Medium, true},
		{"camel", CategoryAnimal, Easy, true},
		{"owl", CategoryAnimal, Easy, true},
		{"peacock", CategoryAnimal, Medium, true},
		{"tiger", CategoryAnimal, Easy, true},
		{"wolf", CategoryAnimal, Medium, true},
		{"rabbit", CategoryAnimal, Easy, true},
		{"snail", CategoryAnimal, Easy, true},

		// =========================
		// OBJECTS
		// =========================
		{"phone", CategoryObject, Easy, true},
		{"guitar", CategoryObject, Medium, true},
		{"microscope", CategoryObject, Hard, true},
		{"umbrella", CategoryObject, Easy, true},
		{"backpack", CategoryObject, Easy, true},
		{"toaster", CategoryObject, Easy, true},
		{"headphones", CategoryObject, Medium, true},
		{"camera", CategoryObject, Easy, true},
		{"skateboard", CategoryObject, Medium, true},
		{"chainsaw", CategoryObject, Medium, true},
		{"flashlight", CategoryObject, Easy, true},
		{"television", CategoryObject, Easy, true},
		{"keyboard", CategoryObject, Easy, true},
		{"rocket", CategoryObject, Medium, true},
		{"drone", CategoryObject, Medium, true},
		{"binoculars", CategoryObject, Medium, true},
		{"traffic light", CategoryObject, Medium, true},
		{"washing machine", CategoryObject, Medium, true},
		{"helmet", CategoryObject, Easy, true},
		{"vacuum cleaner", CategoryObject, Medium, true},

		// =========================
		// ACTIONS
		// =========================
		{"dancing", CategoryAction, Medium, true},
		{"juggling", CategoryAction, Hard, true},
		{"sleeping", CategoryAction, Easy, true},
		{"swimming", CategoryAction, Easy, true},
		{"arguing", CategoryAction, Hard, true},
		{"climbing", CategoryAction, Medium, true},
		{"running", CategoryAction, Easy, true},
		{"crying", CategoryAction, Easy, true},
		{"laughing", CategoryAction, Easy, true},
		{"escaping", CategoryAction, Hard, true},
		{"cooking", CategoryAction, Easy, true},
		{"painting", CategoryAction, Medium, true},
		{"driving", CategoryAction, Easy, true},
		{"typing", CategoryAction, Easy, true},
		{"fishing", CategoryAction, Medium, true},
		{"surfing", CategoryAction, Medium, true},
		{"digging", CategoryAction, Easy, true},
		{"lifting weights", CategoryAction, Medium, true},
		{"sneezing", CategoryAction, Easy, true},
		{"whistling", CategoryAction, Medium, true},

		// =========================
		// FANTASY
		// =========================
		{"dragon", CategoryFantasy, Medium, true},
		{"wizard", CategoryFantasy, Easy, true},
		{"zombie", CategoryFantasy, Easy, true},
		{"vampire", CategoryFantasy, Medium, true},
		{"ghost", CategoryFantasy, Easy, true},
		{"alien invasion", CategoryFantasy, Hard, true},
		{"time machine", CategoryFantasy, Hard, true},
		{"flying castle", CategoryFantasy, Hard, true},
		{"mermaid", CategoryFantasy, Medium, true},
		{"phoenix", CategoryFantasy, Medium, true},
		{"unicorn", CategoryFantasy, Easy, true},
		{"robot ninja", CategoryFantasy, Hard, true},
		{"pirate wizard", CategoryFantasy, Hard, true},
		{"giant spider", CategoryFantasy, Medium, true},
		{"magic portal", CategoryFantasy, Medium, true},

		// =========================
		// PLACES
		// =========================
		{"airport", CategoryPlace, Easy, true},
		{"volcano", CategoryPlace, Medium, true},
		{"desert", CategoryPlace, Easy, true},
		{"school", CategoryPlace, Easy, true},
		{"supermarket", CategoryPlace, Medium, true},
		{"jungle", CategoryPlace, Medium, true},
		{"beach", CategoryPlace, Easy, true},
		{"space station", CategoryPlace, Hard, true},
		{"haunted house", CategoryPlace, Medium, true},
		{"underwater city", CategoryPlace, Hard, true},
		{"mountain", CategoryPlace, Easy, true},
		{"stadium", CategoryPlace, Medium, true},
		{"museum", CategoryPlace, Medium, true},
		{"train station", CategoryPlace, Medium, true},
		{"zoo", CategoryPlace, Easy, true},

		// =========================
		// PROFESSIONS
		// =========================
		{"doctor", CategoryProfession, Easy, true},
		{"chef", CategoryProfession, Easy, true},
		{"astronaut", CategoryProfession, Medium, true},
		{"firefighter", CategoryProfession, Medium, true},
		{"detective", CategoryProfession, Medium, true},
		{"barber", CategoryProfession, Easy, true},
		{"scientist", CategoryProfession, Medium, true},
		{"teacher", CategoryProfession, Easy, true},
		{"pilot", CategoryProfession, Medium, true},
		{"magician", CategoryProfession, Medium, true},

		// =========================
		// SPORTS
		// =========================
		{"football", CategorySports, Easy, true},
		{"basketball", CategorySports, Easy, true},
		{"boxing", CategorySports, Medium, true},
		{"surfing", CategorySports, Medium, true},
		{"skiing", CategorySports, Medium, true},
		{"tennis", CategorySports, Easy, true},
		{"archery", CategorySports, Medium, true},
		{"weightlifting", CategorySports, Medium, true},
		{"karate", CategorySports, Medium, true},
		{"skateboarding", CategorySports, Medium, true},

		// =========================
		// TECHNOLOGY
		// =========================
		{"robot", CategoryTechnology, Easy, true},
		{"computer", CategoryTechnology, Easy, true},
		{"hacker", CategoryTechnology, Medium, true},
		{"drone delivery", CategoryTechnology, Hard, true},
		{"virtual reality", CategoryTechnology, Medium, true},
		{"smartwatch", CategoryTechnology, Easy, true},
		{"cyborg", CategoryTechnology, Medium, true},
		{"server room", CategoryTechnology, Hard, true},
		{"satellite", CategoryTechnology, Medium, true},
		{"spaceship", CategoryTechnology, Medium, true},

		// =========================
		// EMOTIONS / ABSTRACT
		// =========================
		{"panic", CategoryEmotion, Hard, true},
		{"confusion", CategoryEmotion, Hard, true},
		{"jealousy", CategoryEmotion, Hard, true},
		{"awkwardness", CategoryEmotion, Hard, true},
		{"nostalgia", CategoryEmotion, Hard, true},
		{"celebration", CategoryEmotion, Medium, true},
		{"loneliness", CategoryEmotion, Hard, true},
		{"excitement", CategoryEmotion, Medium, true},
		{"fear", CategoryEmotion, Medium, true},
		{"surprise", CategoryEmotion, Medium, true},
	}
}
