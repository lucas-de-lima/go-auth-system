package prisma

import (
	"context"
	"log"

	"github.com/lucas-de-lima/go-auth-system/prisma/db"
)

// DB é uma instância compartilhada do cliente Prisma
var DB *db.PrismaClient

// Init inicializa a conexão com o banco de dados
func Init() {
	DB = db.NewClient()

	if err := DB.Prisma.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
}

// Disconnect fecha a conexão com o banco de dados
func Disconnect() {
	if err := DB.Prisma.Disconnect(); err != nil {
		log.Printf("Erro ao desconectar do banco de dados: %v", err)
	}
}

// Ping verifica se a conexão com o banco de dados está funcionando
func Ping(ctx context.Context) error {
	_, err := DB.User.FindMany().Exec(ctx)
	return err
}
