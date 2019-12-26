package gen

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/markbates/pkger"
)

// Generate create folder structure and generate stubs
func Generate(targetDir string, project Project) error {
	project = prepareProject(project)

	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create project dir: %w", err)
	}

	err = createStructure(targetDir, project)
	if err != nil {
		return fmt.Errorf("create project structure: %w", err)
	}

	err = createModels(filepath.Join(targetDir, "domain"), project)
	if err != nil {
		return fmt.Errorf("generate models err=%s", err)
	}

	cmd := exec.Command("goimports", "-w", "-srcdir", "./", "./")
	cmd.Dir = targetDir
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("err=%s log=%s", err, out)
	}

	return err
}

// prepareProject fill defaults and computed fields
func prepareProject(project Project) Project {
	// fill bi-direction relations
	for i, model := range project.Models {
		for j, relation := range model.Relations {
			if relation.Name == "" {
				panic("relation should have name")
			}
			for k, m := range project.Models {
				if relation.To != m.Name {
					continue
				}
				for l, r := range m.Relations {
					if relation.Ref == r.Name {
						log.Println("set relation ref")
						relation.Relation = newRelationSibling(r)
						if relation.Title == "" {
							relation.Title = r.Title
						} else {
							r.Title = relation.Title
						}
						project.Models[k].Relations[l].Relation = newRelationSibling(relation)
					}
				}
			}
			model.Relations[j] = relation
		}
		project.Models[i] = model
	}
	// add relations fields
	for i, model := range project.Models {
		for _, relation := range model.Relations {
			// many to many
			if relation.Relation.Name != "" {
				fieldName := relation.To + "Ids"
				//TODO: проверить что у модели нет такого поля уже
				if relation.Name[0] != relation.To[0] {
					fieldName = relation.Name + TitleName(relation.To) + "Ids"
				}
				title := relation.Title
				if title == "" {
					title = relation.Relation.Title
				}
				model.Columns = append(model.Columns, Column{
					Name:            fieldName,
					Type:            "[]" + ColumnInt,
					Title:           relation.Title,
					IsRelationField: true,
				})
				if relation.FilledFieldName != "" {
					model.Columns = append(model.Columns, Column{
						Name:            relation.FilledFieldName,
						Type:            "[]" + ColumnModel,
						ModelName:       relation.To,
						Title:           relation.Title,
						IsRelationField: true,
						IDsField:        fieldName,
					})
				}
				continue
			}
			// one to many
			fieldName := relation.Name + "Id"
			if relation.Name != relation.To {
				fieldName = relation.Name + TitleName(relation.To) + "Id"
			}
			model.Columns = append(model.Columns, Column{
				Name:            fieldName,
				Type:            ColumnInt,
				Title:           relation.Title,
				IsRelationField: true,
				IDsFieldTitle:   relation.RelationTitleField,
			}, Column{
				Name:            relation.Name,
				Type:            ColumnModel,
				ModelName:       relation.To,
				IsRelationField: true,
				IDsField:        fieldName,
			})

		}
		project.Models[i] = model
	}
	// add relation model to project
	var relationModelDup = map[string]struct{}{}
	for _, m := range project.Models {
		for _, r := range m.Relations {
			if r.Relation.Name != "" {
				names := []string{r.To, r.Relation.To}
				sort.Strings(names)
				name := names[0] + strings.Title(names[1])
				if _, has := relationModelDup[name]; has {
					continue
				}
				relationModelDup[name] = struct{}{}
				project.Models = append(project.Models, Model{
					Name: name,
					Columns: []Column{
						{
							Name: "id",
							Type: ColumnInt,
						},
						{
							Name: names[0] + "Id",
							Type: ColumnInt,
						},
						{
							Name: names[1] + "Id",
							Type: ColumnInt,
						},
					},
				})

				// add relation model to services
				for j, s := range project.Services {
					for _, m := range s.Models {
						if m == r.To {
							project.Services[j].Models = append(project.Services[j].Models, name)
						}
					}
				}
			}

		}
	}

	return project
}

// GenerateSkeleton generates skeleton from .skeleton file
func GenerateSkeleton() {
	var (
		fname     = ".skeleton"
		targetDir = "."
	)

	f, err := os.Open(fname)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fpath := strings.TrimSpace(scanner.Text())
		if fpath == "" {
			continue
		}
		hasTrailingSlash := fpath[len(fpath)-1] == '/'
		targetPath := filepath.Join(targetDir, fpath)
		if hasTrailingSlash {
			targetPath += "/"
		}
		err = createFileOrDir(targetPath)
		if err != nil {
			log.Fatalln(err)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func createFileOrDir(fpath string) (err error) {
	var (
		isDir   = fpath[len(fpath)-1] == '/'
		dirName = fpath
	)

	if !isDir {
		dirName = filepath.Dir(fpath)
		log.Printf("mkdir all: %s", dirName)
		err = os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return
		}

		//TODO: do not truncate file if it exists
		// _, err = os.Create(fpath)
		// if err != nil {
		// 	return
		// }
		return
	}

	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return
	}

	return
}

