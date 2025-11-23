package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigurationInit(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{
			name:     "successful config initialization with local.toml",
			fileName: "local",
			wantErr:  false,
		},
		{
			name:     "config file not found",
			fileName: "nonexistent",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := Init(tt.fileName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.IsType(t, &viper.Viper{}, config)
			}
		})
	}
}

func TestConfigurationValues(t *testing.T) {
	config, err := Init("local")
	require.NoError(t, err)
	require.NotNil(t, config)

	// Проверяем значения из вашего local.toml
	token := config.GetString("telegram.token")
	maintainerID := config.GetInt64("telegram.maintainer_id")

	assert.NotEmpty(t, token)
	assert.True(t, maintainerID > 0)
}

func TestEnvironmentVariablesOverride(t *testing.T) {
	// Сохраняем оригинальные значения
	originalToken := os.Getenv("DV_TELEGRAM_TOKEN")
	originalMaintainerID := os.Getenv("DV_TELEGRAM_MAINTAINER_ID")

	// Устанавливаем тестовые переменные окружения
	err := os.Setenv("DV_TELEGRAM_TOKEN", "env_test_token_789")
	require.NoError(t, err)
	err = os.Setenv("DV_TELEGRAM_MAINTAINER_ID", "66778899")
	require.NoError(t, err)

	defer func() {
		// Восстанавливаем оригинальные значения
		if originalToken != "" {
			err = os.Setenv("DV_TELEGRAM_TOKEN", originalToken)
			require.NoError(t, err)
		} else {
			err = os.Unsetenv("DV_TELEGRAM_TOKEN")
			require.NoError(t, err)
		}
		if originalMaintainerID != "" {
			err = os.Setenv("DV_TELEGRAM_MAINTAINER_ID", originalMaintainerID)
			require.NoError(t, err)
		} else {
			err = os.Unsetenv("DV_TELEGRAM_MAINTAINER_ID")
			require.NoError(t, err)
		}
	}()

	config, err := Init("local")
	require.NoError(t, err)
	require.NotNil(t, config)

	// Переменные окружения должны иметь приоритет
	assert.Equal(t, "env_test_token_789", config.GetString("telegram.token"))
	assert.Equal(t, int64(66778899), config.GetInt64("telegram.maintainer_id"))
}

func TestConfigStructure(t *testing.T) {
	config, err := Init("local")
	require.NoError(t, err)
	require.NotNil(t, config)

	// Проверяем структуру конфига
	assert.True(t, config.IsSet("telegram"))
	assert.True(t, config.IsSet("telegram.token"))
	assert.True(t, config.IsSet("telegram.maintainer_id"))

	// Проверяем, что значения не пустые
	token := config.GetString("telegram.token")
	maintainerID := config.GetInt64("telegram.maintainer_id")

	assert.NotEmpty(t, token)
	assert.True(t, maintainerID > 0)
}

func TestMultipleConfigCalls(t *testing.T) {
	// Проверяем, что функция работает идемпотентно
	config1, err1 := Init("local")
	config2, err2 := Init("local")

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NotNil(t, config1)
	require.NotNil(t, config2)

	// Оба конфига должны иметь одинаковые значения
	assert.Equal(t,
		config1.GetString("telegram.token"),
		config2.GetString("telegram.token"))
	assert.Equal(t,
		config1.GetInt64("telegram.maintainer_id"),
		config2.GetInt64("telegram.maintainer_id"))
}

func TestConfigWithDifferentEnvironment(t *testing.T) {
	// Тестируем загрузку разных конфигов (если у вас есть production.toml)
	// Можно создать временный тестовый конфиг для этого теста

	t.Run("production config", func(t *testing.T) {
		// Если у вас есть production.toml, раскомментируйте:
		config, err := Init("production")
		if err != nil {
			t.Skip("Production config not available")
			return
		}
		require.NotNil(t, config)
	})
}
