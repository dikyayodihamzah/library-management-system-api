package utils

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	vault "github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"

	_ "github.com/joho/godotenv/autoload"
)

var (
	vaultAddr   = os.Getenv("VAULT_ADDR")
	vaultAuth   = os.Getenv("VAULT_AUTH")
	vaultToken  = os.Getenv("VAULT_TOKEN")
	vaultEngine = os.Getenv("VAULT_ENGINE")
	vaultPath   = os.Getenv("VAULT_PATH")

	envMap = make(map[string]interface{})
)

func LoadEnv(defEnv ...map[string]interface{}) {
	vaultVals := LoadEnvVault()
	dotVals := LoadDotEnv()
	osEnv := LoadEnvFromOS()

	// merge env vars
	if len(defEnv) > 0 {
		for key, val := range defEnv[0] {
			envMap[strings.ToUpper(key)] = val
		}
	}

	for key, val := range osEnv {
		envMap[key] = val
	}

	for key, val := range dotVals {
		envMap[key] = val
	}

	for key, val := range vaultVals {
		envMap[key] = val
	}
}

func LoadDotEnv() map[string]interface{} {
	var dotEnv = make(map[string]interface{})
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil
	}

	for _, key := range os.Environ() {
		kv := strings.Split(key, "=")
		dotEnv[kv[0]] = kv[1]
	}

	return dotEnv
}

func LoadEnvVault() map[string]interface{} {
	var vaultURL string

	if vaultAddr == "" || vaultAuth == "" || vaultToken == "" || vaultPath == "" {
		return nil
	} else {
		vaultURL = vaultAddr

		resp, err := http.Get(vaultURL)
		if err != nil {
			return nil
		}

		defer resp.Body.Close()
	}

	config := vault.DefaultConfig()
	config.Address = vaultURL

	client, err := vault.NewClient(config)
	if err != nil {
		return nil
	}
	client.SetToken(vaultToken)

	secret, err := client.KVv2(vaultEngine).Get(context.Background(), vaultPath)
	if err != nil {
		return nil
	}

	return secret.Data
}

// load env for unit test
func LoadEnvTest(t *testing.T) {
	vaultVals := LoadEnvVault()
	dotVals := LoadDotEnv()

	for key, val := range dotVals {
		t.Setenv(key, val.(string))
	}

	for key, val := range vaultVals {
		t.Setenv(key, val.(string))
	}
}

func LoadEnvFromOS() map[string]interface{} {
	osEnv := make(map[string]interface{})
	for _, key := range os.Environ() {
		kv := strings.SplitN(key, "=", 2)
		osEnv[kv[0]] = kv[1]
	}
	return osEnv
}
