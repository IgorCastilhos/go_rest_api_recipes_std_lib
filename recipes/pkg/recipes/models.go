package recipes

// Recipe - Modelos para as receitas
// Representa uma receita
type Recipe struct {
	Name        string       `json:"name,omitempty"`
	Ingredients []Ingredient `json:"ingredients,omitempty"`
}

// Ingredient - Representa ingredientes individualmente
type Ingredient struct {
	Name string `json:"name,omitempty"`
}
