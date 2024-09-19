package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadConfigPath(t *testing.T) {
	basePath := getBasePath()
	configPath := filepath.Join(basePath, "..", "..", "pkg", "testdata")

	tests := []struct {
		name        string
		configPath  string
		expectedCfg Config
		expectedErr error
	}{
		{
			name:       "config file exists",
			configPath: fmt.Sprintf("%s/config.yaml", configPath),
			expectedCfg: Config{
				Env: "test",
				Conf: Conf{
					TelegramToken: "test",
					PostgresUrl:   "test",
					RedisUrl:      "test",
					MongoDbUrl:    "test",
				},
			},
			expectedErr: nil,
		},
		{
			name:       "config file does not exist",
			configPath: fmt.Sprintf("%s/config_not_exist.yaml", configPath),
			expectedCfg: Config{
				Env: "test",
				Conf: Conf{
					TelegramToken: "test",
					PostgresUrl:   "test",
					RedisUrl:      "test",
					MongoDbUrl:    "test",
				},
			},
			expectedErr: fmt.Errorf("config file does not exist: %s", configPath+"/"+filepath.Join("config_not_exist.yaml")),
		},
		{
			name:       "config path is empty",
			configPath: "",
			expectedCfg: Config{
				Env: "test",
				Conf: Conf{
					TelegramToken: "test",
					PostgresUrl:   "test",
					RedisUrl:      "test",
					MongoDbUrl:    "test",
				},
			},
			expectedErr: fmt.Errorf("config file does not exist: "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadConfigPath(tt.configPath)
			if err != nil {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("LoadConfigPath() error = %v, wantErr %v", err, tt.expectedErr)
					return
				}
			}
			if cfg != nil && cfg.Env != tt.expectedCfg.Env {
				t.Errorf("LoadConfigPath() error = %v, wantErr %v", cfg.Env, tt.expectedCfg.Env)
			}
		})
	}
}

func getBasePath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}
