package main

import (
	"encoding/json"
	"errors"
	recipes "github.com/IgorCastilhos/go_rest_api_recipes_std_lib/pkg/recipes"
	"github.com/gosimple/slug"
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

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
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

// CreateRecipe - Lê os arquivos JSON transportados pelo corpo da requisição HTTP
// e converte em uma instância de recipes.Recipe
func (h *RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	// Objeto de receita que vai ser populado pelos dados JSON
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Converte o nome da receita em uma string URL mais amigável
	resourceID := slug.Make(recipe.Name)
	if err := h.store.Add(resourceID, recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Define o status code para 200
	w.WriteHeader(http.StatusOK)
}
func (h *RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	// Retorna as receitas da loja
	resources, err := h.store.List()
	// Converte a lista retornada em JSON usando a função Marshal
	jsonBytes, err := json.Marshal(resources)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	// Adiciona os dados em JSON para a resposta HTTP usando a função Write
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h *RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	// Recebe o nome do recurso via URl com /recipes/slug-nome-receita
	matches := RecipeReWithID.FindStringSubmatch(r.URL.Path)

	// Espera que as correspondências sejam length >= 2 (full str + 1 grupo)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	// A primeira correspondência ou match ao chamar FindStringSubmatch
	// é sempre a string correspondente completa e, em seguida, todos os
	// subgrupos. Olhando para a regex RecipeReWithID, o primeiro grupo
	// correspondente é o ID do recurso. Só precisamos chamar a função
	// Get da loja com esse ID
	recipe, err := h.store.Get(matches[1])
	if err != nil {
		// caso especial de erro NotFound
		if errors.Is(err, recipes.NotFoundErr) {
			NotFoundHandler(w, r)
			return
		}
		// Qualquer outro erro
		InternalServerErrorHandler(w, r)
		return
	}
	// Converte a struct em dados JSON
	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Adiciona os dados em JSON para a resposta HTTP usando a função Write
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h *RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {}
func (h *RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {}
