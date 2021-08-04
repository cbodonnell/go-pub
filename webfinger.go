package main

import (
	"net/url"
	"errors"
	"fmt"
	"strings"
)

func parseResource(resource string) (string, error) {
	if strings.HasPrefix(resource, "acct:") {
		resource = resource[5:]
	}

	// TODO: Can this be done more elegantly?
	// Maybe separate proto / host
	u, _ := url.Parse(config.ServerName)

	name := resource
	idx := strings.LastIndexByte(name, '/')
	if idx != -1 {
		// TODO: Redo this to account for protocol
		name = name[idx+1:]
		if fmt.Sprintf("%s/%s/%s", u.Host, config.Endpoints.Users, name) != resource {
			return name, errors.New("foreign request rejected")
		}
	} else {
		idx = strings.IndexByte(name, '@')
		if idx != -1 {
			name = name[:idx]
			if !(name+"@"+u.Host == resource) {
				return name, errors.New("foreign request rejected")
			}
		}
	}
	return name, nil
}
