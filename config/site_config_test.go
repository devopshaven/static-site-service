package config_test

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/devopshaven/static-site-service/config"
	"github.com/stretchr/testify/assert"
)

//go:embed example.json
var exampleJsonBytes []byte

func TestParseConfig(t *testing.T) {
	// Remove comments 😂
	re := regexp.MustCompile("(?s)//.*?\n|/\\*.*?\\*/")
	newBytes := re.ReplaceAll(exampleJsonBytes, nil)

	t.Error(string(newBytes))

	var siteConfig config.SiteConfig
	json.Unmarshal(newBytes, &siteConfig)

	assert.NotEmpty(t, siteConfig.Hosting.Public)
}
