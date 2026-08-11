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

==**APAGAR QUANDO TERMINAR:** No momento ainda to aqui, é preciso criar interfaces que recebam os Execute dos usecases e injetar dentro das funções ...Component. Observar /internal/ui/components -> @companies_list.go e @company_details.go para exemplos de implementação
### Atualização:
No momento eu ja passei tudo, ou quase tudo, para usecase e query de verdade. Agora eu preciso fazer as funções que são parecidas ou redundantes se tornarem uma coisa só.
### O que eu to fazendo agora:
O codex criou um widget personalizado para servir de campo clicavel e ja colocou na função showCompanyInformations(), o que é ótimo. Agora é preciso salvar newInlineEditableCompanyField em vars pra os botões Gravar e Cancelar poderem manipula-las. Precisa criar o botão cancelar tb.==

**Próximo passo:** Componentizar os elementos repetidos/semelhantes das telas 1.1 e 1.2.