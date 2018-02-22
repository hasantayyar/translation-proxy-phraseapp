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

func (l *locales) getLocale(projectID string, localeID string, params *phraseapp.LocaleDownloadParams) ([]byte, error) {
	key, err := l.getCacheKey(localeID, params)
	if err == nil {
		locale, err := l.Cache.Get(key)
		if err != nil {
			log.Printf("error: %s", err)
		} else {
			return locale, nil
		}
	}

	locale, err := l.Client.LocaleDownload(projectID, localeID, params)
	if err != nil {
		return nil, err
	}

	if err := l.setLocale(localeID, locale, params); err != nil {
		log.Println(err)
	}

	return locale, nil
}

func (l *locales) setLocale(localeID string, locale []byte, params *phraseapp.LocaleDownloadParams) error {
	key, err := l.getCacheKey(localeID, params)
	if err != nil {
		return err
	}

	err = l.Cache.Set(key, locale)
	return err
}

func (l *locales) getCacheKey(localeID string, params *phraseapp.LocaleDownloadParams) (string, error) {
	hash, err := hashstructure.Hash(params, nil)
	if err != nil {
		log.Fatal(err)
	}

	return strconv.FormatUint(hash, 10), nil
}
