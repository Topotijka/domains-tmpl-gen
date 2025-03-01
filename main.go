package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const usecaseTemplate = `package {{.Entity}}_usecase

type I{{.Entity}}Service interface {}

type {{.Entity}}Usecase struct {
	service I{{.Entity}}Service
}

func New{{.Entity}}Usecase(service I{{.Entity}}Service) *{{.Entity}}Usecase {
	return &{{.Entity}}Usecase{service: service}
}`

const handlerTemplate = `package {{.Entity}}_handler

import (
	"github.com/gofiber/fiber/v2"
	"your_project/storage"
	"your_project/{{.Entity}}_deps"
)

type I{{.Entity}}Usecase interface {}

type {{.Entity}}Handler struct {
	usecase I{{.Entity}}Usecase
}

func AddRoutes(api *fiber.Router, storage *storage.Storage) {
	Deps := {{.Entity}}_deps.NewDeps(storage.DB)
	{{.Entity}}Handler := &{{.Entity}}Handler{usecase: Deps.{{.Entity}}Usecase()}
	api{{.Entity}} := (*api).Group("/{{.EntityLower}}s")
	api{{.Entity}}.Get("/", {{.Entity}}Handler.GetAll)
}`

const repoTemplate = `package {{.Entity}}_repo

import "gorm.io/gorm"

type {{.Entity}}Repo struct {
	db *gorm.DB
}

func New{{.Entity}}Repo(db *gorm.DB) *{{.Entity}}Repo {
	return &{{.Entity}}Repo{db: db}
}`

const depsTemplate = `package {{.Entity}}_deps

import (
	"gorm.io/gorm"
)

type Deps struct {
	db *gorm.DB

	{{.EntityLower}}Repo    *{{.Entity}}_repo.{{.Entity}}Repo
	{{.EntityLower}}Usecase *{{.Entity}}_usecase.{{.Entity}}Usecase
}

func NewDeps(db *gorm.DB) *Deps {
	return &Deps{db: db}
}

func (deps *Deps) {{.Entity}}Repo() *{{.Entity}}_repo.{{.Entity}}Repo {
	if deps.{{.EntityLower}}Repo == nil {
		deps.{{.EntityLower}}Repo = {{.Entity}}_repo.New{{.Entity}}Repo(deps.db)
	}
	return deps.{{.EntityLower}}Repo
}

func (deps *Deps) {{.Entity}}Usecase() *{{.Entity}}_usecase.{{.Entity}}Usecase {
	if deps.{{.EntityLower}}Usecase == nil {
		deps.{{.EntityLower}}Usecase = {{.Entity}}_usecase.New{{.Entity}}Usecase(deps.{{.Entity}}Repo())
	}
	return deps.{{.EntityLower}}Usecase
}`

func createFolders(entity string) error {
	folders := []string{"usecase", "handler", "repo", "deps"}
	for _, folder := range folders {
		basePath := filepath.Join("generated", entity, folder)
		if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func generateFile(entity, folder, templateContent string) error {
	filePath := filepath.Join("generated", entity, folder, fmt.Sprintf("%s_%s.go", entity, folder))
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
		"EntityLower": strings.ToLower(entity),
	})
}

func main() {
	entities := []string{"User", "Order", "Product"}

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
