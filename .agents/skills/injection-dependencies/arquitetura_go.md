# Arquitetura Go com DDD

## Objetivo deste documento

Este arquivo serve como **fonte da verdade** para organizar projetos Go usando:

- DDD (Domain-Driven Design)
- Clean Architecture
- Hexagonal Architecture

O objetivo nao e decorar nomes de pastas. O objetivo e criar uma estrutura em que:

- o dominio fique protegido
- a regra de negocio nao dependa de framework
- detalhes tecnicos fiquem isolados
- cada camada tenha uma responsabilidade clara
- o projeto continue compreensivel mesmo crescendo

Este documento **nao contem tudo o que existe sobre DDD**.
DDD e um assunto muito maior, com temas como bounded contexts, ubiquitous language, aggregates, domain events, sagas, event sourcing e modelagem estrategica.

O foco aqui e outro:

- organizar servicos Go com boas fronteiras
- saber onde cada tipo de codigo deve morar
- evitar acoplamento entre negocio e tecnologia
- tomar decisoes consistentes no dia a dia
- ter exemplos concretos de como isso aparece em HTTP, filas, repositorios e casos de uso

---

## Principio central

O projeto deve ser organizado em torno do **dominio**, nao em torno de tecnologia.

Errado:

- `controllers`
- `models`
- `services`
- `database`
- `utils`

Certo:

- `domain`
- `application`
- `infra`
- `interfaces`

Tecnologia muda. Regra de negocio deveria durar mais.

---

## Regra de dependencias

A direcao das dependencias deve apontar para dentro, em direcao ao dominio.

```text
interfaces -> application -> domain
infra ------^           -> domain
cmd -> monta tudo
```

Regras:

- `domain` nao depende de `application`, `infra`, `interfaces` nem framework
- `application` depende de `domain`
- `infra` pode depender de `domain` para implementar contratos do dominio
- `interfaces` depende de `application` e usa DTOs/adapters de borda
- `cmd` so faz composicao, bootstrap e inicializacao

Se o dominio conhece HTTP, SQL, gRPC, RabbitMQ ou framework web, a arquitetura comecou a vazar.

---

## Papel do `internal`

Em Go, `internal` existe para restringir importacao fora do modulo/pacote permitido.

Estrutura comum:

```text
/cmd
/internal
  /domain
  /application
  /infra
  /interfaces
```

Quando usar:

- quando o codigo e interno da aplicacao e nao deve virar biblioteca publica
- quando voce quer deixar claro que a arquitetura pertence so a este servico

Quando nao usar:

- se o projeto for uma biblioteca Go para ser importada por outros modulos

`internal` nao substitui arquitetura. Ele apenas reforca fronteiras.

---

## Fluxo canonico

```text
HTTP/gRPC/Consumer
  -> handler
  -> usecase
  -> entidades e regras do dominio
  -> contrato de repositorio/servico
  -> implementacao em infra
  -> banco, fila, cache ou API externa
```

O handler nao decide regra de negocio.
O repositorio nao decide caso de uso.
O banco nao modela o dominio.

---

## Diretorios pais

### `/cmd`

Ponto de entrada da aplicacao.

O que entra:

- `main.go`
- leitura inicial de config
- criacao de conexoes
- injecao de dependencias
- montagem da aplicacao
- start do servidor HTTP/gRPC/consumer

O que nao entra:

- regra de negocio
- validacao de entidade
- logica de caso de uso
- SQL de repositorio

Filhos comuns:

- `/cmd/api`
- `/cmd/worker`
- `/cmd/migrate`
- `/cmd/consumer`

Cada subdiretorio em `cmd` representa um executavel diferente.

Exemplo:

```text
/cmd
  /api
    main.go
  /worker
    main.go
```

---

### `/internal`

Raiz da implementacao interna da aplicacao.

O que entra:

- tudo que faz parte do sistema e nao deve ser importado livremente por fora

Filhos canonicos:

- `/internal/domain`
- `/internal/application`
- `/internal/infra`
- `/internal/interfaces`

Se o projeto usar `internal`, o ideal e que as camadas principais fiquem abaixo dele.

---

### `/domain`

Coracao do sistema.

Aqui vivem os conceitos de negocio, invariantes, comportamentos e contratos centrais.

O que entra:

