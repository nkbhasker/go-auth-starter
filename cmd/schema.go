package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/nkbhasker/go-auth-starter/internal/model"
	_ "github.com/nkbhasker/go-auth-starter/internal/storage"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Generate schema",
	Run: func(cmd *cobra.Command, _args []string) {
		Schema()
	},
}

func Schema() {
	sb := &strings.Builder{}
	loadEnums(sb)
	loadModels(sb)

	io.WriteString(os.Stdout, sb.String())
}

func loadEnums(sb *strings.Builder) *strings.Builder {
	enums := []string{
		`CREATE TYPE gender AS ENUM (
			'MALE',
			'FEMALE'
		);`,
		`CREATE TYPE identity_provider AS ENUM (
			'LOCAL',
			'GOOGLE',
			'APPLE'
		);`,
	}
	for _, enum := range enums {
		sb.WriteString(enum)
	}

	return sb
}

func loadModels(sb *strings.Builder) *strings.Builder {
	models := []interface{}{
		&model.User{},
	}
	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	sb.WriteString(stmts)

	return sb
}
