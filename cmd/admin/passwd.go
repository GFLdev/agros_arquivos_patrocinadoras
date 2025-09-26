package main

import (
	"agros_arquivos_patrocinadoras/pkg/app"
	"agros_arquivos_patrocinadoras/pkg/app/config"
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/db"
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
	"time"
)

func getAdmin(ctx *context.Context) (uuid.UUID, error) {
	// Verificar se usuário de administrador está no banco
	schema := ctx.Config.Database.Schema
	query := fmt.Sprintf(
		`SELECT %s FROM %s.%s WHERE %s = :username`,
		schema.UserTable.Columns.UserId,
		schema.Name,
		schema.UserTable.Name,
		schema.UserTable.Columns.Username,
	)

	// Obtenção da linha
	var userId uuid.UUID
	row := ctx.DB.QueryRow(query, sql.Named("username", ctx.Config.AdminUsername))
	err := row.Scan(&userId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("não existe o usuário administrador")
	} else if err != nil {
		ctx.Logger.Error("Erro ao buscar usuário", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível procurar usuário")
	}
	return userId, nil
}

func resetPassword(ctx *context.Context, password string) error {
	var successMsg string
	schema := ctx.Config.Database.Schema
	tx, err := ctx.DB.Begin()
	if err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return fmt.Errorf("não foi possível criar transação")
	}

	// Agendar rollback em caso de erro
	defer func(tx *sql.Tx, err *error) {
		if tx != nil && *err != nil {
			if *err = tx.Rollback(); *err != nil {
				ctx.Logger.Error("Tentativa de rollback falhou", zap.Error(*err))
			}
		}
	}(tx, &err)

	// Criptografar senha e gerar Timestamp
	hash, err := app.HashPassword(ctx, password)
	if err != nil {
		return err
	}
	ts := time.Now().Unix()

	// Verifica se há o usuário como administrador
	adminId, err := getAdmin(ctx)
	if err != nil {
		ctx.Logger.Info("Administrador '" + ctx.Config.AdminUsername + "' não foi encontrado. Fazendo o cadastro.")
		successMsg = "Administrador '" + ctx.Config.AdminUsername + "' foi criado com sucesso."
		adminId = uuid.New()
		insert := fmt.Sprintf(
			`INSERT INTO %s.%s (%s, %s, %s, %s, %s)
			VALUES (:user_id, :username, :name, :password, :updated_at)`,
			schema.Name,
			schema.UserTable.Name,
			schema.UserTable.Columns.UserId,
			schema.UserTable.Columns.Username,
			schema.UserTable.Columns.Name,
			schema.UserTable.Columns.Password,
			schema.UserTable.Columns.UpdatedAt,
		)

		// Criação
		_, err = tx.Exec(
			insert,
			sql.Named("user_id", adminId.String()),
			sql.Named("username", ctx.Config.AdminUsername),
			sql.Named("name", ctx.Config.AdminName),
			sql.Named("password", hash),
			sql.Named("updated_at", ts),
		)
		if err != nil {
			ctx.Logger.Error("Erro ao criar usuário.", zap.Error(err))
			return fmt.Errorf("não foi possível criar usuário")
		}
	} else {
		ctx.Logger.Info("Administrador '" + ctx.Config.AdminUsername + "' foi encontrado. Fazendo atualização.")
		successMsg = "Administrador '" + ctx.Config.AdminUsername + "' foi atualizado com sucesso."
		update := fmt.Sprintf(
			`UPDATE %s.%s
			SET %s = :password, %s = :updated_at
			WHERE %s = :user_id`,
			schema.Name,
			schema.UserTable.Name,
			schema.UserTable.Columns.Password,
			schema.UserTable.Columns.UpdatedAt,
			schema.UserTable.Columns.UserId,
		)

		// Atualização
		_, err = tx.Exec(
			update,
			sql.Named("password", hash),
			sql.Named("updated_at", ts),
			sql.Named("user_id", adminId.String()),
		)
		if err != nil {
			ctx.Logger.Error("Erro ao atualizar usuário.", zap.Error(err))
			return fmt.Errorf("não foi possível atualizar usuário")
		}
	}

	// Commit
	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(err))
		return fmt.Errorf("não foi possível confirmar transação")
	}
	ctx.Logger.Info(successMsg)
	return nil
}

func init() {
	// Criação da pasta logs, caso não exista
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		panic(err)
	}
}

func main() {
	// Logger
	logr := logger.CreateLogger()
	logr.Info("Iniciando aplicação - Repor senha de administrador")

	// Configurações
	cfg, err := config.LoadConfig(logr)
	if err != nil {
		logr.Fatal("Erro ao carregar configurações", zap.Error(err))
	}

	// Banco de dados
	dataBase, err := db.GetSqlDB(&cfg.Database, logr)
	if err != nil {
		logr.Fatal("Erro ao carregar banco de dados", zap.Error(err))
	}
	defer func(dataBase *sql.DB) {
		if err = dataBase.Close(); err != nil {
			logr.Error("Erro ao fechar banco de dados", zap.Error(err))
		}
	}(dataBase)

	// Contexto da aplicação
	ctx := &context.Context{
		Logger: logr,
		Config: cfg,
		DB:     dataBase,
	}

	// Input de senha
	var newPassword, confirmPassword string
	ok := false

	for !ok {
		fmt.Print("\nNova senha do administrador: (>= 4 caracteres) ")
		bytePw, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Print("Senha não pode ser vazia. Finalizando.")
			return
		}
		newPassword = string(bytePw)
		if len(newPassword) < 4 {
			fmt.Print("Senha não pode ter menos de 4 caracteres. Finalizando.")
			return
		}

		// Confirmar senha
		fmt.Print("\nConfirmar senha: ")
		bytePw, err = term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Print("Senha não pode ser vazia. Finalizando.")
			return
		}
		confirmPassword = string(bytePw)
		if newPassword != confirmPassword {
			fmt.Printf("As senhas não coincidem. Tente novamente.\n")
			continue
		}

		// Confirmação
		fmt.Print("\nDeseja continuar? [S/n] (default: n) ")

		var i string
		_, err = fmt.Scanln(&i)
		if err == nil && strings.ToUpper(i) == "S" {
			ok = true
		} else {
			return
		}
	}
	fmt.Println()

	// Cadastrar nova senha
	if err = resetPassword(ctx, newPassword); err != nil {
		logr.Error("Erro ao repor senha de administrador", zap.Error(err))
		return
	}
	logr.Info("Finalizando aplicação.")
}
