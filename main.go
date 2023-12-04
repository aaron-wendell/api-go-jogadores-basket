package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var db *gorm.DB
var err error

// A estrutura do jogador representa um jogador de basquete
type Player struct {
	ID      string `json:"id" gorm:"primary_key"`
	Name    string `json:"name"`
	Team    string `json:"team"`
	Points  int    `json:"points"`
	Assists int    `json:"assists"`
	Rebounds int  `json:"rebounds"`
}

func main() {
	// Conectar no SQLite 
	db, err = gorm.Open("sqlite3", "basketball.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	// O AutoMigrate tentará migrar automaticamente o esquema, sem necessidade de criar a tabela manualmente
	db.AutoMigrate(&Player{})

	// Inicializar Gin router
	r := gin.Default()

	// Definir endpoints da API 
	r.GET("/players", GetPlayers)
	r.GET("/players/:id", GetPlayer)
	r.POST("/players", CreatePlayer)
	r.PUT("/players/:id", UpdatePlayer)
	r.DELETE("/players/:id", DeletePlayer)

	// Configurar um canal para capturar sinais do sistema operacional
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Roda o servidor numa goroutine
	go func() {
		err := r.Run(":8080")
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Aguarda os sinais para encerrar o programa normalmente
	<-stopChan

	fmt.Println("Shutting down the program...")
}

// Retorna todos os jogadores cadastrados
func GetPlayers(c *gin.Context) {
	var players []Player
	if err := db.Find(&players).Error; err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
	} else {
		c.JSON(200, players)
	}
}

// Pesquisar pelo jogador pelo seu ID
func GetPlayer(c *gin.Context) {
	id := c.Params.ByName("id")
	var player Player
	if err := db.Where("id = ?", id).First(&player).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, player)
	}
}

// Criar novo jogador
func CreatePlayer(c *gin.Context) {
	var player Player
	c.BindJSON(&player)

	// Gerar UUID pro jogador
	player.ID = uuid.New().String()

	// Criar o jogador no banco de dados
	db.Create(&player)
	c.JSON(200, player)
}

// Atualizar as informações sobre o jogador
func UpdatePlayer(c *gin.Context) {
	id := c.Params.ByName("id")
	var player Player
	if err := db.Where("id = ?", id).First(&player).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}

	// Vincular as informações atualizadas da solicitação
	c.BindJSON(&player)

	// Salvar as informações no banco de dados
	db.Save(&player)
	c.JSON(200, player)
}

// Remover o jogador do banco de dados 
func DeletePlayer(c *gin.Context) {
	id := c.Params.ByName("id")
	var player Player
	d := db.Where("id = ?", id).Delete(&player)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
