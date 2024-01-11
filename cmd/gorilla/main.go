package main

import (
	"encoding/json"
	"github.com/IgorCastilhos/go_rest_api_recipes_std_lib/pkg/recipes"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

func main() {
	// Cria a Store e o Recipe Handler
	store := recipes.NewMemStore()
	// Cria o roteador
	router := mux.NewRouter()
	//router.HandleFunc("/", &home{})
	//router.Use("/", loggingMiddleware)

	s := router.PathPrefix("/receitas").Subrouter()

	// Registra as rotas
	NewRecipesHandler(store, s)

	// Inicia o servidor
	err := http.ListenAndServe(":8010", router)
	if err != nil {
		return
	}
}

func NewRecipesHandler(s recipeStore, router *mux.Router) *RecipesHandler {
	handler := &RecipesHandler{
		store: s,
	}

	router.HandleFunc("/", handler.ListRecipes).Methods("GET")
	router.HandleFunc("/", handler.CreateRecipe).Methods("POST")
	router.HandleFunc("/{id}", handler.GetRecipe).Methods("GET")
	router.HandleFunc("/{id}", handler.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/{id}", handler.DeleteRecipe).Methods("DELETE")

	return handler
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte("500 Internal Server Error"))
	if err != nil {
		return
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("404 Not Found"))
	if err != nil {
		return
	}
}

type RecipesHandler struct {
	store recipeStore
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Remove(name string) error
}

func (h RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	// objeto da receita que vai ser populado pelo JSON payload
	var recipe recipes.Recipe

	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Cria uma URL mais fácil de entender pra usar como ID
	resourceID := slug.Make(recipe.Name)

	if err := h.store.Add(resourceID, recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
}
func (h RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.store.List()

	jsonBytes, err := json.Marshal(recipes)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	// Quando o ID da receita (slug) é passado como parâmetro, use mux.Vars() com a requisição como parâmetro.
	// Essa função retorna um mapa de parâmetros correspondentes com o padrão da URL definida no router (nesse caso
	// id de /receitas/{id}).
	id := mux.Vars(r)["id"]

	recipe, err := h.store.Get(id)
	if err != nil {
		if err == recipes.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Recebe objeto que vai ser populado pelo JSON
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Update(id, recipe); err != nil {
		if err == recipes.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}
	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.store.Remove(id); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Bem-vindo! Esta é a página inicial."))
	if err != nil {
		return
	}
}
