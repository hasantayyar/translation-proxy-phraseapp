PhraseApp translation proxy to cache translations for unlimited API requests and very faster responses.

## Run
1. `go get github.com/thesoenke/phrase-proxy`
2. `export PHRASEAPP_ACCESS_TOKEN=token`
3. `translation-proxy`

## Usage
The translation proxy replicates the locale download endpoint. The request supports the same parameters as the direct PhraseApp [locale download endpoint](https://phraseapp.com/docs/api/v2/locales/#download). Only the host needs to be adapted. For example when running it locally to `http://localhost:8080/api/v2/projects/:project_id/locales/:id/download`

## Load testing
1. `go get -u github.com/tsenart/vegeta`
2. `cat vegeta/requests.txt | vegeta attack -rate=1000 -duration=30s | vegeta report`