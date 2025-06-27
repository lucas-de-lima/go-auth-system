package logging

import (
	"bytes"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.InfoWriter != os.Stdout {
		t.Errorf("InfoWriter deveria ser os.Stdout, mas foi %v", config.InfoWriter)
	}

	if config.WarningWriter != os.Stdout {
		t.Errorf("WarningWriter deveria ser os.Stdout, mas foi %v", config.WarningWriter)
	}

	if config.ErrorWriter != os.Stderr {
		t.Errorf("ErrorWriter deveria ser os.Stderr, mas foi %v", config.ErrorWriter)
	}

	if config.Prefix != "" {
		t.Errorf("Prefix deveria ser vazio, mas foi %s", config.Prefix)
	}

	expectedFlag := log.LstdFlags | log.Lshortfile
	if config.Flag != expectedFlag {
		t.Errorf("Flag deveria ser %d, mas foi %d", expectedFlag, config.Flag)
	}
}

func TestSetupLogger(t *testing.T) {
	// Reset global variables
	infoLogger = nil
	warningLogger = nil
	errorLogger = nil
	once = sync.Once{}

	var infoBuf, warningBuf, errorBuf bytes.Buffer

	config := Config{
		InfoWriter:    &infoBuf,
		WarningWriter: &warningBuf,
		ErrorWriter:   &errorBuf,
		Prefix:        "TEST: ",
		Flag:          log.LstdFlags,
	}

	SetupLogger(config)

	// Verifica se os loggers foram configurados
	if infoLogger == nil {
		t.Error("infoLogger não deveria ser nil")
	}

	if warningLogger == nil {
		t.Error("warningLogger não deveria ser nil")
	}

	if errorLogger == nil {
		t.Error("errorLogger não deveria ser nil")
	}

	// Testa se a configuração foi aplicada corretamente
	Info("teste info")
	Warning("teste warning")
	Error("teste error")

	infoOutput := infoBuf.String()
	warningOutput := warningBuf.String()
	errorOutput := errorBuf.String()

	if !strings.Contains(infoOutput, "TEST: INFO: ") {
		t.Errorf("Output de info deveria conter 'TEST: INFO: ', mas foi: %s", infoOutput)
	}

	if !strings.Contains(warningOutput, "TEST: WARNING: ") {
		t.Errorf("Output de warning deveria conter 'TEST: WARNING: ', mas foi: %s", warningOutput)
	}

	if !strings.Contains(errorOutput, "TEST: ERROR: ") {
		t.Errorf("Output de error deveria conter 'TEST: ERROR: ', mas foi: %s", errorOutput)
	}

	if !strings.Contains(infoOutput, "teste info") {
		t.Errorf("Output de info deveria conter 'teste info', mas foi: %s", infoOutput)
	}
}

func TestInfo(t *testing.T) {
	// Reset global variables
	infoLogger = nil
	warningLogger = nil
	errorLogger = nil
	once = sync.Once{}

	var buf bytes.Buffer
	config := Config{
		InfoWriter:    &buf,
		WarningWriter: &buf,
		ErrorWriter:   &buf,
		Prefix:        "",
		Flag:          log.LstdFlags,
	}

	SetupLogger(config)

	Info("teste info %s", "mensagem")
	output := buf.String()

	if !strings.Contains(output, "INFO: ") {
		t.Errorf("Output deveria conter 'INFO: ', mas foi: %s", output)
	}

	if !strings.Contains(output, "teste info mensagem") {
		t.Errorf("Output deveria conter 'teste info mensagem', mas foi: %s", output)
	}
}

func TestWarning(t *testing.T) {
	// Reset global variables
	infoLogger = nil
	warningLogger = nil
	errorLogger = nil
	once = sync.Once{}

	var buf bytes.Buffer
	config := Config{
		InfoWriter:    &buf,
		WarningWriter: &buf,
		ErrorWriter:   &buf,
		Prefix:        "",
		Flag:          log.LstdFlags,
	}

	SetupLogger(config)

	Warning("teste warning %d", 123)
	output := buf.String()

	if !strings.Contains(output, "WARNING: ") {
		t.Errorf("Output deveria conter 'WARNING: ', mas foi: %s", output)
	}

	if !strings.Contains(output, "teste warning 123") {
		t.Errorf("Output deveria conter 'teste warning 123', mas foi: %s", output)
	}
}

