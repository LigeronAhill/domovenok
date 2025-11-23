// Package config предоставляет функционал для загрузки и управления конфигурацией приложения.
// Поддерживает загрузку настроек из TOML файлов и переменных окружения с автоматическим
// определением корневой директории проекта.
package config

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

// Init инициализирует и загружает конфигурацию из указанного файла.
//
// Функция автоматически находит корневую директорию проекта (ищущую папку "domovenok")
// и загружает конфигурацию из папки config/ в корне проекта.
//
// Параметры:
//   - fileName: имя конфигурационного файла без расширения (например, "local", "production")
//
// Возвращает:
//   - *viper.Viper: объект конфигурации с загруженными настройками
//   - error: ошибка, если не удалось найти или загрузить конфигурацию
//
// Особенности:
//   - Поддерживает TOML формат конфигурационных файлов
//   - Переменные окружения имеют приоритет над значениями из файла
//   - Префикс переменных окружения: "DV_" (например, DV_TELEGRAM_TOKEN)
//   - Точки в именах настроек заменяются на подчеркивания (telegram.token -> TELEGRAM_TOKEN)
//
// Пример использования:
//
//	config, err := config.Init("local")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	token := config.GetString("telegram.token")
func Init(fileName string) (*viper.Viper, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for {
		if !strings.HasSuffix(currentDir, "domovenok") {
			currentDir = path.Dir(currentDir)
		} else {
			break
		}
	}
	configPath := path.Join(currentDir, "config")
	config := viper.New()
	config.SetConfigName(fileName)
	config.AddConfigPath(configPath)
	config.SetEnvPrefix("dv")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
	err = config.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}
