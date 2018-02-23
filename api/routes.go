package api

import (
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	"github.com/phrase/phraseapp-go/phraseapp"
)

type translationData struct {
	Client *phraseapp.Client
	Cache  *bigcache.BigCache
}

// Run translation proxy API
func Run(client *phraseapp.Client) {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
	if err != nil {
		log.Fatal(err)
	}

	t := translationData{
		Client: client,
		Cache:  cache,
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Translation Proxy is running")
	})

	api := router.Group("/api/v2")
	{
		api.GET("/projects/:project_id/locales/:id/download", t.downloadLocale)
		api.GET("/projects/:project_id/locales", t.projectLocales)
		api.GET("/projects/:project_id/translations", t.projectLocales)
	}

	router.Run(":8080")
}

func (t *translationData) downloadLocale(c *gin.Context) {
	projectID := c.Param("project_id")
	localeID := c.Param("id")

	var params downloadParams
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	locale, cached, err := t.getLocale(projectID, localeID, &params)
	if err != nil {
		log.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cached {
		c.String(http.StatusNotModified, string(locale))
	} else {
		c.String(http.StatusOK, string(locale))
	}
}

func (t *translationData) projectLocales(c *gin.Context) {
	projectID := c.Param("project_id")

	localeList, cached, err := t.getLocaleList(projectID)
	if err != nil {
		log.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cached {
		c.String(http.StatusNotModified, string(localeList))
	} else {
		c.String(http.StatusOK, string(localeList))
	}
}

func (t *translationData) listTranslations(c *gin.Context) {
	projectID := c.Param("project_id")

	var params translationsListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	localeList, cached, err := t.getTranslations(projectID, &params)
	if err != nil {
		log.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cached {
		c.String(http.StatusNotModified, string(localeList))
	} else {
		c.String(http.StatusOK, string(localeList))
	}
}
