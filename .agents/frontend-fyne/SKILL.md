# Skill: Fyne UI/UX Moderna em Go

Sempre que criar telas com Fyne:

## Objetivo
Criar interfaces profissionais, limpas, responsivas e organizadas, evitando aparência padrão “crua”.

## Regras obrigatórias
- Usar `container.NewBorder`, `container.NewVBox`, `container.NewHBox`, `container.NewGridWithColumns`, `container.NewStack` quando fizer sentido.
- Separar layout em funções pequenas:
  - `buildHeader()`
  - `buildSidebar()`
  - `buildContent()`
  - `buildCard()`
  - `buildForm()`
- Nunca jogar todos os widgets direto em um único `VBox`.
- Usar espaçamento, padding e agrupamento visual.
- Criar cards com borda, fundo ou separação clara.
- Usar títulos, subtítulos e hierarquia visual.
- Usar botões primários e secundários com consistência.
- Usar ícones quando fizer sentido.
- Usar validações em formulários.
- A interface precisa parecer um sistema real, não um exemplo de documentação.

## Estilo visual
- Layout com sidebar à esquerda quando for dashboard/sistema administrativo.
- Header superior com título da tela e ações principais.
- Conteúdo central em cards.
- Formulários alinhados e com campos agrupados.
- Tabelas/listas com busca, filtros e ações.
- Evitar telas poluídas.

## Antes de codar
Sempre descreva rapidamente:
1. Estrutura da tela
2. Componentes usados
3. Fluxo do usuário

## Código
- Código Go idiomático.
- Separar UI de lógica.
- Evitar funções gigantes.
- Usar tema customizado quando necessário.
- Comentar apenas o necessário.