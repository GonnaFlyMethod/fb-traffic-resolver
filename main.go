package main

import (
	"context"
	"fmt"
	"github.com/GonnaFlyMethod/fb-traffic-resolver/internal"
	"net/http"
	"net/url"
	"os"
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

	for i := 0; i < 5; i++ {
		log.Info().Msgf("trying to ping %s, attempt %d", apiAddress, i+1)

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Warn().Err(err).Send()

			time.Sleep(3 * time.Second)
			continue
		}

		if err = response.Body.Close(); err != nil {
			log.Fatal().Msg("can't close response body")
		}

		log.Info().Msg("API is alive")
		return
	}

	log.Fatal().Msgf("can't connect to %s", apiAddress)
}

func main() {
	addressOfAPIStr := mustGetEnv("ADDRESS_OF_API")

	addressOfAPIConvertedToURL, err := url.ParseRequestURI(addressOfAPIStr)
	if err != nil {
		log.Fatal().Msg("api host should be composed according to the template: <scheme>://host, example: http://myapi.com")
	}

	pingAPIOnStartStr := mustGetEnv("PING_API_ON_START")

	pingAPIOnStartConvertedToBool, err := strconv.ParseBool(pingAPIOnStartStr)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if pingAPIOnStartConvertedToBool {
		pingAPI(addressOfAPIConvertedToURL.String())
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