- entidades
- value objects
- agregados
- eventos de dominio
- servicos de dominio
- contratos de repositorio
- erros de dominio
- regras puras de negocio

O que nao entra:

- handler HTTP
- DTO de request/response
- ORM
- SQL
- framework web
- client de API externa
- acesso direto a banco

O dominio deve conseguir existir e ser testado sem subir servidor, banco ou fila.

Filhos mais comuns:

- `/domain/entity`
- `/domain/valueobject`
- `/domain/repository`
- `/domain/service`
- `/domain/event`
- `/domain/errors`

---

### `/application`

Camada que orquestra os casos de uso do sistema.

Ela coordena o que precisa acontecer para atender uma acao de negocio, usando o dominio.

O que entra:

- use cases
- comandos e queries
- orquestracao entre repositorios e servicos
- controle de transacao quando fizer parte do caso de uso
- definicao de portas de entrada/saida da aplicacao
- DTOs internos de caso de uso quando necessario

O que nao entra:

- regra de negocio fundamental da entidade
- detalhes de framework web
- SQL
- serializacao HTTP

Filhos mais comuns:

- `/application/usecase`
- `/application/command`
- `/application/query`
- `/application/dto`
- `/application/mapper`
- `/application/port`

Se `domain` responde "o que e valido no negocio", `application` responde "como o sistema executa esse objetivo".

---

### `/infra`

Camada dos detalhes tecnicos.

Aqui ficam implementacoes concretas de tudo aquilo que o dominio e a aplicacao definem como contrato.

O que entra:

- persistencia
- conexao com banco
- clientes HTTP
- publisher e consumer de mensageria
- cache
- configuracao tecnica
- observabilidade
- implementacao de repositorios

O que nao entra:

- decisao de regra de negocio
- validacao de entidade
- semantica de caso de uso

Filhos mais comuns:

- `/infra/persistence`
- `/infra/httpclient`
- `/infra/messaging`
- `/infra/cache`
- `/infra/config`
- `/infra/observability`
- `/infra/security`

`infra` pode conhecer tecnologia. O dominio nao.

---

### `/interfaces`

Camada de entrada e saida do sistema.

E a borda da aplicacao: recebe algo do mundo externo e adapta para o caso de uso.

O que entra:

- handlers HTTP
- controllers
- presenters
- middleware de borda
- routers
- servidores gRPC
- consumers que convertem mensagem em chamada de use case
- serializers/parsers de entrada e saida

O que nao entra:

- regra de negocio central
- implementacao de repositorio
- SQL
- modelagem de entidade por conveniencia de transporte

Filhos mais comuns:

- `/interfaces/http`
- `/interfaces/grpc`
- `/interfaces/consumer`
- `/interfaces/cli`

`interfaces` traduz protocolo. `application` executa o caso de uso.

---

## Filhos comuns e o tipo de codigo que vai em cada um

### `/domain/entity`

Contem entidades com identidade propria.

Vai aqui:

- structs de negocio com identidade
- metodos que alteram estado de forma valida
- invariantes do negocio
- validacoes que pertencem ao conceito de negocio

Exemplo:

- `Author`
- `Book`
- `Order`
- `User`

Pode e deve ter metodos quando fizer sentido.

Nao e so estrutura de dados.

Nao vai aqui:

- tag de JSON pensada para response HTTP
- dependencia de ORM
- codigo de acesso ao banco
- validacao de formato tecnico irrelevante ao negocio

---

### `/domain/valueobject`

Contem objetos definidos por valor, nao por identidade.

Vai aqui:

- tipos imutaveis ou quase imutaveis
- validacoes encapsuladas
- comparacao por valor

Exemplo:

- `Email`
- `Document`
- `Money`
- `Address`
- `Period`

Tambem pode abrigar wrappers fortes para IDs quando isso melhorar o dominio:

- `UserID`
- `OrderID`

Nao e obrigatorio criar VO para todo UUID. Crie quando isso trouxer semantica, seguranca ou legibilidade.

---

### `/domain/repository`

Contem contratos de persistencia que o dominio ou a aplicacao precisam.

Vai aqui:

- interfaces como `AuthorRepository`
- operacoes coerentes com o negocio

Exemplo:

- `Save(author *Author) error`
- `FindByID(id AuthorID) (*Author, error)`
- `ListActiveCustomers() ([]Customer, error)`

Nao vai aqui:

