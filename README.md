PhraseApp translation proxy to cache translations for unlimited API requests and faster responses.

## Run
1. `go get github.com/thesoenke/phrase-proxy`
2. `export PHRASEAPP_ACCESS_TOKEN=token`
3. `translation-proxy`

## Usage
The translation proxy replicates the locale download endpoint. The request supports the same parameters as the direct PhraseApp [locale download endpoint](https://phraseapp.com/docs/api/v2/locales/#download). Only the host needs to be adapted. For example when running it locally to `http://localhost:8080/api/v2/projects/:project_id/locales/:id/download`

## Run
The easiest way to run the translation-proxy is with Docker:

    docker run -i -t -p 8080:8080 -e PHRASEAPP_ACCESS_TOKEN=<access_token> thesoenke/translation-proxy

Without Docker the translation-proxy needs first to be build with go:

    go get
    go build
    ./translation-proxy

## Supported Endpoints
The following endpoints are currently supported:
- [List Locales](https://phraseapp.com/docs/api/v2/locales/#index) `GET /v2/projects/:project_id/locales`
- [Download Locales](https://phraseapp.com/docs/api/v2/locales/#download) `GET /v2/projects/:project_id/locales/:id/download`
- [List Translations](https://phraseapp.com/docs/api/v2/translations/#index) `GET /v2/projects/:project_id/translations
`

## Load testing
1. `go get -u github.com/tsenart/vegeta`
2. `cat vegeta/requests.txt | vegeta attack -rate=1000 -duration=30s | vegeta report`