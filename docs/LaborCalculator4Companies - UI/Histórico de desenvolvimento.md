# Passo 1:
Começar a colocar o esboço de homepage definida ([[excalidraw-sketch.excalidraw]]) em prática. Primeiro resultado: [[homepage.canvas]]. 
**Próximo passo:** criar ícones personalizados que façam sentido com o contexto e ações.

# Passo 2:
Primeiros elementos gráficos criados no figma: 
- Abrir e fechar sidebar.
- Lupa.
Esses elementos compõe o que chamei de sidebar do homepage.
Nesse passo, já implementei o comportamento para o primeiro elemento, agora falta a lupa. Resultado em:  [[homepage.canvas]].

# Passo 3:
Elementos gráficos da sidebar completamente finalizados. Todos tem comportamentos próprios. Pedi ao codex que criasse uma transição/animação, agora eu tenho um "mini" router. Posso trocar o conteúdo da página como um router e fico com uma animaçãozinha de fade de brinde. Resultado em [[homepage.canvas]].
**Próximo passo:** Alterar a lista de empresa de label para botões que mostrem os dados da empresa: nome, CNPJ, funcionários e etc. Essa alteração deve ser refletida, claro, no espaço central do border da homepage.
==*Sugestão para o próximo passo: montar a visão sugerida de forma, mesmo que primitiva, no excalidraw e depois reproduzir no código.==*

# PASSO 4:
No desenho da tela surgiu algo que eu não tinha pensado: criei a entidade de calculadora que fornece formas de chegar à resultados matemáticos, mas não criei a entidade do recibo, que guarda esses resultados. Vou ter que interromper o desenvolvimento da ui para ajustar isso.

OBS PRECISO LEMBRAR: RECIBO PRECISA POSSUIR O ID DO EMPREGADO PARA PODER FAZER UM GET EM TODOS OS RECIBOS. 

# Passo 5:
Implementas as novas features no backend para comportar a lógica dos recibos. Também individualizei as buscas de objetos que possuem vínculo forte pelo elo mais forte da composição, exemplo: a busca de funcionários (get) usa o id da empresa (elo forte) como parâmetro obrigatório, individualizando a busca de forma otimizada. O mesmo ocorre com o recibo em relação ao funcionário. 
**Próximo passo:** Limpar código sujo e implementar a estrutura de router.

# Passo 6:
Usei o codex-5.6-sol no alto pra isso. Aparentemente ele fez um bom trabalho, a documentação do que foi feita está em [[ARQUITETURA_NAVEGACAO_UI]]. Li todo o documento e entendi uns 60%, falta ver no código como de fato foi feito e criar páginas novas, mas tem um porém, o agente não entendeu bem o desenho da tela. Foi arquitetado para o router funcionar multi-telas, e isso até é bom, digamos que futuramente eu queira uma tela de login ou uma tela de full-view dos recibos, porém para o que estou desenhando nunca saio da homepage (vale pensar em mudar para mainpage o nome): a sidebar (left do border) nunca desmonta de verdade, só se transforma; e o conteúdo que hoje são outras páginas, deve ser mostrado dentro do centro (center do border). A IA não entendeu, mas o conteúdo inicial (Acesso rápido, cadastrar empresa e novo recibo) é um tipo de tela. Acho que seria ideal uma nova instância de router morar dentro do centro, sendo controlada pelo que acontece na homepage, que por sua vez, mora em uma outra instância de router, contemplando mudanças futuras como comentei.
**Próximo passo:** Implementar esse roteamento dentro do center do border da "homepage".

# Passo 7:
Houve fiz algumas grandes mudanças e não documentei aqui, então talvez eu esqueça alguma coisa. Usei o codex-5.6-terra pra fazer a refatoração sugerida no passo anterior, agora basicamente a aplicação tem dois routers: externo e interno. O router externo abriga, atualmente, apenas a mainpage (antiga homepage) e possui planos para abrigar outras telas de agregação como tela de login, full view de algum documento ou algo nesse sentido. Já o router interno mora dentro da mainpage e serve como navegador para os ítens clicáveis da mainpage; por exemplo: ao clicar em uma empresa na lista de empresas é - após algumas chamadas e processamentos que não convém descrever no momento - adicionado ao router a visualização das informações da empresa: cod, nome, CNPJ, funcionários e etc. 
**Próximo passo:** desenvolver as telas 1.1 e 1.2, descritas em  [[excalidraw-sketch.excalidraw]].

# Passo 8 :
Desenvolvida as telas 1.1 e 1.2. Porém ainda não as integrei oficialmente na aplicação, apenas criei o design e mock básicos de teste. Uma coisa interessante é que tanto a tela 1.1 quanto a 1.2 tem elementos muito parecidos. Nesse sentido vou componentizar o que for semelhante ou até mesmo reutilizado em ambas as telas.

**Próximo passo:** Componentizar os elementos repetidos/semelhantes das telas 1.1 e 1.2.

