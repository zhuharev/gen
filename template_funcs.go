package gen

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var tplFuncsMap = template.FuncMap{
	"titleName": TitleName,
	"goType": func(in ColumnType) string {
		switch in {
		case ColumnTime:
			return "time.Time"
		case ColumnFile:
			return "string"
		case ColumnBool:
			return "bool"
		case ColumnFloat:
			return "float64"
		case ColumnPassword:
			return "string"
		default:
			return string(in)
		}
	},
	"goType2":       goType2,
	"builtInType":   builtInType,
	"deliveryType":  deliveryType,
	"deliveryType2": deliveryType2,
	"isModel":       isModel,
	"modelByName":   modelByName,
	"hasParam": func(method Method) bool {
		for _, arg := range method.Args {
			if arg.IsRestParam {
				return true
			}
		}
		return false
	},
	"httpRESTMethod": func(method Method) string {
		for _, arg := range method.Args {
			if arg.IsRestParam {
				return "GET"
			}
		}
		return "POST"
	},
	"paramIsInt": func(method Method) bool {
		for _, arg := range method.Args {
			if arg.IsRestParam && arg.Type == ColumnInt {
				return true
			}
		}
		return false
	},
	"paramName": func(method Method) string {
		for _, arg := range method.Args {
			if arg.IsRestParam {
				return jsonName(arg.ModelName, arg.Name)
			}
		}
		return ""
	},
	"paramFieldName": func(method Method) string {
		for _, arg := range method.Args {
			if arg.IsRestParam {
				return TitleName(arg.Name)
			}
		}
		return ""
	},
	"jsonName":     jsonName,
	"isCRUDMethod": isCrud,
	"isImplementedMethod": func(implemented []string, method string) bool {
		for _, imp := range implemented {
			if imp == method {
				return true
			}
		}
		return false
	},
	"nsModels":               nsModels,
	"domainServices":         makeServiceByType(""),
	"infrastructureServices": makeServiceByType("infrastructure"),
	"integrationServices":    makeServiceByType(Integration),
	"parsePeriod":            parsePeriod,
	"backgroundMethods": func(project Project) []Method {
		var res []Method
		for _, s := range project.Services {
			for _, m := range s.Methods {
				m.Name = fmt.Sprintf("%sUc.%s", s.Name, TitleName(m.Name))
				if m.RunInBackground != nil {
					res = append(res, m)
				}
			}
		}
		return res
	},
	"serviceImport": func(project Project, serviceName string) string {
		for _, s := range project.Services {
			if s.Name == serviceName {
				if s.Type == "" {
					return project.GoModule + "/domain/" + serviceName
				} else if s.Type == "infrastructure" {
					return project.GoModule + "/infrastructure"
				} else if s.Type == Integration {
					return project.GoModule + "/integration/" + serviceName
				}
			}
		}
		panic(fmt.Sprintf("cannot find service with name: %s", serviceName))
	},
	"servicePkg": func(project Project, serviceName string) string {
		for _, s := range project.Services {
			if s.Name == serviceName {
				if s.Type == "" {
					return serviceName
				} else if s.Type == "infrastructure" {
					return "infrastructure"
				} else if s.Type == Integration {
					return serviceName
				}
			}
		}
		panic(fmt.Sprintf("cannot find service with name: %s", serviceName))
	},
	"usecaseName": func(project Project, serviceName string) string {
		for _, s := range project.Services {
			if s.Name == serviceName {
				if s.Type == "" {
					return "Usecase"
				} else if s.Type == "infrastructure" {
					return TitleName(serviceName)
				} else if s.Type == Integration {
					return "Usecase"
				}
			}
		}
		panic("unreacheble")
	},
	"convertRepresentation": convertRepresentation,
	"removePackage":         removePackage,
	"isArray":               isArray,
	"usecaseBody":           ucBody,
	"isServiceWithRepo":     isServiceWithRepo,
	"enumColumns":           enumColumns,
	"hasEnumColumns":        func(m Model) bool { return len(enumColumns(m)) > 0 },
	"hasID": func(m Model) bool {
		for _, c := range m.Columns {
			if c.Name == "id" {
				return true
			}
		}
		return false
	},
}

func enumColumns(model Model) []Column {
	var (
		enums []Column
	)

	for _, c := range model.Columns {
		if c.Type == ColumnEnum {
			enums = append(enums, c)
		}
	}
	return enums
}

