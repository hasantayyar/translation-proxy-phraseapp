package api

import (
	"encoding/json"
	"log"

	"github.com/phrase/phraseapp-go/phraseapp"
)

type translationsListParams struct {
	Order *string `form:"order,omitempty"`
	Q     *string `form:"q,omitempty"`
	Sort  *string `form:"sort,omitempty"`
}

func (t *translationData) getTranslations(projectID string, params *translationsListParams) ([]byte, bool, error) {
	key, err := t.getCacheKey(projectID, params)
	if err == nil {
		locale, err := t.Cache.Get(key)
		if err != nil {
			log.Printf("error: %s", err)
		} else {
			return locale, true, nil
		}
	}

	translations, err := t.Client.TranslationsList(projectID, 0, 1000, &phraseapp.TranslationsListParams{
		Order: params.Order,
		Q:     params.Q,
		Sort:  params.Sort,
	})
	if err != nil {
		return nil, false, err
	}

	translationMarshal, err := json.Marshal(translations)
	if err != nil {
		return nil, false, err
	}

	if err := t.setTranslations(projectID, translationMarshal, params); err != nil {
		log.Println(err)
	}

	return translationMarshal, false, nil
}

func (t *translationData) setTranslations(projectID string, data []byte, params *translationsListParams) error {
	key, err := t.getCacheKey(projectID, params)
	if err != nil {
		return err
	}

	err = t.Cache.Set(key, data)
	return err
}
