package main

import (
	"net/http"
	"regexp"
)

func main() {
	// Cria um multiplexador de requisições
	// Recebe solicitações HTTP e as envia para os handlers correspondentes
	mux := http.NewServeMux()
	// Registra as rotas e os handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/receitas", &RecipesHandler{})
	mux.Handle("/receitas/", &RecipesHandler{})
	// Executa o servidor
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

type RecipesHandler struct{}

func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bem-vindo à página de receitas!"))
}

type homeHandler struct{}

// Na STD lib, um handler é uma interface que define a assinatura do método
// ServeHTTP(w http.ResponseWriter, r *http.Request). Portanto, para criar um
// handler, é necessário criar uma estrutura (struct) e implementar o método ServeHTTP
func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bem-vindo à página inicial!"))
}

// As duas regexes diferenciam os dois possíveis URIs (/recipes vs. /recipes/<id>)
var (
	RecipeRe       = regexp.MustCompile(`^/receitas/*$`)
	RecipeReWithID = regexp.MustCompile(`^/receitas/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

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