func nsModels(project Project, ns string) []Model {
	var models []Model
	var dup = map[string]struct{}{}
	var service Service
	for _, s := range project.Services {
		if s.Name == ns {
			service = s
			break
		}
	}
	if service.Name == "" {
		panic("service not found")
	}

	for _, dep := range service.Models {
		if _, has := dup[dep]; has {
			continue
		}
		dup[dep] = struct{}{}
		models = append(models, modelByName(project, dep))
		for _, col := range modelByName(project, dep).Columns {
			if col.ModelName != "" && col.Type.IsModel() {
				if _, has := dup[col.ModelName]; has {
					continue
				}
				dup[col.ModelName] = struct{}{}
				models = append(models, modelByName(project, col.ModelName))
				for _, col := range modelByName(project, col.ModelName).Columns {
					if col.ModelName != "" && col.Type.IsModel() {
						if _, has := dup[col.ModelName]; has {
							continue
						}
						dup[col.ModelName] = struct{}{}
						models = append(models, modelByName(project, col.ModelName))
					}
				}
			}
		}
	}

	for _, m := range service.Methods {
		for _, ret := range m.Return {
			if _, has := dup[ret.ModelName]; has {
				continue
			}

			if ret.ModelName == "" {
				continue
			}
			models = append(models, modelByName(project, ret.ModelName))
			dup[ret.ModelName] = struct{}{}

			for _, col := range modelByName(project, ret.ModelName).Columns {
				if col.ModelName != "" && col.Type.IsModel() {
					if _, has := dup[col.ModelName]; has {
						continue
					}
					dup[col.ModelName] = struct{}{}
					models = append(models, modelByName(project, col.ModelName))
					for _, col := range modelByName(project, col.ModelName).Columns {
						if col.ModelName != "" && col.Type.IsModel() {
							if _, has := dup[col.ModelName]; has {
								continue
							}
							dup[col.ModelName] = struct{}{}
							models = append(models, modelByName(project, col.ModelName))
						}
					}
				}
			}
		}
	}

	for _, m := range models {
		for _, r := range m.Relations {
			if _, has := dup[r.To]; has {
				continue
			}
			dup[r.To] = struct{}{}
			models = append(models, modelByName(project, r.To))
		}
	}

	var names []string
	for _, m := range models {
		names = append(names, m.Name)
	}
	log.Printf("models: %s: %+v %+v", ns, names, service.Models)
	return models
}

func convertRepresentation(column Column) string {
	if isArray(column.Type) {
		return fmt.Sprintf("convert%[1]sList(%[2]s)", TitleName(string(column.ModelName)), column.Name)
	}
	return fmt.Sprintf("convert%[1]s(%[2]s)", TitleName(string(column.ModelName)), column.Name)
}

func makeServiceByType(typ string) func(project Project) []Service {
	return func(project Project) []Service {
		var ret []Service
		for _, s := range project.Services {
			if s.Type == typ {
				ret = append(ret, s)
			}
		}
		return ret
	}
}

func isCrud(method Method) bool {
	if len(method.Args) != 1 || len(method.Return) != 1 {
		return false
	}
	return method.Args[0].ModelName == method.Return[0].ModelName
}

func ucBody(method Method) string {
	switch method.Name {
	case "create":
		if len(method.Args) != 1 {
			return `panic("unimplemented")`
		}
		return fmt.Sprintf(`	err := uc.repo.Create%[1]s(ctx, %[2]s)
	return %[2]s.ID, err`, TitleName(method.Args[0].ModelName), method.Args[0].Name)
	case "update":
		if len(method.Args) != 1 {
			return `panic("unimplemented")`
		}
		return fmt.Sprintf(`return uc.repo.Update%[1]s(ctx, %[2]s)`, TitleName(method.Args[0].ModelName), method.Args[0].Name)
	case "list":
		if len(method.Return) != 1 {
			return `panic("unimplemented")`
		}
		return fmt.Sprintf(`return uc.repo.List%[1]s(ctx)`, TitleName(method.Return[0].ModelName))
	}
	if len(method.Args) == 1 && len(method.Return) == 1 && method.Args[0].ModelName == method.Return[0].ModelName {
		return fmt.Sprintf(`return uc.repo.Get%[1]s%[2]s(ctx, %[3]s)`, TitleName(method.Args[0].ModelName), TitleName(method.Name), method.Args[0].Name)
	}
	return `panic("unimplemented")`
}

func isServiceWithRepo(s Service) bool {
	return s.Type == ""
}

