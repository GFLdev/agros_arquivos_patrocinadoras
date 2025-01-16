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
)

// LoadConfig carrega as configurações do servidor.
func LoadConfig() {
	log.Println("Carregando arquivo de configurações")
	viper.SetDefault("environment", "development")

	// Caminho do arquivo config.json
	viper.AddConfigPath(".")
	// Nome do arquivo
	viper.SetConfigName("config.json")
	// Extensão do arquivos
	viper.SetConfigType("json")
	// Leitura
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Erro fatal no arquivo de configuração: %s\n", err)
	}

	env := viper.GetString("environment")
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
		AllowOrigins:     viper.GetStringSlice("origins"),
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

	// Redirecionamento para HTTPS
	e.Pre(middleware.HTTPSRedirect())

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
		// Sistema de arquivos estáticos
		middleware.StaticWithConfig(middleware.StaticConfig{
			Filesystem: http.Dir("./frontend/dist"),
		}),
	)
}

// Serve inicializa o servidor echo.
func Serve(e *echo.Echo) {
	env := viper.GetString("environment")
	port := viper.GetInt("port")
	if env == "production" {
		log.Printf("Iniciando servidor de produção na porta %d\n", port)

		err := e.StartAutoTLS(":" + strconv.Itoa(port))
		if err != nil {
			log.Fatalf("Erro na execução do servidor, %s\n", err)
		}
	} else if env == "development" {
		log.Printf("Iniciando servidor de desenvolvimento na porta %d\n", port)
		certFile := viper.GetString("cert_file")
		keyFile := viper.GetString("key_file")

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
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	ConfigMiddleware(e)

	// Grupos
	//adminGroup := e.Group("/admin")

	// Rotas
	e.GET("/", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/html")

		return nil
	})
	e.POST("/login", handlers.LoginHandler)
	e.POST("/download", handlers.DownloadFileHandler)

	// Inicializar servidor
	Serve(e)
}
