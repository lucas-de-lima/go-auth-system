package logging

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	once          sync.Once
)

// Config contém as configurações do logger
type Config struct {
	InfoWriter    io.Writer
	WarningWriter io.Writer
	ErrorWriter   io.Writer
	Prefix        string
	Flag          int
}

// DefaultConfig retorna a configuração padrão para o logger
func DefaultConfig() Config {
	return Config{
		InfoWriter:    os.Stdout,
		WarningWriter: os.Stdout,
		ErrorWriter:   os.Stderr,
		Prefix:        "",
		Flag:          log.LstdFlags | log.Lshortfile,
	}
}

// SetupLogger configura os loggers com a configuração fornecida
func SetupLogger(config Config) {
	once.Do(func() {
		infoLogger = log.New(config.InfoWriter, config.Prefix+"INFO: ", config.Flag)
		warningLogger = log.New(config.WarningWriter, config.Prefix+"WARNING: ", config.Flag)
		errorLogger = log.New(config.ErrorWriter, config.Prefix+"ERROR: ", config.Flag)
	})
}

// Info registra uma mensagem de informação
func Info(format string, v ...interface{}) {
	setupIfNeeded()
	infoLogger.Printf(format, v...)
}

// Warning registra uma mensagem de aviso
func Warning(format string, v ...interface{}) {
	setupIfNeeded()
	warningLogger.Printf(format, v...)
}

// Error registra uma mensagem de erro
func Error(format string, v ...interface{}) {
	setupIfNeeded()
	errorLogger.Printf(format, v...)
}

// Fatal registra uma mensagem de erro e encerra o programa
func Fatal(format string, v ...interface{}) {
	setupIfNeeded()
	errorLogger.Fatalf(format, v...)
}

// setupIfNeeded configura os loggers com a configuração padrão se ainda não foram configurados
func setupIfNeeded() {
	once.Do(func() {
		config := DefaultConfig()
		infoLogger = log.New(config.InfoWriter, config.Prefix+"INFO: ", config.Flag)
		warningLogger = log.New(config.WarningWriter, config.Prefix+"WARNING: ", config.Flag)
		errorLogger = log.New(config.ErrorWriter, config.Prefix+"ERROR: ", config.Flag)
	})
}
