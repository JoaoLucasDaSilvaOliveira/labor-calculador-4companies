# Curso prático: testes unitários do Router com Fyne — Sprint 1

Este guia ensina os testes necessários para concluir a Sprint 1: a operação
**Router.Reset** e a segurança do histórico. Ele não altera código de produção.
Use-o para entender o contrato, escrever os testes e deixar que eles indiquem
se a implementação o cumpre.

## 1. O que é um teste unitário aqui?

Um teste unitário verifica uma unidade pequena de comportamento automaticamente.
A unidade deste caso é o **Router**: ele recebe comandos de navegação, mantém o
histórico e manda uma view ao AnimatedContent.

Não estamos testando pixels, uma janela real do sistema operacional ou cliques
de mouse. Estamos testando regras observáveis:

- qual é a rota atual;
- quais entradas ficaram no histórico;
- se é possível voltar;
- qual view terminou renderizada;
- o que acontece quando uma factory falha.

Como o Router usa objetos Fyne reais, isso é um teste de unidade com uma
integração pequena e controlada ao Fyne — não é um teste manual de interface.

## 2. Vocabulário essencial

| Termo | Significado |
| --- | --- |
| SUT (System Under Test) | O objeto testado: router. |
| Arrange | Preparar cenário: app, router, rotas e histórico inicial. |
| Act | Executar a ação, por exemplo Router.Reset. |
| Assert | Conferir o resultado com if e t.Fatal ou t.Fatalf. |
| Fixture | Preparação reutilizável, como newTestRouter. |
| Factory | Função que constrói a view de uma rota (PageFactory). |
| Invariante | Regra sempre verdadeira; após Reset, há uma entrada no histórico. |

O modelo mental de qualquer teste é:

~~~text
Arrange → Act → Assert
preparar → agir → conferir
~~~

## 3. O papel do Fyne no teste

As páginas do router são fyne.CanvasObject. Por isso, as factories usadas nos
testes devolvem objetos reais e simples, como widget.NewLabel("Main").

No início de testes que usam widgets ou animação, crie a aplicação de teste:

~~~go
fyneApp := test.NewApp()
t.Cleanup(fyneApp.Quit)
~~~

- test.NewApp() cria uma aplicação apropriada para teste, sem abrir o app
  desktop normal.
- t.Cleanup(fyneApp.Quit) garante que ela será encerrada ao fim do teste,
  inclusive quando alguma asserção falhar.

O teste de registro duplicado pode dispensar isso porque ele não chega a
construir nem exibir widgets.

## 4. Contrato da Sprint 1

Reset(route, params) inicia uma sessão de navegação nova:

~~~text
histórico antes:  [A, B, C]
Reset(D)
histórico depois: [D]
~~~

Depois de um Reset bem-sucedido:

1. Current() retorna D.
2. CanGoBack() retorna false.
3. Back() retorna false.
4. Nenhuma entrada de A, B ou C permanece acessível.

### Reset não é apenas Replace

Replace troca somente a última entrada:

~~~text
[A, B, C] + Replace(D) = [A, B, D]
~~~

Logo, Back() ainda retorna para B. Isso está correto para uma troca local de
tela, mas não encerra a sessão anterior.

~~~text
[A, B, C] + Reset(D) = [D]
~~~

Portanto, é perfeitamente razoável que Reset **reaproveite** a lógica de
Replace depois de preparar uma sessão nova. O que não implementa Reset é chamar
somente Replace sobre o histórico antigo. A diferença é semântica: uma operação
substitui o topo; a outra zera a pilha.

## 5. Como aproveitar os testes existentes

O arquivo-base é internal/ui/navigation/router_test.go. Ele já contém uma
fixture útil:

~~~go
router := newTestRouter(t)
if err := router.Replace(RouteMain, nil); err != nil {
    t.Fatalf("Replace(RouteMain) returned an error: %v", err)
}
~~~

- newTestRouter(t) registra factories previsíveis para RouteMain e
  routeTestDetails.
