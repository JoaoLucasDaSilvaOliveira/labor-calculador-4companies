# Arquitetura de navegação da interface Fyne

Este documento explica, ponto a ponto, a implementação do router da interface,
a refatoração da homepage e as decisões arquiteturais tomadas.

O objetivo é permitir que a solução seja:

- Entendida sem depender do histórico do chat.
- Reproduzida em outras páginas e projetos Fyne.
- Evoluída sem concentrar responsabilidades na homepage.
- Testada sem abrir uma janela gráfica real.

## Estado vigente

As seções que mencionam `HomePage`, `RouteHome` ou os primeiros placeholders
descrevem etapas históricas da refatoração. O código vigente usa `MainPage` e
`RouteMain`.

Atualmente a composição é:

```text
RootRouter
└── MainPage
    ├── Sidebar persistente
    └── WorkspaceRouter.View()
        ├── QuickAccessPage
        ├── CompanyPage
        └── EmployeePage
```

As telas de empresa e funcionário são somente leitura nesta etapa. Elas
consultam os use cases reais, filtram funcionários por empresa e recibos por
funcionário e distinguem carregamento, estado vazio e erro recuperável. A
composição visual segue os protótipos de `cmd/main/testes_ui/1.1` e
`cmd/main/testes_ui/1.2`, incluindo os painéis customizados, dimensões e ações
de rodapé. A edição inline está visualmente preparada, mas ainda não persiste
alterações.

## Sumário

