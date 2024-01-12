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
	router.Run()
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
	// Pega o corpo da requisição e converte em recipes.Recipe
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := slug.Make(recipe.Name)

	h.store.Add(id, recipe)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h RecipesHandler) ListRecipes(c *gin.Context) {
	r, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, r)
}
func (h RecipesHandler) GetRecipe(c *gin.Context) {
	id := c.Param("id")

	recipe, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(200, recipe)
}
func (h RecipesHandler) UpdateRecipe(c *gin.Context) {
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")

	err := h.store.Update(id, recipe)
	if err != nil {
		if err == recipes.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h RecipesHandler) DeleteRecipe(c *gin.Context) {
	id := c.Param("id")

	err := h.store.Remove(id)
	if err != nil {
		if err == recipes.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
