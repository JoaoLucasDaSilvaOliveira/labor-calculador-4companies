# Skill: Go DDD Dependency Injection via `main.go`

## Objetivo

Use esta skill quando o usuário pedir para criar, refatorar ou corrigir **injeção de dependências em projetos Go** organizados com:

- DDD
- Clean Architecture
- Hexagonal Architecture
- `/cmd`, `/internal/domain`, `/internal/application`, `/internal/infra`, `/internal/interfaces`

A regra principal desta skill é:

> O `main.go`, dentro de `/cmd/<executavel>`, é o local correto para composição, bootstrap e injeção das dependências do serviço.

Não remova a injeção de dependências do `main` sem necessidade.  
Não transforme o `main` em regra de negócio.  
O `main` deve montar a aplicação.

---

## Fonte da verdade

Sempre siga o documento de arquitetura fornecido pelo usuário:

- `arquitetura_go.md`

Esse documento define:

```text
interfaces -> application -> domain
infra ------^           -> domain
cmd -> monta tudo
```

Regras obrigatórias:

- `domain` não depende de `application`, `infra`, `interfaces` nem framework
- `application` depende de `domain`
- `infra` pode depender de `domain` para implementar contratos
- `interfaces` depende de `application`
- `cmd` faz composição, bootstrap e inicialização
- detalhes técnicos ficam em `infra`
- handlers, routers, middlewares e adapters de protocolo ficam em `interfaces`
- use cases ficam em `application`
- entidades, value objects, contratos de repositório e erros de negócio ficam em `domain`

---

## Arquivo modelo obrigatório

Use o arquivo abaixo como referência concreta de como o projeto atual faz a montagem:

- `main.example.go`

Esse arquivo é o exemplo-base para entender o estilo esperado de composition root.

Ele mostra o `main()` realizando:

1. load de configuração
2. abertura da conexão com banco
3. criação da factory de repositórios
4. instanciação de repositories
5. instanciação de use cases
6. instanciação de handlers
7. instanciação de route providers
8. criação do router
9. start do servidor HTTP

O Codex deve preservar essa ideia arquitetural.

---

## Quando aplicar

Aplique esta skill quando o usuário pedir coisas como:

- “crie a injeção de dependências”
- “organize o main”
- “monte o bootstrap do serviço”
- “crie o composition root”
- “injete repositories/usecases/handlers”
- “ajuste o main.go”
- “adicione uma nova entidade/usecase/handler no serviço”
- “adicione novas rotas seguindo o padrão”
- “crie um container manual”
- “faça DI sem framework”

---

## Modelo mental

Em Go, neste projeto, a injeção é **manual e explícita**.

Não use automaticamente frameworks como:

- Wire
- Dig
- Fx
- Spring-like containers

Só use ferramenta externa se o usuário pedir explicitamente.

A forma padrão é:

```go
func main() {
    // config
    // conexões técnicas
    // factories
    // repositories
    // use cases
    // handlers
    // route providers
    // router
    // server
}
```

---

## Regra de ouro

O `main.go` pode conhecer todas as camadas porque ele está fora do fluxo de negócio e serve para montar o executável.

Isso é permitido:

```go
package main

import (
    "meuapp/internal/application/usecase/book"
    "meuapp/internal/infra/persistence/sqlite"
    bookHTTP "meuapp/internal/interfaces/http/book"
)
```

Isso é proibido:

```go
// domain/entity/book.go
import "meuapp/internal/infra/persistence/sqlite"
```

```go
// application/usecase/create_book.go
import "meuapp/internal/interfaces/http/book"
```

```go
// domain/entity/author.go
import "net/http"
```

---

## Responsabilidade do `main.go`

O `main.go` deve conter apenas:

- carregar `.env`
- carregar config
- criar conexões externas
- criar factories técnicas
- criar repositories concretos
- criar use cases com seus contratos
- criar handlers com seus use cases
- criar route providers
- criar router
- subir HTTP server, worker, consumer ou CLI
- configurar shutdown/defer/close quando necessário

O `main.go` não deve conter:

- regra de negócio
- validação de entidade
- SQL
- parse detalhado de request HTTP
- resposta HTTP
- mapeamento de status code
- DTO de request/response
- algoritmo de caso de uso

---

## Ordem padrão de montagem

Siga esta ordem sempre que possível:

