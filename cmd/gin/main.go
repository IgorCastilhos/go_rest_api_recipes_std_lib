package main

import (
	"github.com/IgorCastilhos/go_rest_api_recipes_std_lib/pkg/recipes"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"net/http"
)

func main() {
	// Cria um roteador Gin
	router := gin.Default()
	// Instancia o recipe handler e provisiona uma implementação da store de dados
	store := recipes.NewMemStore()
	recipesHandler := NewRecipeHandler(store)

	// Registra Rotas
	router.GET("/", homePage)
	router.GET("/receitas", recipesHandler.ListRecipes)
	router.POST("/receitas", recipesHandler.CreateRecipe)
	router.GET("/receitas/:id", recipesHandler.GetRecipe)
	router.PUT("/receitas/:id", recipesHandler.UpdateRecipe)
	router.DELETE("/receitas/:id", recipesHandler.DeleteRecipe)

	// Inicia o servidor
	err := router.Run()
	if err != nil {
		return
	}
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "Bem-vindo(a)! Esta é a página inicial.")
}

type RecipesHandler struct {
	store recipeStore
}

func NewRecipeHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Remove(name string) error
}

// Definindo a assinatura das funções handler

func (h RecipesHandler) CreateRecipe(c *gin.Context) {
	// *gin.Context provê uma func ShouldBindJSON que converte o corpo do HTTP em uma struct.
	// ShouldBindJSON é essencialmente um wrapper em torno da função JSON.Marshal usada nas outras APIs

	// Pega o corpo da requisição e converte em recipes.Recipe
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := slug.Make(recipe.Name)

	err := h.store.Add(id, recipe)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h RecipesHandler) ListRecipes(c *gin.Context) {
	// Chama a store para pegar a lista de receitas
	r, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, r)
}
func (h RecipesHandler) GetRecipe(c *gin.Context)    {}
func (h RecipesHandler) UpdateRecipe(c *gin.Context) {}
func (h RecipesHandler) DeleteRecipe(c *gin.Context) {}

// NewRecipesHandler é o construtor de RecipesHandler
func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}
