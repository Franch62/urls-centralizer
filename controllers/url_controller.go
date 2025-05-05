package controllers

import (
	"io"
	"net/http"
	"urls-centralizer/config"
	"urls-centralizer/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func CreateURL(c *gin.Context) {
	var url models.URL
	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, url)
}

func GetAllURLs(c *gin.Context) {
	var urls []models.URL
	config.DB.Find(&urls)
	c.JSON(http.StatusOK, urls)
}

func FetchYAMLFromURL(c *gin.Context) {
	var url models.URL
	id := c.Param("id")

	if err := config.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL não encontrada"})
		return
	}

	resp, err := http.Get(url.URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Erro ao acessar a URL remota"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler conteúdo da resposta"})
		return
	}

	// Retorna o conteúdo como texto (YAML bruto)
	c.Data(http.StatusOK, "text/yaml", body)
}

func UpdateURL(c *gin.Context) {
	id := c.Param("id")
	var existing models.URL

	if err := config.DB.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL não encontrada"})
		return
	}

	var updatedData models.URL
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing.Source = updatedData.Source
	existing.URL = updatedData.URL

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func DeleteURL(c *gin.Context) {
	id := c.Param("id")
	var url models.URL

	if err := config.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL não encontrada"})
		return
	}

	if err := config.DB.Delete(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar a URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deletada com sucesso"})
}

func ServeSwaggerUI(c *gin.Context) {
	id := c.Param("id")
	var url models.URL

	if err := config.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Redireciona para o editor online usando sua própria API como origem do YAML
	redirectURL := "https://editor.swagger.io/?url=http://localhost:8080/api/urls/" + id + "/fetch"
	c.Redirect(http.StatusFound, redirectURL)
}

// ProxyYAML busca o conteúdo de uma URL .yaml e retorna com headers de CORS
// ProxyYAML busca o conteúdo de uma URL .yaml e retorna com headers de CORS
func ProxyYAML(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	// Faz o download do YAML da URL fornecida
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch the yaml file"})
		return
	}
	defer resp.Body.Close()

	// Adiciona os headers para permitir o acesso via CORS
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "text/yaml")
	c.Status(http.StatusOK)
	io.Copy(c.Writer, resp.Body)
}

func GetURLEndpoints(c *gin.Context) {
	id := c.Param("id")

	// Busca a URL no banco
	var url models.URL
	if err := config.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL não encontrada"})
		return
	}

	// Busca o YAML remoto
	resp, err := http.Get(url.URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Erro ao acessar a URL remota"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler o conteúdo da resposta"})
		return
	}

	// Estrutura para decodificar apenas a parte de "paths"
	var parsed struct {
		Paths map[string]interface{} `yaml:"paths"`
	}

	if err := yaml.Unmarshal(body, &parsed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao parsear o YAML"})
		return
	}

	// Extrai os caminhos (endpoints)
	endpoints := make([]string, 0, len(parsed.Paths))
	for path := range parsed.Paths {
		endpoints = append(endpoints, path)
	}

	c.JSON(http.StatusOK, gin.H{"endpoints": endpoints})
}
