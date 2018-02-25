# translation-proxy-phraseapp

[PhraseApp](https://phraseapp.com) translation proxy to cache translations for unlimited API requests and faster responses.

## Run
### Docker

    docker build -t translation-proxy-phraseapp .
    docker run -it -p 8080:8080 -e PHRASEAPP_ACCESS_TOKEN=<access_token> translation-proxy-phraseapp

### Build from source

    go get
    go build
    export PHRASEAPP_ACCESS_TOKEN=<access_token>
    ./translation-proxy

## Supported Endpoints
The translation proxy replicates the API from PhraseApp of multiple GET endpoints. The following endpoints are supported:
- [List Locales](https://phraseapp.com/docs/api/v2/locales/#index) `GET /v2/projects/:project_id/locales`
- [Download Locales](https://phraseapp.com/docs/api/v2/locales/#download) `GET /v2/projects/:project_id/locales/:id/download`
- [List Translations](https://phraseapp.com/docs/api/v2/translations/#index) `GET /v2/projects/:project_id/translations
`

### Webhooks
By default the translations will be cached for 10 minutes. To immediately clear the cache and get the newest translations the translation proxy supports webhooks.

The translation proxy will print the webhooks URL containing a secret on start: `/webhooks/:secret`. The secret will change after every restart. The webhooks URL then has to be added in the project on PhraseApp with events that should trigger a cache reset like: `translations:update`

## Load Test
1. `go get -u github.com/tsenart/vegeta`
2. `cat vegeta/requests.txt | vegeta attack -rate=1000 -duration=30s | vegeta report`