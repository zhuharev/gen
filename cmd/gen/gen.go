package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
	"github.com/zhuharev/gen"
	"gopkg.in/yaml.v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			CmdGenerate,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

// CmdGenerate run generator
var CmdGenerate = &cli.Command{
	Name: "generate",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "c",
		},
		&cli.StringFlag{
			Name: "o",
		},
	},
	Action: runGenerator,
}

func runGenerator(ctx *cli.Context) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open(ctx.String("c"))
	if err != nil {
		return err
	}
	defer f.Close()
	var p gen.Project
	err = yaml.NewDecoder(f).Decode(&p)
	if err != nil {
		return err
	}
	err = gen.Generate(ctx.String("o"), p)
	if err != nil {
		return err
	}
	return nil

}

// var testProject = gen.Project{
// 	Name:     "testdata",
// 	GoModule: "github.com/zhuharev/test",
// 	Models: []gen.Model{
// 		UserModel,
// 		DialogModel,
// 		DialogMemberModel,
// 	},
// 	Services: []gen.Service{
// 		{
// 			Name:   "auth",
// 			Plural: "auth",
// 			Deps:   []string{"user"},
// 			Methods: []gen.Method{
// 				{
// 					Name: "login",
// 					Args: []gen.MethodArg{
// 						{
// 							Column: gen.Column{
// 								Name: "phone",
// 								Type: gen.ColumnString,
// 							},
// 						},
// 						{
// 							Column: gen.Column{
// 								Name: "password",
// 								Type: gen.ColumnString,
// 							},
// 						},
// 					},
// 				},
// 				{
// 					Name: "sendCode",
// 				},
// 			},
// 		},
// 		{
// 			Name:   "user",
// 			Plural: "users",
// 			Models: []string{"user"},
// 			Methods: []gen.Method{
// 				{
// 					Name: "byId",
// 					Args: []gen.MethodArg{
// 						{
// 							IsRestParam: true,
// 							Column: gen.Column{
// 								Name:      "id",
// 								Type:      gen.ColumnInt,
// 								ModelName: "user",
// 							},
// 						},
// 					},
// 					Return: []gen.Column{
// 						{
// 							Name:      "user",
// 							Type:      gen.ColumnModel,
// 							ModelName: "user",
// 						},
// 					},
// 				},
// 				{
// 					Name: "byPhone",
// 					Args: []gen.MethodArg{
// 						{
// 							Column: gen.Column{
// 								Name:      "phone",
// 								Type:      gen.ColumnString,
// 								ModelName: "user",
// 							},
// 						},
// 					},
// 					Return: []gen.Column{
// 						{
// 							Name:      "user",
// 							Type:      gen.ColumnModel,
// 							ModelName: "user",
// 						},
// 					},
// 				},
// 				{
// 					Name: "update",
// 					Args: []gen.MethodArg{
// 						{
// 							Column: gen.Column{
// 								Name:      "user",
// 								Type:      gen.ColumnModel,
// 								ModelName: "user",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Name:   "messenger",
// 			Plural: "messenger",
// 			Models: []string{"dialog", "dialogMember"},
// 			Methods: []gen.Method{
// 				{
// 					Name: "byId",
// 					Args: []gen.MethodArg{
// 						{
// 							IsRestParam: true,
// 							Column: gen.Column{
// 								Name:      "id",
// 								Type:      gen.ColumnInt,
// 								ModelName: "dialog",
// 							},
// 						},
// 					},
// 					Return: []gen.Column{
// 						{
// 							Name:      "dialog",
// 							Type:      gen.ColumnModel,
// 							ModelName: "dialog",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// }

// var (
// 	UserModel = gen.Model{
// 		Name: "user",
// 		Type: gen.ModelEntity,
// 		Columns: []gen.Column{
// 			{
// 				Name: "id",
// 				Type: gen.ColumnInt,
// 			},
// 			{
// 				Name: "username",
// 				Type: gen.ColumnString,
// 			},
// 			{
// 				Name: "firstName",
// 				Type: gen.ColumnString,
// 			},
// 			{
// 				Name: "lastName",
// 				Type: gen.ColumnString,
// 			},
// 			{
// 				Name: "phone",
// 				Type: gen.ColumnString,
// 			},
// 			{
// 				Name: "createdAt",
// 				Type: gen.ColumnTime,
// 			},
// 			{
// 				Name: "updatedAt",
// 				Type: gen.ColumnTime,
// 			},
// 		},
// 	}
// 	DialogModel = gen.Model{
// 		Name: "dialog",
// 		Type: gen.ModelEntity,
// 		Columns: []gen.Column{
// 			{
// 				Name: "id",
// 				Type: gen.ColumnInt,
// 			},
// 			{
// 				Name: "createdAt",
// 				Type: gen.ColumnTime,
// 			},
// 		},
// 	}
// 	DialogMemberModel = gen.Model{
// 		Name: "dialogMember",
// 		Type: gen.ModelRelation,
// 		Columns: []gen.Column{
// 			{
// 				Name: "userId",
// 				Type: gen.ColumnInt,
// 			},
// 			{
// 				Name: "dialogId",
// 				Type: gen.ColumnInt,
// 			},
// 			{
// 				Name: "createdAt",
// 				Type: gen.ColumnTime,
// 			},
// 		},
// 	}
// )
