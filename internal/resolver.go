package internal

import (
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type file struct {
	path        string
	content     []byte
	contentType string
}

type TrafficResolver struct {
	addressOfAPI          *url.URL
	prefixOfAPIInPath     string
	proxyForRequestsToAPI *httputil.ReverseProxy

	indexHTML             file
	filesWithoutIndexHTML []file
}

func NewTrafficResolver(addressOfAPI *url.URL, prefixOfAPIInPath, buildFolderPath string) *TrafficResolver {
	var indexHTML file

	var filesWithoutIndexHTML []file
	err := filepath.Walk(buildFolderPath,
		func(path string, fileInfo os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if fileInfo.IsDir() {
				return nil
			}

			f := file{
				path: path,
			}

			fileData, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contentType := mime.TypeByExtension(filepath.Ext(path))
			if contentType == "" {
				contentType = http.DetectContentType(fileData)
			}

			f.content = fileData
			f.contentType = contentType

			if strings.HasSuffix(f.path, "index.html") {
				indexHTML = f
				return nil
			}

			filesWithoutIndexHTML = append(filesWithoutIndexHTML, f)
			return nil
		})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return &TrafficResolver{
		addressOfAPI:          addressOfAPI,
		prefixOfAPIInPath:     prefixOfAPIInPath,
		proxyForRequestsToAPI: httputil.NewSingleHostReverseProxy(addressOfAPI),

		indexHTML:             indexHTML,
		filesWithoutIndexHTML: filesWithoutIndexHTML,
	}
}

func (tr *TrafficResolver) Resolve(res http.ResponseWriter, req *http.Request) {
	switch {
	case strings.HasPrefix(req.URL.Path, tr.prefixOfAPIInPath):
		req.URL.Host = tr.addressOfAPI.Host
		req.URL.Scheme = tr.addressOfAPI.Scheme
		req.Host = tr.addressOfAPI.Host

		tr.proxyForRequestsToAPI.ServeHTTP(res, req)
	default:
		f := tr.determineWhichFileToReturn(req)

		res.Header().Set("Content-Type", f.contentType)
		if _, err := res.Write(f.content); err != nil {
			log.Info().Msg("can't write response")
		}
	}
}

func (tr *TrafficResolver) determineWhichFileToReturn(req *http.Request) *file {
	if req.URL.Path != "/" && req.URL.Path != "/index.html" {
		for _, f := range tr.filesWithoutIndexHTML {
			if strings.HasSuffix(f.path, req.URL.Path) {
				return &f
			}
		}
	}

	return &tr.indexHTML
}