- `sql.DB`
- query SQL
- implementacao concreta
- detalhe de Postgres, MySQL, Redis

Repositorio em `domain` e contrato. Implementacao vai em `infra`.

---

### `/domain/service`

Contem servicos de dominio.

Use quando uma regra de negocio:

- nao cabe naturalmente em uma unica entidade
- envolve colaboracao entre multiplos conceitos do dominio
- precisa de linguagem explicita do negocio

Vai aqui:

- regras puras e semanticamente de negocio

Exemplo:

- politica de desconto
- calculo de elegibilidade
- regra de fechamento de folha

Nao use `service` como pasta coringa para qualquer coisa.

Se virou deposito de funcoes soltas, a modelagem esta fraca.

---

### `/domain/event`

Contem eventos de dominio.

Vai aqui:

- fatos importantes do negocio que ja aconteceram

Exemplo:

- `OrderPaid`
- `EmployeeClockedIn`
- `InvoiceClosed`

Eventos de dominio nao sao a mesma coisa que mensagem de broker.
Primeiro existe o fato no dominio. Depois, se necessario, isso pode ser traduzido para integracao em outra camada.

---

### `/domain/errors`

Contem erros semanticos do negocio.

Vai aqui:

- `ErrBookUnavailable`
- `ErrInvalidStatusTransition`
- `ErrAuthorInactive`

Ajuda a manter a linguagem do dominio explicita e evita espalhar strings pelo projeto.

---

### `/application/usecase`

Contem a implementacao dos casos de uso.

Vai aqui:

- `CreateAuthor`
- `RegisterBook`
- `BorrowBook`
- `CloseInvoice`

Responsabilidade:

- receber dados de entrada
- carregar entidades
- chamar metodos de dominio
- persistir alteracoes
- coordenar dependencias
- devolver resultado apropriado

Nao vai aqui:

- parse de HTTP
- query SQL
- regra que deveria estar na entidade

Use case e orquestrador, nao repositorio nem handler.

---

### `/application/command`

Contem comandos de entrada para operacoes que mudam estado.

Vai aqui:

- structs como `CreateUserCommand`
- dados necessarios para executar uma acao de escrita

Serve para deixar a entrada do caso de uso explicita.

---

### `/application/query`

Contem consultas orientadas a leitura.

Vai aqui:

- contratos e modelos para buscas
- query handlers
- respostas projetadas para leitura, quando fizer sentido separar leitura de escrita

Em cenarios simples, pode nem existir pasta propria.

---

### `/application/dto`

Contem DTOs internos da camada de aplicacao.

Vai aqui:

- input/output de casos de uso
- estruturas de transferencia entre borda e aplicacao

Nao vai aqui:

- entidade de dominio disfarcada
- payload acoplado a framework

DTO existe para transporte. Entidade existe para negocio.

---

### `/application/mapper`

Contem conversoes entre DTOs, entidades e modelos de borda.

Vai aqui:

- funcoes de traducao
- montagem de resposta

Use com moderacao.
Se todo fluxo exige dezenas de mapeamentos desnecessarios, pode haver burocracia demais na arquitetura.

---

### `/application/port`

Contem portas usadas pela aplicacao quando a interface faz mais sentido aqui do que no dominio.

Exemplo:

- notificacao
- clock
- geracao de token
- fila de jobs

Regra pratica:

- se o contrato expressa linguagem central do negocio, tende a ficar em `domain`
- se o contrato atende uma necessidade operacional do caso de uso, pode ficar em `application`

---

### `/infra/persistence`

Contem persistencia concreta.

Filhos comuns:

- `/infra/persistence/postgres`
- `/infra/persistence/mysql`
- `/infra/persistence/sqlite`
- `/infra/persistence/memory`

Vai aqui:

- implementacao de repositorios
- SQL
- uso de ORM
- transacoes concretas
- mapeamento para schema

Nao vai aqui:

- regra de negocio
- decisao de fluxo do caso de uso

---

### `/infra/httpclient`

Contem clients de APIs externas.

Vai aqui:

- chamadas REST/SOAP/gRPC externas
- autenticacao tecnica
- retry tecnico
- serializacao de integracao

O dominio nao deve saber que isso e um `http.Client`.

---

### `/infra/messaging`

Contem implementacoes de integracao com mensageria.

Filhos comuns:

