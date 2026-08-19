# Edição inline e guarda de navegação no workspace

**Status:** Accepted  
**Data:** 2026-08-17

**Superseded by:** [[0003-dirty-state-e-reload-da-sidebar]]

## Contexto

As telas de empresa e funcionário apresentavam campos como labels e não
possuíam fluxo integrado de edição, cancelamento ou persistência. O workspace
usa um router interno persistente, portanto uma navegação poderia desmontar a
tela atual e descartar valores temporários sem confirmação.

## Decisão

As telas usarão uma sessão de edição em nível de composição. O clique em um
campo de negócio ativa todos os campos editáveis da tela, habilita `GRAVAR` e
`CANCELAR` e registra os valores originais nos componentes inline.

`GRAVAR` valida os campos, chama o use case de update injetado e somente encerra
a sessão após persistência bem-sucedida. `CANCELAR` restaura os snapshots sem
chamar persistência. Códigos de empresa e funcionário permanecem somente
leitura; o funcionário expõe primeiro nome e sobrenome em campos separados.

O router interno recebe uma guarda de navegação compartilhada pela MainPage e
pelas telas. Enquanto a sessão estiver ativa, `Push`, `Reset`, `Replace` e
`Back` apresentam um diálogo. Continuar descarta a sessão e executa a operação
pendente; voltar à edição mantém a tela e seus valores.

## Consequências

- Páginas continuam dependendo da interface `Navigator`, enquanto a
  configuração da guarda permanece no router concreto da composição.
- A validação permanece próxima ao campo e não usa o diálogo de navegação.
- O router não conhece entidades nem regras de persistência; a página monta os
  comandos e usa as interfaces de update injetadas.
- A construção visual das páginas ocorre na UI thread após o carregamento dos
  dados, evitando criar widgets Fyne dentro da goroutine de consulta.
