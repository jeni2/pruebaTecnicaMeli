package controller

import (
	"fmt"
	"meli_golang_gin_basic_app/cmd/api/model"
	"meli_golang_gin_basic_app/cmd/api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Database struct {
	dService service.IDatabase
}

func NewDatabase(dService service.IDatabase) *Database {
	return &Database{
		dService: dService,
	}
}

func (d *Database) Persist(c *gin.Context) {
	var requestBody model.PersistRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if requestBody.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'host' es requerido"})
		return
	}
	if requestBody.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'pasword' es requerido"})
		return
	}
	if requestBody.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'username' es requerido"})
		return
	}
	if requestBody.Port == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'port' es requerido y no puede ser 0"})
		return
	}

	response := model.PersistResponse{Id: service.NewDatabase().Persist(&requestBody)}
	if response.Id == -1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Se presento un error almacenando los datos"})
	}
	c.JSON(http.StatusCreated, response)
}

func (d *Database) Scan(c *gin.Context) {
	id := c.Param("id")
	fmt.Print(id)
	service.NewDatabase().Scan(id)
	c.Status(http.StatusCreated)
}

func (d *Database) GetClassification(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, service.NewDatabase().GetClassification(id))
}

func (d *Database) GetWordList(c *gin.Context) {
	c.JSON(http.StatusOK, service.NewDatabase().GetWordList())
}

func (d *Database) AddNewWord(c *gin.Context) {
	word := c.Param("word")
	c.JSON(http.StatusOK, model.PersistResponse{Id: service.NewDatabase().AddNewWord(word)})
}