- A checagem do erro impede que o teste continue com cenário inválido.
- t.Fatalf encerra o teste e mostra o erro com contexto.
- t.Helper(), dentro de newTestRouter, faz a falha apontar para quem chamou o
  helper, e não para dentro dele.

Os testes estão no pacote navigation, e não em navigation_test. Por isso podem
verificar router.history. Para esta Sprint, isso é justificável: o requisito
exige explicitamente que o histórico termine com uma única entrada.

## 6. Os quatro testes da Sprint

### 6.1 TestRouterResetClearsHistory

Este teste satisfaz o CA-01. Primeiro monte uma pilha com três entradas; depois
faça Reset e confirme que existe apenas uma:

~~~go
func TestRouterResetClearsHistory(t *testing.T) {
    fyneApp := test.NewApp()
    t.Cleanup(fyneApp.Quit)

    router := newTestRouter(t)
    // Monte [Main, Details, Details].

    if err := router.Reset(RouteMain, nil); err != nil {
        t.Fatalf("Reset(RouteMain) returned an error: %v", err)
    }

    if got := len(router.history); got != 1 {
        t.Fatalf("history length after Reset = %d, want 1", got)
    }
    if got := router.Current(); got != RouteMain {
        t.Fatalf("Current() after Reset = %q, want %q", got, RouteMain)
    }
}
~~~

Use rotas já registradas na fixture. Por exemplo, faça Replace(RouteMain, nil)
e duas chamadas a Push(routeTestDetails, ...). Montar mais de uma entrada é
crucial: um teste com histórico vazio não prova que Reset limpa contexto
anterior.

### 6.2 TestRouterBackReturnsFalseAfterReset

Este é o CA-02: depois de iniciar nova sessão, o usuário não pode voltar à
sessão anterior.

~~~go
func TestRouterBackReturnsFalseAfterReset(t *testing.T) {
    fyneApp := test.NewApp()
    t.Cleanup(fyneApp.Quit)

    router := newTestRouter(t)
    // Monte [Main, Details].

    if err := router.Reset(RouteMain, nil); err != nil {
        t.Fatalf("Reset(RouteMain) returned an error: %v", err)
    }

    if router.CanGoBack() {
        t.Fatal("CanGoBack() after Reset = true, want false")
    }
    if router.Back() {
        t.Fatal("Back() after Reset = true, want false")
    }
}
~~~

As duas asserts não são redundantes. CanGoBack é a consulta usada pela UI para
habilitar ou desabilitar o botão; Back é o comando que efetivamente muda o
estado.

### 6.3 Falha da factory e política de erro

O CA-03 exige que uma factory com erro não deixe estado parcial. O plano da
Sprint recomenda esta política:

1. Tentar construir a rota solicitada.
2. Se falhar, construir o fallback.
3. Se o fallback funcionar, o histórico final é [fallback].
4. Retornar o erro original para quem chamou Reset.

Registre uma factory que falha intencionalmente:

~~~go
const routeFailing RouteID = "test.failing"

err := router.Register(routeFailing, func(any) (fyne.CanvasObject, error) {
    return nil, errors.New("falha intencional da factory")
})
if err != nil {
    t.Fatalf("Register(routeFailing) returned an error: %v", err)
}
~~~

Depois de partir de [Main, Details], faça Reset(routeFailing, nil) e valide:

~~~go
if err == nil {
    t.Fatal("Reset(routeFailing) returned nil error, want an error")
}
if got := len(router.history); got != 1 {
    t.Fatalf("history length after failed Reset = %d, want 1", got)
}
if got := router.Current(); got != RouteMain {
    t.Fatalf("Current() after failed Reset = %q, want fallback %q", got, RouteMain)
}
~~~

Nome recomendado: TestRouterResetUsesFallbackWhenFactoryFails. Ele é mais
preciso que “preserves state”, pois a política recomendada não preserva o
histórico antigo: abre uma sessão limpa no fallback.

> Atenção ao estado atual do código: Reset limpa history antes de chamar
> Replace. Em uma falha, o fallback é construído/exibido, mas o histórico
> anterior é restaurado depois. Isso não coincide com a política recomendada e
> pode deixar histórico e view divergentes. O teste deve expor essa decisão;
> este guia não altera o método.

