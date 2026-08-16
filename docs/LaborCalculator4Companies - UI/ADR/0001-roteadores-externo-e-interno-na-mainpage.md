# Roteadores externo e interno na mainpage

A aplicação usará um router externo para alternar entre telas de agregação e um router interno, persistente no centro da mainpage, para o fluxo de trabalho de acesso rápido, empresa e funcionário. Esta separação preserva a sidebar durante a navegação contextual e evita que telas globais compartilhem indevidamente o histórico do workspace.

## Consequências

As rotas de empresa e funcionário pertencem ao router interno; a rota externa da mainpage permanece como a unidade que contém esse workspace. O acesso rápido é a rota inicial e o fallback do router interno, cujo histórico segue o fluxo acesso rápido → empresa → funcionário.

As rotas internas carregam somente o identificador da entidade. Cada tela consulta a fonte de verdade ao ser construída e deve renová-la quando seu fluxo permitir, antes de apresentar ou executar ações dependentes desses dados.

O workspace reutiliza o contrato genérico `Navigator` e  a implementação pronta `Router`; não será criada, nesta etapa, uma abstração de navegação com métodos específicos de empresa ou funcionário. A composição do workspace registra as rotas e injeta dependências, enquanto as páginas solicitam navegação por rotas tipadas.

Uma falha de dados de uma tela interna é apresentada como um estado visual recuperável da própria rota, com retorno pelo histórico, e não como erro da factory do router. Estados vazios e validações de formulário continuam sendo apresentados no contexto da lista ou do campo afetado.
