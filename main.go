package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Verse representa um versículo da Bíblia
type Verse struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Livro     string `json:"livro" binding:"required,len=2"`
	Capitulo  int    `json:"capitulo" binding:"required,gte=1"`
	Versiculo int    `json:"versiculo" binding:"required,gte=1"`
	Texto     string `json:"texto" binding:"required"`
}

var db *gorm.DB

func initDB() {
	dsn := "host=db user=postgres password=mysecretpassword dbname=verses sslmode=disable port=5432"
	var err error

	for i := 0; i < 10; i++ { // Tentativas de conexão
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Tentando conectar ao banco de dados... (tentativa %d)\n", i+1)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		panic("Falha ao conectar com o banco de dados: " + err.Error())
	}

	db.AutoMigrate(&Verse{})
}

func main() {
	// Inicializa o banco de dados
	initDB()

	r := gin.Default()

	// Rota para criar um versículo
	r.POST("/verses", createVerse)
	// Rota para pegar um versículo aleatório
	r.GET("/verses/random", getRandomVerse)

	// Obtém a porta da variável de ambiente ou usa a 8080 como fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicia o servidor na porta dinâmica
	r.Run(":" + port) // Escuta na porta fornecida pelo Cloud Run ou na 8080 como fallback
}

func createVerse(c *gin.Context) {
	var verse Verse
	// Bind dos dados do JSON recebidos na requisição
	if err := c.ShouldBindJSON(&verse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Criação do versículo no banco de dados
	if err := db.Create(&verse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar versículo"})
		return
	}

	c.JSON(http.StatusCreated, verse)
}

func getRandomVerse(c *gin.Context) {
	var verses []Verse
	// Busca todos os versículos no banco
	if err := db.Find(&verses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar versículos"})
		return
	}

	if len(verses) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhum versículo encontrado"})
		return
	}

	// Gera um índice aleatório para retornar um versículo aleatório
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(verses))
	randomVerse := verses[randomIndex]

	c.JSON(http.StatusOK, randomVerse)
}