### 6.4 TestRouterLatestResetWinsDuringTransition

Este é o CA-04. Ele protege contra uma animação antiga terminar depois de uma
navegação mais recente e sobrescrever a tela correta.

~~~text
Replace(Main) → Push(Details) → Reset(Main) → aguardar → conferir Main
~~~

O padrão já existe em TestRouterLatestNavigationWinsDuringTransition:

~~~go
time.Sleep(canvas.DurationShort * 3)
~~~

- canvas.DurationShort é uma duração curta padrão do Fyne.
- O multiplicador 3 deixa callbacks de fade em andamento ocorrerem.
- Esperar tempo não é ideal para regra de negócio, mas aqui a propriedade
  testada depende deliberadamente da transição assíncrona.

Depois, valide a rota e a view. Validar somente Current() não basta: histórico
correto e view errada ainda seria um bug.

~~~go
if got := router.Current(); got != RouteMain {
    t.Fatalf("Current() = %q, want %q", got, RouteMain)
}

outlet := router.View().(*fyne.Container)
label, ok := outlet.Objects[0].(*widget.Label)
if !ok {
    t.Fatalf("rendered object type = %T, want *widget.Label", outlet.Objects[0])
}
if label.Text != "Main" {
    t.Fatalf("rendered label = %q, want Main", label.Text)
}
~~~

## 7. Quando extrair helpers?

Se a preparação se repetir, extraia helpers pequenos e com nome de domínio:

~~~go
func pushDetails(t *testing.T, router *Router, companyID int) {
    t.Helper()
    if err := router.Push(routeTestDetails, CompanyDetailsParams{CompanyID: companyID}); err != nil {
        t.Fatalf("Push(routeTestDetails) returned an error: %v", err)
    }
}
~~~

Evite helpers genéricos demais que escondem o cenário. “Abrir detalhes” deixa
o teste legível; uma lista abstrata de passos torna o teste mais difícil de
entender do que as chamadas diretas.

## 8. Executando os testes

Comece pelo pacote do Router:

~~~bash
go test ./internal/ui/navigation -count=1
~~~

- go test compila o pacote e executa funções TestXxx.
- ./internal/ui/navigation limita a execução à unidade alterada.
- -count=1 desativa o cache de resultados, garantindo execução atual.

Para executar um teste específico:

~~~bash
go test ./internal/ui/navigation -run '^TestRouterResetClearsHistory$' -count=1
~~~

- -run recebe uma expressão regular (regex).
- ^ representa o início e $ o fim do nome: isso evita rodar nomes parecidos.

Quando os quatro estiverem verdes, amplie:

~~~bash
go test ./internal/ui/... -count=1
go vet ./internal/ui/...
~~~

go vet não executa testes; ele procura construções suspeitas que podem compilar
mesmo contendo erros prováveis.

Se o cache global do Go estiver somente leitura:

~~~bash
env GOCACHE=/tmp/labor-calculator-go-cache go test ./internal/ui/navigation -count=1
~~~

env VAR=valor comando define uma variável apenas naquele processo. GOCACHE
aponta o cache temporário da compilação Go.

## 9. Checklist de entrega

- [ ] TestRouterResetClearsHistory cobre [A, B, C] → [D].
- [ ] TestRouterBackReturnsFalseAfterReset prova que não há retorno.
- [ ] Um teste de factory com falha define e protege a política de fallback.
- [ ] TestRouterLatestResetWinsDuringTransition valida rota e view.
- [ ] Os testes passam isoladamente.
- [ ] go test ./internal/ui/... -count=1 passa.
- [ ] go vet ./internal/ui/... passa.

## 10. Próximo passo

Escreva primeiro os quatro testes sem mudar router.go. Se algum falhar, compare
a falha com o contrato deste documento e com a política de erro escolhida.
Assim, o teste protege a regra de navegação da Sprint 1, em vez de apenas
confirmar o comportamento atual da implementação.
