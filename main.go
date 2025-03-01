package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const usecaseTemplate = `package {{.EntityLower}}_usecase

type I{{.Entity}}Repo interface {}

type {{.Entity}}Usecase struct {
	repo I{{.Entity}}Repo
}

func New{{.Entity}}Usecase(repo I{{.Entity}}Repo) *{{.Entity}}Usecase {
	return &{{.Entity}}Usecase{repo: repo}
}`

const handlerTemplate = `package {{.EntityLower}}_handler

import (
	"github.com/gofiber/fiber/v2"
)

type I{{.Entity}}Usecase interface {}

type {{.Entity}}Handler struct {
	usecase I{{.Entity}}Usecase
}

func AddRoutes(api *fiber.Router, storage *storage.Storage) {
	deps := {{.EntityLower}}_deps.NewDeps(storage.DB)
	{{.Entity}}Handler := &{{.Entity}}Handler{usecase: deps.{{.Entity}}Usecase()}

	api{{.Entity}} := (*api).Group("/{{.EntityLower}}s")
	api{{.Entity}}.Get("/", {{.Entity}}Handler.GetAll)
}`

const repoTemplate = `package {{.EntityLower}}_repo

import "gorm.io/gorm"

type {{.Entity}}Repo struct {
	db *gorm.DB
}

func New{{.Entity}}Repo(db *gorm.DB) *{{.Entity}}Repo {
	return &{{.Entity}}Repo{db: db}
}`

const depsTemplate = `package {{.EntityLower}}_deps

import (
	"gorm.io/gorm"
	"your_project/{{.EntityLower}}_repo"
	"your_project/{{.EntityLower}}_usecase"
)

type Deps struct {
	db *gorm.DB

	{{.EntityLower}}Repo    *{{.EntityLower}}_repo.{{.Entity}}Repo
	{{.EntityLower}}Usecase *{{.EntityLower}}_usecase.{{.Entity}}Usecase
}

func NewDeps(db *gorm.DB) *Deps {
	return &Deps{db: db}
}

func (deps *Deps) {{.Entity}}Repo() *{{.EntityLower}}_repo.{{.Entity}}Repo {
	if deps.{{.EntityLower}}Repo == nil {
		deps.{{.EntityLower}}Repo = {{.EntityLower}}_repo.New{{.Entity}}Repo(deps.db)
	}
	return deps.{{.EntityLower}}Repo
}

func (deps *Deps) {{.Entity}}Usecase() *{{.EntityLower}}_usecase.{{.Entity}}Usecase {
	if deps.{{.EntityLower}}Usecase == nil {
		deps.{{.EntityLower}}Usecase = {{.EntityLower}}_usecase.New{{.Entity}}Usecase(deps.{{.Entity}}Repo())
	}
	return deps.{{.EntityLower}}Usecase
}`

func createFolders(entity string) error {
	entityLower := strings.ToLower(entity)

	folders := []string{"usecase", "handler", "repo", "deps"}
	for _, folder := range folders {
		basePath := filepath.Join("generated", entityLower, folder)
		if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func generateFile(entity, folder, templateContent string) error {
	entityLower := strings.ToLower(entity)

	filePath := filepath.Join("generated", entityLower, folder, fmt.Sprintf("%s_%s.go", entityLower, folder))
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.New(folder).Parse(templateContent)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, map[string]string{
		"Entity":      entity,
		"EntityLower": entityLower,
	})
}

func main() {
	entities := []string{"Actual"}

	templates := map[string]string{
		"usecase": usecaseTemplate,
		"handler": handlerTemplate,
		"repo":    repoTemplate,
		"deps":    depsTemplate,
	}

	for _, entity := range entities {
		if err := createFolders(entity); err != nil {
			fmt.Println("Ошибка создания папок:", err)
			continue
		}

		for folder, templateContent := range templates {
			if err := generateFile(entity, folder, templateContent); err != nil {
				fmt.Println("Ошибка генерации файла:", err)
				continue
			}
		}
	}

	fmt.Println("Генерация завершена!")
}
