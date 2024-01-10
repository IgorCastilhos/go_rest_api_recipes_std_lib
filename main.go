package main

import (
	"github.com/IgorCastilhos/go_rest_api_recipes_std_lib/recipes/pkg/recipes"
	"net/http"
	"regexp"
)

// As duas regexes diferenciam os dois possíveis URIs (/recipes vs. /recipes/<id>)
var (
	RecipeRe       = regexp.MustCompile(`^/receitas/*$`)
	RecipeReWithID = regexp.MustCompile(`^/receitas/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func main() {
	// Cria a Store e o Recipe Handler
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	// Cria um multiplexador de requisições
	// Recebe solicitações HTTP e as envia para os handlers correspondentes
	mux := http.NewServeMux()
	// Registra as rotas e os handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/receitas", recipesHandler)
	mux.Handle("/receitas/", recipesHandler)
	// Executa o servidor
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

type homeHandler struct{}

// Na STD lib, um handler é uma interface que define a assinatura do método
// ServeHTTP(w http.ResponseWriter, r *http.Request). Portanto, para criar um
// handler, é necessário criar uma estrutura (struct) e implementar o método ServeHTTP
func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bem-vindo à página inicial!"))
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	List() (map[string]recipes.Recipe, error)
	Remove(name string) error
}

// RecipesHandler - implementa http.Handler e despacha requisições para a loja
type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}

// Roteamento (Routing)
func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.CreateRecipe(w, r)
		return
	case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
		h.ListRecipes(w, r)
		return
	case r.Method == http.MethodGet && RecipeReWithID.MatchString(r.URL.Path):
		h.GetRecipe(w, r)
		return
	case r.Method == http.MethodPut && RecipeReWithID.MatchString(r.URL.Path):
		h.UpdateRecipe(w, r)
		return
	case r.Method == http.MethodDelete && RecipeReWithID.MatchString(r.URL.Path):
		h.DeleteRecipe(w, r)
		return
	default:
		return
	}
}

func (h *RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {}
func (h *RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request)  {}
func (h *RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request)    {}
func (h *RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {}
func (h *RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {}
