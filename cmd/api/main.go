package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JuanCS-Dev/typecraft/internal/api/handlers"
	"github.com/JuanCS-Dev/typecraft/internal/config"
	"github.com/JuanCS-Dev/typecraft/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Banner
	fmt.Println(`
╔════════════════════════════════════════════════════════════════════╗
║                                                                    ║
║   ████████╗██╗   ██╗██████╗ ███████╗ ██████╗██████╗  █████╗ ███████╗████████╗
║   ╚══██╔══╝╚██╗ ██╔╝██╔══██╗██╔════╝██╔════╝██╔══██╗██╔══██╗██╔════╝╚══██╔══╝
║      ██║    ╚████╔╝ ██████╔╝█████╗  ██║     ██████╔╝███████║█████╗     ██║   
║      ██║     ╚██╔╝  ██╔═══╝ ██╔══╝  ██║     ██╔══██╗██╔══██║██╔══╝     ██║   
║      ██║      ██║   ██║     ███████╗╚██████╗██║  ██║██║  ██║██║        ██║   
║      ╚═╝      ╚═╝   ╚═╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝        ╚═╝   
║                                                                    ║
║                   AI-Powered Book Production Engine               ║
║                         Servidor API v0.1.0                        ║
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝
	`)

	// Carregar configurações
	log.Println("📋 Carregando configurações...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Erro ao carregar configurações: %v", err)
	}
	log.Println("✅ Configurações carregadas")

	// Conectar ao banco de dados
	log.Println("🔌 Conectando ao banco de dados...")
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco: %v", err)
	}

	// Executar migrations
	log.Println("🔄 Executando migrations...")
	if err := database.Migrate(); err != nil {
		log.Fatalf("❌ Erro nas migrations: %v", err)
	}

	// Configurar Gin
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Criar router
	router := gin.Default()

	// Middleware de CORS
	router.Use(corsMiddleware(cfg.AllowedOrigins))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		if err := database.Health(); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "Typecraft API",
			"version": "0.1.0",
			"status":  "running",
		})
	})

	// Inicializar handlers
	projectHandler := handlers.NewProjectHandler()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Projects
		projects := v1.Group("/projects")
		{
			projects.POST("", projectHandler.CreateProject)
			projects.GET("", projectHandler.ListProjects)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PATCH("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
			projects.POST("/:id/upload", projectHandler.UploadManuscript)
			projects.POST("/:id/process", projectHandler.ProcessProject)
			projects.GET("/:id/jobs", projectHandler.GetProjectJobs)
		}
	}

	// Iniciar servidor
	addr := fmt.Sprintf(":%d", cfg.APIPort)
	log.Printf("🚀 Servidor iniciando na porta %d", cfg.APIPort)
	log.Printf("📍 http://localhost:%d", cfg.APIPort)
	log.Printf("📍 http://localhost:%d/health", cfg.APIPort)

	// Graceful shutdown
	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Desligando servidor...")
	if err := database.Close(); err != nil {
		log.Printf("⚠️  Erro ao fechar banco: %v", err)
	}
	log.Println("✅ Servidor desligado")
}

// corsMiddleware adiciona headers CORS
func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Verificar se origin é permitido
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
