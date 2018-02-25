package api

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	"github.com/phrase/phraseapp-go/phraseapp"
)

type translationData struct {
	Client        *phraseapp.Client
	Cache         *bigcache.BigCache
	WebhookSecret int64
}

// Run translation proxy API
func Run(client *phraseapp.Client) {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         10 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		Verbose:            true,
		HardMaxCacheSize:   131072,
	}
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		log.Fatal(err)
	}

	secretID := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := translationData{
		Client:        client,
		Cache:         cache,
		WebhookSecret: secretID.Int63(),
	}

	fmt.Printf("Webhooks URL: /webhooks/%d\n", t.WebhookSecret)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Translation Proxy is running")
	})

	router.POST("/webhooks/:secret_id", t.webhooks)
	api := router.Group("/api/v2")
	{
		api.GET("/projects/:project_id/locales/:id/download", t.downloadLocale)
		api.GET("/projects/:project_id/locales", t.projectLocales)
		api.GET("/projects/:project_id/translations", t.listTranslations)
		// api.GET("/projects/:project_id/locales/:locale_id/translations", t.listTranslationsByLocale) // TODO find a solution for routing conflict with downloadLocale
	}

	router.Run(":8080")
}

func (t *translationData) webhooks(c *gin.Context) {
	secret := c.Param("secret_id")
	id, err := strconv.ParseInt(secret, 10, 64)
	if err != nil {
		log.Printf("error: invalid secret: '%s'", secret)
	} else if id == t.WebhookSecret {
		err = t.Cache.Reset()
		if err != nil {
			log.Printf("error: %s", err)
		}
	}
	c.String(http.StatusOK, "")
}

func (t *translationData) downloadLocale(c *gin.Context) {
	projectID := c.Param("project_id")
	localeID := c.Param("id")

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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

func (t *translationData) listTranslationsByLocale(c *gin.Context) {
	projectID := c.Param("project_id")
	localeID := c.Param("locale_id")

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var params translationsListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	localeList, cached, err := t.getTranslationsByLocale(projectID, localeID, &params)
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
