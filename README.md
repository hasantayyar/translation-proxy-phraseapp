# translation-proxy-phraseapp

[PhraseApp](https://phraseapp.com) translation proxy to cache translations for unlimited API requests and faster responses.

## Run
The easiest way to run the translation-proxy is with Docker:

    docker run -i -t -p 8080:8080 -e PHRASEAPP_ACCESS_TOKEN=<access_token> thesoenke/translation-proxy

Without Docker the translation-proxy needs first to be built with go:

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

## Load Test
1. `go get -u github.com/tsenart/vegeta`
2. `cat vegeta/requests.txt | vegeta attack -rate=1000 -duration=30s | vegeta report`