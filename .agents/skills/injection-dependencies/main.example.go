package main

import (
	authorUC "estudos-microservicos/crud/internal/application/usecase/author"
	bookUC "estudos-microservicos/crud/internal/application/usecase/book"
	"estudos-microservicos/crud/internal/infra/config"
	"estudos-microservicos/crud/internal/infra/persistence/sqlite"
	authorHandlerLib "estudos-microservicos/crud/internal/interfaces/http/author"
	bookHandlerLib "estudos-microservicos/crud/internal/interfaces/http/book"
	"estudos-microservicos/crud/internal/interfaces/http/shared"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/joho/godotenv"
)

func main() {
	// LOAD DAS CONFIGURAÇÕES BASICAS
	// ----------------------------------------------------------------------------
		//load .env $HOME
		slog.Info("Carregando .env")
		
		baseDir, existsEnv := os.LookupEnv("HOME")
		if !existsEnv {

			log.Fatal("Erro ao carregar a $HOME")
		}

		//join .env $HOME+ProjectPath
		envPath := path.Join(baseDir, "Estudos Golang/estudos microserviços/crud/.env")

		//load main .env
		if err := godotenv.Load(envPath); err != nil {
			log.Fatal("Ocorreu um erro ao carrgar o arquivo .env principal: "+err.Error())
		}
		
		slog.Info("Carregamento finalizado")
	// ----------------------------------------------------------------------------

	// LOAD DA CONFIGURAÇÃO DE BANCO DE DADOS
	// ----------------------------------------------------------------------------
		slog.Info("Iniciando a configuração do banco de dados")

		dbConfig := config.LoadDBConfig()
		// inicia o pool de conexões com o banco de dados e sobe o esquema do banco de dados atual

		slog.Info("Iniciando o pool de conexões com o banco de dados")
		database, err := sqlite.OpenDB(dbConfig)

		if err != nil {
			log.Fatal("Erro de banco de dados: "+err.Error())
		}
		
		//garante que feche o pool quando a aplicação finalizar
		defer database.Close()

		// instancia a fábrica de repositórios
		repositoryFactory := sqlite.NewRepositoryFactory(database)
		
		slog.Info("Configuração do banco de dados finalizada")
	// ----------------------------------------------------------------------------

	// INICIALIZAÇÃO DOS REPOSITORIOS
	// ----------------------------------------------------------------------------
		slog.Info("Instanciando repository factories")
		authorRepository := repositoryFactory.NewAuthorRepository()
		bookRepository := repositoryFactory.NewBookRepository()
	// ----------------------------------------------------------------------------

	// INICIALIZAÇÃO DOS USECASES
	// ----------------------------------------------------------------------------
		slog.Info("Instanciando usecases")
		//autor
		authorCreateUseCase := authorUC.CreateUseCase{Repo: authorRepository}
		authorListUseCase := authorUC.ListUseCase{Repo: authorRepository}
		authorListByIdUseCase := authorUC.ListByIdUseCase{Repo: authorRepository}
		authorUpdateUseCase := authorUC.UpdateUseCase{Repo: authorRepository}
		authorDeleteUseCase := authorUC.DeleteUseCase{Repo: authorRepository}

		//livro
		bookCreateUseCase := bookUC.CreateUseCase{Repo: bookRepository}
		bookListUseCase := bookUC.ListUseCase{Repo: bookRepository}
		bookListByIdUseCase := bookUC.ListByIdUseCase{Repo: bookRepository}
		bookUpdateUseCase := bookUC.UpdateUseCase{Repo: bookRepository}
		bookDeleteUseCase := bookUC.DeleteUseCase{Repo: bookRepository}
	// ----------------------------------------------------------------------------
	
	// INICIALIZAÇÃO DOS HANDLERS
	// ----------------------------------------------------------------------------
		slog.Info("Instanciando handlers")
		//autor
		authorHandler := authorHandlerLib.NewHandler(&authorCreateUseCase, &authorListUseCase, &authorListByIdUseCase, &authorUpdateUseCase, &authorDeleteUseCase)

		//livro
		bookHandler := bookHandlerLib.NewHandler(&bookCreateUseCase, &bookListUseCase, &bookListByIdUseCase, &bookUpdateUseCase, &bookDeleteUseCase)
	// ----------------------------------------------------------------------------
	
	// INICIALIZAÇÃO DAS ROTAS
		slog.Info("Instanciando route providers")
		//autor - provider
		authorRouteProvider := authorHandlerLib.NewAuthorRouteProvider(authorHandler)

		// livro - provider
		bookRouteProvider := bookHandlerLib.NewBookRouteProvider(bookHandler)

		slog.Info("Instanciando router")
		//instancia de um router de fato
		mux := shared.NewRouter(authorRouteProvider, bookRouteProvider)
	// ----------------------------------------------------------------------------
	
	// INICIALIZAÇÃO DO SERVIDOR
	// ----------------------------------------------------------------------------
		slog.Info("Subindo o servidor")
		port := os.Getenv("SERVER_PORT")
		
		slog.Info("Ouvindo na porta: "+port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatal("erro ao iniciar servidor: "+err.Error())
		}
	// ----------------------------------------------------------------------------
}