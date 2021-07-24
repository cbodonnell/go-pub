package main

import (
	"errors"
	"fmt"
	"strings"
)

func parseResource(resource string) (string, error) {
	if strings.HasPrefix(resource, "acct:") {
		resource = resource[5:]
	}

	name := resource
	idx := strings.LastIndexByte(name, '/')
	if idx != -1 {
		name = name[idx+1:]
		if fmt.Sprintf("%s/%s/%s", config.ServerName, config.Endpoints.Users, name) != resource {
			return name, errors.New("foreign request rejected")
		}
	} else {
		idx = strings.IndexByte(name, '@')
		if idx != -1 {
			name = name[:idx]
			if !(name+"@"+config.ServerName == resource) {
				return name, errors.New("foreign request rejected")
			}
		}
	}
	return name, nil
}
