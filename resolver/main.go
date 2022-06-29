package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sccsmanager/resolver/internal"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func mustGetEnv(envName string) string {
	res := os.Getenv(envName)
	if res == "" {
		log.Fatal().Msgf("can't read env %s", envName)
	}
	return res
}

func pingAPI(apiAddress string) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiAddress, http.NoBody)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	const timeToSleep = 3 * time.Second

	for i := 0; i < 5; i++ {
		log.Info().Msgf("trying to ping %s, attempt %d", apiAddress, i+1)

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Warn().Err(err).Send()

			time.Sleep(timeToSleep)
			continue
		}

		if err = response.Body.Close(); err != nil {
			log.Fatal().Msg("can't close response body")
		}

		if response.StatusCode == http.StatusOK {
			log.Info().Msg("API is alive")
			return
		}

		log.Warn().Msgf("response status code from api is not 2xx, actual response status code: %d", response.StatusCode)
		time.Sleep(timeToSleep)
	}

	log.Fatal().Msgf("can't connect to %s", apiAddress)
}

func main() {
	pingAPIOnStartStr := mustGetEnv("PING_API_ON_START")
	pingAPIOnStartConvertedToBool, err := strconv.ParseBool(pingAPIOnStartStr)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	addressOfAPIStr := mustGetEnv("ADDRESS_OF_API")
	if pingAPIOnStartConvertedToBool {
		pingAPI(addressOfAPIStr)
	}

	addressOfAPIConvertedToURL, err := url.Parse(addressOfAPIStr)
	if err != nil {
		log.Fatal().Msg("api host should be composed according to the template: <scheme>://host, example: http://myapi.com")
	}

	const (
		prefixOfAPIInPath = "/api"
		pathToBuildFolder = "./build"
	)
	tr := internal.NewTrafficResolver(addressOfAPIConvertedToURL, prefixOfAPIInPath, pathToBuildFolder)

	http.HandleFunc("/", tr.Resolve)

	resolverPort := mustGetEnv("RESOLVER_PORT")
	resolverAddress := fmt.Sprintf(":%s", resolverPort)

	startMsg := fmt.Sprintf("starting resolver on %s", resolverAddress)
	log.Info().Msg(startMsg)

	err = http.ListenAndServe(resolverAddress, nil)
	log.Fatal().Err(err).Send()
}
