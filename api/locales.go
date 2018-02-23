package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"

	"github.com/mitchellh/hashstructure"
	"github.com/phrase/phraseapp-go/phraseapp"
)

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

func (t *translationData) getLocaleList(projectID string) ([]byte, bool, error) {
	localesData, err := t.Cache.Get(projectID)
	if err != nil {
		log.Printf("error: %s", err)
	} else {
		return localesData, true, nil
	}

	locales, err := t.Client.LocalesList(projectID, 0, 100)
	if err != nil {
		return nil, false, err
	}

	localesData, err = json.Marshal(locales)
	if err != nil {
		return nil, false, err
	}

	if err := t.Cache.Set(projectID, localesData); err != nil {
		log.Println(err)
	}

	return localesData, false, nil
}

func (t *translationData) getLocale(projectID string, localeID string, params *downloadParams) ([]byte, bool, error) {
	key, err := t.getCacheKey(localeID, params)
	if err == nil {
		locale, err := t.Cache.Get(key)
		if err != nil {
			log.Printf("error: %s", err)
		} else {
			return locale, true, nil
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

	locale, err := t.Client.LocaleDownload(projectID, localeID, localeParams)
	if err != nil {
		return nil, false, err
	}

	if err := t.setLocale(localeID, locale, params); err != nil {
		log.Println(err)
	}

	return locale, false, nil
}

func (t *translationData) setLocale(localeID string, locale []byte, params *downloadParams) error {
	key, err := t.getCacheKey(localeID, params)
	if err != nil {
		return err
	}

	err = t.Cache.Set(key, locale)
	return err
}

func (t *translationData) getCacheKey(data string, params interface{}) (string, error) {
	hashParams, err := hashstructure.Hash(params, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}

	key := []byte(data + strconv.FormatUint(hashParams, 10))
	digest := md5.New()
	digest.Write(key)
	hash := digest.Sum(nil)
	return hex.EncodeToString(hash), nil
}
