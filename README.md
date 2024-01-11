### Construindo uma API REST com a Standard Library

* Implementar uma **API** que ajuda os usuários a encontrar receitas que eles podem fazer com os ingredientes na geladeira deles. 

| Ação      | Verbo  | Caminho        | Descrição                                         |
|-----------|--------|----------------|---------------------------------------------------|
| Criar     | POST   | /receitas      | Criar uma entidade representada pelo payload JSON |
| Listar    | GET    | /receitas      | Obter todas as entidades do recurso               |
| Ler       | GET    | /receitas/<id> | Obter uma única entidade                          |
| Atualizar | PUT    | /receitas/<id> | Atualizar uma entidade com o payload JSON         |
| Excluir   | DELETE | /receitas/<id> | Excluir uma entidade                              |
### Todo

1. [x]  Routing
2. [x]  Criar
3. [x]  Listar
4. [x]  Ler
5. [x]  Atualizar
6. [x]  Excluir

### Construindo uma API REST com o pacote de roteamento gorilla/mux

* O roteador combina URLs com funções handler e as chama adequadamente, o que ajuda a reduzir o código necessário no multiplexador criado para usar a Biblioteca Padrão
