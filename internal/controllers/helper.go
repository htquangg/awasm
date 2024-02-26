package controllers

import (
	"net/url"
	"strings"
)

func trimmedEndpointFromURL(url *url.URL) string {
	path := strings.TrimPrefix(url.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) == 0 {
		return "/"
	}
	// path: /api/v<1/2>/<preview/live>/<deploymentID/endpointID>/*
	return "/" + strings.Join(pathParts[4:], "/")
}
