package api

import (
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	"github.com/phrase/phraseapp-go/phraseapp"
)

// Run translation proxy API
func Run(client *phraseapp.Client) {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
	if err != nil {
		log.Fatal(err)
	}

	l := locales{
		Client: client,
		Cache:  cache,
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Translation Proxy is running")
	})

	api := router.Group("/api/v2")
	{
		api.GET("/projects/:project_id/locales/:id/download", l.downloadLocale)
	}

	router.Run(":8080")
}

func (l *locales) downloadLocale(c *gin.Context) {
	projectID := c.Param("project_id")
	localeID := c.Param("id")

	type Params struct {
		ConvertEmoji               bool              `form:"convert_emoji,omitempty"`
		Encoding                   string            `form:"encoding,omitempty"`
		FallbackLocaleID           string            `form:"fallback_locale_id,omitempty"`
		FileFormat                 string            `form:"file_format" binding:"required"`
		FormatOptions              map[string]string `form:"format_options,omitempty"`
		IncludeEmptyTranslations   bool              `form:"include_empty_translations,omitempty"`
		KeepNotranslateTags        bool              `form:"keep_notranslate_tags,omitempty"`
		SkipUnverifiedTranslations bool              `form:"skip_unverified_translations,omitempty"`
		Tag                        string            `form:"tag,omitempty"`
	}

	var params Params
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	localeParams := phraseapp.LocaleDownloadParams{
		ConvertEmoji:               params.ConvertEmoji,
		Encoding:                   &params.Encoding,
		FallbackLocaleID:           &params.FallbackLocaleID,
		FileFormat:                 &params.FileFormat,
		IncludeEmptyTranslations:   params.IncludeEmptyTranslations,
		KeepNotranslateTags:        params.KeepNotranslateTags,
		SkipUnverifiedTranslations: params.SkipUnverifiedTranslations,
		Tag: &params.Tag,
	}

	locale, err := l.getLocale(projectID, localeID, &localeParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, string(locale))
}