# Passo 9:
Consolidado o desenho para integrar as telas 1.1 (empresa) e 1.2 (funcionário): a mainpage terá um router interno persistente no conteúdo central, enquanto o router externo continuará responsável apenas pelas telas de agregação. O workspace inicia e faz fallback em Acesso rápido, preserva o histórico `acesso rápido → empresa → funcionário` e recebe rotas nomeadas pelo padrão `<local>.<ação>`. A decisão está registrada em [[0001-roteadores-externo-e-interno-na-mainpage]].

As telas manterão composições próprias e usarão somente componentes visuais pequenos e reutilizáveis. Nesta integração, dados serão somente leitura; navegação interna, filtros por vínculo e estados de lista vazia serão funcionais. Falhas reais de consulta terão estado recuperável com retorno pelo histórico; validações de formulário permanecem próximas aos campos.
**Próximo passo:** implementar a integração do workspace e das telas 1.1/1.2; em seguida, retomar o fluxo de edição inline (`GRAVAR`/`CANCELAR`/validações) e a criação de novo recibo.

# Passo 10:
Integradas as telas 1.1 e 1.2 ao fluxo real da aplicação. A `MainPage` agora
preserva a sidebar e hospeda o router interno no conteúdo central, enquanto o
router externo permanece responsável pela `MainPage`. A seleção de empresa
inicia uma sessão limpa com `Reset`; a seleção de funcionário usa `Push`.

As telas consultam os use cases reais e apresentam dados somente leitura,
listas filtradas pelo vínculo obrigatório, estados de carregamento, estados
vazios e erros recuperáveis. Campos de detalhe, células de listas e estados
visuais comuns foram extraídos como componentes pequenos, mantendo as
composições de empresa e funcionário separadas.

Também foram adicionados testes para o isolamento do histórico e para as
consultas por `CompanyID` e `EmployeeID`. O fluxo de edição inline, CRUD de
recibos e formulário real de cadastro continuam como próximos incrementos.

**Próximo passo:** retomar a edição inline com `GRAVAR`/`CANCELAR` e validações;
depois implementar a criação e edição de recibos.

# Passo 11:
Revisada a integração visual das telas 1.1 e 1.2 para preservar a composição
dos protótipos existentes. Os blocos de informações, linhas customizadas,
painéis com moldura e legenda, ações de recibo, larguras e espaçamentos foram
componentizados sem substituir a estrutura original por uma grade genérica.

O botão `VOLTAR` da Tela 1.1 foi movido para o rodapé, junto de `GRAVAR`, e
agora encerra a sessão do workspace executando `Reset` seguido de `Replace` em
`workspace.quick-access`. Na Tela 1.2, `VOLTAR` permanece no mesmo rodapé e
retorna pela pilha para a empresa.

**Próximo passo:** implementar o comportamento de edição e persistência dos
campos mantendo os componentes visuais atuais.

# Passo 12:
Corrigido o dimensionamento das telas 1.1 e 1.2 sem alterar a composição dos
protótipos. A janela inicia em `1280x800`, mantendo espaço para a sidebar e os
campos de largura fixa antes de o `HSplit` calcular suas restrições.

Os painéis de funcionários e recibos agora calculam a altura a partir da
linha visual real, aplicam o limite máximo de `300` para a lista e preservam o
acréscimo de `60` do protótipo. Um layout visual pequeno mantém o frame no
tamanho calculado mesmo quando o `Border` recebe altura extra; o espaço
restante fica fora do painel. Foram adicionados testes para a altura normal e
para o limite máximo das duas listas.

**Próximo passo:** implementar o comportamento de edição e persistência dos
campos mantendo os componentes visuais atuais.

# Passo 13:
Implementado o fluxo de edição inline das Telas 1.1 e 1.2. Os campos de
negócio alternam entre label e `Entry`, `GRAVAR` persiste pelos use cases de
update e `CANCELAR` restaura os valores originais. A Tela 1.2 passou a expor
primeiro nome e sobrenome separadamente; os códigos continuam somente leitura.

O router interno agora compartilha uma sessão de edição e bloqueia
`Push`/`Reset`/`Replace`/`Back` com diálogo de alterações não salvas. A
composição passou a construir widgets na UI thread depois das consultas
assíncronas, e foram adicionados testes para campos e guarda de navegação,
mantendo os testes existentes das consultas das páginas.

**Próximo passo:** implementar a criação e edição de recibos, reaproveitando a
mesma política de persistência e proteção contra navegação.

# Passo 14:
Refinado o estado de edição inline para separar o modo visual `Entry` de uma
alteração pendente. `GRAVAR`, `CANCELAR` e a proteção de navegação só são
ativados quando algum valor diverge do original; retornar ao valor inicial
desfaz o estado pendente.

Os estados label e edição agora usam apresentações exclusivas, com destaque
visual azul no modo editável. Após salvar uma empresa, a lista atualmente
montada na sidebar repete a consulta em goroutine e atualiza seus itens na UI
thread. Foram adicionados testes para dirty state, troca visual e reload da
lista.

**Próximo passo:** implementar a criação e edição de recibos, reaproveitando a
mesma política de persistência e proteção contra navegação.
