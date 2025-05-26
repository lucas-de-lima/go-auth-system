package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Obtém o diretório de trabalho atual
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Erro ao obter o diretório de trabalho: %v", err)
	}
	log.Printf("Diretório de trabalho: %s", wd)

	// Lista de possíveis localizações do arquivo de configuração
	// A ordem é importante - tentamos do mais específico para o mais genérico
	configPaths := []string{
		filepath.Join(wd, "..", "configs", "app.env"),         // Se executado de prisma/
		filepath.Join(wd, "configs", "app.env"),               // Se executado da raiz
		filepath.Join(filepath.Dir(wd), "configs", "app.env"), // Um nível acima
		"../configs/app.env",                                  // Relativo a prisma/
		"../../configs/app.env",                               // Relativo a prisma/cmd/
	}

	// Tenta carregar o arquivo de configuração de cada localização possível
	var loadedPath string
	configLoaded := false

	for _, path := range configPaths {
		absPath, _ := filepath.Abs(path)
		log.Printf("Tentando carregar configuração de: %s", absPath)

		if _, err := os.Stat(absPath); err == nil {
			if err := godotenv.Load(absPath); err == nil {
				log.Printf("✅ Configurações carregadas com sucesso de: %s", absPath)
				configLoaded = true
				loadedPath = absPath
				break
			} else {
				log.Printf("❌ Erro ao carregar o arquivo %s: %v", absPath, err)
			}
		} else {
			log.Printf("❌ Arquivo não encontrado: %s", absPath)
		}
	}

	if !configLoaded {
		log.Println("⚠️ ATENÇÃO: Não foi possível carregar o arquivo de configuração de nenhum local conhecido.")

		// Como último recurso, tenta encontrar recursivamente na árvore de diretórios
		log.Println("🔍 Procurando o arquivo de configuração recursivamente...")
		foundPath := findConfigFile(filepath.Dir(wd), "app.env")

		if foundPath != "" {
			log.Printf("🔎 Arquivo de configuração encontrado em: %s", foundPath)
			if err := godotenv.Load(foundPath); err == nil {
				log.Printf("✅ Configurações carregadas com sucesso de: %s", foundPath)
				loadedPath = foundPath
				// Não precisamos atualizar configLoaded aqui, pois não é usado após este ponto
			} else {
				log.Printf("❌ Erro ao carregar o arquivo %s: %v", foundPath, err)
			}
		} else {
			log.Println("❌ Não foi possível encontrar o arquivo de configuração em nenhum lugar.")
		}
	}

	// Verificar se a DATABASE_URL está definida
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		// Definir uma URL padrão se não estiver definida
		defaultDbUrl := "postgresql://postgres:postgres@localhost:5432/auth_system?sslmode=disable"
		os.Setenv("DATABASE_URL", defaultDbUrl)
		log.Println("⚠️ ATENÇÃO: DATABASE_URL não encontrada no arquivo de configuração!")
		log.Printf("⚠️ Usando URL de banco de dados padrão: %s", defaultDbUrl)
	} else {
		log.Printf("✅ Usando DATABASE_URL do arquivo de configuração: %s", dbUrl)
		if loadedPath != "" {
			log.Printf("✅ Carregado de: %s", loadedPath)
		}
	}

	// Obter os argumentos para o comando Prisma
	prismaArgs := []string{"github.com/steebchen/prisma-client-go"}
	if len(os.Args) > 1 {
		prismaArgs = append(prismaArgs, os.Args[1:]...)
	}

	// Executar o comando Prisma
	cmd := exec.Command("go", append([]string{"run"}, prismaArgs...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ() // Importante: passar as variáveis de ambiente para o comando

	log.Printf("▶️ Executando: go %s", strings.Join(append([]string{"run"}, prismaArgs...), " "))

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro ao executar o comando Prisma: %v\n", err)
		os.Exit(1)
	}
}

// Função auxiliar para procurar o arquivo de configuração recursivamente
func findConfigFile(dir, fileName string) string {
	// Limite de profundidade para evitar loops infinitos
	maxDepth := 5
	currentDepth := 0

	for currentDepth < maxDepth {
		configPath := filepath.Join(dir, "configs", fileName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		// Tenta no diretório atual
		configPath = filepath.Join(dir, fileName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		// Sobe um nível
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// Chegamos à raiz do sistema de arquivos
			break
		}
		dir = parentDir
		currentDepth++
	}

	return ""
}