```go
func main() {
    loadEnv()
    cfg := loadConfig()

    database := openDatabase(cfg.DB)
    defer database.Close()

    repositoryFactory := sqlite.NewRepositoryFactory(database)

    authorRepository := repositoryFactory.NewAuthorRepository()
    bookRepository := repositoryFactory.NewBookRepository()

    authorCreateUseCase := authorUC.CreateUseCase{Repo: authorRepository}
    authorListUseCase := authorUC.ListUseCase{Repo: authorRepository}

    bookCreateUseCase := bookUC.CreateUseCase{Repo: bookRepository}
    bookListUseCase := bookUC.ListUseCase{Repo: bookRepository}

    authorHandler := authorHandlerLib.NewHandler(
        &authorCreateUseCase,
        &authorListUseCase,
    )

    bookHandler := bookHandlerLib.NewHandler(
        &bookCreateUseCase,
        &bookListUseCase,
    )

    authorRouteProvider := authorHandlerLib.NewAuthorRouteProvider(authorHandler)
    bookRouteProvider := bookHandlerLib.NewBookRouteProvider(bookHandler)

    mux := shared.NewRouter(authorRouteProvider, bookRouteProvider)

    http.ListenAndServe(":"+cfg.Server.Port, mux)
}
```

Adapte nomes ao projeto real.

---

## Padrão esperado ao adicionar uma nova feature

Quando o usuário pedir para adicionar uma nova entidade, por exemplo `Publisher`, faça o encadeamento inteiro.

### 1. Domain

Criar ou usar:

```text
/internal/domain/entity/publisher.go
/internal/domain/repository/publisher_repository.go
/internal/domain/errors/...
```

Contratos ficam no domínio se representarem persistência essencial ao negócio.

### 2. Infra

Criar implementação concreta:

```text
/internal/infra/persistence/sqlite/publisher_repository.go
```

Atualizar factory:

```go
func (f *RepositoryFactory) NewPublisherRepository() domainrepo.PublisherRepository {
    return NewPublisherRepository(f.db)
}
```

### 3. Application

Criar use cases:

```text
/internal/application/usecase/publisher/create.go
/internal/application/usecase/publisher/list.go
/internal/application/usecase/publisher/list_by_id.go
/internal/application/usecase/publisher/update.go
/internal/application/usecase/publisher/delete.go
```

Cada use case recebe contratos por campo ou construtor:

```go
type CreateUseCase struct {
    Repo repository.PublisherRepository
}
```

### 4. Interfaces

Criar handler e rotas:

```text
/internal/interfaces/http/publisher/handler.go
/internal/interfaces/http/publisher/routes.go
/internal/interfaces/http/publisher/request.go
/internal/interfaces/http/publisher/response.go
```

Handler recebe use cases:

```go
func NewHandler(
    create *publisherUC.CreateUseCase,
    list *publisherUC.ListUseCase,
    listByID *publisherUC.ListByIDUseCase,
    update *publisherUC.UpdateUseCase,
    delete *publisherUC.DeleteUseCase,
) *Handler {
    return &Handler{
        create: create,
        list: list,
        listByID: listByID,
        update: update,
        delete: delete,
    }
}
```

### 5. Cmd / main

Atualizar o `main.go`, pois é ali que a feature entra no executável:

```go
publisherRepository := repositoryFactory.NewPublisherRepository()

publisherCreateUseCase := publisherUC.CreateUseCase{Repo: publisherRepository}
publisherListUseCase := publisherUC.ListUseCase{Repo: publisherRepository}
publisherListByIdUseCase := publisherUC.ListByIdUseCase{Repo: publisherRepository}
publisherUpdateUseCase := publisherUC.UpdateUseCase{Repo: publisherRepository}
publisherDeleteUseCase := publisherUC.DeleteUseCase{Repo: publisherRepository}

publisherHandler := publisherHandlerLib.NewHandler(
    &publisherCreateUseCase,
    &publisherListUseCase,
    &publisherListByIdUseCase,
    &publisherUpdateUseCase,
    &publisherDeleteUseCase,
)

publisherRouteProvider := publisherHandlerLib.NewPublisherRouteProvider(publisherHandler)

mux := shared.NewRouter(
    authorRouteProvider,
    bookRouteProvider,
    publisherRouteProvider,
)
```

---

## Padrão de DI para use cases

Prefira use cases explícitos.

Aceito:

```go
type CreateUseCase struct {
    Repo repository.AuthorRepository
}
```

Também aceito:

```go
func NewCreateUseCase(repo repository.AuthorRepository) *CreateUseCase {
    return &CreateUseCase{Repo: repo}
}
```

Evite:

```go
type CreateUseCase struct {
    DB *sql.DB
}
```

Use case deve depender de contrato, não de banco concreto.

---

## Padrão de DI para handlers

Handler deve receber use cases prontos.

Aceito:

```go
type Handler struct {
    create *authorUC.CreateUseCase
    list   *authorUC.ListUseCase
}
```

Evite handler criando dependências:

```go
func NewHandler() *Handler {
    db := sqlite.OpenDB()
    repo := sqlite.NewAuthorRepository(db)
    uc := authorUC.CreateUseCase{Repo: repo}
    return &Handler{create: &uc}
}
```

Handler não monta infra.  
Handler adapta HTTP para application.

---

## Padrão de DI para repositories