- `/infra/messaging/rabbitmq`
- `/infra/messaging/kafka`
- `/infra/messaging/nats`

Vai aqui:

- publishers
- subscribers tecnicos
- adapters de broker
- serializacao de mensagens

Lembrete:

- evento de dominio pertence ao `domain`
- publicacao em RabbitMQ pertence ao `infra`

---

### `/infra/cache`

Contem detalhes de cache.

Vai aqui:

- Redis
- cache em memoria
- estrategias tecnicas de invalidacao

Cache e detalhe de implementacao, salvo quando ele altera explicitamente regra de negocio, o que e raro.

---

### `/infra/config`

Contem carregamento de configuracao tecnica.

Vai aqui:

- leitura de env
- parsing de config
- defaults tecnicos

Nao vai aqui:

- regra de negocio
- validacao de entidade

---

### `/infra/observability`

Contem preocupacoes transversais tecnicas.

Vai aqui:

- logger
- metrics
- tracing
- correlation IDs

O dominio pode receber um contrato abstrato se realmente precisar registrar algo, mas o mecanismo concreto fica em `infra`.

---

### `/interfaces/http`

Contem a interface HTTP.

Filhos comuns:

- `/interfaces/http/handler`
- `/interfaces/http/middleware`
- `/interfaces/http/request`
- `/interfaces/http/response`
- `/interfaces/http/router`

Vai aqui:

- bind de request
- validacao de payload de transporte
- chamada do use case
- traducao de erro para status code
- serializacao JSON

Nao vai aqui:

- regra de negocio central
- SQL
- mutacao de entidade sem passar por use case

Exemplo de responsabilidade:

- ler `path params`
- ler `query params`
- ler headers
- fazer bind do body JSON
- validar formato do request
- montar input do use case
- traduzir erro para resposta HTTP

---

### `/interfaces/grpc`

Contem a interface gRPC.

Vai aqui:

- server gRPC
- adaptacao de protobuf para DTO/use case
- tratamento de codigos de erro gRPC

Mesmo principio do HTTP: traduz protocolo e delega para aplicacao.

---

### `/interfaces/consumer`

Contem consumidores de fila do ponto de vista da entrada do sistema.

Vai aqui:

- leitura de mensagem
- parse do payload
- chamada de use case
- tratamento de retry na borda quando apropriado

Se o consumer esta decidindo regra de negocio, ele esta gordo demais.

---

### `/interfaces/cli`

Contem comandos de linha de comando como interface de entrada.

Vai aqui:

- parse de flags
- composicao de input
- chamada de use case

CLI e interface. Nao e lugar para regra de dominio.

---

## Onde cada tipo de codigo deve ficar

### Validacoes

Depende do tipo de validacao:

- validacao de formato de entrada HTTP -> `interfaces`
- validacao de contrato de caso de uso -> `application`
- validacao de regra de negocio -> `domain`

Exemplo:

- "campo obrigatorio no JSON" -> `interfaces/http/request`
- "nao pode criar pedido sem itens" -> `domain`
- "use case exige tenantID no contexto" -> `application`

---

### Query params, path variables, headers e body

Esses elementos pertencem a **interface de transporte**, entao nascem em `interfaces/http`.

#### Path variables

Exemplos:

- `GET /authors/{id}`
- `PATCH /books/{id}`

Uso correto:

- identificar recurso especifico
- carregar um agregado ou entidade por ID

Onde tratar:

- ler em `interfaces/http/handler`
- validar formato basico em `interfaces/http/request`
- converter para input do use case
- se virar conceito forte, transformar em VO/tipo forte antes de chegar ao dominio

O que nao fazer:

- deixar a entidade ler `chi.URLParam` ou `gin.Param`
- acoplar o dominio ao roteador

#### Query params

Exemplos:

- `GET /books?page=2&limit=20`
- `GET /books?author_id=123`
- `GET /books?status=available&sort=title`

Uso correto:

- filtros
- ordenacao
- paginacao
- busca textual
- flags de expansao ou visao

Onde tratar:

- leitura e parse em `interfaces/http`
- normalizacao de valores em DTO/request model
- envio para `application/query` ou `application/usecase`

O que nao fazer:

- colocar query param dentro de entidade de dominio so porque veio do HTTP
- misturar parametro de transporte com regra de negocio

#### Headers

