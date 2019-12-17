// Code generated by assets compiler. DO NOT EDIT.
	
package gen

var (
	assets = map[string]string{
"/models.tmpl": "// Code generated by https://github.com/zhuharev/gen DO NOT EDIT.\n\npackage domain\n\nimport (\n  \"time\"\n  \"regexp\"\n\n  \"github.com/nbutton23/zxcvbn-go\"\n)\n\ntype FieldError struct {\n  FieldName string `json:\"fieldName\"`\n  Description string `json:\"description\"`\n}\n\nvar (\n  emailRe = regexp.MustCompile(`(?i)^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,5})$`)\n  phoneRe = regexp.MustCompile(`^[0-9]{10}$`)\n  cyrillicRe = regexp.MustCompile(`^\\p{Cyrillic}+$`)\n  aliasRe = regexp.MustCompile(`^[a-z][\\w\\.]+$`)\n)\n\nvar Columns = struct {\n  {{ range .project.Models }}{{ $model := . }}\n  {{ titleName $model.Name }} struct { {{ range .Columns }}\n    {{ titleName .Name }} string\n    {{ end }}\n  }\n  {{ end }}\n} {\n {{ range .project.Models }}{{ $model := . }}\n  {{ titleName $model.Name }}: struct { {{ range .Columns }}\n    {{ titleName .Name }} string\n    {{ end }}\n  } {\n  {{ range .Columns }}\n    {{ titleName .Name }}: \"{{ titleName .Name }}\",\n    {{ end }}  \n  },\n  {{ end }}\n}\n\n{{ range .project.Models }}{{ $model := . }}\n\n{{ range enumColumns $model }}{{ $column := . }}\n\ntype {{ titleName $model.Name }}{{ titleName .Name }} int\n\nconst (\n  {{ range $i, $e := .Enums }}{{ titleName $model.Name }}{{ titleName $column.Name }}{{ titleName $e.Name }}{{ if eq $i 0}} {{ titleName $model.Name }}{{ titleName $column.Name }} = iota + 1 {{ end }}\n  {{ end }}\n)\n\nfunc (v {{ titleName $model.Name }}{{ titleName .Name }}) Int() int {\n  return int(v)\n}\n\n{{ end }}\n\n// {{ titleName .Name }} describe model\ntype {{ titleName .Name }} struct { {{ range .Columns }}\n  {{ titleName .Name }}     {{ removePackage (goType2 $.project .) \"domain\" }} {{ if  eq (titleName .Name) \"ID\"}} `storm:\"id,increment\"` {{ else if .Validation}}{{if .Validation.Unique}}`storm:\"unique\"`{{end}}{{else if (and .IsRelationField .IDsField)}}`json:\"-\"`{{ end }}{{ end }}\n}\n\nfunc (s {{ titleName $model.Name }}) Validate() []FieldError {\n  var res []FieldError\n  {{ range .Columns }}\n  if isValid, description := s.Validate{{ titleName .Name }}(); !isValid {\n    res = append(res, FieldError{FieldName: \"{{ jsonName .ModelName .Name }}\", Description: description})\n  }{{ end }}\n  return res\n}\n\n{{ range .Columns }}func (s {{ titleName $model.Name }}) Validate{{ titleName .Name }}() (isValid bool, errDescription string) {\n  {{ if .Validation }} {{ $column := . }}\n  {{ range .Validation.CustomFunctions }}\n  isValid, errDescription = {{ .Name }}(s.{{ titleName $column.Name }})\n  if !isValid {\n    return\n  }\n  {{ end }}\n  {{ end }}\n  {{ if eq .Type \"enum\" }}\n  {{ if .Validation }}\n  {{ if .Validation.Required}}\n  if s.{{ titleName .Name }}.Int() == 0 {\n    errDescription = \"Required\"\n    isValid = false\n    return\n  }\n  {{end}}\n  {{ end }}\n  {{ end }}\n  {{ if eq .Type \"password\" }}\n  passwordStrength := zxcvbn.PasswordStrength(s.{{ titleName .Name }}, nil)\n  if s.{{ titleName .Name }} != \"\" && passwordStrength.Score < 3 {\n    errDescription = \"Weak password\"\n    isValid = false\n    return\n  } else if s.{{ titleName .Name }} != \"\" && len(s.{{ titleName .Name }}) < 6 {\n    errDescription = \"Weak password\"\n    isValid = false\n    return    \n  }\n  {{ end }}\n{{ if .Validation }}\n  {{ if eq .Type \"int\" }}\n    {{- if .Validation.Required -}}\n  if s.{{ titleName .Name }} == 0 {\n    errDescription = \"Required\"\n    isValid = false\n    return\n  }\n    {{ end }}\n  {{ else if eq .Type \"string\" }}\n    {{ if .Validation.Required }}\n  if s.{{ titleName .Name }} == \"\" {\n    errDescription = \"Required\"\n    isValid = false\n    return\n  }\n    {{ end }}\n    {{ if .Validation.IsCyrillic }}\n  if s.{{ titleName .Name }} != \"\" && !cyrillicRe.MatchString(s.{{ titleName .Name }}) {\n    errDescription = \"Must be cyrillic\"\n    isValid = false\n    return    \n  }\n    {{ end }}\n    {{ if .Validation.IsAlias }}\n  if s.{{ titleName .Name }} != \"\" && !aliasRe.MatchString(s.{{ titleName .Name }}) {\n    errDescription = \"Must be lowercase word. Must starts with letter. Can contain _ and . (dot)\"\n    isValid = false\n    return    \n  }\n    {{ end }}\n    {{ if .Validation.IsPhone }}\n  if s.{{ titleName .Name }} != \"\" && !phoneRe.MatchString(s.{{ titleName .Name }}) || (s.{{ titleName .Name }}[0] != '7' && s.{{ titleName .Name }}[0] != '9' ) {\n    errDescription = \"10 digits of phone number without spaces or separators (e.g. 9001234567)\"\n    isValid = false\n    return\n  }\n    {{ end }}\n    {{ if .Validation.IsEmail }}\n  if s.{{ titleName .Name }} != \"\" && emailRe.FindString(s.{{ titleName .Name }}) == \"\" {\n    errDescription = \"Bad email\"\n    isValid = false\n    return\n  }\n    {{ end }}\n    {{ if .Validation.Min }}\n  if s.{{ titleName .Name }} != \"\" && len([]rune(s.{{ titleName .Name }})) < {{ .Validation.Min.Value }} {\n    errDescription = \"{{ .Validation.Min.Description }}\"\n    isValid = false\n    return    \n  }\n    {{ end }}\n    {{ if .Validation.Max }}\n  if s.{{ titleName .Name }} != \"\" && len([]rune(s.{{ titleName .Name }})) > {{ .Validation.Max.Value }} {\n    errDescription = \"{{ .Validation.Max.Description }}\"\n    isValid = false\n    return    \n  }\n    {{ end }}\n  {{ else if or (eq .Type \"time\") (eq .Type \"date\") }}\n    {{ if .Validation.Required }}\n  if s.{{ titleName .Name }}.IsZero() {\n    errDescription = \"Required\"\n    isValid = false\n    return\n  }\n    {{ end }}\n    {{ if .Validation.MinAge }}\n  minAge := time.Duration({{ parsePeriod .Validation.MinAge.Value }})\n  minAgeReached := s.{{ titleName .Name }}.Add(minAge)\n  if time.Now().Before(minAgeReached) {\n    errDescription = \"{{ .Validation.MinAge.Description }}\"\n    isValid = false\n    return\n  }\n    {{ end }}\n    {{ if .Validation.MaxAge }}\n  maxAge := time.Duration({{ parsePeriod .Validation.MaxAge.Value}})\n  maxAgeReached := s.{{ titleName .Name }}.Add(maxAge)\n  if time.Now().After(maxAgeReached) {\n    errDescription = \"{{ .Validation.MaxAge.Description }}\"\n    isValid = false\n    return\n  }\n    {{ end }}\n  {{ end }}\n{{ end }}\n  isValid = true\n  return \n}\n{{end}}\n{{ end }}",
"/repo.tmpl": "// Code generated by https://github.com/zhuharev/gen DO NOT EDIT.\n\npackage {{ .service.Name }}\n\nimport (\n    \"fmt\"\n    \"context\"\n\n    \"{{ .project.GoModule }}/domain\"\n    \"{{ .project.GoModule }}/infrastructure\"\n)\n\ntype ReadAccessRulesApplier interface {\n    ApplyReadAccessRules(context.Context) error\n}\n\ntype Repo struct {\n    infrastructure.Database\n\n    {{ range .service.Models }}\n{{ $model := modelByName $.project . }}\n    afterReadList{{ titleName $model.Name }} func(context.Context, *Repo, []domain.{{ titleName $model.Name }}) ([]domain.{{ titleName $model.Name }},error)\n    afterRead{{ titleName $model.Name }} func(context.Context, *Repo, *domain.{{ titleName $model.Name }}) (error)\n{{ end }}\n}\n\nfunc NewRepo(db infrastructure.Database) *Repo {\n    return &Repo{\n        Database: db,\n    }\n}\n\nfunc (r *Repo) RunInTransaction(ctx context.Context, fn func(context.Context, *Repo) error) error {\n    return r.Database.RunInTransaction(ctx, func(ctx context.Context, db infrastructure.Database) error {\n        repo := NewRepo(db)\n            {{ range .service.Models }}\n{{ $model := modelByName $.project . }}\n        repo.afterReadList{{ titleName $model.Name }} = r.afterReadList{{ titleName $model.Name }}\n        repo.afterRead{{ titleName $model.Name }} = r.afterRead{{ titleName $model.Name }}\n        {{ end }}\n        return fn(ctx, repo)\n    })\n}\n\n{{ range .service.Models }}\n{{ $model := modelByName $.project . }}\n{{ range $model.Columns }}\n{{ if isModel $.project .}}\n{{ else }}\nfunc (r *Repo) Get{{ titleName $model.Name }}By{{ titleName .Name }}(ctx context.Context, value {{ goType2 $.project . }}) (*domain.{{ titleName $model.Name }}, error) {\n    var u domain.{{ titleName $model.Name }}\n    err := r.Database.GetByField(ctx, \"{{ titleName .Name }}\", value, &u)\n    if err != nil {\n        return nil, fmt.Errorf(\"get {{$model.Name}} by {{ .Name }}: %w\", err)\n    }\n\n    if rulesApplier, ok := interface{}(&u).(ReadAccessRulesApplier); ok {\n        err := rulesApplier.ApplyReadAccessRules(ctx)\n        if err !=nil {\n            return nil, err\n        }\n    }\n\n    if r.afterRead{{ titleName $model.Name }} != nil {\n        err = r.afterRead{{ titleName $model.Name }}(ctx, r, &u)\n        if err != nil {\n            return nil, fmt.Errorf(\"get callback: %w\", err)\n        }\n    }\n\n    return &u, nil\n}\nfunc (r *Repo) List{{ titleName $model.Name }}By{{ titleName .Name }}(ctx context.Context, value {{ goType2 $.project . }}, opts ...infrastructure.ListOptions) ([]domain.{{ titleName $model.Name }}, error) {\n    var u []domain.{{ titleName $model.Name }}\n    err := r.Database.ListByField(ctx, \"{{ titleName .Name }}\", value, &u, opts...)\n    if err == infrastructure.ErrNotFound {\n        return nil, nil\n    }\n    if err != nil {\n        return nil, fmt.Errorf(\"get list {{$model.Name}} by {{ .Name }}: %w\", err)\n    }\n\n    if len(u) > 0 {\n        if _, ok := interface{}(&u[0]).(ReadAccessRulesApplier); ok {\n            for i := range u {\n                err := interface{}(&u[i]).(ReadAccessRulesApplier).ApplyReadAccessRules(ctx)\n                if err != nil {\n                    return nil, err\n                }\n            }\n        }\n    }\n\n    if r.afterReadList{{ titleName $model.Name }} != nil {\n        u, err = r.afterReadList{{ titleName $model.Name }}(ctx, r, u)\n        if err != nil {\n            return nil, fmt.Errorf(\"get list callback: %w\", err)\n        }\n    }\n\n    return u, nil\n}\n{{end}}\n{{ end }}\n\n\n\nfunc (r *Repo) Create{{ titleName . }}(ctx context.Context, model *domain.{{ titleName . }}) error {\n    err := r.Database.Create(ctx, model)\n    if err != nil {\n        return fmt.Errorf(\"create {{ titleName . }}: %w\", err)\n    }\n    return nil\n}\n\n{{ if hasID $model }}\nfunc (r *Repo) Delete{{ titleName . }}(ctx context.Context, id int ) error {\n    model := &domain.{{ titleName . }}{ID: id}\n    err := r.Database.Delete(ctx, model)\n    if err != nil {\n        return fmt.Errorf(\"delete {{ titleName . }}: %w\", err)\n    }\n    return nil\n}\n{{ end }}\n\nfunc (r *Repo) Update{{ titleName . }}(ctx context.Context, model *domain.{{ titleName . }}) error {\n    err := r.Database.Update(ctx, model)\n    if err != nil {\n        return fmt.Errorf(\"update {{ titleName . }}: %w\", err)\n    }\n    return nil\n}\n\nfunc (r *Repo) Update{{ titleName . }}Field(ctx context.Context, model *domain.{{ titleName . }}, field infrastructure.Field, fields ...infrastructure.Field) error {\n    err := r.Database.UpdateField(ctx, model, field, fields...)\n    if err != nil {\n        return fmt.Errorf(\"update {{ titleName . }}: %w\", err)\n    }\n    return nil\n}\n\nfunc (r *Repo) List{{ titleName . }}(ctx context.Context,opts ...infrastructure.ListOptions) ([]domain.{{ titleName . }}, error) {\n    var to []domain.{{ titleName . }}\n    err := r.Database.List(ctx, &to, opts...)\n    if err == infrastructure.ErrNotFound {\n        return nil, nil\n    }\n    if err != nil {\n        return nil, fmt.Errorf(\"get list of {{ titleName . }} from db: %w\", err)\n    }\n\n    if len(to) > 0 {\n        if _, ok := interface{}(&to[0]).(ReadAccessRulesApplier); ok {\n            for i := range to {\n                err := interface{}(&to[i]).(ReadAccessRulesApplier).ApplyReadAccessRules(ctx)\n                if err != nil {\n                    return nil, err\n                }\n            }\n        }\n    }\n\n    if r.afterReadList{{ titleName $model.Name }} != nil {\n        to, err = r.afterReadList{{ titleName $model.Name }}(ctx, r, to)\n        if err != nil {\n            return nil, fmt.Errorf(\"get list callback: %w\", err)\n        }\n    }\n\n    return to, nil\n}\n\nfunc (r *Repo) List{{ titleName . }}ByIDs(ctx context.Context, ids []int, opts ...infrastructure.ListOptions) ([]domain.{{ titleName . }}, error) {\n    list, err := r.List{{ titleName .}}ByFields(ctx, []infrastructure.Field{\n        {\n            Name: r.ColumnByName(\"ID\"),\n            Value: ids,\n            CompareType: infrastructure.CompareTypeIn,\n        },\n    })\n    if err != nil {\n        return nil, err \n    }\n\n    if len(list) > 0 {\n        if _, ok := interface{}(&list[0]).(ReadAccessRulesApplier); ok {\n            for i := range list {\n                err := interface{}(&list[i]).(ReadAccessRulesApplier).ApplyReadAccessRules(ctx)\n                if err != nil {\n                    return nil, err\n                }\n            }\n        }\n    }\n\n    if r.afterReadList{{ titleName $model.Name }} != nil {\n        var err error\n        list, err = r.afterReadList{{ titleName $model.Name }}(ctx, r, list)\n        if err != nil {\n            return nil, fmt.Errorf(\"get list callback: %w\", err)\n        }\n    }    \n\n    return list, nil\n}\n\nfunc (r *Repo) List{{ titleName . }}ByFields(ctx context.Context, fields []infrastructure.Field, opts ...infrastructure.ListOptions) ([]domain.{{ titleName . }}, error) {\n    var to []domain.{{ titleName . }}\n    err := r.Database.ListByFields(ctx, fields, &to, opts...)\n    if err == infrastructure.ErrNotFound {\n        return nil, nil\n    }\n    if err != nil {\n        return nil, fmt.Errorf(\"get list of {{ titleName . }} from db: %w\", err)\n    }\n\n    if len(to) > 0 {\n        if _, ok := interface{}(&to[0]).(ReadAccessRulesApplier); ok {\n            for i := range to {\n                err := interface{}(&to[i]).(ReadAccessRulesApplier).ApplyReadAccessRules(ctx)\n                if err != nil {\n                    return nil, err\n                }\n            }\n        }\n    }\n\n    if r.afterReadList{{ titleName $model.Name }} != nil {\n        to, err = r.afterReadList{{ titleName $model.Name }}(ctx, r, to)\n        if err != nil {\n            return nil, fmt.Errorf(\"get list callback: %w\", err)\n        }\n    }\n\n    return to, nil\n}\n\nfunc (r *Repo) Count{{ titleName . }}ByFields(ctx context.Context, fields []infrastructure.Field) (int, error) {\n    var to domain.{{ titleName . }}\n    return r.Database.CountByFields(ctx, fields, &to)\n}\n\nfunc (r *Repo) Get{{ titleName . }}ByFields(ctx context.Context, fields []infrastructure.Field) (*domain.{{ titleName . }}, error) {\n    var to domain.{{ titleName . }}\n    err := r.Database.ListByFields(ctx, fields, &to)\n    if err == infrastructure.ErrNotFound {\n        return nil, err\n    }\n    if err != nil {\n        return nil, fmt.Errorf(\"get {{ titleName . }} from db by fields: %w\", err)\n    }\n\n    if rulesApplier, ok := interface{}(&to).(ReadAccessRulesApplier); ok {\n        err := rulesApplier.ApplyReadAccessRules(ctx)\n        if err != nil {\n            return nil, err\n        }\n    }\n\n    if r.afterRead{{ titleName $model.Name }} != nil {\n        err = r.afterRead{{ titleName $model.Name }}(ctx, r, &to)\n        if err != nil {\n            return nil, fmt.Errorf(\"get callback: %w\", err)\n        }\n    }\n\n    return &to, nil\n}\n{{ end }}\n\n",
"/rest.tmpl": "// Code generated by https://github.com/zhuharev/gen DO NOT EDIT.\n\npackage rest\n\nimport (\n\t\"fmt\"\n\t\"net/http\"\n\t\"strconv\"\n\n\t\"{{ .project.GoModule }}/domain\"\n\t\"{{ .project.GoModule }}/domain/{{ .service.Name }}\"\n\n\t\"github.com/labstack/echo/v4\"\n)\n\n// Date is date representation in format: 31.02.2009\ntype Date time.Time\n\nconst dateTimeLayout = \"02.01.2006\"\n\nfunc (d *Date) UnmarshalJSON(data []byte) error {\n\tvar s string\n\terr := json.Unmarshal(data, &s)\n\tif err != nil {\n\t\treturn err \n\t}\n\tt, err := time.Parse(dateTimeLayout, s)\n\tif err != nil {\n\t\treturn err\n\t}\n\t*d = Date(t)\n\treturn nil\n}\n\nfunc (d *Date) MarshalJSON() ([]byte, error) {\n\tif d == nil {\n\t\treturn nil, nil\n\t}\n\ts := time.Time(*d).Format(dateTimeLayout)\n\treturn []byte(`\"` + s + `\"`), nil\n}\n\nfunc (d Date) String() string {\n\treturn time.Time(d).Format(dateTimeLayout)\n}\n\ntype Validator interface {\n\tValidate() []domain.FieldError\n}\n\n// EchoDelivery can registre http handler on echo server\ntype EchoDelivery struct {\n\tusecase {{ .service.Name }}.Usecase\n}\n\n// NewEchoDelivery returns new delivery\nfunc New(usecase {{ .service.Name }}.Usecase) *EchoDelivery {\n\treturn &EchoDelivery{\n\t\tusecase: usecase,\n\t}\n}\n\ntype UserIDer interface {\n\tUserID() (int,bool)\n}\n\n// RegistreHandlers registre handlers to echo.Group\nfunc (e *EchoDelivery) RegistreHandlers(group *echo.Group) {\n{{ range .service.Methods }}\n{{ if (hasParam . )}}\n\tgroup.GET(\"/:{{ paramName . }}\", e.{{ titleName .Name }})\n{{else}}\n\tgroup.Any(\"/{{ .Name }}\", e.{{ titleName .Name }})\n{{end}}   \n{{ end }}\n\t\n}\n\ntype Error struct {\n\tCode    int    `json:\"code\"`\n\tMessage string `json:\"message\"`\n\tValidationErrors []domain.FieldError `json:\"validationErrors,omitempty\"`\n}\n\nvar ErrBadRequest = fmt.Errorf(\"bad request\")\n\n{{ range (nsModels .project .service.Name)}}\n{{ $model := .}}\ntype {{ titleName .Name }} struct { {{ range .Columns }}\n  {{ titleName .Name }}     {{ removePackage (deliveryType2 $.project .) \"domain\" }} {{ if eq .Type.String \"password\"}} `json:\"password\"` {{else}} `json:\"{{ jsonName .ModelName .Name }}\"`{{end}} {{ end }}\n}\n\nfunc convert{{ titleName .Name }}(model *domain.{{ titleName .Name }}) {{ titleName .Name }} {\n\treturn {{ titleName .Name }}{\n\t\t{{ range .Columns }} \n\t\t{{ titleName .Name }}: {{ if isModel $.project .}}{{ if isArray (.Type) }}convert{{ titleName .ModelName}}List(model.{{ titleName .Name }}){{ else }}convert{{ titleName .ModelName}}(&model.{{ titleName .Name }}){{end}}{{ else if (eq .Type \"enum\") }}model.{{ titleName .Name }}.Int(){{ else if eq .Type \"date\" }}Date(model.{{ titleName .Name }}){{ else if eq .Type \"password\" }}\"\"{{ else }}  model.{{ titleName .Name }}{{end}},   {{ end }}\n\t}\n}\n\nfunc convert{{ titleName .Name }}ToDomain(model {{ titleName .Name }}) *domain.{{ titleName .Name }} {\n\treturn &domain.{{ titleName .Name }}{\n\t\t{{ range .Columns }} \n\t\t{{ titleName .Name }}: {{ if isModel $.project .}}{{ if isArray (.Type) }}convert{{ titleName .ModelName}}ListToDomain(model.{{ titleName .Name }}){{ else }}*convert{{ titleName .ModelName}}ToDomain(model.{{ titleName .Name }}){{end}}{{ else if (eq .Type \"enum\") }}domain.{{ titleName $model.Name }}{{ titleName .Name }}(model.{{ titleName .Name }}){{ else if eq .Type \"date\" }}time.Time(model.{{ titleName .Name }}){{ else }}  model.{{ titleName .Name }}{{end}},   {{ end }}\n\t}\n}\n\nfunc convert{{ titleName .Name }}List(in []domain.{{ titleName .Name }}) []{{ titleName .Name }} {\n\tres := make([]{{ titleName .Name }}, len(in))\n\tfor i, v := range in {\n\t\tres[i] = convert{{ titleName .Name }}(&v)\n\t}\n\treturn res\n}\n\nfunc convert{{ titleName .Name }}ListToDomain(in []{{ titleName .Name }}) []domain.{{ titleName .Name }} {\n\tres := make([]domain.{{ titleName .Name }}, len(in))\n\tfor i, v := range in {\n\t\tres[i] = *convert{{ titleName .Name }}ToDomain(v)\n\t}\n\treturn res\n}\n{{ end }}\n\n{{ range .service.Methods }}\n\n{{ if (hasParam . )}}\nfunc set{{ titleName .Name }}Param(ctx echo.Context, r *{{ titleName .Name }}Request) error {\n\t{{ if (paramIsInt .) }}\n\tvar err error\n\tr.{{ paramFieldName . }}, err =strconv.Atoi(ctx.Param(\"{{ paramName . }}\"))\n\tif err != nil || r.{{ paramFieldName . }} < 1 {\n\t\treturn ctx.JSON(http.StatusBadRequest, Error{Code: 400, Message: \"param should be positive int\"})\n\t}{{else}}\n\tr.{{ paramFieldName . }} = ctx.Param(\"{{ paramName . }}\")\t\n\t{{end}}\n\treturn nil\n}\n{{ end }}\n\ntype {{ titleName .Name }}Request struct {\n\t{{ range .Args}}\n\t{{ titleName .Name }} {{ deliveryType2 $.project .Column }} `json:\"{{ jsonName .ModelName .Name }}\"`\n\t{{ end }}\n}\n\ntype {{ titleName .Name }}Response struct {\n\t{{ range .Return }}\n\t{{if isModel $.project .}}\n\t{{ titleName .Name }} {{ deliveryType $.project . }} `json:\"{{ jsonName .ModelName .Name }}\"`\n\t{{ else }}\n\t{{ titleName .Name }} {{ deliveryType2 $.project . }} `json:\"{{ jsonName .ModelName .Name }}\"`\n\t{{ end }}\n\t{{ end }}\n}\n\nfunc (e *EchoDelivery) {{ titleName .Name }}(c echo.Context) error {\n\tctx := c.Request().Context()\n\tif ider,ok:=c.(UserIDer); ok {\n\t\tif id,has := ider.UserID(); has {\n\t\t\tlog.Println(\"user found in context, id:\", id)\n\t\t\tctx = context.WithValue(ctx, \"userID\", id)\n\t\t} else {\n\t\t\tlog.Println(\"user not found in context\")\n\t\t}\n\t} else {\n\t\tpanic(\"context not ider\")\n\t}\n\n\tvar r {{ titleName .Name }}Request\n\tif err := c.Bind(&r); err != nil {\n\t\treturn c.JSON(http.StatusBadRequest, Error{Code: 400, Message: fmt.Sprintf(\"bad params: %s\", err)})\n\t}\n\t{{ if (hasParam . )}}\n\tif err := set{{ titleName .Name }}Param(c, &r); err != nil {\n\t\treturn nil\n\t}\n\t{{ end }}\n\t{{ range .Args }}\n\t{{ if isModel $.project .Column }}\n\tif validationErrors := convert{{ titleName .Name }}ToDomain(r.{{ titleName .Name }}).Validate(); len(validationErrors) > 0 {\n\t\treturn c.JSON(500, Error{Code: 400, Message: \"validation errors\", ValidationErrors: validationErrors})\n\t}{{ end }}{{end}}\n\t{{ range .Return }} {{ .Name }} ,{{ end }}err := e.usecase.{{ titleName .Name }}(ctx, {{ range .Args }}{{ if isModel $.project .Column}} convert{{ titleName .Name }}ToDomain(r.{{ titleName .Name }}){{ else if eq .Type \"date\" }}time.Time(r.{{ titleName .Name }}){{ else }}r.{{ titleName .Name }}{{ end }},{{ end }})\n\tif err != nil {\n\t\tlog.Printf(\"err call method err=%s\", err)\n\t\treturn c.JSON(500, Error{Code: 500, Message: err.Error()})\n\t}\n\tvar resp {{ titleName .Name }}Response\n\t{{ range .Return }}{{if isModel $.project .}}resp.{{ titleName .Name }} = {{ convertRepresentation .}}{{else}}resp.{{ titleName .Name }} = {{ .Name }}{{end}}\n\t{{ end }}\n   err = c.JSON(200, resp)\n   if err != nil {\n\t   log.Printf(\"err writing response err=%s\", err)\n   }\n   return err\n}\n{{ end }}",
"/usecase.tmpl": "// Code generated by https://github.com/zhuharev/gen DO NOT EDIT.\n\npackage {{ .service.Name }}\n\nimport (\n    \"context\"\n\n    \"{{ .project.GoModule }}/domain\"\n)\n\ntype Usecase interface{\n    {{ range .service.Methods }}{{ titleName .Name }}(ctx context.Context, {{range  .Args }}{{ .Name }} {{ goType2 $.project .Column }}, {{end}}) ({{ range .Return }} {{ goType2 $.project . }} ,{{ end }}error)\n    {{end}}\n}\n\n",
"/usecase_imp.tmpl": "package {{ .service.Name }}\n\nimport (\n    \"context\"\n\n    \"{{ .project.GoModule }}/domain\"\n    {{ range .service.Deps }}\n    \"{{ serviceImport $.project .Name }}\"{{end}}\n\n    \"go.uber.org/zap\"\n)\n\ntype UsecaseImp struct{\n  {{ if (isServiceWithRepo .service)}}repo *Repo{{ end }}\n  {{ range .service.Deps }}\n  {{.Name}}Uc {{servicePkg $.project .Name}}.{{usecaseName $.project .Name}}\n  {{end}}\n\n  *zap.Logger\n}\n\nfunc NewUsecase(logger *zap.Logger, {{ if (isServiceWithRepo .service)}}repo *Repo,{{end}} {{ range .service.Deps }}{{.Name}}Uc {{servicePkg $.project .Name}}.{{usecaseName $.project .Name}},{{end}}) Usecase {\n    return &UsecaseImp{\n        {{ if (isServiceWithRepo .service)}}repo: repo,{{end}}\n        {{ range .service.Deps }}\n        {{.Name}}Uc:  {{.Name}}Uc,\n        {{end}}\n        Logger: logger,\n    }\n}\n\n{{ range .service.Methods }}\nfunc (uc *UsecaseImp) {{ titleName .Name }}(ctx context.Context, {{range  .Args }}{{ .Name }} {{ goType2 $.project .Column }}, {{end}}) ({{ range .Return }} {{ goType2 $.project . }} ,{{ end }}error) {\n    {{ usecaseBody .}}\n}\n{{end}}\n\n",
"/usecase_imp_append.tmpl": "\n{{ range .service.Methods }}\n{{ if (isImplementedMethod $.implemented (titleName .Name))}}{{else}}\nfunc (uc *UsecaseImp) {{ titleName .Name }}(ctx context.Context, {{range  .Args }}{{ .Name }} {{ goType2 $.project .Column }}, {{end}}) ({{ range .Return }} {{ goType2 $.project . }} ,{{ end }}error) {\n    {{ usecaseBody .}}\n}\n{{end}}\n{{end}}\n\n",
}
)