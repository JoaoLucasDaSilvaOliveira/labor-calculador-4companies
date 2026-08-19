# LaborCalculator4Companies — UI

Vocabulário do contexto de interface da calculadora trabalhista, responsável
por apresentar e editar empresas, funcionários e seus recibos.

## Telas

**Tela 1.1 — informações da empresa**:
A composição de interface que apresenta os dados de uma empresa e a sua lista
de funcionários. Preserva responsabilidades e fluxo próprios, mesmo usando
componentes visuais compartilhados.
_Evitar_: tela genérica de entidade

**Tela 1.2 — informações do funcionário**:
A composição de interface que apresenta os dados de um funcionário e a sua
lista de recibos. Preserva responsabilidades e fluxo próprios, mesmo usando
componentes visuais compartilhados.
_Evitar_: tela genérica de entidade

**Componente visual reutilizável**:
Um bloco pequeno e semântico de apresentação ou interação, usado por mais de
uma tela sem assumir o fluxo completo de empresa ou funcionário.
_Evitar_: componente de tela configurável

## Navegação

**Router externo**:
O navegador entre telas de agregação da aplicação, como a mainpage, login ou
visualizações em tela cheia. Ele não navega entre os conteúdos de trabalho da
mainpage.
_Evitar_: router da homepage

**Router interno**:
O navegador persistente do conteúdo central da mainpage. Ele contém as telas
de trabalho, incluindo as telas 1.1 de empresa e 1.2 de funcionário, sem
desmontar a sidebar.
_Evitar_: router global

## Estados de apresentação

**Estado vazio**:
Uma apresentação normal de uma lista sem entidades associadas, como uma
empresa sem funcionários ou um funcionário sem recibos. Não representa uma
falha do sistema.
_Evitar_: erro de lista vazia

**Erro recuperável de tela**:
Uma falha real ao carregar ou consultar dados para uma tela do workspace. A
tela explica a falha e permite retornar pelo histórico do router interno.
_Evitar_: estado vazio

**Validação de preenchimento**:
O feedback sobre um valor inválido ou obrigatório em um formulário. É exibido
no contexto do campo correspondente e não altera a navegação atual.
_Evitar_: erro recuperável de tela, pop-up de erro

## Edição e persistência

**Sessão de edição**:
O estado de alteração pendente da tela atual. Clicar em um campo inicia o modo
`Entry` para todos os campos de negócio editáveis, mas a sessão só é ativada
quando algum valor diverge do snapshot original.

**Campo editável inline**:
Um campo que alterna entre label de visualização e `Entry`, preservando um
snapshot do valor original para cancelamento. Códigos identificadores não
fazem parte dos campos editáveis de empresa ou funcionário.

**Alteração não salva**:
Um valor de campo que diverge do snapshot original. `GRAVAR` e `CANCELAR`
ficam habilitados somente nesse estado; retornar ao valor original encerra a
alteração pendente e não bloqueia a navegação.

**Guarda de navegação**:
O mecanismo do router interno que apresenta o diálogo de alterações não
salvas. `Voltar à edição` preserva os valores; `Continuar` executa o
cancelamento da sessão antes de navegar.

**Reload da lista de empresas**:
Após a persistência de uma empresa, o componente de lista atualmente montado
na sidebar repete a consulta em uma goroutine e aplica o resultado na UI
thread, mantendo o nome atualizado sem reconstruir a MainPage.