Exemplos:

- `Authorization`
- `X-Request-ID`
- `X-Tenant-ID`
- `Idempotency-Key`

Uso correto:

- autenticacao/autorizacao
- correlacao
- contexto multitenant
- idempotencia
- metadados de protocolo

Onde tratar:

- leitura em middleware ou handler de `interfaces/http`
- traducao para contexto de aplicacao quando necessario

Header nao e conceito de dominio por padrao.
Ele so vira informacao de negocio se o caso de uso realmente precisar dele de forma semantica.

#### Body

Exemplos:

- `POST /authors`
- `PUT /books/{id}`

Uso correto:

- dados de entrada da operacao
- payloads de criacao ou atualizacao

Onde tratar:

- bind e validacao estrutural em `interfaces/http/request`
- conversao para DTO/input de `application`

O body nao deve ser passado bruto para o dominio.

---

### Paginacao, filtro e ordenacao

Esses conceitos geralmente vivem entre `interfaces` e `application`, nao em `domain`.

Exemplo:

```text
GET /books?page=2&limit=10&sort=title&status=available
```

Distribuicao recomendada:

- `interfaces/http/request`: parse de `page`, `limit`, `sort`, `status`
- `application/query`: definicao da consulta e das regras de uso da consulta
- `infra/persistence`: traducao disso para SQL, ORM ou outro mecanismo de busca

Regras praticas:

- limite maximo de pagina pode ficar em `application`
- parse de string para inteiro fica em `interfaces`
- traducao para `ORDER BY` fica em `infra`

Se `sort` aceita somente certos campos, essa whitelist costuma nascer em `application` ou no modelo de query, nao direto no banco.

---

### UUID e IDs

Opcoes validas:

- usar `uuid.UUID` diretamente quando simplicidade for suficiente
- criar value objects ou tipos fortes quando o dominio ganhar clareza com isso

Use tipo forte quando quiser evitar confusao como:

- `UserID` sendo passado onde deveria ser `OrderID`

---

### DTOs

DTO nao deve morar no `domain`.

Regra pratica:

- DTO de request/response de protocolo -> `interfaces`
- DTO de input/output de caso de uso -> `application`
- entidade de negocio -> `domain`

---

### Contratos

Nem toda interface precisa ir no `domain`.

Use esta heuristica:

- contrato essencial ao negocio -> `domain`
- contrato operacional do caso de uso -> `application`
- detalhe concreto -> `infra`

---

### Integracoes externas

Sempre fora do dominio.

Vai normalmente em:

- `infra/httpclient`
- `infra/messaging`
- `infra/security`

Se a aplicacao precisar usar isso, injete por contrato.

---

### Erros e status codes

Separacao recomendada:

- `domain`: erros semanticos de negocio
- `application`: erros de orquestracao, autorizacao de caso de uso ou pre-condicoes
- `interfaces/http`: mapeamento desses erros para status code e payload

Exemplo:

- `domain.ErrBookUnavailable`
- `domain.ErrInvalidStatusTransition`
- `application.ErrUnauthorizedTenantAccess`

Mapeamento HTTP tipico:

- erro de validacao de request -> `400 Bad Request`
- recurso nao encontrado -> `404 Not Found`
- conflito de regra de negocio -> `409 Conflict`
- autenticacao ausente/invalida -> `401 Unauthorized`
- sem permissao -> `403 Forbidden`
- erro inesperado -> `500 Internal Server Error`

O dominio nao retorna status code.
O dominio retorna significado.
Quem traduz isso para HTTP e a camada de `interfaces`.

---

### Idempotencia

Idempotencia costuma aparecer em endpoints como:

- `POST /payments`
- `POST /clock-in`
- `POST /orders`

Normalmente a chave de idempotencia entra por header, por exemplo:

- `Idempotency-Key`

Distribuicao recomendada:

- `interfaces/http`: le o header
- `application`: decide a politica de idempotencia do caso de uso
- `infra`: persiste ou consulta a chave

Nao coloque a logica de idempotencia inteira no handler.
Ela costuma ser parte da execucao segura do caso de uso.

---

### Autenticacao e autorizacao

Separar bem isso evita muita confusao.

Autenticacao:

- verifica quem esta chamando
- normalmente acontece na borda, em middleware ou adapter

Autorizacao:

