# Project instructions

## Documentation

Este projeto utiliza Markdown compatível com Obsidian como documentação persistente.

Antes de propor alterações arquiteturais relevantes, leia:

* `docs/LaborCalculator4Companies - UI/CONTEXT.md`
* `docs/LaborCalculator4Companies - UI/Histórico de desenvolvimento.md`
* documentação relacionada à área alterada;
* ADRs existentes em `docs/LaborCalculator4Companies - UI/ADR/`.

Os arquivos Markdown podem conter links no formato Obsidian `[[arquivo]]`. Preserve esse formato quando apropriado.

## Source of truth

O código implementado é a fonte de verdade sobre o comportamento atual do sistema.

Quando houver divergência entre código e documentação:

* não assuma silenciosamente que a documentação está correta;
* investigue a divergência antes de realizar alterações;
* se a documentação estiver desatualizada, atualize-a para refletir o código;
* se a divergência indicar possível erro de implementação ou violação de uma decisão arquitetural registrada, informe o usuário antes de alterar o comportamento ou a decisão;
* não altere uma decisão arquitetural vigente apenas para fazer a documentação coincidir com o código.

ADRs registram decisões e suas justificativas. Portanto, uma divergência entre código e um ADR `Accepted` deve ser tratada como algo a investigar, e não automaticamente como documentação desatualizada.

## Grill workflow

Sempre que o usuário solicitar `grill-with-docs`:

1. Leia `docs/LaborCalculator4Companies - UI/CONTEXT.md` e os ADRs relevantes existentes.
2. Investigue o codebase e a documentação relacionada antes de fazer perguntas.
3. Não pergunte ao usuário algo que possa ser determinado pelo código ou pela documentação.
4. Faça uma pergunta por vez para eliminar ambiguidades.
5. Trate ADRs com status `Accepted` como decisões arquiteturais vigentes.
6. Não reverta silenciosamente uma decisão registrada em ADR. Se uma nova decisão entrar em conflito com um ADR existente, apresente o conflito ao usuário.
7. Registre decisões consolidadas em `docs/LaborCalculator4Companies - UI/CONTEXT.md`.
8. Para novas decisões arquiteturais relevantes, crie um ADR em `docs/LaborCalculator4Companies - UI/ADR/`.
9. Ao final da sessão, adicione uma entrada curta em `docs/LaborCalculator4Companies - UI/Histórico de desenvolvimento.md`.

Não implemente a feature durante o grill, salvo solicitação explícita do usuário.

## Development log

Após concluir uma alteração relevante:

* atualize `docs/LaborCalculator4Companies - UI/Histórico de desenvolvimento.md`;
* registre apenas decisões, alterações relevantes e próximo passo;
* não registre detalhes triviais de implementação;
* preserve o estilo cronológico existente no documento.

O histórico deve servir como registro resumido da evolução do projeto, e não como substituto da documentação técnica, dos ADRs ou do `CONTEXT.md`.

## Architecture decisions

Uma decisão merece ADR quando:

* altera limites entre módulos;
* cria uma abstração relevante;
* muda fluxo de dados;
* muda persistência;
* adiciona dependência estrutural;
* altera contratos importantes;
* possui trade-offs relevantes.

Não crie ADR para pequenas refatorações locais.

Para novas decisões arquiteturais relevantes, crie um ADR.

Não reescreva retrospectivamente ADRs aceitos para representar uma nova decisão.

Quando uma decisão substituir outra:

1. crie um novo ADR;
2. preserve o ADR original;
3. marque no novo ADR a relação `Supersedes`;
4. marque no ADR antigo a relação `Superseded by`.

ADRs com status `Accepted` devem ser considerados vigentes até que sejam explicitamente substituídos.

## After implementation

Após implementar trabalho previamente definido ou planejado:

1. Revise `docs/LaborCalculator4Companies - UI/CONTEXT.md`.
2. Atualize o contexto para refletir o estado real do sistema após a implementação.
3. Remova hipóteses temporárias ou planos que não representam mais o estado atual.
4. Preserve conceitos, invariantes, terminologia e decisões que continuam válidos.
5. Verifique se os ADRs relacionados continuam representando as decisões adotadas.
6. Não altere retrospectivamente a justificativa de ADRs aceitos.
7. Caso a implementação tenha substituído uma decisão anterior, registre um novo ADR e relacione-o ao anterior usando `Supersedes` / `Superseded by`.
8. Atualize `docs/LaborCalculator4Companies - UI/Histórico de desenvolvimento.md` com:

   * o que foi implementado;
   * decisões relevantes;
   * problemas ou desvios importantes;
   * próximo passo.
9. Atualize qualquer documentação técnica afetada pela implementação.
10. Verifique se documentação, ADRs e comportamento implementado não possuem divergências não explicadas.

Ao final da tarefa, a documentação deve refletir o comportamento real do código, enquanto os ADRs devem preservar o histórico e a justificativa das decisões arquiteturais.