func createModels(targetDir string, project Project) error {
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create domain dir: %w", err)
	}

	targetFile := filepath.Join(targetDir, "models.go")

	f, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create models.go: %w", err)
	}

	err = executeModelsToWriter(f, "models.tmpl", project)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	return nil
}

func createStructure(targetDir string, project Project) error {
	const layout = "layout/clean-arch"

	err := pkger.Walk("/layout", func(path string, info os.FileInfo, err error) error {
		fpath := strings.ReplaceAll(path, "github.com/zhuharev/gen:/", "")
		if err != nil {
			return err
		}
		fpath = strings.TrimPrefix(fpath, layout)
		targetPath := strings.TrimSuffix(filepath.Join(targetDir, fpath), ".tmpl")

		fmt.Printf(
			"%s %s\n",
			path,
			targetPath,
		)
		//return nil

		err = createFileOrDir(targetPath)
		if err != nil {
			return fmt.Errorf("create file or dir: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		f, err := pkger.Open(path)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(f)
		if strings.HasSuffix(fpath, ".tmpl") {
			tpl, err := template.New(fpath).Funcs(tplFuncsMap).Parse(string(data))
			if err != nil {
				return fmt.Errorf("parse templates: %w", err)
			}

			projectData, err := json.MarshalIndent(project, "  ", "  ")
			if err != nil {
				panic("marhsal project schema")
			}

			var buf bytes.Buffer
			err = tpl.Execute(&buf, map[string]interface{}{
				"project": project,
				"schema":  string(projectData),
			})
			if err != nil {
				return fmt.Errorf("execute template: %w", err)
			}

			// Do not overwrite existsing md files
			if (strings.HasSuffix(fpath, ".md") || strings.HasSuffix(fpath, ".md.tmpl")) && fpath != "/API.md.tmpl" {
				if exists, _ := isFileExist(targetPath); exists {
					log.Printf("file %s exists, skip create", targetPath)
					return nil
				}
			}

			err = ioutil.WriteFile(targetPath, buf.Bytes(), os.ModePerm)
			if err != nil {
				return err
			}
			return nil
		} else if strings.HasSuffix(targetPath, ".md") {
			// Do not overwrite .md files
			if exists, _ := isFileExist(targetPath); exists {
				//log.Printf("file %s exists, skip create", targetPath)
				return nil
			}
		}

		// Do not overwrite config
		if fpath == "/app/config.go" || fpath == "/conf/config.yml" {
			if exists, _ := isFileExist(targetPath); exists {
				log.Printf("file %s exists, skip create", targetPath)
				return nil
			}
		}

		err = ioutil.WriteFile(targetPath, data, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	return generateServices(targetDir, project)
}

func execTemplate(name string, project *Project, service *Service) (_ []byte, err error) {
	var buf bytes.Buffer
	err = execTemplateToWriter(&buf, name, project, service, nil)
	return buf.Bytes(), err
}

func executeModelsToWriter(wr io.Writer, name string, project Project) error {
	t, err := template.New(name).Funcs(tplFuncsMap).Parse(assets["/"+name])
	if err != nil {
		return fmt.Errorf("parse %s template: %w", name, err)
	}

	var tmplData = map[string]interface{}{
		"project": project,
	}
	err = t.Execute(wr, tmplData)
	return err
}

func execTemplateToWriter(wr io.Writer, name string, project *Project, service *Service, implemented []string) error {
	t, err := template.New(name).Funcs(tplFuncsMap).Parse(assets["/"+name])
	if err != nil {
		return fmt.Errorf("parse %s template: %w", name, err)
	}

	var tmplData = make(map[string]interface{})
	if project != nil {
		tmplData["project"] = project
	}
	if service != nil {
		tmplData["service"] = service
	}
	tmplData["implemented"] = implemented
	err = t.Execute(wr, tmplData)
	return err
}

func generateDomainServices(targetDir string, project Project, s Service) error {
	domainDir := filepath.Join(targetDir, "domain")
	err := os.MkdirAll(filepath.Join(domainDir, s.Name), os.ModePerm)
	if err != nil {
		return err
	}

	// delivery
	err = createFileOrDir(filepath.Join(domainDir, s.Name, "delivery", "rest", "rest.go"))
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(domainDir, s.Name, "delivery", "rest", fmt.Sprintf("rest.go")))
	if err != nil {
		return err
	}
	defer f.Close()

	err = execTemplateToWriter(f, "rest.tmpl", &project, &s, nil)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	// err = createFileOrDir(filepath.Join(domainDir, s.Name, "delivery", "rpc", "rpc.go"))
	// if err != nil {
	// 	return err
	// }
	// err = createFileOrDir(filepath.Join(domainDir, s.Name, "delivery", "rpc", "convert.go"))
	// if err != nil {
	// 	return err
	// }

	f, err = os.Create(filepath.Join(domainDir, s.Name, fmt.Sprintf("%s.go", s.Name)))
	if err != nil {
		return err
	}
	defer f.Close()

	err = execTemplateToWriter(f, "usecase.tmpl", &project, &s, nil)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	usecaseImpPath := filepath.Join(domainDir, s.Name, fmt.Sprintf("%s_usecase.go", s.Name))
	var implemented []string
	if _, err := os.Stat(usecaseImpPath); os.IsNotExist(err) {

	} else if err == nil {
		implemented, err = implementedMethods(usecaseImpPath)
		if err != nil {
			return fmt.Errorf("get implemented methods from file: %w", err)
		}
	}

	flags := os.O_CREATE | os.O_WRONLY
	tmpl := "usecase_imp.tmpl"
	if len(implemented) > 0 {
		tmpl = "usecase_imp_append.tmpl"
		flags = flags | os.O_APPEND
	}

	f, err = os.OpenFile(usecaseImpPath, flags, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open usecase imp file: %w", err)
	}
	defer f.Close()

	err = execTemplateToWriter(f, tmpl, &project, &s, implemented)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	f, err = os.Create(filepath.Join(domainDir, s.Name, fmt.Sprintf("%s_repo.go", s.Name)))
	if err != nil {
		return err
	}
	err = execTemplateToWriter(f, "repo.tmpl", &project, &s, nil)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	f.Close()
	return nil
}

func isFileExist(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func generateInfrastructureService(targetDir string, project Project, s Service) error {
	log.Println("generate infrastructure service", s.Name)
	integrationDir := filepath.Join(targetDir, "infrastructure", s.Name)
	err := os.MkdirAll(integrationDir, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(integrationDir, fmt.Sprintf("%s.go", s.Name)))
	if err != nil {
		return err
	}
	defer f.Close()

	err = execTemplateToWriter(f, "usecase.tmpl", &project, &s, nil)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	f, err = os.Create(filepath.Join(integrationDir, fmt.Sprintf("%s_usecase.go", s.Name)))
	if err != nil {
		return err
	}
	defer f.Close()

	err = execTemplateToWriter(f, "usecase_imp.tmpl", &project, &s, nil)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

const (
	Integration = "integration"
)

func generateIntegrationService(targetDir string, project Project, s Service) error {
	integrationDir := filepath.Join(targetDir, "integration", s.Name)
	err := os.MkdirAll(integrationDir, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(integrationDir, fmt.Sprintf("%s.go", s.Name)))
	if err != nil {
		return err
	}
	defer f.Close()

	err = execTemplateToWriter(f, "usecase.tmpl", &project, &s, nil)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	usecaseImpPath := filepath.Join(integrationDir, fmt.Sprintf("%s_usecase.go", s.Name))
	var implemented []string
	if _, err := os.Stat(usecaseImpPath); os.IsNotExist(err) {

	} else if err == nil {
		implemented, err = implementedMethods(usecaseImpPath)
		if err != nil {
			return fmt.Errorf("get implemented methods from file: %w", err)
		}
	}

	flags := os.O_CREATE | os.O_WRONLY
	tmpl := "usecase_imp.tmpl"
	if len(implemented) > 0 {
		tmpl = "usecase_imp_append.tmpl"
		flags = flags | os.O_APPEND
	}

	f, err = os.OpenFile(usecaseImpPath, flags, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	err = execTemplateToWriter(f, tmpl, &project, &s, implemented)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

func generateServices(targetDir string, project Project) error {
	for _, s := range project.Services {
		var err error
		switch s.Type {
		case "":
			err = generateDomainServices(targetDir, project, s)
		case "infrastructure":
			err = generateInfrastructureService(targetDir, project, s)
		case Integration:
			err = generateIntegrationService(targetDir, project, s)
		default:
			panic("unknown type: " + s.Type)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func implementedMethods(fpath string) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fpath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	var res []string
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Recv.NumFields() > 0 {
			res = append(res, fn.Name.Name)
		}
	}
	return res, nil
}