- verifica se aquele ator pode executar aquela acao
- pode acontecer em `application`
- em alguns casos, pode depender de regra do `domain`

Exemplo:

- middleware extrai usuario do JWT
- handler monta contexto autenticado
- use case verifica se aquele usuario pode alterar aquele recurso
- entidade pode reforcar uma regra de estado ou ownership se isso fizer parte do negocio

---

### Transacoes

A decisao de iniciar/encerrar transacao geralmente nasce na camada de `application`, porque isso faz parte da execucao do caso de uso.

A implementacao concreta da transacao fica em `infra`.

---

### Banco de dados

Banco nao modela o dominio. Banco persiste o dominio.

Por isso:

- schema, SQL e ORM ficam em `infra`
- entidades e comportamentos ficam em `domain`

Nao modele entidade pensando primeiro em tabela.

---

## Interface HTTP na pratica

### Exemplo 1: buscar um livro por ID

Endpoint:

```text
GET /books/{id}
```

Fluxo:

1. `interfaces/http/handler` le `id` da rota
2. valida formato basico do identificador
3. monta input do caso de uso
4. chama `application/usecase/GetBookByID`
5. use case chama contrato de repositorio
6. `infra/persistence/postgres` busca no banco
7. use case devolve DTO de saida
8. handler serializa JSON e responde

O que nao fazer:

- rodar SQL no handler
- devolver a entidade de dominio diretamente como response model por conveniencia

---

### Exemplo 2: listar livros com filtro e paginacao

Endpoint:

```text
GET /books?author_id=123&status=available&page=1&limit=20&sort=title
```

Leitura correta:

- `author_id` -> filtro
- `status` -> filtro
- `page` -> paginacao
- `limit` -> paginacao
- `sort` -> ordenacao

Onde cada parte mora:

- parse e validacao de formato -> `interfaces/http/request`
- definicao da consulta -> `application/query`
- traducao para query SQL -> `infra/persistence`

O dominio nem sempre precisa participar de listagens simples.
Em cenarios de consulta, a `application/query` pode ser suficiente.

---

### Exemplo 3: criar autor

Endpoint:

```text
POST /authors
Content-Type: application/json

{
  "name": "Machado de Assis",
  "email": "autor@dominio.com"
}
```

Fluxo:

1. handler faz bind do body
2. request model valida estrutura minima
3. handler monta `CreateAuthorInput`
4. use case cria entidades/value objects necessarios
5. dominio valida invariantes
6. repositorio persiste
7. handler traduz resultado para `201 Created`

Exemplo de distribuicao:

- email malformado no payload -> pode falhar cedo na borda
- regra "nao pode criar autor inativo sem motivo" -> pertence ao dominio
- verificacao de duplicidade -> geralmente passa por application + repository

---

### Exemplo 4: atualizar parcialmente um recurso

Endpoint:

```text
PATCH /books/{id}
```

Regra pratica:

- `PATCH` trata alteracoes parciais
- `PUT` tende a representar substituicao completa

Cuidados:

- o handler nao deve sair alterando struct do dominio campo por campo sem criterio
- o ideal e chamar um use case explicito, por exemplo `RenameBook`, `ArchiveBook`, `UpdateBookMetadata`

Quanto mais explicito o caso de uso, menos ambiguidade entra no sistema.

---

## Comandos e queries

Separar leitura de escrita costuma ajudar muito.

Comandos:

- alteram estado
- passam por regras de negocio
- costumam usar entidades e repositorios de escrita

Queries:

- leem estado
- podem usar modelos otimizados de leitura
- nem sempre precisam materializar entidade completa

Exemplo:

- `CreateBook` -> comando
- `ListAvailableBooks` -> query
- `GetAuthorProfile` -> query
- `CloseInvoice` -> comando

Nem todo projeto precisa de CQRS formal.
Mas separar mentalmente `command` de `query` ja melhora bastante o desenho.

---

## Testes por camada

### Testes de `domain`

Devem ser os mais puros e mais rapidos.

Teste aqui:

- invariantes
- transicoes de estado
- metodos de entidade
- value objects
- servicos de dominio

Nao dependa de banco, HTTP ou framework.

### Testes de `application`

Teste aqui:

- orquestracao de caso de uso
- chamadas esperadas para contratos
- cenarios de sucesso e falha
- controle de transacao

