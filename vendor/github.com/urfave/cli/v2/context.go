package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
)

// Context is a type that is passed through to
// each Handler action in a cli application. Context
// can be used to retrieve context-specific args and
// parsed command-line options.
type Context struct {
	context.Context
	App           *App
	Command       *Command
	shellComplete bool
	setFlags      map[string]bool
	flagSet       *flag.FlagSet
	parentContext *Context
}

// NewContext creates a new context. For use in when invoking an App or Command action.
func NewContext(app *App, set *flag.FlagSet, parentCtx *Context) *Context {
	c := &Context{App: app, flagSet: set, parentContext: parentCtx}
	if parentCtx != nil {
		c.Context = parentCtx.Context
		c.shellComplete = parentCtx.shellComplete
	}

	c.Command = &Command{}

	if c.Context == nil {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			defer cancel()
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			<-sigs
		}()
		c.Context = ctx
	}

	return c
}

// NumFlags returns the number of flags set
func (c *Context) NumFlags() int {
	return c.flagSet.NFlag()
}

// Set sets a context flag to a value.
func (c *Context) Set(name, value string) error {
	return c.flagSet.Set(name, value)
}

// IsSet determines if the flag was actually set
func (c *Context) IsSet(name string) bool {
	if fs := lookupFlagSet(name, c); fs != nil {
		if fs := lookupFlagSet(name, c); fs != nil {
			isSet := false
			fs.Visit(func(f *flag.Flag) {
				if f.Name == name {
					isSet = true
				}
			})
			if isSet {
				return true
			}
		}

		// XXX hack to support IsSet for flags with EnvVar
		//
		// There isn't an easy way to do this with the current implementation since
		// whether a flag was set via an environment variable is very difficult to
		// determine here. Instead, we intend to introduce a backwards incompatible
		// change in version 2 to add `IsSet` to the Flag interface to push the
		// responsibility closer to where the information required to determine
		// whether a flag is set by non-standard means such as environment
		// variables is available.
		//
		// See https://github.com/urfave/cli/issues/294 for additional discussion
		f := lookupFlag(name, c)
		if f == nil {
			return false
		}

		val := reflect.ValueOf(f)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		filePathValue := val.FieldByName("FilePath")
		if !filePathValue.IsValid() {
			return false
		}

		envVarValues := val.FieldByName("EnvVars")
		if !envVarValues.IsValid() {
			return false
		}

		_, ok := flagFromEnvOrFile(envVarValues.Interface().([]string), filePathValue.Interface().(string))
		return ok
	}

	return false
}

// LocalFlagNames returns a slice of flag names used in this context.
func (c *Context) LocalFlagNames() []string {
	var names []string
	c.flagSet.Visit(makeFlagNameVisitor(&names))
	return names
}

// FlagNames returns a slice of flag names used by the this context and all of
// its parent contexts.
func (c *Context) FlagNames() []string {
	var names []string
	for _, ctx := range c.Lineage() {
		ctx.flagSet.Visit(makeFlagNameVisitor(&names))

	}
	return names
}

// FlagNames returns a slice of flag names used in this context.
//func (c *Context) FlagNames() (names []string) {
//	for _, f := range c.Command.Flags {
//		name := strings.Split(f.GetName(), ",")[0]
//		if name == "help" {
//			continue
//		}
//		names = append(names, name)
//	}
//	return
//}

// GlobalFlagNames returns a slice of global flag names used by the app.
//func (c *Context) GlobalFlagNames() (names []string) {
//	for _, f := range c.App.Flags {
//		name := strings.Split(f.GetName(), ",")[0]
//		if name == "help" || name == "version" {
//			continue
//		}
//		names = append(names, name)
//	}
//	return names
//}

// Lineage returns *this* context and all of its ancestor contexts in order from
// child to parent
func (c *Context) Lineage() []*Context {
	var lineage []*Context

	for cur := c; cur != nil; cur = cur.parentContext {
		lineage = append(lineage, cur)
	}

	return lineage
}

// value returns the value of the flag corresponding to `name`
func (c *Context) value(name string) interface{} {
	return c.flagSet.Lookup(name).Value.(flag.Getter).Get()
}

// Args returns the command line arguments associated with the context.
func (c *Context) Args() Args {
	ret := args(c.flagSet.Args())
	return &ret
}

// NArg returns the number of the command line arguments.
func (c *Context) NArg() int {
	return c.Args().Len()
}

func lookupFlag(name string, ctx *Context) Flag {
	for _, c := range ctx.Lineage() {
		if c.Command == nil {
			continue
		}

		for _, f := range c.Command.Flags {
			for _, n := range f.Names() {
				if n == name {
					return f
				}
			}
		}
	}

	if ctx.App != nil {
		for _, f := range ctx.App.Flags {
			for _, n := range f.Names() {
				if n == name {
					return f
				}
			}
		}
	}

	return nil
}

func lookupFlagSet(name string, ctx *Context) *flag.FlagSet {
	for _, c := range ctx.Lineage() {
		if f := c.flagSet.Lookup(name); f != nil {
			return c.flagSet
		}
	}

	return nil
}

func copyFlag(name string, ff *flag.Flag, set *flag.FlagSet) {
	switch ff.Value.(type) {
	case Serializer:
		_ = set.Set(name, ff.Value.(Serializer).Serialize())
	default:
		_ = set.Set(name, ff.Value.String())
	}
}

func normalizeFlags(flags []Flag, set *flag.FlagSet) error {
	visited := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})
	for _, f := range flags {
		parts := f.Names()
		if len(parts) == 1 {
			continue
		}
		var ff *flag.Flag
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if visited[name] {
				if ff != nil {
					return errors.New("Cannot use two forms of the same flag: " + name + " " + ff.Name)
				}
				ff = set.Lookup(name)
			}
		}
		if ff == nil {
			continue
		}
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if !visited[name] {
				copyFlag(name, ff, set)
			}
		}
	}
	return nil
}

func makeFlagNameVisitor(names *[]string) func(*flag.Flag) {
	return func(f *flag.Flag) {
		nameParts := strings.Split(f.Name, ",")
		name := strings.TrimSpace(nameParts[0])

		for _, part := range nameParts {
			part = strings.TrimSpace(part)
			if len(part) > len(name) {
				name = part
			}
		}

		if name != "" {
			*names = append(*names, name)
		}
	}
}

type requiredFlagsErr interface {
	error
	getMissingFlags() []string
}

type errRequiredFlags struct {
	missingFlags []string
}

func (e *errRequiredFlags) Error() string {
	numberOfMissingFlags := len(e.missingFlags)
	if numberOfMissingFlags == 1 {
		return fmt.Sprintf("Required flag %q not set", e.missingFlags[0])
	}
	joinedMissingFlags := strings.Join(e.missingFlags, ", ")
	return fmt.Sprintf("Required flags %q not set", joinedMissingFlags)
}

func (e *errRequiredFlags) getMissingFlags() []string {
	return e.missingFlags
}

func checkRequiredFlags(flags []Flag, context *Context) requiredFlagsErr {
	var missingFlags []string
	for _, f := range flags {
		if rf, ok := f.(RequiredFlag); ok && rf.IsRequired() {
			var flagPresent bool
			var flagName string
			for _, key := range f.Names() {
				if len(key) > 1 {
					flagName = key
				}

				if context.IsSet(strings.TrimSpace(key)) {
					flagPresent = true
				}
			}

			if !flagPresent && flagName != "" {
				missingFlags = append(missingFlags, flagName)
			}
		}
	}

	if len(missingFlags) != 0 {
		return &errRequiredFlags{missingFlags: missingFlags}
	}

	return nil
}
