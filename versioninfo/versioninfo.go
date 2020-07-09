package versioninfo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//BuildMeta build meta
type BuildMeta struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	BuildDate  string `json:"buildDate"`
}

//RegisterMetaEndpoint register an endpoint for build meta
func RegisterMetaEndpoint(basePath string, meta *BuildMeta) {
	http.HandleFunc(
		basePath+"/version/",
		func(writer http.ResponseWriter, request *http.Request) {
			metaBytes, _ := json.Marshal(meta)
			fmt.Fprint(writer, string(metaBytes))
		},
	)
}
