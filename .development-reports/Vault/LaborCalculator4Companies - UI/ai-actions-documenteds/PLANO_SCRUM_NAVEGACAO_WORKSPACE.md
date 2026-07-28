# Plano Scrum — navegação persistente da MainPage

Este documento organiza em sprints a evolução da navegação da interface do
LaborCalculator4Companies.

O plano parte das decisões registradas no Passo 6 do histórico de
desenvolvimento:

- A sidebar deve permanecer montada durante o trabalho normal.
- As telas de empresa, funcionário, recibos e cálculos devem ocupar somente a
  área central da MainPage.
- Deve existir um router externo para telas excepcionais de janela inteira.
- Deve existir um router interno para o workspace central.
- Selecionar uma empresa inicia uma nova sessão de navegação no workspace.
- Trocar de empresa deve eliminar completamente o histórico da empresa
  anterior.
- Navegações internas da empresa devem usar histórico e permitir `Back`.
- Um futuro botão `Close All` deve encerrar a sessão atual e voltar ao acesso
  rápido.

## Sumário

1. [Visão do produto](#1-visão-do-produto)
2. [Modelo de navegação](#2-modelo-de-navegação)
3. [Regras invariantes de UX](#3-regras-invariantes-de-ux)
4. [Conceitos Scrum usados](#4-conceitos-scrum-usados)
5. [Definition of Ready](#5-definition-of-ready)
6. [Definition of Done](#6-definition-of-done)
7. [Product Backlog priorizado](#7-product-backlog-priorizado)
8. [Sprint 1 — operação Reset e segurança do histórico](#8-sprint-1--operação-reset-e-segurança-do-histórico)
9. [Sprint 2 — controller da sessão do workspace](#9-sprint-2--controller-da-sessão-do-workspace)
10. [Sprint 3 — MainPage persistente e QuickAccessPage](#10-sprint-3--mainpage-persistente-e-quickaccesspage)
11. [Sprint 4 — dois routers e providers separados](#11-sprint-4--dois-routers-e-providers-separados)
12. [Sprint 5 — fluxo Empresa → Funcionário](#12-sprint-5--fluxo-empresa--funcionário)
13. [Sprint 6 — fluxo Funcionário → Recibos → Recibo](#13-sprint-6--fluxo-funcionário--recibos--recibo)
14. [Sprint 7 — Close All e controles do workspace](#14-sprint-7--close-all-e-controles-do-workspace)
15. [Sprint 8 — hardening, lifecycle e documentação](#15-sprint-8--hardening-lifecycle-e-documentação)
16. [Cenários de aceitação ponta a ponta](#16-cenários-de-aceitação-ponta-a-ponta)
17. [Estratégia de testes](#17-estratégia-de-testes)
18. [Riscos e respostas](#18-riscos-e-respostas)
19. [Backlog futuro](#19-backlog-futuro)
20. [Cerimônias Scrum sugeridas](#20-cerimônias-scrum-sugeridas)
21. [Métricas úteis](#21-métricas-úteis)
22. [Ordem final das entregas](#22-ordem-final-das-entregas)

---

## 1. Visão do produto

### Product Goal

Entregar uma MainPage com experiência semelhante a um editor de código:

- A sidebar funciona como explorador permanente.
- A área central funciona como workspace.
- Selecionar uma empresa abre uma nova sessão de trabalho.
- Funcionários, recibos e cálculos são abertos dentro dessa sessão.
- O usuário pode retornar pelas telas da mesma empresa.
- O usuário nunca retorna acidentalmente para dados de outra empresa.
- A sessão pode ser encerrada de forma explícita por `Close All`.
- Telas excepcionais continuam podendo ocupar a janela inteira.

### Resultado visual pretendido

```text
MasterWindow
└── RootRouter
    ├── LoginPage                         futuro
    ├── MainPage
    │   └── HSplit
    │       ├── Sidebar persistente
    │       └── WorkspaceRouter
    │           ├── QuickAccessPage
    │           ├── CompanyPage
    │           ├── EmployeePage
    │           ├── ReceiptListPage
    │           ├── ReceiptPage
    │           └── CalculationPage
    └── ReceiptFullViewPage               futuro
```

### O que significa “persistente”

Durante a navegação normal:

- `MainPage` não é substituída.
- O `HSplit` não é recriado.
- A sidebar não é desmontada.
- A busca da sidebar pode continuar existindo.
- Apenas `WorkspaceRouter.View()` troca seu conteúdo.

O `RootRouter` troca a MainPage somente em fluxos de janela inteira.

---

## 2. Modelo de navegação

## 2.1 RootRouter

Escopo: janela inteira.

Rotas previstas:

```text
RouteMain
RouteLogin                    futuro
RouteReceiptFullView          futuro
RouteFatalError               opcional
```

O `RootRouter` não deve receber rotas comuns de:

- Empresa.
- Funcionário.
- Lista de recibos.
- Cadastro normal.
- Cálculo dentro da operação diária.

## 2.2 WorkspaceRouter

Escopo: painel central da MainPage.

Rotas previstas:

```text
RouteQuickAccess
RouteCompanyDetails
RouteCompanyCreate
RouteEmployeeDetails
RouteEmployeeCreate
RouteReceiptList
RouteReceiptDetails
RouteReceiptCreate
RouteCalculationCreate
```

## 2.3 Sessões do workspace

Existem dois tipos de sessão.

### Sessão sem empresa

Histórico:

```text
[QuickAccess]
```

### Sessão de empresa

Exemplo:

```text
[Company(10), Employee(25), ReceiptList(25), Receipt(98)]
```

O primeiro item é a raiz da sessão.

`Back` pode remover itens até chegar a:

```text
[Company(10)]
```

`Back` não deve sair da empresa e não deve chegar ao QuickAccess.

Para encerrar a sessão da empresa, o usuário usa:

```text
Close All
```

Resultado:

```text
[QuickAccess]
```

## 2.4 Troca de empresa

Estado inicial:

```text
[Company(10), Employee(25), Receipt(98)]
```

O usuário seleciona Empresa 20 na sidebar.

Resultado obrigatório:

```text
[Company(20)]
```

Resultado proibido:

```text
[Company(10), Employee(25), Receipt(98), Company(20)]
```

Também é proibido:

```text
[Company(10), Employee(25), Company(20)]
```

Em ambos os casos proibidos, `Back` poderia expor dados da empresa anterior.

---

## 3. Regras invariantes de UX

Invariante é uma regra que deve permanecer verdadeira em todos os estados
válidos do sistema.

### INV-01 — empresa ativa única

O workspace pode ter no máximo uma empresa ativa.

### INV-02 — histórico pertence à empresa ativa

Se existe uma empresa ativa, todas as entradas dependentes do histórico devem
pertencer a ela.

### INV-03 — troca de empresa é Reset

Selecionar qualquer nova empresa deve limpar o histórico antes de instalar a
página raiz da nova empresa.

A implementação deve oferecer uma operação atômica de `Reset`; não deve
depender de chamadas públicas separadas como:

```go
router.Clear()
router.Replace(...)
```

Se a segunda chamada falhar, o workspace ficaria vazio. `Reset` deve construir
o destino antes de confirmar a limpeza.

### INV-04 — Push somente dentro da sessão

Depois de abrir a empresa, as navegações dependentes usam `Push`.

Exemplo:

```text
Company -> Employee -> ReceiptList -> Receipt
```

### INV-05 — Back respeita a raiz da sessão

Depois de `Reset(Company)`, `Back()` deve retornar `false` na página da empresa.

### INV-06 — QuickAccess não fica atrás da empresa

Selecionar empresa não deve produzir:

```text
[QuickAccess, Company]
```

Deve produzir:

```text
[Company]
```

Isso impede que o botão voltar seja usado para encerrar implicitamente a
sessão.

### INV-07 — Close All é explícito

`Close All` deve produzir:

```text
[QuickAccess]
```

e remover:

- Empresa ativa.
- Funcionário atual.
- Recibo atual.
- Páginas preservadas no histórico.
- Estado visual associado à sessão encerrada.

### INV-08 — sidebar permanece montada

`Push`, `Back`, `Reset` e `Close All` no workspace não podem substituir a
MainPage no `RootRouter`.

### INV-09 — rotas fullscreen pertencem ao RootRouter

Uma página fullscreen pode esconder a MainPage, mas ao voltar deve restaurar a
mesma instância do workspace, quando esse for o comportamento desejado.

---

## 4. Conceitos Scrum usados

## 4.1 Sprint

Período fixo em que um incremento utilizável é produzido.

Sugestão para este projeto:

```text
1 sprint = 1 semana ou 5 dias de trabalho focado
```

Se o desenvolvimento ocorrer somente em horários livres, manter o objetivo da
sprint e adaptar a quantidade de itens, sem aumentar indefinidamente a sprint.

## 4.2 Sprint Goal

Objetivo único que explica por que a sprint existe.

Itens podem ser ajustados durante a sprint, mas o objetivo deve continuar
estável.

## 4.3 Product Backlog

Lista priorizada de tudo que pode ser necessário para o produto.

O backlog é vivo. Novas descobertas podem:

- Adicionar histórias.
- Dividir histórias grandes.
- Repriorizar trabalho futuro.

## 4.4 Sprint Backlog

Histórias e tarefas selecionadas para a sprint atual.

## 4.5 User Story

Formato:

```text
Como <tipo de usuário>,
quero <capacidade>,
para <benefício>.
```

## 4.6 Story Points

Estimativa relativa de complexidade, risco e volume.

Escala sugerida:

```text
1, 2, 3, 5, 8, 13
```

Story points não são horas.

Uma história de 8 pontos não significa oito horas; significa que ela parece
mais complexa ou incerta que uma história de 3 pontos.

## 4.7 Incremento

Resultado integrado e verificável ao final da sprint.

O incremento não precisa significar que toda a feature está pronta. Precisa
representar uma melhoria utilizável e coerente.

---

## 5. Definition of Ready

Uma história está pronta para entrar em sprint quando:

- [ ] O comportamento esperado está descrito.
- [ ] Os parâmetros da rota são conhecidos.
- [ ] O escopo é RootRouter ou WorkspaceRouter.
- [ ] Dependências necessárias foram identificadas.
- [ ] Critérios de aceitação são verificáveis.
- [ ] Não existe decisão arquitetural crítica pendente.
- [ ] A história cabe em uma sprint.
- [ ] Dependências de histórias anteriores estão concluídas.

Se uma história não cumprir esses itens, ela deve passar por refinement antes
da Sprint Planning.

---

## 6. Definition of Done

Uma história de navegação é considerada concluída quando:

- [ ] O código está formatado com `gofmt`.
- [ ] A solução compila.
- [ ] `go vet` dos pacotes alterados passa.
- [ ] Testes focados passam.
- [ ] Não existe acesso direto a repository a partir da UI.
- [ ] Dependências são injetadas explicitamente.
- [ ] A página não chama `MasterWindow.SetContent`.
- [ ] A rota está registrada no router correto.
- [ ] Parâmetros da rota são validados.
- [ ] Erros de navegação são tratados.
- [ ] A sidebar continua montada nas rotas internas.
- [ ] `Back` respeita a raiz da sessão.
- [ ] Troca de empresa não preserva histórico anterior.
- [ ] A documentação técnica foi atualizada.
- [ ] O cenário foi validado manualmente na interface.
- [ ] Não foram introduzidas dependências ou abstrações sem necessidade.

Comandos mínimos:

```bash
gofmt -w $(rg --files cmd/main internal/ui -g '*.go')
go vet ./internal/ui/... ./cmd/main
go test ./internal/ui/... ./cmd/main -count=1
go build -o /tmp/labor-calculator-main-build ./cmd/main
```

---

## 7. Product Backlog priorizado

| ID | Item | Prioridade | Pontos | Sprint prevista |
|---|---|---:|---:|---:|
| PB-01 | Adicionar operação atômica `Reset` ao router | P0 | 5 | 1 |
| PB-02 | Testar limpeza integral do histórico | P0 | 3 | 1 |
| PB-03 | Definir política de falha durante `Reset` | P0 | 3 | 1 |
| PB-04 | Criar controller/facade do workspace | P0 | 8 | 2 |
| PB-05 | Controlar empresa ativa no workspace | P0 | 5 | 2 |
| PB-06 | Extrair `QuickAccessPage` da homepage | P0 | 5 | 3 |
| PB-07 | Transformar Homepage em MainPage persistente | P0 | 8 | 3 |
| PB-08 | Montar `WorkspaceRouter.View()` no centro | P0 | 5 | 3 |
| PB-09 | Separar RootRouter e WorkspaceRouter | P0 | 8 | 4 |
| PB-10 | Separar providers de rotas | P0 | 5 | 4 |
| PB-11 | Registrar fallbacks independentes | P0 | 3 | 4 |
| PB-12 | Seleção de empresa executar Reset | P0 | 5 | 5 |
| PB-13 | Abrir funcionário usando Push | P0 | 5 | 5 |
| PB-14 | Voltar até a raiz da empresa | P0 | 3 | 5 |
| PB-15 | Abrir lista de recibos por funcionário | P0 | 5 | 6 |
| PB-16 | Abrir recibo usando Push | P0 | 5 | 6 |
| PB-17 | Validar pertencimento ao contexto ativo | P0 | 8 | 6 |
| PB-18 | Implementar `Close All` | P1 | 5 | 7 |
| PB-19 | Adicionar estado visual do workspace | P1 | 3 | 7 |
| PB-20 | Preparar confirmação de conteúdo não salvo | P2 | 8 | Futuro |
| PB-21 | Adicionar lifecycle `OnShow`/`OnHide` | P2 | 8 | 8 |
| PB-22 | Criar teste ponta a ponta da navegação | P1 | 8 | 8 |
| PB-23 | Atualizar documentação arquitetural | P1 | 3 | Todas |
| PB-24 | Criar rota fullscreen de recibo | P2 | 8 | Futuro |
| PB-25 | Adicionar tabs de documentos | P2 | 13 | Futuro |

Prioridades:

- `P0`: necessário para o fluxo correto e seguro.
- `P1`: importante para experiência e manutenção.
- `P2`: evolução posterior; não bloqueia o MVP do workspace.

---

## 8. Sprint 1 — operação Reset e segurança do histórico

### Sprint Goal

Dar ao router uma operação segura capaz de iniciar uma nova sessão sem manter
qualquer entrada anterior.

### User Story US-01

```text
Como usuário,
quero que uma nova sessão substitua integralmente a anterior,
para nunca voltar acidentalmente a um contexto que já foi encerrado.
```

### Escopo

Adicionar ao router:

```go
Reset(route RouteID, params any) error
```

Semântica:

1. Construir a nova entrada.
2. Se a construção falhar, executar a política de erro definida.
3. Somente após sucesso, substituir todo o histórico.
4. Exibir a nova view.

Resultado:

```go
r.history = []historyEntry{entry}
```

### Por que não usar `Replace`

`Replace` substitui apenas a última entrada.

Estado:

```text
[CompanyA, EmployeeA, ReceiptA]
```

Executar:

```text
Replace(CompanyB)
```

produziria:

```text
[CompanyA, EmployeeA, CompanyB]
```

`Back` ainda retornaria para EmployeeA.

`Reset` deve produzir:

```text
[CompanyB]
```

### Tarefas

- [ ] Especificar `Reset` na documentação do router.
- [ ] Decidir se `Reset` entra em `Navigator` ou em interface especializada.
- [ ] Implementar construção da entrada antes da mutação.
- [ ] Limpar o histórico de forma atômica.
- [ ] Reutilizar `AnimatedContent`.
- [ ] Manter a política `latest wins`.
- [ ] Definir comportamento quando a rota de Reset falha.
- [ ] Adicionar testes.

### Decisão recomendada de interface

Manter o router genérico com:

```go
func (r *Router) Reset(route RouteID, params any) error
```

E introduzir posteriormente um contrato de workspace:

```go
type WorkspaceNavigator interface {
    SelectCompany(companyID int) error
    Push(route RouteID, params any) error
    Back() bool
    CloseAll() error
}
```

Assim, páginas comuns não precisam conhecer `Reset` diretamente.

### Critérios de aceitação

#### CA-01

```text
Dado history = [A, B, C]
Quando Reset(D) é executado
Então history = [D]
```

#### CA-02

```text
Dado history = [A, B]
Quando Reset(D) termina
Então Back() retorna false
```

#### CA-03

```text
Dado history = [A, B]
Quando a factory de D falha
Então nenhuma entrada parcial é adicionada
```

#### CA-04

A view renderizada deve corresponder à rota atual depois de uma animação
interrompida.

### Testes

- `TestRouterResetClearsHistory`
- `TestRouterBackReturnsFalseAfterReset`
- `TestRouterResetPreservesStateWhenFactoryFails`, conforme política escolhida.
- `TestRouterLatestResetWinsDuringTransition`

### Incremento

Router capaz de iniciar uma nova pilha de navegação com segurança.

### Risco principal

O fallback atual também modifica o histórico. A política de falha do Reset deve
ser explícita:

Opção recomendada:

- Construir a rota solicitada.
- Se falhar, construir o fallback.
- Se o fallback funcionar, substituir o histórico por `[fallback]`.
- Retornar o erro original.

---

## 9. Sprint 2 — controller da sessão do workspace

### Sprint Goal

Impedir que páginas e componentes manipulem diretamente as regras de sessão da
empresa.

### Motivação

O router conhece:

- Rotas.
- Histórico.
- Factories.
- Views.

Ele não deve conhecer:

- Empresa ativa.
- Relação empresa → funcionário.
- Relação funcionário → recibo.
- Significado de `Close All`.

Essas regras pertencem a um controller/facade de navegação do workspace.

### Estrutura sugerida

```text
internal/ui/workspace/
├── controller.go
└── controller_test.go
```

Exemplo conceitual:

```go
type Controller struct {
    router          *navigation.Router
    activeCompanyID int
}
```

### User Story US-02

```text
Como usuário,
quero que todas as telas abertas pertençam à empresa selecionada,
para não misturar funcionários ou recibos entre empresas.
```

### Responsabilidades do controller

- Selecionar empresa.
- Guardar o ID da empresa ativa.
- Iniciar a sessão com `Reset`.
- Abrir rotas dependentes com `Push`.
- Encerrar a sessão com `CloseAll`.
- Expor `Back`.
- Informar se existe empresa ativa.
- Validar pré-condições de contexto.

### API conceitual

```go
type Navigator interface {
    SelectCompany(companyID int) error
    OpenEmployee(employeeID int) error
    OpenReceiptList(employeeID int) error
    OpenReceipt(receiptID int) error
    OpenNewReceipt(employeeID int) error
    OpenNewCalculation(employeeID int) error
    Back() bool
    CloseAll() error
    ActiveCompanyID() (int, bool)
}
```

Não é obrigatório implementar todos os métodos nesta sprint. A interface deve
ser mantida proporcional às páginas existentes.

### Seleção da empresa

```text
SelectCompany(10)
  ├── valida ID
  ├── Router.Reset(RouteCompanyDetails, CompanyDetailsParams{10})
  └── após sucesso, define activeCompanyID = 10
```

O ID ativo só deve mudar depois que o Reset for aceito.

### Seleção da mesma empresa

Política recomendada:

```text
clicar na empresa já ativa
  -> Reset para a raiz da mesma empresa
```

Isso funciona como “voltar ao arquivo raiz/projeto” e elimina páginas
dependentes abertas.

Se futuramente a UX indicar que o clique deve apenas manter a tela atual, essa
política pode ser revisada.

### Tarefas

- [ ] Criar pacote `workspace`.
- [ ] Definir interface pública pequena.
- [ ] Criar controller com router injetado.
- [ ] Implementar `SelectCompany`.
- [ ] Implementar `Back`.
- [ ] Implementar consulta da empresa ativa.
- [ ] Implementar limpeza do contexto.
- [ ] Criar fake router ou interface mínima para testes.
- [ ] Documentar invariantes.

### Critérios de aceitação

#### CA-01

Selecionar empresa válida chama `Reset`, não `Push`.

#### CA-02

Empresa ativa é atualizada somente depois de Reset bem-sucedido.

#### CA-03

Selecionar nova empresa substitui o ID ativo.

#### CA-04

Nenhum componente visual precisa acessar `Router.history`.

#### CA-05

O controller não consulta banco e não contém regra de negócio de domínio.

### Incremento

Uma facade testável que transforma ações semânticas em operações de router.

### Pattern aplicado

Facade/Controller:

- Esconde operações genéricas do router.
- Oferece métodos relacionados à UX.
- Centraliza invariantes do workspace.

---

## 10. Sprint 3 — MainPage persistente e QuickAccessPage

### Sprint Goal

Transformar a homepage em um shell permanente e transformar o acesso rápido em
uma página real do workspace.

### User Story US-03

```text
Como usuário,
quero que a sidebar permaneça visível enquanto troco o conteúdo central,
para manter orientação e acesso rápido às empresas.
```

### Mudança conceitual

Antes:

```text
HomePage
├── Sidebar
└── QuickAccess hardcoded
```

Depois:

```text
MainPage
├── Sidebar
└── WorkspaceView injetada

WorkspaceRouter
└── QuickAccessPage
```

### Tarefas

- [ ] Criar `pages/quick_access.go`.
- [ ] Mover atalhos para `NewQuickAccessPage`.
- [ ] Remover `buildQuickAccess` da MainPage.
- [ ] Renomear `homePage` para `mainPage`, se aprovado.
- [ ] Renomear `NewHomePage` para `NewMainPage`, se aprovado.
- [ ] Criar `MainPageDeps`.
- [ ] Injetar `WorkspaceView fyne.CanvasObject`.
- [ ] Colocar `WorkspaceView` no painel direito do `HSplit`.
- [ ] Preservar os três estados da sidebar.
- [ ] Preservar animação de largura da sidebar.
- [ ] Confirmar que a sidebar não é recriada ao trocar rota interna.
- [ ] Registrar `RouteQuickAccess`.

### Dependências sugeridas

```go
type MainPageDeps struct {
    Companies         components.CompanyFinder
    WorkspaceView     fyne.CanvasObject
    WorkspaceNavigator workspace.Navigator
}
```

O nome exato pode ser ajustado durante a sprint.

### HSplit versus Border

Manter o `HSplit` é recomendado porque ele oferece divisor redimensionável,
comportamento semelhante a uma IDE.

Uma estrutura futura pode ser:

```text
MainPage Border
├── Top: toolbar
├── Bottom: status bar
└── Center: HSplit(sidebar, workspace)
```

Não é necessário substituir `HSplit` nesta sprint.

### Critérios de aceitação

#### CA-01

O acesso rápido é criado por uma factory registrada.

#### CA-02

O centro da MainPage recebe `WorkspaceRouter.View()`.

#### CA-03

Trocar de QuickAccess para outra rota não altera o objeto da sidebar.

#### CA-04

Abrir e fechar a sidebar continua funcionando.

#### CA-05

Busca de empresas continua funcionando.

### Incremento

MainPage visualmente igual à atual, porém com outlet central independente.

### Risco principal

Recriar a MainPage dentro de uma factory a cada fallback pode reconstruir o
workspace. A composição deve decidir quando a instância do shell é preservada.

---

## 11. Sprint 4 — dois routers e providers separados

### Sprint Goal

Compor RootRouter e WorkspaceRouter sem misturar rotas ou dependências.

### User Story US-04

```text
Como mantenedor,
quero separar rotas fullscreen de rotas do workspace,
para evitar que uma navegação comum desmonte a MainPage.
```

### Estrutura pretendida

```text
RootRouter
└── MainPage
    └── WorkspaceRouter.View()
```

### Providers

Separar:

```text
RootRouteProvider
├── RouteMain
└── rotas fullscreen futuras

WorkspaceRouteProvider
├── RouteQuickAccess
├── RouteCompanyDetails
├── RouteCompanyCreate
├── RouteEmployeeDetails
├── RouteReceiptList
└── outras rotas internas
```

### Namespaces de rota

Mesmo que cada router tenha mapa próprio, usar nomes claros:

```go
RootRouteMain
RootRouteReceiptFullView

WorkspaceRouteQuickAccess
WorkspaceRouteCompanyDetails
WorkspaceRouteEmployeeDetails
```

Alternativa:

```text
root.main
root.receipt.full-view

workspace.quick-access
workspace.company.details
workspace.employee.details
```

### Ordem de composição sugerida

```text
1. Criar aplicação Fyne.
2. Criar RootRouter.
3. Criar WorkspaceRouter.
4. Criar WorkspaceController.
5. Criar WorkspaceRouteProvider.
6. Registrar rotas do workspace.
7. Reset do workspace para QuickAccess.
8. Criar MainPage usando WorkspaceRouter.View().
9. Criar RootRouteProvider usando MainPage.
10. Registrar rotas root.
11. Abrir RouteMain no RootRouter.
12. Executar application.Run(RootRouter.View()).
```

### Tarefas

- [ ] Renomear provider atual ou dividi-lo.
- [ ] Registrar rotas internas somente no WorkspaceRouter.
- [ ] Registrar MainPage somente no RootRouter.
- [ ] Definir fallback do RootRouter.
- [ ] Definir fallback do WorkspaceRouter.
- [ ] Injetar navigator correto em cada página.
- [ ] Remover registros antigos no router único.
- [ ] Atualizar bootstrap de `cmd/main`.
- [ ] Criar testes de escopo.

### Fallbacks

```text
RootRouter fallback      -> RootRouteMain
WorkspaceRouter fallback -> WorkspaceRouteQuickAccess
```

Quando autenticação existir:

```text
RootRouter fallback -> RootRouteLogin ou RootRouteMain
```

dependendo do estado autenticado.

### Critérios de aceitação

#### CA-01

CompanyPage não pode ser aberta pelo RootRouter.

#### CA-02

MainPage não precisa ser registrada no WorkspaceRouter.

#### CA-03

Fallback interno mantém a sidebar montada.

#### CA-04

Fallback externo continua capaz de substituir a janela inteira.

#### CA-05

Ao voltar de uma rota fullscreen, o estado interno da MainPage é preservado,
quando a mesma instância estiver no histórico root.

### Incremento

Dois escopos de navegação operando independentemente.

---

## 12. Sprint 5 — fluxo Empresa → Funcionário

### Sprint Goal

Entregar a primeira navegação vertical completa dentro de uma sessão de
empresa.

### User Story US-05

```text
Como usuário,
quero selecionar uma empresa e abrir seus funcionários,
para navegar nos dados sem perder o contexto da empresa.
```

### Fluxo

```text
QuickAccess
  │
  └── sidebar seleciona CompanyA
      │
      └── Reset -> [CompanyA]
          │
          └── seleciona EmployeeA
              │
              └── Push -> [CompanyA, EmployeeA]
```

### Seleção de nova empresa

```text
[CompanyA, EmployeeA]
  │
  └── sidebar seleciona CompanyB
      │
      └── Reset -> [CompanyB]
```

### Tarefas

- [ ] Conectar lista de empresas ao WorkspaceController.
- [ ] Remover chamada genérica `Push(Company)` da MainPage.
- [ ] Implementar `SelectCompany`.
- [ ] Criar/carregar CompanyPage dentro do workspace.
- [ ] Apresentar lista de funcionários vinculados à empresa.
- [ ] Garantir query obrigatória pelo ID da empresa.
- [ ] Criar callback `OnEmployeeSelected`.
- [ ] Implementar `OpenEmployee`.
- [ ] Registrar EmployeePage.
- [ ] Adicionar ação de voltar.
- [ ] Adicionar loading, empty e error states.
- [ ] Testar troca de empresa durante EmployeePage.

### Regra de dados

Conforme o Passo 5 do histórico:

- Busca de funcionários deve exigir o ID da empresa.
- O elo forte da composição individualiza a consulta.

A UI deve passar:

```text
CompanyID -> GetEmployees
```

e não realizar uma busca global seguida de filtro visual.

### Critérios de aceitação

#### CA-01

Selecionar CompanyA produz histórico:

```text
[CompanyA]
```

#### CA-02

Selecionar EmployeeA produz:

```text
[CompanyA, EmployeeA]
```

#### CA-03

Back em EmployeeA retorna CompanyA.

#### CA-04

Back em CompanyA retorna `false`.

#### CA-05

Selecionar CompanyB enquanto EmployeeA está aberto produz:

```text
[CompanyB]
```

#### CA-06

Nenhum dado de CompanyA permanece visível depois da conclusão do Reset.

### Incremento

Primeira sessão empresarial navegável e segura.

---

## 13. Sprint 6 — fluxo Funcionário → Recibos → Recibo

### Sprint Goal

Completar a árvore principal de navegação planejada no Excalidraw.

### User Story US-06

```text
Como usuário,
quero acessar os recibos de um funcionário e abrir um recibo,
para consultar e continuar o trabalho dentro da empresa selecionada.
```

### Fluxo

```text
[Company]
  -> Push(Employee)
  -> Push(ReceiptList)
  -> Push(Receipt)
```

Histórico:

```text
[Company(10), Employee(25), ReceiptList(25), Receipt(98)]
```

Back:

```text
Receipt(98)
  -> ReceiptList(25)
  -> Employee(25)
  -> Company(10)
```

### Tarefas

- [ ] Implementar/ajustar ReceiptListPage.
- [ ] Consultar recibos por EmployeeID obrigatório.
- [ ] Criar estado vazio com ação `Novo recibo`.
- [ ] Criar callback de seleção do recibo.
- [ ] Registrar ReceiptPage.
- [ ] Implementar `OpenReceiptList`.
- [ ] Implementar `OpenReceipt`.
- [ ] Implementar `OpenNewReceipt`.
- [ ] Validar IDs nas factories.
- [ ] Preservar a empresa ativa.
- [ ] Testar Back em todos os níveis.
- [ ] Testar troca de empresa a partir de ReceiptPage.

### Regra de dados

Conforme Passos 4 e 5:

- Recibo possui vínculo com funcionário.
- Busca de recibos usa EmployeeID como parâmetro obrigatório.
- Funcionário pertence à empresa ativa.

### Validação de contexto

No mínimo, a navegação deve exigir empresa ativa antes de abrir:

- Funcionário.
- Lista de recibos.
- Recibo.
- Novo recibo.

Validações fortes de pertencimento devem continuar no application/domain.

A UI não deve ser a única barreira de segurança ou consistência.

### Critérios de aceitação

#### CA-01

ReceiptList só abre dentro de uma sessão empresarial.

#### CA-02

A busca recebe EmployeeID.

#### CA-03

Back percorre cada nível na ordem inversa.

#### CA-04

Back para na CompanyPage.

#### CA-05

Selecionar outra empresa em qualquer nível elimina todo o histórico anterior.

#### CA-06

Erros de carregamento produzem estado visual compreensível.

### Incremento

Árvore Empresa → Funcionário → Recibos navegável dentro do workspace.

---

## 14. Sprint 7 — Close All e controles do workspace

### Sprint Goal

Dar ao usuário uma forma explícita de encerrar a sessão atual e retornar ao
estado neutro.

### User Story US-07

```text
Como usuário,
quero fechar todo o contexto aberto,
para voltar ao acesso rápido sem usar repetidamente o botão voltar.
```

### Semântica

Antes:

```text
[Company(10), Employee(25), Receipt(98)]
activeCompanyID = 10
```

Depois de `CloseAll`:

```text
[QuickAccess]
activeCompanyID = none
```

### Implementação recomendada

No WorkspaceController:

```go
func (c *Controller) CloseAll() error {
    err := c.router.Reset(navigation.RouteQuickAccess, nil)
    if err != nil {
        return err
    }

    c.activeCompanyID = 0
    return nil
}
```

O estado ativo só deve ser limpo depois de Reset bem-sucedido, conforme a
política escolhida.

### Localização visual

Opções:

1. Toolbar global da MainPage.
2. Barra superior do workspace.
3. Menu de contexto do workspace.

Recomendação:

```text
barra superior do workspace
```

porque `Close All` opera sobre o conteúdo central, não sobre a janela inteira.

### Estado do botão

No QuickAccess:

- Oculto; ou
- Desabilitado.

Durante sessão de empresa:

- Habilitado.

### Tarefas

- [ ] Implementar `CloseAll` no controller.
- [ ] Adicionar botão/ícone.
- [ ] Definir tooltip.
- [ ] Atualizar estado habilitado/desabilitado.
- [ ] Limpar empresa ativa.
- [ ] Limpar seleção visual da empresa, se aplicável.
- [ ] Resetar workspace para QuickAccess.
- [ ] Testar CloseAll em todos os níveis.
- [ ] Garantir que a sidebar permaneça montada.

### Critérios de aceitação

#### CA-01

CloseAll sempre produz:

```text
[QuickAccess]
```

#### CA-02

Back retorna `false` depois de CloseAll.

#### CA-03

Nenhuma empresa permanece ativa.

#### CA-04

Selecionar uma empresa depois de CloseAll inicia uma sessão limpa.

#### CA-05

A sidebar não é recriada.

### Conteúdo não salvo

Nesta sprint, implementar confirmação somente se já existirem formulários com
estado “dirty”.

Caso contrário, registrar a proteção como backlog futuro para não criar uma
abstração sem consumidor real.

### Incremento

Workspace com início e encerramento explícitos de sessão.

---

## 15. Sprint 8 — hardening, lifecycle e documentação

### Sprint Goal

Consolidar a arquitetura, eliminar inconsistências e preparar o workspace para
evolução.

### User Story US-08

```text
Como mantenedor,
quero testes e contratos claros para a navegação,
para adicionar novas páginas sem quebrar isolamento ou histórico.
```

### Tarefas

- [ ] Revisar nomes Root/Workspace.
- [ ] Revisar interfaces injetadas.
- [ ] Evitar dependência em router concreto nas páginas.
- [ ] Revisar fallback de cada escopo.
- [ ] Adicionar teste ponta a ponta de navegação.
- [ ] Adicionar teste de troca rápida entre empresas.
- [ ] Adicionar teste de CloseAll durante animação.
- [ ] Verificar latest wins nos dois routers.
- [ ] Avaliar lifecycle `OnShow`.
- [ ] Recarregar dados preservados quando necessário.
- [ ] Revisar crescimento do histórico.
- [ ] Atualizar `ARQUITETURA_NAVEGACAO_UI.md`.
- [ ] Atualizar histórico de desenvolvimento.
- [ ] Registrar exemplos de criação de rota root e workspace.
- [ ] Rodar format, vet, testes e build.

### Lifecycle opcional

Se houver necessidade concreta:

```go
type Page interface {
    View() fyne.CanvasObject
    OnShow()
    OnHide()
}
```

Uso:

- `OnShow`: recarregar dados ao voltar.
- `OnHide`: interromper carregamento ou persistir estado visual.

Não implementar se as páginas existentes ainda não precisarem.

### Critérios de aceitação

#### CA-01

Todos os cenários da seção 16 passam.

#### CA-02

Nenhuma rota interna desmonta a MainPage.

#### CA-03

Não há caminho de Back entre empresas.

#### CA-04

Documentação permite criar nova rota sem consultar o histórico do chat.

#### CA-05

Build do executável passa.

### Incremento

Arquitetura de workspace documentada, testada e pronta para novas features.

---

## 16. Cenários de aceitação ponta a ponta

## Cenário A — inicialização

```text
Dado que a aplicação foi aberta
Quando a MainPage é montada
Então a sidebar está visível
E o workspace mostra QuickAccess
E Back retorna false
```

Histórico:

```text
[QuickAccess]
```

## Cenário B — seleção de empresa

```text
Dado QuickAccess
Quando CompanyA é selecionada
Então CompanyPage(A) aparece no centro
E a sidebar permanece visível
E Back retorna false
```

Histórico:

```text
[CompanyA]
```

## Cenário C — navegação interna

```text
Dado CompanyA
Quando EmployeeA é selecionado
E ReceiptA é selecionado
Então Back retorna primeiro ao funcionário
E depois retorna à empresa
E não retorna ao QuickAccess
```

Histórico máximo:

```text
[CompanyA, EmployeeA, ReceiptListA, ReceiptA]
```

## Cenário D — troca de empresa na CompanyPage

```text
Dado CompanyA
Quando CompanyB é selecionada
Então o histórico contém somente CompanyB
E Back retorna false
```

## Cenário E — troca de empresa em tela profunda

```text
Dado CompanyA -> EmployeeA -> ReceiptA
Quando CompanyB é selecionada
Então CompanyB aparece no centro
E nenhum Back retorna a CompanyA
E nenhum widget de CompanyA permanece ativo
```

## Cenário F — Close All

```text
Dado CompanyA -> EmployeeA -> ReceiptA
Quando Close All é acionado
Então QuickAccess aparece no centro
E empresa ativa é removida
E Back retorna false
E sidebar permanece montada
```

## Cenário G — nova sessão após Close All

```text
Dado QuickAccess depois de Close All
Quando CompanyB é selecionada
Então o histórico contém somente CompanyB
```

## Cenário H — rota interna inválida

```text
Dado que uma rota interna não pode ser construída
Quando o fallback é acionado
Então QuickAccess aparece no centro
E a MainPage permanece montada
```

## Cenário I — fullscreen futuro

```text
Dado CompanyA -> ReceiptA no workspace
Quando ReceiptFullView é aberto pelo RootRouter
E o usuário volta
Então a mesma MainPage é restaurada
E ReceiptA continua ativo no workspace
```

---

## 17. Estratégia de testes

## 17.1 Testes do Router

Responsáveis por mecânica genérica:

- Register.
- Push.
- Replace.
- Reset.
- Back.
- Current.
- Fallback.
- Latest wins.

Não devem conhecer CompanyID ou ReceiptID.

## 17.2 Testes do WorkspaceController

Responsáveis por semântica da UX:

- Selecionar empresa usa Reset.
- Abrir funcionário usa Push.
- Empresa ativa muda corretamente.
- CloseAll volta ao QuickAccess.
- Ações dependentes exigem contexto.

## 17.3 Testes das factories

Responsáveis por:

- Type assertion de params.
- ID maior que zero.
- Injeção das dependências corretas.
- Retorno de erro quando inválido.

## 17.4 Testes dos componentes

Responsáveis por:

- Seleção em `widget.List`.
- Callback com ID correto.
- Limpeza da seleção.
- Empty state.
- Error state.
- Busca com espaços.

## 17.5 Testes de integração da UI

Responsáveis por:

- MainPage mantém o mesmo objeto da sidebar.
- Workspace troca somente o centro.
- CloseAll não desmonta shell.
- RootRouter e WorkspaceRouter não compartilham histórico.

## 17.6 Teste manual

Checklist:

- [ ] Abrir aplicação.
- [ ] Abrir/fechar sidebar.
- [ ] Pesquisar empresa.
- [ ] Selecionar empresa.
- [ ] Abrir funcionário.
- [ ] Abrir recibo.
- [ ] Voltar nível por nível.
- [ ] Confirmar parada na empresa.
- [ ] Abrir novamente funcionário.
- [ ] Trocar de empresa no nível mais profundo.
- [ ] Confirmar que Back não retorna à anterior.
- [ ] Executar Close All.
- [ ] Confirmar QuickAccess.
- [ ] Confirmar que sidebar não piscou nem foi recriada.

---

## 18. Riscos e respostas

## Risco R-01 — Replace usado no lugar de Reset

Impacto:

- Histórico de empresa anterior permanece.

Resposta:

- Controller deve ser o único responsável por `SelectCompany`.
- Teste obrigatório de troca profunda de empresa.

## Risco R-02 — páginas recebem router errado

Impacto:

- Uma página interna pode substituir a janela inteira.

Resposta:

- Interfaces `RootNavigator` e `WorkspaceNavigator`.
- Providers separados.
- Rotas com namespace explícito.

## Risco R-03 — MainPage recriada

Impacto:

- Sidebar e workspace perdem estado.

Resposta:

- Revisar lifecycle do RootRouter.
- Preservar a instância no histórico.
- Testar fullscreen + Back.

## Risco R-04 — dados antigos em views preservadas

Impacto:

- Back mostra informação desatualizada.

Resposta:

- Adicionar `OnShow` somente quando necessário.
- Recarregar listas após alterações.

## Risco R-05 — consulta bloqueia UI

Impacto:

- Janela congela ao carregar funcionários/recibos.

Resposta:

- Loading state.
- Goroutine para I/O.
- `fyne.Do` para atualização visual.

## Risco R-06 — animação antiga conclui depois do Reset

Impacto:

- Empresa errada aparece.

Resposta:

- Preservar `transitionID`.
- Testar Reset durante animação.

## Risco R-07 — Close All perde formulário não salvo

Impacto:

- Perda de dados digitados.

Resposta futura:

- Dirty state.
- Confirmação.
- Contrato `CanClose`.

## Risco R-08 — controller vira regra de negócio

Impacto:

- UI passa a substituir application/domain.

Resposta:

- Controller conhece apenas contexto de navegação.
- Validação de pertencimento real continua no use case/domain.

---

## 19. Backlog futuro

Itens não necessários para o primeiro incremento do workspace.

## 19.1 Dirty state

Representa páginas com alterações não salvas.

Possível contrato:

```go
type ClosablePage interface {
    HasUnsavedChanges() bool
}
```

Antes de Reset ou CloseAll:

```text
existem alterações?
  -> pedir confirmação
  -> cancelar ou continuar
```

## 19.2 Tabs

Para aproximar ainda mais de uma IDE:

```text
Workspace
├── Tab Company
├── Tab Employee
└── Tab Receipt *
```

O asterisco pode indicar alteração não salva.

Tabs exigem separar:

- Documentos abertos.
- Documento ativo.
- Histórico de navegação.
- Ordem visual.
- Fechamento.

Não misturar tabs prematuramente com a pilha atual.

## 19.3 Back e Forward

Hoje existe `Back`.

Forward exigiria uma segunda pilha ou cursor:

```text
history + currentIndex
```

## 19.4 Full-view de recibo

Rota do RootRouter.

Deve preservar o workspace ao retornar.

## 19.5 Login

Rota root inicial.

Após login:

```text
RootRouter.Reset(RouteMain)
```

Isso impede voltar ao login.

## 19.6 Persistência de sessão

No futuro, a aplicação poderia restaurar:

- Última empresa.
- Última página.
- Tamanho da sidebar.

Não implementar antes de o fluxo básico estar estável.

---

## 20. Cerimônias Scrum sugeridas

Como o projeto pode ter apenas um desenvolvedor, as cerimônias devem ser leves.

## Sprint Planning

Duração sugerida:

```text
30 a 60 minutos
```

Perguntas:

1. Qual é o Sprint Goal?
2. Quais histórias estão Ready?
3. Qual capacidade real existe nesta semana?
4. Quais riscos podem impedir o objetivo?

## Daily Scrum

Duração:

```text
5 minutos
```

Registrar no report:

1. O que avancei?
2. O que farei agora?
3. Existe bloqueio?
4. Descobri algo que muda o backlog?

O Daily não é relatório para chefe. É inspeção do plano.

## Backlog Refinement

Duração:

```text
20 a 40 minutos por semana
```

Atividades:

- Dividir histórias grandes.
- Esclarecer critérios.
- Reestimar.
- Repriorizar.
- Remover itens obsoletos.

## Sprint Review

Demonstrar o incremento funcionando.

Exemplo Sprint 5:

```text
Selecionar CompanyA
Abrir EmployeeA
Voltar
Trocar para CompanyB
Confirmar histórico limpo
```

## Sprint Retrospective

Perguntas:

1. O que funcionou?
2. O que atrasou?
3. Qual decisão gerou retrabalho?
4. O que será mudado na próxima sprint?

Escolher no máximo uma ou duas ações concretas de melhoria.

---

## 21. Métricas úteis

Evitar métricas de vaidade, como quantidade de linhas.

### Cenários de aceitação passando

Principal indicador funcional.

### Bugs de contexto

Quantidade de ocorrências em que:

- Empresa errada aparece.
- Back atravessa sessões.
- Dados de outra empresa permanecem.

Meta:

```text
zero
```

### Lead time de nova rota

Tempo entre:

```text
história Ready -> rota integrada e Done
```

### Retrabalho por dependência errada

Exemplo:

- Página recebeu RootRouter quando precisava WorkspaceRouter.

### Velocidade

Soma de story points concluídos por sprint.

Usar somente depois de algumas sprints. Não comparar velocidade com outras
pessoas ou projetos.

---

## 22. Ordem final das entregas

```text
Sprint 1
└── Router.Reset

Sprint 2
└── WorkspaceController e empresa ativa

Sprint 3
└── MainPage persistente + QuickAccessPage

Sprint 4
└── RootRouter + WorkspaceRouter + providers separados

Sprint 5
└── Company -> Employee

Sprint 6
└── Employee -> Receipts -> Receipt

Sprint 7
└── Close All

Sprint 8
└── Hardening, lifecycle e documentação
```

### Marco 1 — fundação

Sprints 1 e 2:

- Reset seguro.
- Regras de sessão centralizadas.

### Marco 2 — shell persistente

Sprints 3 e 4:

- Dois routers.
- MainPage nunca desmontada durante navegação comum.

### Marco 3 — fluxo de negócio navegável

Sprints 5 e 6:

- Empresa.
- Funcionário.
- Recibos.

### Marco 4 — encerramento e estabilidade

Sprints 7 e 8:

- Close All.
- Testes ponta a ponta.
- Documentação consolidada.

Ao final, o comportamento central deve ser:

```text
QuickAccess
  -> selecionar CompanyA
  -> navegar dentro de CompanyA com Push/Back
  -> selecionar CompanyB
  -> histórico de CompanyA é eliminado
  -> navegar dentro de CompanyB
  -> Close All
  -> QuickAccess
```

Esse fluxo é a referência funcional para todas as decisões futuras do
workspace.
