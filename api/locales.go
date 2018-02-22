package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/allegro/bigcache"
	"github.com/phrase/phraseapp-go/phraseapp"
)

type locales struct {
	Client *phraseapp.Client
	Cache  *bigcache.BigCache
}

func (l *locales) getLocale(projectID string, localeID string, params *phraseapp.LocaleDownloadParams) ([]byte, error) {
	key := l.getCacheKey(localeID, params)
	locale, err := l.Cache.Get(key)
	if err != nil {
		log.Printf("cache error: %s", err)
	} else {
		return locale, nil
	}

	locale, err = l.Client.LocaleDownload(projectID, localeID, params)
	return locale, err
}

func (l *locales) getCacheKey(localeID string, params *phraseapp.LocaleDownloadParams) string {
	key := []byte(fmt.Sprintf("%s %v", localeID, params))
	digest := md5.New()
	digest.Write(key)
	hash := digest.Sum(nil)
	return hex.EncodeToString(hash)
}
