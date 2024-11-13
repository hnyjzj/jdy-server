package utils

import (
	"strings"
)

func IsWechat(agent string) bool {
	if agent == "" {
		return false
	}

	if strings.Contains(agent, "MicroMessenger") || strings.Contains(agent, "wxwork") {
		return true
	}

	return false
}