Use mocks, fakes ou implementacoes em memoria quando fizer sentido.

### Testes de `infra`

Teste aqui:

- repositorios concretos
- integracoes externas
- serializacao
- adaptadores de broker

Frequentemente fazem sentido testes de integracao.

### Testes de `interfaces`

Teste aqui:

- binding de request
- parse de query/path/header/body
- mapeamento de erro para status code
- contrato HTTP/gRPC exposto

Esses testes garantem que a borda conversa corretamente com a aplicacao.

---

## Anti-patterns comuns

### 1. `service` virar pasta lixo

Sintoma:

- toda regra vai para `service`
- entidades ficam anemicas

Correcao:

- mova comportamento para entidade ou VO quando fizer sentido
- deixe `domain/service` apenas para regras que realmente cruzam conceitos

### 2. Handler gordo

Sintoma:

- handler valida negocio
- abre transacao
- chama banco direto
- monta resposta complexa sozinho

Correcao:

- handler adapta protocolo
- use case executa a acao

### 3. Repositorio com regra de negocio

Sintoma:

- repositorio decide status
- repositorio aplica politica de negocio

Correcao:

- repositorio persiste e recupera
- regra fica em `domain` ou `application`

### 4. Entidade anemica

Sintoma:

- entidade e so struct com getters/setters
- toda regra mora fora

Correcao:

- coloque comportamento e invariantes dentro da entidade

### 5. DTO virando entidade

Sintoma:

- request/response structs vazam para o dominio

Correcao:

- DTO serve transporte
- entidade serve negocio

### 6. Dominio acoplado a framework

Sintoma:

- entidade importa pacote HTTP
- value object depende de ORM

Correcao:

- remova detalhes tecnicos do `domain`

### 7. Pasta `utils` crescendo sem controle

Sintoma:

- tudo que nao tem lugar vai para `utils`

Correcao:

- separe por responsabilidade real
- helper de HTTP fica perto de HTTP
- helper de persistencia fica perto de persistencia
- regra de negocio nao vira utilitario generico

---

## Como pensar quando surgir uma duvida

Quando nao souber onde colocar um codigo, faca estas perguntas:

1. Isso expressa regra de negocio ou detalhe tecnico?
2. Isso existiria mesmo se eu trocasse HTTP por gRPC?
3. Isso existiria mesmo se eu trocasse Postgres por Redis?
4. Isso tem semantica de dominio ou e so transporte?
5. Isso esta orquestrando um objetivo ou implementando uma tecnologia?

Se a resposta apontar para negocio, aproxime do `domain`.
Se apontar para execucao do caso de uso, aproxime de `application`.
Se apontar para protocolo ou tecnologia, aproxime de `interfaces` ou `infra`.

---

## Estrutura canonica de referencia

```text
/cmd
  /api
    main.go
  /worker
    main.go

/internal
  /domain
    /entity
    /valueobject
    /repository
    /service
    /event
    /errors

  /application
    /usecase
    /command
    /query
    /dto
    /mapper
    /port

  /infra
    /persistence
      /postgres
      /sqlite
      /memory
    /httpclient
    /messaging
      /rabbitmq
      /kafka
    /cache
    /config
    /observability
    /security

  /interfaces
    /http
      /handler
      /middleware
      /request
      /response
      /router
    /grpc
    /consumer
    /cli
```

Essa estrutura e canonica, nao obrigatoria em 100 por cento dos projetos.

O importante e preservar os principios:

- dominio isolado
- use cases orquestrando
- detalhes tecnicos nas bordas
- dependencias apontando para dentro

---

## FAQ curto

### Pode colocar metodo dentro de `entity`?

Sim. Em geral, deve.

Entidade rica protege regra de negocio melhor do que struct anemica com regra espalhada.

### Entity e Value Object sao a mesma coisa?

Nao.

- `Entity` tem identidade
- `Value Object` e definido por valor

### Preciso criar Value Object para todo UUID?

Nao.

Crie quando isso deixar o dominio mais seguro e mais expressivo.

### Onde fica funcao de validacao?

Depende do tipo:

- validacao de transporte -> `interfaces`
- validacao de caso de uso -> `application`
- validacao de negocio -> `domain`

### `internal` e obrigatorio?

Nao.

Mas em aplicacoes Go ele costuma ser uma boa escolha para reforcar encapsulamento arquitetural.
