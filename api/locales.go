package api

import (
	"log"
	"strconv"

	"github.com/allegro/bigcache"
	"github.com/mitchellh/hashstructure"
	"github.com/phrase/phraseapp-go/phraseapp"
)

type locales struct {
	Client *phraseapp.Client
	Cache  *bigcache.BigCache
}

type downloadParams struct {
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

func (l *locales) getLocale(projectID string, localeID string, params *downloadParams) ([]byte, error) {
	key, err := l.getCacheKey(localeID, params)
	if err == nil {
		locale, err := l.Cache.Get(key)
		if err != nil {
			log.Printf("error: %s", err)
		} else {
			return locale, nil
		}
	}

	localeParams := &phraseapp.LocaleDownloadParams{
		ConvertEmoji:               params.ConvertEmoji,
		Encoding:                   &params.Encoding,
		FallbackLocaleID:           &params.FallbackLocaleID,
		FileFormat:                 &params.FileFormat,
		IncludeEmptyTranslations:   params.IncludeEmptyTranslations,
		KeepNotranslateTags:        params.KeepNotranslateTags,
		SkipUnverifiedTranslations: params.SkipUnverifiedTranslations,
		Tag: &params.Tag,
	}
	locale, err := l.Client.LocaleDownload(projectID, localeID, localeParams)
	if err != nil {
		return nil, err
	}

	if err := l.setLocale(localeID, locale, params); err != nil {
		log.Println(err)
	}

	return locale, nil
}

func (l *locales) setLocale(localeID string, locale []byte, params *downloadParams) error {
	key, err := l.getCacheKey(localeID, params)
	if err != nil {
		return err
	}

	err = l.Cache.Set(key, locale)
	return err
}

func (l *locales) getCacheKey(localeID string, params *downloadParams) (string, error) {
	hash, err := hashstructure.Hash(params, nil)
	if err != nil {
		log.Fatal(err)
	}

	return strconv.FormatUint(hash, 10), nil
}