func TestError(t *testing.T) {
	// Reset global variables
	infoLogger = nil
	warningLogger = nil
	errorLogger = nil
	once = sync.Once{}

	var buf bytes.Buffer
	config := Config{
		InfoWriter:    &buf,
		WarningWriter: &buf,
		ErrorWriter:   &buf,
		Prefix:        "",
		Flag:          log.LstdFlags,
	}

	SetupLogger(config)

	Error("teste error %v", "erro")
	output := buf.String()

	if !strings.Contains(output, "ERROR: ") {
		t.Errorf("Output deveria conter 'ERROR: ', mas foi: %s", output)
	}

	if !strings.Contains(output, "teste error erro") {
		t.Errorf("Output deveria conter 'teste error erro', mas foi: %s", output)
	}
}

func TestFatal(t *testing.T) {
	// Reset global variables
	infoLogger = nil
	warningLogger = nil
	errorLogger = nil
	once = sync.Once{}

	var buf bytes.Buffer
	config := Config{
		InfoWriter:    &buf,
		WarningWriter: &buf,
		ErrorWriter:   &buf,
		Prefix:        "",
		Flag:          log.LstdFlags,
	}

	SetupLogger(config)

	// Fatal chama os.Exit(1), então não podemos testar diretamente
	// Mas podemos verificar se o logger foi configurado corretamente
	if errorLogger == nil {
		t.Error("errorLogger deveria estar configurado")
	}
}

func TestSetupIfNeeded(t *testing.T) {
	// Reset global variables
	infoLogger = nil
	warningLogger = nil
	errorLogger = nil
	once = sync.Once{}

	// Chama Info sem configurar o logger primeiro
	// Isso deve chamar setupIfNeeded automaticamente
	Info("teste setup automático")

	// Verifica se os loggers foram configurados com valores padrão
	if infoLogger == nil {
		t.Error("infoLogger deveria ter sido configurado automaticamente")
	}

	if warningLogger == nil {
		t.Error("warningLogger deveria ter sido configurado automaticamente")
	}

	if errorLogger == nil {
		t.Error("errorLogger deveria ter sido configurado automaticamente")
	}
}

func TestMultipleSetupCalls(t *testing.T) {
	// Reset global variables
	infoLogger = nil
	warningLogger = nil
	errorLogger = nil
	once = sync.Once{}

	var buf1, buf2 bytes.Buffer

	config1 := Config{
		InfoWriter:    &buf1,
		WarningWriter: &buf1,
		ErrorWriter:   &buf1,
		Prefix:        "FIRST: ",
		Flag:          log.LstdFlags,
	}

	config2 := Config{
		InfoWriter:    &buf2,
		WarningWriter: &buf2,
		ErrorWriter:   &buf2,
		Prefix:        "SECOND: ",
		Flag:          log.LstdFlags,
	}

	// Primeira configuração
	SetupLogger(config1)
	Info("primeira mensagem")

	// Segunda configuração (não deve sobrescrever a primeira)
	SetupLogger(config2)
	Info("segunda mensagem")

	// Verifica se apenas a primeira configuração foi usada
	output1 := buf1.String()
	output2 := buf2.String()

	if !strings.Contains(output1, "FIRST: INFO: ") {
		t.Errorf("Primeira configuração deveria ter sido usada: %s", output1)
	}

	if !strings.Contains(output1, "primeira mensagem") {
		t.Errorf("Primeira mensagem deveria estar na primeira configuração: %s", output1)
	}

	if !strings.Contains(output1, "segunda mensagem") {
		t.Errorf("Segunda mensagem deveria estar na primeira configuração: %s", output1)
	}

	if output2 != "" {
		t.Errorf("Segunda configuração não deveria ter sido usada: %s", output2)
	}
}
