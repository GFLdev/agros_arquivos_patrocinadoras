package main

import (
	"agros_arquivos_patrocinadoras/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// LoadConfig carrega as configurações nas variáveis de ambiente.
func LoadConfig() {
	log.Println("Carregando variáveis de ambiente")
	viper.SetDefault("ENVIRONMENT", "development")

	// Caminho do arquivo .env
	viper.AddConfigPath(".")
	// Nome do arquivo
	viper.SetConfigName(".env")
	// Extensão do arquivos
	viper.SetConfigType("env")
	// Leitura
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Erro fatal no arquivo de configuração: %s\n", err)
	}

	env := viper.GetString("ENVIRONMENT")
	if env == "production" {
		log.Println("Configurando servidor de produção")
	} else if env == "development" {
		log.Println("Configurando servidor de desenvolvimento")
	} else {
		log.Println("Ambiente não definido. Fallback para development")
	}
}

// ConfigMiddleware configura os middlewares a serem utilizados pelo servidor.
func ConfigMiddleware(e *echo.Echo) {
	log.Println("Configurando middlewares")

	corsConfig := middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     strings.Split(viper.GetString("ORIGINS"), ","),
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderAccessControlAllowOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	}

	e.Use(
		// Limitações de requisições IP/segundo
		middleware.RateLimiter(
			middleware.NewRateLimiterMemoryStore(rate.Limit(5)),
		),
		// XSSProtection
		middleware.Secure(),
		middleware.Recover(),
		// GZip
		middleware.Gzip(),
		// CORS
		middleware.CORSWithConfig(corsConfig),
	)
}

// Serve inicializa o servidor echo.
func Serve(e *echo.Echo) {
	env := viper.GetString("ENVIRONMENT")
	port := viper.GetInt("PORT")
	if env == "production" {
		log.Printf("Iniciando servidor de produção na porta %d\n", port)

		err := e.StartAutoTLS(":" + strconv.Itoa(port))
		if err != nil {
			log.Fatalf("Erro na execução do servidor, %s\n", err)
		}
	} else if env == "development" {
		log.Printf("Iniciando servidor de desenvolvimento na porta %d\n", port)
		certFile := viper.GetString("CERT_FILE")
		keyFile := viper.GetString("KEY_FILE")

		err := e.StartTLS(":"+strconv.Itoa(port), certFile, keyFile)
		if err != nil {
			log.Fatalf("Erro na execução do servidor:\n\t%s\n", err)
		}
	}
}

func main() {
	log.Println("Iniciando aplicação...")

	// Configurações
	LoadConfig()

	// Servidor Echo
	e := echo.New()
	e.HideBanner = true
	ConfigMiddleware(e)

	// Grupos
	adminGroup := e.Group("/admin")

	// SPA
	if viper.GetString("ENVIRONMENT") == "production" {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Filesystem: http.Dir("./frontend/dist"),
			HTML5:      true,
		}))
	} else {
		res, err := http.Get("http: //localhost:5173/")
		if err != nil {
			log.Fatalf("Erro na obtenção da página estática:\n\t%s\n", err)
		}

		e.Static("/", "")
	}

	// Rotas
	e.POST("/login", handlers.LoginHandler)
	e.POST("/download", handlers.DownloadFileHandler)

	// Inicializar servidor
	Serve(e)
}