Repository concreto mora em `infra`.

Aceito:

```go
type AuthorRepository struct {
    db *sql.DB
}
```

Ele implementa contrato do domínio ou da aplicação.

Factory concreta pode ficar em `infra/persistence/<driver>`:

```go
type RepositoryFactory struct {
    db *sql.DB
}

func NewRepositoryFactory(db *sql.DB) *RepositoryFactory {
    return &RepositoryFactory{db: db}
}

func (f *RepositoryFactory) NewAuthorRepository() repository.AuthorRepository {
    return NewAuthorRepository(f.db)
}
```

---

## Padrão de DI para router

O router deve receber route providers ou handlers já montados.

Aceito:

```go
mux := shared.NewRouter(authorRouteProvider, bookRouteProvider)
```

Evite router criando use cases ou repositories:

```go
func NewRouter() http.Handler {
    db := sqlite.OpenDB()
    repo := sqlite.NewAuthorRepository(db)
    uc := authorUC.CreateUseCase{Repo: repo}
    handler := authorHTTP.NewHandler(&uc)
    ...
}
```

Router organiza rotas.  
Não monta aplicação inteira.

---

## Checklist obrigatório antes de finalizar uma alteração

Antes de responder ou aplicar patch, verifique:

- [ ] O `main.go` continua sendo a composition root?
- [ ] A injeção de dependências acontece no `cmd`?
- [ ] O domínio continua sem imports de infra/interfaces/framework?
- [ ] Use cases dependem de contratos, não de banco concreto?
- [ ] Handlers recebem use cases, não criam repositories?
- [ ] Router recebe handlers/providers, não cria banco?
- [ ] SQL ficou em `infra/persistence`?
- [ ] Request/response ficou em `interfaces/http`?
- [ ] Regras de negócio ficaram em `domain` ou `application`?
- [ ] O padrão segue o `main.example.go`?

---

## Como responder ao usuário

Quando gerar ou corrigir código, explique de forma direta:

1. o que foi montado no `main`
2. quais dependências foram injetadas
3. quais arquivos precisaram ser alterados
4. se alguma camada estava violando a arquitetura
5. como testar rapidamente

Exemplo de resposta:

```text
Ajustei o DI seguindo o main.example.go:

- repositoryFactory cria os repositories concretos
- use cases recebem repositories por contrato
- handlers recebem use cases
- route providers recebem handlers
- shared.NewRouter recebe os providers
- main.go continua como composition root

Não movi a injeção para dentro dos handlers porque isso acoplaria interfaces com infra.
```

---

## Proibições fortes

Não faça:

- criar repository dentro de handler
- criar use case dentro de router
- abrir banco dentro de use case
- colocar SQL no application
- colocar DTO HTTP no domain
- colocar status code HTTP no domain
- remover DI do main só para “limpar”
- esconder tudo em um container mágico sem necessidade
- transformar `cmd` em camada de regra de negócio

---

## Quando o `main` ficar grande

Se o `main` crescer muito, pode extrair funções privadas no próprio pacote `main`.

Aceito:

```go
func main() {
    app := buildApp()
    app.Start()
}

func buildApp() *App {
    // composição e DI continuam no pacote main/cmd
}
```

Ou:

```go
func buildRepositories(factory *sqlite.RepositoryFactory) Repositories
func buildUseCases(repos Repositories) UseCases
func buildHandlers(ucs UseCases) Handlers
```

Mas a responsabilidade continua sendo de `/cmd`.

Não mova a composição para `domain`, `application`, `infra` ou `interfaces`.

Se criar um pacote de bootstrap, ele deve continuar sendo claramente composição do executável, por exemplo:

```text
/cmd/api/main.go
/cmd/api/bootstrap.go
```

ou, se o projeto aceitar:

```text
/internal/bootstrap
```

Mas só use `/internal/bootstrap` se o usuário pedir ou se o projeto já tiver esse padrão. O padrão desta skill é manter em `/cmd`.

---

## Exemplo de estrutura esperada para o arquivo principal

```text
/cmd
  /api
    main.go
    bootstrap.go optional
```

O arquivo `main.go` pode ficar parecido com o `main.example.go`:

```go
package main

func main() {
    // LOAD CONFIG
    // DB
    // REPOSITORIES
    // USE CASES
    // HANDLERS
    // ROUTES
    // SERVER
}
```

Essa separação por blocos com comentários é aceitável e coerente com o arquivo modelo.

---

## Tom de trabalho do Codex

Ao aplicar esta skill:

- seja conservador com arquitetura
- preserve o padrão existente do projeto
- use nomes já presentes no código
- siga o `main.example.go` como referência concreta
- não invente framework de DI
- não “corrija” o main removendo a composição dele
- prefira DI explícita, legível e manual

A meta é um código simples, explícito e alinhado com DDD/Clean/Hexagonal em Go.