func removePackage(in string, pkg string) string {
	if !strings.Contains(in, ".") || !strings.Contains(in, pkg+".") {
		return in
	}
	if isArray(ColumnType(in)) {
		return "[]" + strings.SplitN(in, ".", 2)[1]
	}

	return strings.SplitN(in, ".", 2)[1]
}

func isArray(columnType ColumnType) bool {
	return strings.HasPrefix(string(columnType), "[]")
}

func deliveryType(project Project, column Column) string {
	if isArray(column.Type) {
		return strings.ReplaceAll(goType2(project, column), "[]domain.", "[]")
	}
	if column.Type == ColumnModel {
		typ := goType2(project, column)
		typ = strings.TrimPrefix(typ, "domain.")
		typ = strings.TrimPrefix(typ, "*domain.")
		return typ
	}
	return builtInType(project, column) // TitleName(column.Name)
}

func deliveryType2(project Project, column Column) string {
	if column.Type == ColumnDate {
		return "Date"
	}
	if column.Type == ColumnEnum {
		return "int"
	}
	return deliveryType(project, column)
}

func builtInType(project Project, column Column) string {
	if column.Type == ColumnEnum {
		return "int"
	}
	if column.Type == ColumnDate {
		return "string"
	}
	return goType2(project, column)
}

func goType2(project Project, column Column) string {
	var (
		columnType = column.Type
		isArray    bool
	)

	if strings.HasPrefix(string(column.Type), "[]") {
		columnType = column.Type[2:]
		isArray = true
	}

	var typ string
COLUMN:
	switch columnType {
	case ColumnDate:
		typ = "time.Time"
	case ColumnInt:
		typ = "int"
	case ColumnTime:
		typ = "time.Time"
	case ColumnFile:
		typ = "string"
	case ColumnBool:
		typ = "bool"
	case ColumnFloat:
		typ = "float64"
	case ColumnPassword:
		typ = "string"
	case ColumnString:
		typ = "string"
	case ColumnEnum:
		for _, model := range project.Models {
			if model.Name == column.ModelName {
				for _, c := range model.Columns {
					if c.Name == column.Name {
						typ = "domain." + TitleName(model.Name) + TitleName(c.Name)
						break COLUMN
					}
				}
			}
		}
		panic("unknown type " + typ + " " + column.Name)
	case ColumnModel:
		for _, model := range project.Models {
			if model.Name == column.ModelName {
				typ = "domain." + TitleName(model.Name)
				if !isArray {
					typ = "*" + typ
				}
				break COLUMN
			}
		}
		panic("unknown type " + typ + " " + column.Name)
	default:
		panic(fmt.Sprintf("unknown column type %+v ct:%s isArray:%v", column, columnType, isArray))
	}
	if isArray {
		return "[]" + typ
	}
	return typ
}

func isModel(project Project, column Column) bool {
	columnType := column.Type
	if strings.HasPrefix(string(columnType), "[]") {
		columnType = columnType[2:]
	}
	return columnType == ColumnModel
}

func TitleName(in string) string {
	if strings.HasSuffix(strings.ToLower(string(in)), "id") {
		nn := in[:len(in)-2] + "ID"
		return strings.Title(string(nn))
	}
	if strings.HasSuffix(strings.ToLower(string(in)), "ids") {
		nn := in[:len(in)-3] + "IDs"
		return strings.Title(string(nn))
	}
	return strings.Title(string(in))
}

func jsonName(modelName string, fieldName string) string {
	if strings.ToLower(fieldName) == "id" && modelName != "" {
		return modelName + "Id"
	}
	return fieldName
}

func modelByName(project Project, name string) Model {
	for _, m := range project.Models {
		if m.Name == name {
			return m
		}
	}
	panic(fmt.Sprintf("model %s not found", name))
}

func parsePeriod(period string) int {
	arr := strings.Fields(period)
	if len(arr) != 2 {
		panic("bad period format")
	}
	value, err := strconv.Atoi(arr[0])
	if err != nil {
		panic("bad period format")
	}
	var multiplier time.Duration
	switch arr[1] {
	case "years", "year":
		multiplier = time.Hour * 24 * 365
	case "day", "days":
		multiplier = time.Hour * 24
	case "hour", "hours":
		multiplier = time.Hour
	case "minute", "minutes":
		multiplier = time.Minute
	case "second", "seconds":
		multiplier = time.Second
	default:
		panic("unknown period: " + arr[1])
	}
	return int(multiplier * time.Duration(value))
}