1. [Problema original](#1-problema-original)
2. [Decisão arquitetural](#2-decisão-arquitetural)
3. [Estrutura final](#3-estrutura-final)
4. [Fluxo completo da aplicação](#4-fluxo-completo-da-aplicação)
5. [IDs de rota e parâmetros](#5-ids-de-rota-e-parâmetros)
6. [Contrato Navigator](#6-contrato-navigator)
7. [Estrutura interna do Router](#7-estrutura-interna-do-router)
8. [Registro e construção de páginas](#8-registro-e-construção-de-páginas)
9. [Push, Replace, Back e histórico](#9-push-replace-back-e-histórico)
10. [Fallback e tratamento de erro](#10-fallback-e-tratamento-de-erro)
11. [AnimatedContent e transições](#11-animatedcontent-e-transições)
12. [Refatoração da homepage](#12-refatoração-da-homepage)
13. [Estados e transições da sidebar](#13-estados-e-transições-da-sidebar)
14. [Comunicação entre lista, homepage e router](#14-comunicação-entre-lista-homepage-e-router)
15. [Injeção de dependências e composition root](#15-injeção-de-dependências-e-composition-root)
16. [Páginas de destino](#16-páginas-de-destino)
17. [Patterns utilizados](#17-patterns-utilizados)
18. [Como adicionar uma nova rota](#18-como-adicionar-uma-nova-rota)
19. [Como tratar mudanças globais](#19-como-tratar-mudanças-globais)
20. [Testes implementados](#20-testes-implementados)
21. [Comandos de verificação](#21-comandos-de-verificação)
22. [Limitações e próximos passos](#22-limitações-e-próximos-passos)
23. [Glossário](#23-glossário)
24. [Ordem recomendada de leitura do código](# 24.%20Ordem%20recomendada%20de%20leitura do código)
25. [Checklist para novas páginas](#25-checklist-para-novas-páginas)

---

## 1. Problema original

O arquivo `internal/ui/pages/homepage.go` concentrava responsabilidades demais
dentro de `NewHomePage`.

O mesmo construtor:

1. Criava a estrutura visual da homepage.
2. Criava os componentes de acesso rápido.
3. Mantinha o estado da sidebar.
4. Declarava callbacks para abrir, fechar e pesquisar.
5. Controlava a animação do `HSplit`.
6. Substituía o conteúdo variável da sidebar.
7. Instalava diretamente a homepage na janela principal.

Além disso, havia funções declaradas dentro de outras funções. Algumas eram
callbacks necessários do Fyne, mas outras continham comportamento relevante da
página. Isso dificultava:

- Ler o fluxo de cima para baixo.
- Testar partes isoladas.
- Reutilizar a janela para outras páginas.
- Navegar entre telas.
- Descobrir quem era responsável por cada estado.
- Adicionar comportamento de voltar.
- Preservar uma página anterior.

O construtor terminava com:

```go
app.MasterWindow.SetContent(homepageBox)
```

Isso fazia a própria página controlar a janela. Se todas as páginas seguissem
esse modelo, cada uma precisaria conhecer `MainApplication` e substituir o root
da janela por conta própria.

Também existiam callbacks provisórios:

```go
func() { fmt.Println("tapped1") }
func() { fmt.Println("tapped2") }
```

Eles foram substituídos por navegação real.

---

## 2. Decisão arquitetural

A decisão foi criar um router específico para navegação e manter
`AnimatedContent` como um componente de transição independente.

A separação é:

```text
Router
├── conhece rotas
├── constrói páginas por factories
├── mantém histórico
├── implementa voltar
├── valida fallback
└── pede ao AnimatedContent para exibir uma view

AnimatedContent
├── mantém um container estável
├── executa fade-out
├── substitui o CanvasObject
├── executa fade-in
└── ignora callbacks de transições antigas
```

O `AnimatedContent` não foi transformado em router porque ele também é útil em
situações que não são navegação:

- Resultado de pesquisa.
- Conteúdo variável de uma sidebar.
- Etapas visuais de um formulário.
- Estados de carregamento, erro e sucesso.
- Trocas internas de um componente.

Se o código de navegação fosse colocado dentro dele, um componente visual
passaria a conhecer rota, histórico, parâmetros e fallback. Isso violaria o
princípio de responsabilidade única.

O router, por outro lado, usa `AnimatedContent` por composição.

---

## 3. Estrutura final

```text
cmd/main/
├── main.go
└── ui_routes.go

internal/ui/
├── application/
│   └── app.go
├── components/
│   ├── animated_content.go
│   ├── companies_list.go
│   └── homepage_sidebar.go
├── navigation/
│   ├── router.go
│   ├── router_test.go
│   └── routes.go
└── pages/
    ├── homepage.go
    ├── company.go
    └── calculation.go
```

Responsabilidades:

| Arquivo                                      | Responsabilidade                                              |
| -------------------------------------------- | ------------------------------------------------------------- |
| `cmd/main/main.go`                           | Montar banco, repository, use case, aplicação, router e rotas |
| `cmd/main/ui_routes.go`                      | Associar IDs de rota aos construtores de página               |
| `internal/ui/application/app.go`             | Criar a aplicação Fyne, a janela e iniciar o event loop       |
| `internal/ui/navigation/routes.go`           | Declarar IDs e parâmetros de rotas                            |
| `internal/ui/navigation/router.go`           | Controlar registro, histórico, fallback e página atual        |
| `internal/ui/components/animated_content.go` | Trocar conteúdo usando fade                                   |
| `internal/ui/pages/homepage.go`              | Controlar layout e estado local da homepage                   |
| `internal/ui/components/homepage_sidebar.go` | Construir os três estados visuais da sidebar                  |
| `internal/ui/components/companies_list.go`   | Consultar e apresentar empresas, emitindo seleção             |
| `internal/ui/pages/company.go`               | Destino inicial das rotas de empresa                          |
| `internal/ui/pages/calculation.go`           | Destino inicial da rota de novo cálculo                       |
| `internal/ui/navigation/router_test.go`      | Verificar os comportamentos do router                         |

---

## 4. Fluxo completo da aplicação

O fluxo de inicialização é:

```text
main()
  │
  ├── abre o SQLite
  │
  ├── cria CompanyRepository
  │
  ├── cria GetCompanyUsecase
  │
  ├── cria MainApplication
  │
  ├── cria Router com fallback RouteHome
  │
  ├── cria uiRouteProvider
  │
  ├── registra as factories
  │
  ├── abre RouteHome usando Replace
  │
  └── monta Router.View() na janela e inicia o event loop
```

Depois da inicialização, a janela não troca mais seu root:

```text
MainWindow
└── Router.View()                  <- root estável
    └── AnimatedContent.View()     <- outlet estável
        └── página atual           <- conteúdo variável
```

O termo `outlet` representa a região em que as páginas são renderizadas. É uma
ideia semelhante ao espaço reservado por routers web, mas implementada com um
`fyne.Container`.

---

## 5. IDs de rota e parâmetros

Arquivo:

```text
internal/ui/navigation/routes.go
```

O tipo:

```go
type RouteID string
```

cria um tipo semântico para identificar rotas. Poderia ser usado `string`
diretamente, mas um tipo próprio:

- Documenta a intenção do valor.
- Evita misturar qualquer string com uma rota por acidente.
- Facilita mudar a representação no futuro.
- Faz assinaturas como `Push(route RouteID, ...)` ficarem explícitas.

As rotas atuais são:

```go
const (
    RouteHome              RouteID = "home"
    RouteCompanyDetails    RouteID = "company.details"
    RouteCompanyCreate     RouteID = "company.create"
    RouteCalculationCreate RouteID = "calculation.create"
)
```

Foi utilizada a convenção:

```text
recurso.ação
```

Exemplos:

- `company.details`
- `company.create`
- `calculation.create`

Não é obrigatório usar essa convenção, mas ela deixa as rotas previsíveis.

### Parâmetros da empresa

```go
type CompanyDetailsParams struct {
    CompanyID int
}
```

A rota transporta apenas o ID.

Não foi transportado `*entity.Company` porque uma entidade guardada no
histórico pode ficar desatualizada. Quando a página real de empresa carregar
seus dados, ela poderá consultar o caso de uso usando `CompanyID`.

---

## 6. Contrato Navigator

O contrato injetado nas páginas é:

```go
type Navigator interface {
    Push(route RouteID, params any) error
    Replace(route RouteID, params any) error
    Back() bool
    CanGoBack() bool
    Current() RouteID
}
```

As páginas dependem da interface, e não de `*Router`.

Isso aplica inversão de dependência:

```text
Homepage -> Navigator <- Router
```

A homepage conhece o comportamento de navegação de que precisa, mas não conhece:

- O mapa interno de factories.
- O container usado pelo router.
- A implementação do histórico.
- A animação utilizada.
- A política de fallback.

### Verificação de implementação

```go
var _ Navigator = (*Router)(nil)
```

Essa linha não cria um router.

Ela pede ao compilador que confirme que `*Router` pode ser atribuído a uma
variável do tipo `Navigator`. Se um método for removido ou tiver sua assinatura
alterada, o projeto para de compilar naquele ponto.

O `_` é o identificador em branco do Go. Ele descarta o valor porque o objetivo
é somente validar o tipo.

---

## 7. Estrutura interna do Router

```go
type Router struct {
    outlet        *components.AnimatedContent
    fallbackRoute RouteID
    factories     map[RouteID]PageFactory
    history       []historyEntry
}
```

### `outlet`

É o host visual que permanece montado na janela.

O router não modifica diretamente `fyne.Container.Objects`. Ele delega a troca
para `AnimatedContent`, preservando o fade.

### `fallbackRoute`

É a rota segura usada quando uma navegação não pode ser concluída.

Atualmente:

```go
navigation.NewRouter(navigation.RouteHome)
```

### `factories`

```go
map[RouteID]PageFactory
```

É um mapa em que:

- A chave é o ID da rota.
- O valor é a função que constrói a página.

Exemplo mental:

```text
home               -> homePage
company.details    -> companyPage
company.create     -> companyRegistrationPage
calculation.create -> calculationCreationPage
```

### `history`

```go
[]historyEntry
```

É um slice usado como pilha.

Em Go, um `slice` é uma visão dinâmica sobre um array. Ele possui comprimento e
capacidade e pode crescer com `append`.

---

## 8. Registro e construção de páginas

A assinatura de uma factory é:

```go
type PageFactory func(params any) (fyne.CanvasObject, error)
```

Uma factory:

1. Recebe os parâmetros da rota.
2. Valida e converte esses parâmetros.
3. Injeta as dependências da página.
4. Retorna o `fyne.CanvasObject`.
5. Retorna erro se não puder construir a página.

### Por que retornar erro

Uma factory pode falhar porque:

- Recebeu parâmetros do tipo errado.
- Recebeu um ID inválido.
- Não possui uma dependência obrigatória.
- No futuro, pode precisar carregar configuração inicial.

Sem `error`, uma falha tenderia a virar `panic` ou uma página incompleta.

### Registro

```go
func (r *Router) Register(route RouteID, factory PageFactory) error
```

As validações são:

1. O ID não pode ser vazio.
2. A factory não pode ser `nil`.
3. A mesma rota não pode ser registrada duas vezes.

Erros sentinela foram declarados:

```go
ErrEmptyRouteID
ErrNilPageFactory
ErrRouteAlreadyAdded
ErrRouteNotFound
ErrNilPage
```

Um erro sentinela é uma variável de erro que pode ser reconhecida com:

```go
errors.Is(err, ErrRouteNotFound)
```

Isso é mais seguro do que comparar a mensagem textual.

### Construção de uma entrada

`buildEntry` executa:

```text
RouteID
  │
  ├── procura factory
  ├── chama factory(params)
  ├── valida erro
  ├── valida view não nula
  └── monta historyEntry
```

Uma entrada contém:

```go
type historyEntry struct {
    route  RouteID
    params any
    view   fyne.CanvasObject
}
```

`params` é mantido para identificar o contexto da entrada e permitir futuras
funcionalidades, como lifecycle e reconstrução.

`view` é mantida para que `Back` restaure a mesma instância.

---

## 9. Push, Replace, Back e histórico

### `Push`

```go
func (r *Router) Push(route RouteID, params any) error
```

Fluxo:

```text
Push
  ├── constrói historyEntry
  ├── se falhar, abre fallback
  ├── adiciona entrada ao final do slice
  └── mostra a nova view
```

Exemplo:

```text
antes: [Home]
Push(CompanyDetails)
depois: [Home, CompanyDetails]
```

Use `Push` quando o usuário deve poder voltar.

### `Replace`

```go
func (r *Router) Replace(route RouteID, params any) error
```

Fluxo:

```text
Replace
  ├── constrói historyEntry
  ├── se o histórico estiver vazio, adiciona
  ├── caso contrário, substitui a última entrada
  └── mostra a nova view
```

Exemplo:

```text
antes: [Home, Login]
Replace(Home)
depois: [Home, Home]
```

No caso inicial:

```text
antes: []
Replace(Home)
depois: [Home]
```

Use `Replace` quando a tela anterior não deve ser acessada por `Back`.

Exemplos comuns:

- Inicialização da aplicação.
- Login substituído por home após autenticação.
- Tela temporária substituída pelo resultado final.
- Redirecionamento de fallback.

### `Back`

```go
func (r *Router) Back() bool
```

Fluxo:

```text
Back
  ├── verifica CanGoBack
  ├── remove a última entrada
  ├── recupera a nova última entrada
  └── mostra a view já existente
```

Exemplo:

```text
antes: [Home, CompanyDetails, Calculation]
Back()
depois: [Home, CompanyDetails]
```

O retorno é:

- `true`: uma página anterior foi restaurada.
- `false`: não havia página anterior.

### `CanGoBack`

```go
return len(r.history) > 1
```

Se existe somente uma entrada, ela é a raiz atual. Removê-la deixaria o outlet
sem página.

### `Current`

`Current` não armazena uma cópia separada da rota atual.

Ele deriva a informação:

```go
return r.history[len(r.history)-1].route
```

Isso evita inconsistência entre:

- `history`.
- `currentView`.
- `currentRoute`.
- `lastView`.

Há uma única fonte da verdade: o histórico.

---

## 10. Fallback e tratamento de erro

Quando `Push` ou `Replace` não consegue construir a rota solicitada,
`handleNavigationError`:

1. Tenta construir a rota fallback.
2. Substitui a entrada atual pelo fallback.
3. Mostra a view fallback.
4. Retorna o erro original.

O erro original não é engolido. A página segura é exibida, mas o chamador ainda
pode registrar o problema:

```go
if err := navigator.Push(route, params); err != nil {
    fyne.LogError("Não foi possível navegar.", err)
}
```

Se o próprio fallback também falhar, os erros são unidos com:

```go
errors.Join(...)
```

`errors.Join` cria um erro que representa múltiplas causas e continua
compatível com `errors.Is` e `errors.As`.

### Comportamento atual do fallback

O fallback substitui a posição atual. Ele não é adicionado por `Push`.

Isso evita que o botão voltar retorne imediatamente para uma rota inválida.

---

## 11. AnimatedContent e transições

Arquivo:

```text
internal/ui/components/animated_content.go
```

Estrutura:

```go
type AnimatedContent struct {
    container    *fyne.Container
    overlay      *canvas.Rectangle
    animation    *fyne.Animation
    hasContent   bool
    transitionID uint64
}
```

### Container estável

O container é criado com:

```go
container.NewStack(overlay)
```

`NewStack` posiciona objetos no mesmo espaço, um sobre o outro.

Depois da primeira página, os objetos são:

```go
[]fyne.CanvasObject{
    content,
    overlay,
}
```

O overlay fica acima do conteúdo.

### Fade

Para esconder o conteúdo, o overlay recebe a cor de fundo com alpha crescente.

```text
alpha 0   = transparente
alpha 255 = opaco
```

Fluxo da troca:

```text
conteúdo antigo
  │
  ├── fade do overlay até alpha 255
  ├── substitui o conteúdo atrás do overlay
  └── fade do overlay até alpha 0
      │
      └── conteúdo novo visível
```

### `transitionID`

`transitionID` é um contador `uint64`.

`uint64` significa inteiro sem sinal de 64 bits:

- Não aceita valores negativos.
- Possui espaço suficiente para um número extremamente alto de transições.

Cada chamada de `SetContent` faz:

```go
a.transitionID++
currentTransitionID := a.transitionID
```

Os callbacks verificam:

```go
transitionID == a.transitionID
```

Exemplo:

```text
SetContent(Empresa) recebe ID 10
SetContent(Home) recebe ID 11

callback da Empresa:
10 == 11? não
callback ignorado

callback da Home:
11 == 11? sim
callback executado
```

Essa política é chamada de “latest wins”: a solicitação mais recente vence.

Ela evita que uma animação antiga substitua uma página mais nova.

### Por que ainda verificar `a.animation == animation`

O ID identifica a solicitação de conteúdo.

A comparação de ponteiro identifica a instância exata da animação ativa.

As duas verificações protegem níveis diferentes:

- `transitionID`: solicitação de troca.
- `a.animation == animation`: animação concreta.

---

## 12. Refatoração da homepage

Arquivo:

```text
internal/ui/pages/homepage.go
```

O construtor anterior recebia:

```go
func NewHomePage(app *application.MainApplication)
```

O novo construtor recebe dependências específicas:

```go
type HomePageDeps struct {
    Companies components.CompanyFinder
    Navigator navigation.Navigator
}
```

e retorna:

```go
func NewHomePage(deps HomePageDeps) fyne.CanvasObject
```

### Por que retornar `fyne.CanvasObject`

A página deixa de decidir onde será montada.

Ela apenas constrói e devolve sua representação visual. O router decide quando
e onde mostrá-la.

Esse formato também facilita:

- Testes.
- Pré-visualização.
- Reutilização em outro container.
- Troca do mecanismo de navegação.

### Struct `homePage`

```go
type homePage struct {
    deps HomePageDeps

    sidebarContent *components.AnimatedContent
    sidebarHost    fyne.CanvasObject
    split          *container.Split
    splitAnimation *fyne.Animation

    sidebarState           homePageSidebarState
    lastOpenedSidebarState homePageSidebarState
    renderedSidebarState   homePageSidebarState
    hasRenderedSidebar     bool
}
```

A struct guarda apenas:

- Dependências da página.
- Objetos visuais que precisam sobreviver entre métodos.
- Estado local da sidebar.
- Animação local do split.

### Ordem de leitura da homepage

O arquivo foi organizado para ser lido nesta sequência:

1. Dependências.
2. Estado.
3. Construtor.
4. Montagem do layout.
5. Montagem do acesso rápido.
6. Transições de estado.
7. Escolha da sidebar.
8. Animações.
9. Navegação global.

### Métodos extraídos

O comportamento antes aninhado foi dividido em métodos:

```go
build
buildQuickAccess
sidebarDeps
closeSidebar
openSidebar
showSearch
showCompanies
renderSidebar
buildSidebar
setInitialSidebar
expandSidebar
collapseSidebar
swapSidebar
animateSplitOffset
minimumOffset
openCompanyDetails
openCompanyRegistration
openCalculationCreation
navigate
```

Cada nome descreve uma ação única.

### Cores do tema

Os atalhos antes usavam:

```go
color.Black
utils.NewColor(0, 105, 204, 255)
```

Agora usam:

```go
defaultColor := theme.Color(theme.ColorNameForeground)
hoverColor := theme.Color(theme.ColorNamePrimary)
```

Isso evita assumir que o fundo sempre será claro. No tema escuro,
`ColorNameForeground` fornece uma cor de texto apropriada, e
`ColorNamePrimary` mantém o hover alinhado à cor primária configurada no tema.

Também foi removida dessa página a dependência do helper externo usado apenas
para construir a cor fixa.

### Closures que permanecem

Ainda existem funções anônimas em locais como:

```go
p.animateSplitOffset(..., func() {
    p.sidebarContent.SetContent(..., func() {
        // próxima fase
    })
})
```

Elas permanecem porque representam callbacks curtos de conclusão.

Nesse caso, as fases precisam ocorrer em ordem:

```text
animação de largura termina
  -> conteúdo é substituído
     -> próxima animação começa
```

A regra principal de cada fase continua em métodos nomeados. As closures apenas
encadeiam eventos assíncronos do Fyne.

---

## 13. Estados e transições da sidebar

Estados:

```go
const (
    homePageSidebarCompanies homePageSidebarState = iota
    homePageSidebarSearch
    homePageSidebarClosed
)
```

### `iota`

`iota` é um contador do Go usado em blocos de constantes.

Nesse bloco, os valores são:

```text
homePageSidebarCompanies = 0
homePageSidebarSearch    = 1
homePageSidebarClosed    = 2
```

O tipo definido:

```go
type homePageSidebarState uint8
```

impede usar qualquer inteiro comum sem conversão e documenta que o valor
representa um estado da sidebar.

### Por que esses estados não são rotas

Sidebar aberta, busca e sidebar fechada são modos internos da mesma homepage.

Eles não representam destinos globais porque:

- Não substituem a página principal.
- Não precisam entrar no histórico global.
- Não deveriam aparecer ao voltar de uma página de empresa.
- São relevantes somente enquanto a homepage existe.

Portanto:

```text
Home, CompanyDetails, CompanyCreate = rotas globais
Companies, Search, Closed           = estado local
```

### Quatro campos de controle

#### `sidebarState`

Estado solicitado pelo usuário.

#### `renderedSidebarState`

Estado que terminou de ser colocado no container.

Durante uma animação, o solicitado e o renderizado podem ser diferentes.

#### `lastOpenedSidebarState`

Guarda se a sidebar estava mostrando empresas ou pesquisa antes de fechar.

Assim:

```text
Search -> Closed -> Open = Search
Companies -> Closed -> Open = Companies
```

#### `hasRenderedSidebar`

Diferencia:

- Primeira montagem, sem necessidade de fade.
- Trocas posteriores, que devem ser animadas.

Isso remove a dependência implícita do valor zero de
`renderedSidebarState`.

### `renderSidebar`

O método funciona como dispatcher:

```text
primeira renderização?
  -> setInitialSidebar

estava fechada e vai abrir?
  -> expandSidebar

estava aberta e vai fechar?
  -> collapseSidebar

outra troca?
  -> swapSidebar
```

### `minimumOffset`

O `HSplit` usa um offset proporcional, mas o Fyne fornece largura mínima em
pixels.

O cálculo transforma largura em proporção:

```go
availableWidth := splitWidth - dividerWidth
offset := viewMinWidth / availableWidth
```

O divisor é descontado porque ele também ocupa espaço no `HSplit`.

Se a largura disponível ainda for zero, retorna `0` para evitar divisão por
zero.

---

## 14. Comunicação entre lista, homepage e router

O componente de lista não conhece o router.

As duas funções anteriores de lista repetiam a criação de `widget.List`.

A construção comum foi concentrada em:

```go
newCompaniesList(
    finder,
    filter,
    emptyMessage,
    onCompanySelected,
)
```

As funções públicas agora apenas escolhem:

- O filtro.
- A mensagem de resultado vazio.
- O callback de seleção.

### Contrato de consulta

```go
type CompanyFinder interface {
    Execute(query.GetCompanyWithFilter) ([]*entity.Company, error)
}
```

`GetCompanyUsecase` implementa esse contrato implicitamente. Em Go não é
necessário escrever `implements`.

Isso permite substituir o use case por um fake em testes sem entregar à lista
uma aplicação inteira.

### Estados apresentados pela lista

A lista agora diferencia:

- Falha ao carregar: `Não foi possível carregar as empresas.`
- Nenhuma empresa cadastrada.
- Nenhum resultado para a busca.
- Busca ainda vazia: `Digite um nome para buscar.`

O erro técnico não é mostrado diretamente ao usuário. A UI apresenta mensagem
compreensível, enquanto uma evolução futura pode registrar a causa em log ou
observabilidade.

O texto pesquisado é processado com:

```go
strings.TrimSpace(companyName)
```

`TrimSpace` remove espaços, tabs e quebras de linha no começo e no fim. Uma
busca contendo apenas espaços volta ao estado de instrução.

Ele recebe:

```go
onCompanySelected func(companyID int)
```

Quando um item é selecionado:

```text
widget.List.OnSelected
  │
  ├── limpa a seleção
  └── chama onCompanySelected(companyID)
```

A homepage conecta esse callback:

```go
OnCompanySelected: p.openCompanyDetails
```

Então:

```go
func (p *homePage) openCompanyDetails(companyID int) {
    p.navigate(
        navigation.RouteCompanyDetails,
        navigation.CompanyDetailsParams{CompanyID: companyID},
    )
}
```

Fluxo completo:

```text
usuário seleciona empresa
  -> CompaniesList obtém o ID
  -> chama callback da homepage
  -> homepage cria CompanyDetailsParams
  -> Navigator.Push
  -> Router localiza a factory
  -> uiRouteProvider valida os parâmetros
  -> NewCompanyPage constrói a view
  -> Router adiciona ao histórico
  -> AnimatedContent executa o fade
```

### Por que limpar a seleção

O router preserva a homepage e sua lista no histórico.

Se o usuário:

1. Seleciona Empresa 10.
2. Abre a página da Empresa 10.
3. Volta.
4. Seleciona novamente Empresa 10.

o Fyne pode considerar que o item já estava selecionado e não chamar novamente
`OnSelected`.

Por isso:

```go
companiesList.Unselect(id)
```

é chamado antes da navegação.

### Ícone de cadastrar empresa

O botão com `assets.AddCompany`, que antes possuía callback vazio, passou a
receber:

```go
OnAddCompany: p.openCompanyRegistration
```

Assim, o atalho central e o ícone da sidebar usam a mesma rota e o mesmo método
de navegação.

### `time.AfterFunc` e `fyne.Do`

Os ícones da toolbar preservam um pequeno atraso antes de substituir a sidebar:

```go
time.AfterFunc(canvas.DurationShort, func() {
    fyne.Do(tapped)
})
```

Motivo:

1. O botão recebe o clique.
2. O Fyne tem tempo de apresentar o feedback visual do toque.
3. `AfterFunc` executa o callback depois da duração informada.
4. Como o callback do timer não está necessariamente na UI thread, `fyne.Do`
   agenda a alteração visual no contexto correto do Fyne.

---

## 15. Injeção de dependências e composition root

O `main.go` continua sendo o composition root.

Composition root é o ponto em que as implementações concretas são criadas e
conectadas.

Fluxo:

```go
companyRepository := sqlite.NewCompanyRepository(database)
getCompaniesUseCase := companyUC.NewGetCompanyUsecase(companyRepository)

application := uiApplication.NewApplication()
router := navigation.NewRouter(navigation.RouteHome)
routeProvider := newUIRouteProvider(getCompaniesUseCase, router)
```

As dependências apontam:

```text
SQLite Repository -> GetCompanyUsecase -> UI components
Router             -> Navigator        -> Pages
```

### `uiRouteProvider`

O provider recebe dependências prontas:

```go
type uiRouteProvider struct {
    getCompanies *companyUC.GetCompanyUsecase
    navigator    navigation.Navigator
}
```

Ele não:

- Abre banco.
- Cria repository.
- Cria use case.
- Executa regra de negócio.
- Controla o event loop.

Ele adapta rotas para páginas.

Exemplo:

```go
func (p *uiRouteProvider) homePage(_ any) (fyne.CanvasObject, error) {
    return pages.NewHomePage(pages.HomePageDeps{
        Companies: p.getCompanies,
        Navigator: p.navigator,
    }), nil
}
```

### Type assertion da rota de empresa

```go
companyParams, ok := params.(navigation.CompanyDetailsParams)
```

`any` é um alias para `interface{}`.

Uma interface vazia pode guardar valores de qualquer tipo. A type assertion
extrai o tipo concreto esperado.

Se `ok` for `false`, a rota recebeu parâmetros errados e a factory retorna erro.

O ID também é validado:

```go
if companyParams.CompanyID <= 0 {
    return nil, fmt.Errorf(...)
}
```

### `MainApplication`

`MainApplication` não guarda mais use cases.

Ele possui somente a infraestrutura do Fyne:

```go
type MainApplication struct {
    application  fyne.App
    masterWindow fyne.Window
}
```

E o método:

```go
func (a *MainApplication) Run(rootView fyne.CanvasObject)
```

monta o root uma vez e inicia:

```go
a.masterWindow.ShowAndRun()
```

`ShowAndRun` mostra a janela e inicia o event loop do Fyne. O event loop recebe
eventos de teclado, mouse, desenho e animações até a aplicação terminar.

---

## 16. Páginas de destino

Foram criados destinos iniciais para confirmar a navegação.

### Página de empresa — Tela 1.1

```text
internal/ui/pages/company.go
```

Dependências:

```go
type CompanyPageDeps struct {
    CompanyID int
    Company   CompanyFinderByID
    Employees EmployeeFinder
    Navigator navigation.Navigator
}
```

Ela consulta a empresa por ID e os funcionários por `CompanyID`. A seleção de
um funcionário navega para a Tela 1.2 usando `Push`. Falhas de consulta são
apresentadas dentro da própria rota com retorno pelo histórico interno.

### Página de funcionário — Tela 1.2

```text
internal/ui/pages/employee.go
```

A página consulta o funcionário por ID e seus recibos por `EmployeeID`. A lista
de recibos possui estado vazio próprio; falhas de consulta são apresentadas
como erro recuperável. Edição, exclusão e criação de recibos permanecem fora
do incremento atual.

### Cadastro de empresa

Também está em `company.go`.

É um destino inicial para os dois botões de cadastrar empresa:

- Atalho central.
- Ícone da sidebar.

O formulário real ainda deve ser implementado.

### Novo cálculo

```text
internal/ui/pages/calculation.go
```

É o destino inicial do atalho “NOVO CÁLCULO”.

Ele valida que a navegação está conectada sem introduzir prematuramente a regra
e o formulário completo de cálculo.

### Componente compartilhado de voltar

```go
newPageWithBackAction(...)
```

Esse helper evita duplicar:

- Botão com `theme.NavigateBackIcon()`.
- Header.
- Padding.
- Chamada de `navigator.Back()`.

Não foi criada uma abstração maior porque ainda existem poucas páginas.

---

## 17. Patterns utilizados

### Router Pattern

Centraliza a decisão de qual página deve ser apresentada.

Sem router:

```text
Página -> Window.SetContent(outra página)
```

Com router:

```text
Página -> Navigator -> Router -> outlet
```

### Factory Method

`PageFactory` posterga a criação da página até a navegação.

Isso permite:

- Criar uma instância por `Push`.
- Validar parâmetros.
- Injetar dependências.
- Retornar erro.

### Dependency Injection

Dependências são passadas por construtores e structs:

```go
HomePageDeps
CompanyPageDeps
HomePageSidebarDeps
```

Nenhuma página abre banco ou cria repository.

### Dependency Inversion

Páginas dependem de:

```go
navigation.Navigator
components.CompanyFinder
```

e não das implementações concretas.

### Composition

O router contém um `AnimatedContent`.

Ele não herda, copia nem incorpora a responsabilidade da animação.

### State

A sidebar usa um estado enumerado e métodos de transição.

Não foram criadas classes/structs diferentes para cada estado porque seriam
abstrações excessivas para apenas três modos simples.

É uma aplicação proporcional da ideia do State Pattern, sem o peso completo do
pattern GoF.

### Callback / Inversion of Control

A lista emite:

```go
onCompanySelected(companyID)
```

Ela não decide o que acontece depois. O componente pai injeta o comportamento.

### Provider

`uiRouteProvider` agrupa o registro e a criação de páginas.

Isso mantém `main.go` como composition root sem colocar layout detalhado nele.

### Latest Wins

`transitionID` garante que somente a transição mais recente conclua seus
callbacks.

---

## 18. Como adicionar uma nova rota

Exemplo: página de detalhes de funcionário.

### Passo 1: declarar a rota

Em `internal/ui/navigation/routes.go`:

```go
const (
    // rotas existentes...
    RouteEmployeeDetails RouteID = "employee.details"
)
```

### Passo 2: declarar os parâmetros

```go
type EmployeeDetailsParams struct {
    EmployeeID int
}
```

Use somente os valores necessários para identificar o destino.

Evite transportar widgets, repositories ou entidades mutáveis.

### Passo 3: criar a página

Em `internal/ui/pages/employee.go`:

```go
package pages

import (
    "labor-calculador-4companies/internal/ui/navigation"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

type EmployeePageDeps struct {
    EmployeeID int
    Navigator  navigation.Navigator
}

func NewEmployeePage(deps EmployeePageDeps) fyne.CanvasObject {
    title := widget.NewLabel("Funcionário")
    id := widget.NewLabel(fmt.Sprintf("ID: %d", deps.EmployeeID))

    back := widget.NewButton("Voltar", func() {
        deps.Navigator.Back()
    })

    return container.NewVBox(back, title, id)
}
```

Em uma página real, o use case de consulta também deve entrar em
`EmployeePageDeps`.

### Passo 4: adicionar a dependência no provider

Se a página precisar de um use case:

```go
type uiRouteProvider struct {
    getCompanies *companyUC.GetCompanyUsecase
    getEmployee  *employeeUC.GetEmployeeUsecase
    navigator    navigation.Navigator
}
```

O use case deve ser criado no `main.go` e passado ao provider.

### Passo 5: criar a factory

Em `cmd/main/ui_routes.go`:

```go
func (p *uiRouteProvider) employeePage(params any) (fyne.CanvasObject, error) {
    employeeParams, ok := params.(navigation.EmployeeDetailsParams)
    if !ok {
        return nil, fmt.Errorf(
            "parâmetros inválidos para a página do funcionário",
        )
    }

    if employeeParams.EmployeeID <= 0 {
        return nil, fmt.Errorf(
            "ID do funcionário inválido: %d",
            employeeParams.EmployeeID,
        )
    }

    return pages.NewEmployeePage(pages.EmployeePageDeps{
        EmployeeID: employeeParams.EmployeeID,
        Navigator:  p.navigator,
    }), nil
}
```

### Passo 6: registrar

No método `Register`:

```go
if err := router.Register(
    navigation.RouteEmployeeDetails,
    p.employeePage,
); err != nil {
    return err
}
```

### Passo 7: navegar

No callback que conhece o ID:

```go
err := navigator.Push(
    navigation.RouteEmployeeDetails,
    navigation.EmployeeDetailsParams{
        EmployeeID: employeeID,
    },
)
if err != nil {
    fyne.LogError("Não foi possível abrir o funcionário.", err)
}
```

### Passo 8: testar

Adicione testes para:

- Registro.
- Parâmetros inválidos.
- `Push`.
- `Back`.
- Fallback.

Depois execute:

```bash
go test ./internal/ui/... ./cmd/main -count=1
go vet ./internal/ui/... ./cmd/main
go build -o /tmp/labor-calculator-main-build ./cmd/main
```

---

## 19. Como tratar mudanças globais

O router deve controlar navegação, não sincronização de dados.

Exemplo:

```text
empresa é editada
  -> nome muda
  -> lista da homepage precisa refletir a mudança
```

Não coloque no router:

```go
router.NotifyCompanyChanged(...)
```

Isso faria o router virar:

- Navegador.
- Event bus.
- Store global.
- Cache de domínio.

### Opção 1: recarregar ao mostrar

É a primeira recomendação.

Quando a homepage voltar a ficar visível, ela consulta novamente as empresas.

Para isso, uma evolução possível é:

```go
type Page interface {
    View() fyne.CanvasObject
    OnShow()
}
```

O router chamaria `OnShow` ao restaurar a página.

Use quando:

- Dados podem mudar fora da página.
- A consulta é barata.
- Não é preciso atualização em tempo real.

### Opção 2: callback específico

Uma página de edição pode receber:

```go
OnCompanySaved func(companyID int)
```

Use quando:

- Existe um único interessado.
- O fluxo é local e previsível.
- Não vale criar infraestrutura de eventos.

### Opção 3: Observer/eventos tipados

Quando vários componentes precisarem reagir:

```go
type CompanyUpdated struct {
    CompanyID int
}
```

Interessados se inscrevem e recebem o evento.

Use Observer quando:

- Existem vários consumidores.
- Atualizações precisam ocorrer enquanto telas continuam vivas.
- O acoplamento direto por callbacks começou a crescer.

Não foi criado agora porque ainda não existe necessidade concreta suficiente.

---

## 20. Testes implementados

Arquivo:

```text
internal/ui/navigation/router_test.go
```

Os testes usam:

```go
fyneApp := test.NewApp()
t.Cleanup(fyneApp.Quit)
```

`test.NewApp` cria uma aplicação Fyne apropriada para testes, sem depender de
abrir a aplicação desktop real.

`t.Cleanup` registra uma função executada ao final do teste, mesmo se o teste
falhar.

### `TestRouterPushAndBackRestoreHistory`

Verifica:

1. Home é aberta com `Replace`.
2. Detalhes são abertos com `Push`.
3. `Current` aponta para detalhes.
4. `CanGoBack` retorna `true`.
5. Os parâmetros ficam guardados.
6. `Back` retorna para Home.
7. Não existe outra página anterior.

### `TestRouterReplaceDoesNotAddHistory`

Verifica que `Replace` não aumenta o histórico.

Depois de substituir a home por detalhes:

```text
len(history) == 1
Back() == false
```

### `TestRouterUsesFallbackForUnknownRoute`

Tenta abrir uma rota inexistente.

Verifica:

- O erro pode ser reconhecido como `ErrRouteNotFound`.
- A rota atual volta a ser Home.
- O histórico não recebe a rota inválida.

### `TestRouterLatestNavigationWinsDuringTransition`

Sequência:

1. Abre Home.
2. Executa `Push` para detalhes.
3. Executa `Back` imediatamente.
4. Aguarda as animações.
5. Confirma que a rota e o objeto renderizado são Home.

Esse teste protege contra callbacks antigos de animação.

### `TestRouterRejectsDuplicateRegistration`

Registra Home duas vezes e confirma:

```go
errors.Is(err, ErrRouteAlreadyAdded)
```

---

## 21. Comandos de verificação

### Formatação

```bash
gofmt -w $(rg --files cmd/main internal/ui -g '*.go')
```

`gofmt` aplica a formatação oficial do Go.

Elementos:

- `-w`: escreve o resultado formatado de volta nos arquivos.
- `rg --files`: lista arquivos rapidamente.
- `-g '*.go'`: glob que mantém somente arquivos com extensão `.go`.
- `$(...)`: command substitution do shell; a saída de `rg` vira argumentos de
  `gofmt`.

### Análise estática

```bash
go vet ./internal/ui/... ./cmd/main
```

`go vet` procura construções suspeitas que ainda podem compilar, como:

- Formatação incompatível em algumas chamadas.
- Uso incorreto de determinados recursos da biblioteca padrão.
- Erros estruturais reconhecidos pelos analyzers ativos.

`./internal/ui/...` significa todos os pacotes abaixo de `internal/ui`.

### Testes focados

```bash
go test ./internal/ui/... ./cmd/main -count=1
```

Flags e argumentos:

- `./internal/ui/...`: testa recursivamente a UI.
- `./cmd/main`: compila e testa o pacote do executável.
- `-count=1`: não reutiliza resultado de teste em cache.

### Build

```bash
go build -o /tmp/labor-calculator-main-build ./cmd/main
```

Elementos:

- `go build`: compila o pacote e suas dependências.
- `-o`: define o caminho do arquivo de saída.
- `/tmp/labor-calculator-main-build`: mantém o binário fora do repositório.
- `./cmd/main`: pacote do executável.

### Cache temporário

No ambiente de execução usado durante a alteração, o cache global do Go ficou
somente leitura. A verificação foi repetida usando:

```bash
env GOCACHE=/tmp/labor-calculator-go-cache go test ...
```

Partes:

- `env`: executa um comando com variáveis de ambiente específicas.
- `GOCACHE=...`: muda o diretório de cache de compilação do Go somente para
  aquele comando.

Isso não altera permanentemente a configuração do sistema.

### Resultado

Passaram:

```text
go vet da UI e cmd/main
testes de internal/ui e cmd/main
build de cmd/main
```

### Falha preexistente na suíte completa

O comando:

```bash
go test ./... -count=1
```

continua falhando em:

```text
internal/domain/tests/calculator_test.go
TestCalculateVacation
returns_error_when_start_is_not_before_end
```

Mensagem:

```text
expected error início do período aquisitivo deve ser anterior ao fim do
período aquisitivo, got <nil>
```

Essa falha pertence ao cálculo de férias e não foi causada nem alterada pela
implementação do router.

---

## 22. Limitações e próximos passos

### Pesquisa síncrona

Atualmente a busca executa o caso de uso dentro de `Entry.OnChanged`.

Isso significa uma consulta a cada alteração do texto.

Quando a base crescer, implementar:

1. Debounce.
2. Estado de carregamento.
3. Consulta fora da UI thread.
4. Atualização visual com `fyne.Do`.
5. Cancelamento/ID da busca anterior.

Debounce é uma espera curta após a última digitação. Se o usuário continuar
digitando, o temporizador reinicia.

### Histórico sem limite

Cada `Push` preserva uma view.

Uma navegação muito longa pode manter muitos widgets em memória.

Possíveis evoluções:

- Limite máximo de histórico.
- `PopUntil`.
- Rotas que não preservam view.
- Reconstrução da página durante `Back`.

Não foi implementado agora porque aumentaria a complexidade sem necessidade
atual.

### Dados preservados podem ficar antigos

Restaurar a mesma view mantém formulários e scroll, mas também mantém dados
visuais antigos.

Uma evolução adequada é lifecycle:

```go
type Page interface {
    View() fyne.CanvasObject
    OnShow()
    OnHide()
}
```

Não introduza lifecycle antes de existir uma página que realmente precise.

### Parâmetros com `any`

`any` mantém o router pequeno, mas a validação ocorre em runtime.

Com muitas rotas, opções futuras:

- Helpers genéricos para factories tipadas.
- Métodos semânticos em um navigator de aplicação.
- Uma facade como `OpenCompany(companyID int)`.

Exemplo de facade:

```go
type AppNavigator interface {
    OpenHome()
    OpenCompany(companyID int)
    OpenNewCalculation()
    Back() bool
}
```

Trade-off:

- Mais segurança e semântica.
- Mais métodos e código para cada nova rota.

Para o tamanho atual, `RouteID + params` é proporcional.

### Router não é thread-safe

O histórico e o mapa não usam mutex.

As operações de navegação devem ocorrer na UI thread do Fyne.

Se uma goroutine precisar navegar:

```go
fyne.Do(func() {
    navigator.Push(route, params)
})
```

`fyne.Do` agenda a função na thread/event loop apropriado do Fyne.

### Feedback visual de erro

Hoje a homepage usa `fyne.LogError` quando a navegação falha.

Em uma aplicação mais madura, pode ser criado:

- Dialog de erro.
- Página de erro.
- Serviço central de notificações visuais.

O router deve continuar retornando erros; a camada visual decide como
apresentá-los.

### Páginas ainda iniciais

As páginas de empresa, cadastro e cálculo são destinos estruturais.

Ainda faltam:

- Carregamento dos detalhes da empresa.
- Formulário de cadastro.
- Validação.
- Loading state.
- Error state detalhado.
- Success feedback.
- Fluxo real de cálculo.

---

## 23. Glossário

### Alpha

Canal que representa transparência de uma cor.

- `0`: transparente.
- `255`: completamente opaco.

### `any`

Alias do Go para `interface{}`. Pode guardar valores de qualquer tipo.

### Back

Operação que remove a entrada atual do histórico e restaura a anterior.

### Callback

Função passada para ser chamada quando um evento ocorrer.

Exemplo:

```go
button := widget.NewButton("Abrir", onOpen)
```

### CanvasObject

Interface base dos objetos visuais que podem ser desenhados pelo Fyne.

Widgets e containers implementam `fyne.CanvasObject`.

### Closure

Função que captura variáveis do escopo externo.

Exemplo:

```go
companyID := 10
callback := func() {
    openCompany(companyID)
}
```

### Composition root

Ponto do programa que instancia e conecta dependências concretas.

Neste projeto:

```text
cmd/main
```

### Factory

Função responsável por construir um objeto.

Neste router, uma factory constrói uma página.

### Fade

Transição gradual de visibilidade.

### Fallback

Destino seguro usado quando a rota solicitada não pode ser aberta.

### History stack

Pilha que registra a ordem de navegação.

O último item é a página atual.

### Outlet

Container estável em que a página atual é apresentada.

### Params

Dados necessários para construir uma rota.

### Push

Adiciona uma nova entrada ao histórico.

### Replace

Substitui a entrada atual sem criar um novo nível para `Back`.

### RouteID

Identificador semântico de uma rota.

### Type assertion

Operação que extrai um tipo concreto de uma interface:

```go
value, ok := interfaceValue.(ConcreteType)
```

### UI thread

Contexto em que o framework processa eventos e atualizações visuais.

Operações visuais disparadas por goroutines devem voltar para esse contexto com
`fyne.Do`.

---

## 24. Ordem recomendada de leitura do código

Para compreender sem pular dependências:

- [x] `internal/ui/navigation/routes.go` 
- [x]  `internal/ui/navigation/router.go` 
- [x]  `internal/ui/components/animated_content.go`
- [x]  `cmd/main/ui_routes.go`
- [x]  `cmd/main/main.go`
- [x]  `internal/ui/application/app.go`
- [x]  `internal/ui/pages/homepage.go`
- [x]  `internal/ui/components/homepage_sidebar.go`
- [x]  `internal/ui/components/companies_list.go`
- [x]  `internal/ui/pages/company.go`
- [x]  `internal/ui/pages/calculation.go`
- [x]  `internal/ui/navigation/router_test.go`

Perguntas para responder durante a leitura:

1. Onde uma rota é identificada?
2. Onde uma rota é registrada?
3. Quem constrói a página?
4. Quem guarda a página anterior?
5. Quem executa a animação?
6. Quem controla a sidebar?
7. Como o ID da empresa chega à página?
8. Por que a página não acessa a janela?
9. Como uma rota inválida volta para Home?
10. Como os testes verificam uma animação interrompida?

---

## 25. Checklist para novas páginas

Antes de considerar uma nova rota pronta:

- [ ] Existe um `RouteID` estável?
- [ ] Os parâmetros possuem struct própria quando necessário?
- [ ] Somente IDs e dados mínimos são transportados?
- [ ] A página retorna `fyne.CanvasObject`?
- [ ] A página não chama `Window.SetContent`?
- [ ] A página depende de `Navigator`, não de `*Router`?
- [ ] Use cases entram por dependência explícita?
- [ ] A factory valida o tipo de `params`?
- [ ] A factory valida IDs obrigatórios?
- [ ] A rota foi registrada no provider?
- [ ] O chamador trata o erro de `Push` ou `Replace`?
- [ ] O botão voltar usa `Navigator.Back`?
- [ ] Estado local ficou na página, e não no router?
- [ ] Atualização de dados ficou fora do router?
- [ ] Há teste de navegação ou factory?
- [ ] `go vet` passou?
- [ ] Os testes focados passaram?
- [ ] O build passou?

Seguindo essa lista, o router permanece pequeno e as páginas continuam
independentes da janela e da implementação concreta da navegação.
