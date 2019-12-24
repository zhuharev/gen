package gen

import (
	"regexp"
	"strconv"
	"strings"
)

// Project is main app container
type Project struct {
	Name     string `json:"name,omitempty"`
	GoModule string `yaml:"goModule" json:"goModule,omitempty"`
	Version  string `json:"version,omitempty" yaml:"version"`

	Services      []Service      `json:"services,omitempty"`
	Models        []Model        `json:"models,omitempty"`
	Forms         []Form         `json:"forms,omitempty"`
	UploadServers []UploadServer `json:"upload_servers,omitempty" yaml:"uploadServers"`
}

type UploadServer struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	FileName    string `json:"fileName" yaml:"fileName"`

	Preprocess *Preprocess `json:"preprocess,omitempty"`
	Validation *Validation `json:"validation,omitempty"`
}

type Preprocess struct {
	Thumbnail *struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"thumbnail"`
}

// Service container contains all info about service
type Service struct {
	Name   string `json:"name,omitempty"`
	Plural string `json:"plural,omitempty"`
	Type   string `json:"type,omitempty"`

	Models  []string `json:"models,omitempty"`
	Deps    []Dep    `json:"deps,omitempty"`
	Methods []Method `json:"methods,omitempty"`
}

type Dep struct {
	Name string `json:"name" yaml:"name"`
}

// Model is buisness-logic model used in services
type Model struct {
	Name       string     `json:"name,omitempty"`
	Plural     string     `json:"plural"`
	Type       ModelType  `json:"type,omitempty"`
	Columns    []Column   `json:"columns,omitempty"`
	Relations  []Relation `json:"relations,omitempty"`
	Extensions KVs        `json:"extensions,omitempty"`
}

type KV struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type KVs []KV

// ValueByName search and return value by name
func (kvs KVs) ValueByName(name string) string {
	for _, v := range kvs {
		if v.Name == name {
			return v.Value
		}
	}
	return ""
}

// Relation describe model relation
type Relation struct {
	To                 string `json:"to"`
	Name               string `json:"name"`
	Ref                string `json:"ref"`
	Title              string `json:"title"`
	FilledFieldName    string `json:"filledFieldName" yaml:"filledFieldName"`
	RelationTitleField string `json:"relationTitleField" yaml:"relationTitleField"`

	Relation RelationSibling `json:"-"`
}

type RelationSibling struct {
	Name            string `json:"name"`
	To              string `json:"to"`
	Title           string `json:"title"`
	FilledFieldName string `json:"filledFieldName" yaml:"filledFieldName"`
}

func newRelationSibling(relation Relation) RelationSibling {
	return RelationSibling{
		Name:            relation.Name,
		FilledFieldName: relation.FilledFieldName,
		To:              relation.To,
		Title:           relation.Title,
	}
}

// Form for vt
type Form struct {
	Name                string               `json:"name,omitempty"`
	Title               string               `json:"title,omitempty"`
	Type                string               `json:"type"`
	CreateMethod        string               `json:"createMethod,omitempty" yaml:"createMethod"`
	UpdateMethod        string               `json:"updateMethod,omitempty" yaml:"updateMethod"`
	ListMethod          string               `json:"listMethod,omitempty" yaml:"listMethod"`
	DeleteMethod        string               `json:"deleteMethod,omitempty" yaml:"deleteMethod"`
	AutocompleteMethods []AutocompleteMethod `json:"autocompleteMethods" yaml:"autocompleteMethods"`
	ListColumns         []Column             `json:"listColumns,omitempty" yaml:"listColumns"`
}

// AutocompleteMethod
type AutocompleteMethod struct {
	ColumnName string `json:"columnName" yaml:"columnName"`
	Method     string `json:"method"`
}

// ModelType describe model
type ModelType string

const (
	ModelEntity   ModelType = "entity"
	ModelRelation ModelType = "relation"
)

// Column is model column description
type Column struct {
	Name            string        `json:"name,omitempty"`
	Title           string        `json:"title,omitempty" yaml:"title,omitempty"`
	Type            ColumnType    `json:"type,omitempty"`
	ModelName       string        `yaml:"modelName" json:"modelName,omitempty"`
	Enums           []Enum        `yaml:"enums,omitempty" json:"enums"`
	Validation      *Validation   `json:"validation"`
	Display         Display       `json:"display"`
	UploadServer    string        `json:"uploadServer" yaml:"uploadServer"`
	IsRelationField bool          `json:"isRelationField" yaml:"isRelationField"`
	IDsField        string        `json:"idsField" yaml:"idsField"`
	IDsFieldTitle   string        `json:"idsFieldTitle" yaml:"idsFieldTitle"`
	Control         ColumnControl `json:"control" yaml:"control"`
	ReadOnly        bool          `json:"readonly,omitempty" yaml:"readonly,omitempty"`
}

type ColumnControl string

const (
	ColumnControlTextarea = "textarea"
)

type Display struct {
	Prefix string `json:"prefix,omitempty"`
}

type Validation struct {
	Required        bool             `json:"required,omitempty"`
	Unique          bool             `json:"unique,omitempty"`
	IsPhone         bool             `json:"isPhone,omitempty" yaml:"isPhone"`
	IsEmail         bool             `json:"isEmail,omitempty" yaml:"isEmail"`
	IsCyrillic      bool             `json:"isCyrillic,omitempty" yaml:"isCyrillic"`
	IsAlias         bool             `json:"isAlias,omitempty" yaml:"isAlias"`
	CustomFunctions []CustomFunction `json:"customFunctions" yaml:"customFunctions"`
	Max             *struct {
		Description string `json:"description"`
		Value       int    `json:"value"`
	}
	Min *struct {
		Description string `json:"description"`
		Value       int    `json:"value"`
	}
	MinAge *struct {
		Description string `json:"description"`
		Value       string `json:"value"`
	} `yaml:"minAge" json:"minAge"`
	MaxAge *struct {
		Description string `json:"description"`
		Value       string `json:"value"`
	} `yaml:"maxAge" json:"maxAge"`
	Mime  string `json:"mime"`
	Image *struct {
		AspectRatio float64 `json:"aspectRatio" yaml:"aspectRatio"`
		MinWidth    int     `json:"minWidth" yaml:"minWidth"`
		MaxWidth    int     `json:"maxWidth" yaml:"maxWidth"`
		MinHeight   int     `json:"minHeight" yaml:"minHeight"`
		MaxHeight   int     `json:"maxHeight" yaml:"maxHeight"`
	} `yaml:"image" json:"image,omitempty"`
}

type CustomFunction struct {
	Name string `json:"name"`
}

type ValidationError struct {
	FieldName  string
	Len        *int
	NeedNumber *bool
	Max        *int
}

var numberRegex = regexp.MustCompile("^[0-9]+$")

func (c Column) IsValid(fieldData string) (bool, *ValidationError) {
	if c.Validation == nil {
		return true, nil
	}

	fieldData = strings.TrimSpace(fieldData)

	if c.Validation.Required && len(fieldData) == 0 {
		return false, &ValidationError{FieldName: c.Name}
	}

	if c.Validation.IsPhone {
		var expectedLen = 10
		if len(fieldData) != expectedLen {
			return false, &ValidationError{FieldName: c.Name, Len: &expectedLen}
		}
		if !numberRegex.MatchString(fieldData) {
			var needNumber = true
			return false, &ValidationError{FieldName: c.Name, NeedNumber: &needNumber}
		}
	}

	if c.Validation.Max != nil {
		//TODO: check column type and validate string
		intVal, _ := strconv.Atoi(fieldData)
		if intVal > c.Validation.Max.Value {
			return false, &ValidationError{FieldName: c.Name, Max: &c.Validation.Max.Value}
		}

	}

	return true, nil
}

// Enum value description for enum
type Enum struct {
	Name  string `yaml:"name" json:"name"`
	Title string `yaml:"title" json:"title"`
}

type ColumnType string

const (
	ColumnString   ColumnType = "string"
	ColumnInt      ColumnType = "int"
	ColumnTime     ColumnType = "time"
	ColumnModel    ColumnType = "model"
	ColumnFloat    ColumnType = "float"
	ColumnBool     ColumnType = "bool"
	ColumnPassword ColumnType = "password"
	ColumnEnum     ColumnType = "enum"
	ColumnFile     ColumnType = "file"
	ColumnDate     ColumnType = "date"
)

func (t ColumnType) String() string {
	return string(t)
}

func (t ColumnType) IsModel() bool {
	if isArray(t) {
		return ColumnModel == t[2:]
	}
	return t == ColumnModel
}

func (t ColumnType) IsEnum() bool {
	return t == ColumnEnum
}

func (t ColumnType) IsArray() bool {
	return isArray(t)
}

// Method describe buisness-logic operation
type Method struct {
	Name            string         `json:"name,omitempty"`
	Description     string         `json:"description"`
	Args            []MethodArg    `json:"args,omitempty"`
	Return          []Column       `json:"return,omitempty"`
	Menus           []string       `json:"menus,omitempty"`
	RunInBackground *BackgroundJob `json:"runInBackground" yaml:"runInBackground"`
}

type BackgroundJob struct {
	Config string `json:"config" yaml:"config"`
}

type MethodArg struct {
	Column      `json:"column,omitempty"`
	IsRestParam bool `yaml:"isRestParam" json:"isRestParam,omitempty"`
}
