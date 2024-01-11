package main

import (
	"bytes"
	"github.com/IgorCastilhos/go_rest_api_recipes_std_lib/pkg/recipes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile("../../testdata/" + name)
	if err != nil {
		t.Errorf("Não conseguiu ler %v", name)
	}
	return content
}

func TestRecipesHandlerCRUD_Integration(t *testing.T) {

	//	Cria uma MemStore e um Recipe Handler
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	//	Testa os dados
	queijoEPresunto := readTestData(t, "receita_queijo_e_presunto.json")
	queijoEPresuntoReader := bytes.NewReader(queijoEPresunto)

	queijoPresuntoComManteiga := readTestData(t, "receita_queijo_presunto_com_manteiga.json")
	queijoPresuntoComManteigaReader := bytes.NewReader(queijoPresuntoComManteiga)

	//	CREATE - adiciona uma nova receita
	req := httptest.NewRequest(http.MethodPost, "/receitas", queijoEPresuntoReader)
	w := httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()
	assert.Equal(t, 200, result.StatusCode)

	saved, _ := store.List()
	assert.Len(t, saved, 1)

	// GET - Encontra o registro criado no CREATE
	req = httptest.NewRequest(http.MethodGet, "/receitas/torrada-de-queijo-e-presunto", queijoEPresuntoReader)
	w = httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	result = w.Result()
	defer result.Body.Close()
	assert.Equal(t, 200, result.StatusCode)

	data, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Erro inesperado: %v", err)
	}

	assert.JSONEq(t, string(queijoEPresunto), string(data))

	// UPDATE - adiciona manteiga à receita
	req = httptest.NewRequest(http.MethodPut, "/receitas/torrada-de-queijo-e-presunto", queijoPresuntoComManteigaReader)
	w = httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	result = w.Result()
	defer result.Body.Close()
	assert.Equal(t, 200, result.StatusCode)

	updatePresuntoEQueijo, err := store.Get("torrada-de-queijo-e-presunto")
	assert.NoError(t, err)

	assert.Contains(t, updatePresuntoEQueijo.Ingredients, recipes.Ingredient{Name: "manteiga"})

	//DELETE - remove a receita da torrada
	req = httptest.NewRequest(http.MethodDelete, "/receitas/torrada-de-queijo-e-presunto", nil)
	w = httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	result = w.Result()
	defer result.Body.Close()
	assert.Equal(t, 200, result.StatusCode)

	saved, _ = store.List()
	assert.Len(t, saved, 0)

}
