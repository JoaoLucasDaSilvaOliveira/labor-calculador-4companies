# Dirty state da edição inline e reload da sidebar

**Status:** Accepted  
**Data:** 2026-08-19  
**Supersedes:** [[0002-edicao-inline-e-guarda-de-navegacao]]

## Contexto

A primeira implementação ativava `GRAVAR`, `CANCELAR` e a guarda de
navegação apenas com o clique no campo. Além disso, o label da sidebar era
carregado somente na construção inicial do componente e não refletia a
alteração de nome persistida na Tela 1.1.

## Decisão

O modo visual `Entry` e o estado de alteração pendente serão separados. O
clique em qualquer campo editável troca todos os campos de negócio da tela
para `Entry`, mas a sessão de edição só começa quando algum valor diverge do
snapshot original. Se o usuário retornar ao valor original, os botões
`GRAVAR` e `CANCELAR` são desabilitados novamente e a navegação não é
bloqueada.

Cada `InlineEditableField` mantém uma apresentação exclusiva para leitura e
outra para edição. A troca substitui a apresentação montada, evitando que
elementos ocultos do estado anterior permaneçam na composição visual. O
estado de edição usa fundo azul claro e borda da cor primária para tornar a
mudança inequívoca.

Após atualizar uma empresa, a Tela 1.1 solicita o reload do componente de
lista atualmente montado na sidebar. A lista executa a consulta em goroutine,
descarta respostas antigas e aplica o resultado na UI thread. Um handle
estável é compartilhado pela MainPage e pelas páginas de empresa para manter
essa integração sem acoplar a página à estrutura da sidebar.

## Consequências

- A guarda protege somente alterações reais, não uma simples entrada no modo
  de edição.
- O componente de lista permanece responsável por consulta, loading, vazio,
  erro e atualização do seu próprio conteúdo.
- A MainPage não precisa ser reconstruída para refletir o novo nome da
  empresa.
- A decisão anterior sobre sessão baseada apenas no modo visual fica
  preservada no ADR 0002 como histórico e é substituída por esta decisão.
