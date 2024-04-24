// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package flags

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
)

// usecolor defines whether the CLI help should use colored output or normal dumb
// colorless terminal formatting.
var usecolor = (isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())) && os.Getenv("TERM") != "dumb"

func init() {
	if usecolor {
		// Annotate all help categories with colors
		cli.AppHelpTemplate = regexp.MustCompile("[A-Z ]+:").ReplaceAllString(cli.AppHelpTemplate, "\u001B[33m$0\u001B[0m")

		// Annotate flag categories with colors (private template, so need to
		// copy-paste the entire thing here...)
		cli.AppHelpTemplate = strings.ReplaceAll(cli.AppHelpTemplate, "{{template \"visibleFlagCategoryTemplate\" .}}", "{{range .VisibleFlagCategories}}\n   {{if .Name}}\u001B[33m{{.Name}}\u001B[0m\n\n   {{end}}{{$flglen := len .Flags}}{{range $i, $e := .Flags}}{{if eq (subtract $flglen $i) 1}}{{$e}}\n{{else}}{{$e}}\n   {{end}}{{end}}{{end}}")
	}
	cli.FlagStringer = FlagString
}

// FlagString prints a single flag in help.
func FlagString(f cli.Flag) string {
	df, ok := f.(cli.DocGenerationFlag)
	if !ok {
		return ""
	}
	needsPlaceholder := df.TakesValue()
	placeholder := ""
	if needsPlaceholder {
		placeholder = "value"
	}

	namesText := cli.FlagNamePrefixer(df.Names(), placeholder)

	defaultValueString := ""
	if s := df.GetDefaultText(); s != "" {
		defaultValueString = " (default: " + s + ")"
	}
	envHint := strings.TrimSpace(cli.FlagEnvHinter(df.GetEnvVars(), ""))
	if envHint != "" {
		envHint = " (" + envHint[1:len(envHint)-1] + ")"
	}
	usage := strings.TrimSpace(df.GetUsage())
	usage = wordWrap(usage, 80)
	usage = indent(usage, 10)

	if usecolor {
		return fmt.Sprintf("\n    \u001B[32m%-35s%-35s\u001B[0m%s\n%s", namesText, defaultValueString, envHint, usage)
	} else {
		return fmt.Sprintf("\n    %-35s%-35s%s\n%s", namesText, defaultValueString, envHint, usage)
	}
}

func indent(s string, nspace int) string {
	ind := strings.Repeat(" ", nspace)
	return ind + strings.ReplaceAll(s, "\n", "\n"+ind)
}

func wordWrap(s string, width int) string {
	var (
		output     strings.Builder
		lineLength = 0
	)

	for {
		sp := strings.IndexByte(s, ' ')
		var word string
		if sp == -1 {
			word = s
		} else {
			word = s[:sp]
		}
		wlen := len(word)
		over := lineLength+wlen >= width
		if over {
			output.WriteByte('\n')
			lineLength = 0
		} else {
			if lineLength != 0 {
				output.WriteByte(' ')
				lineLength++
			}
		}

		output.WriteString(word)
		lineLength += wlen

		if sp == -1 {
			break
		}
		s = s[wlen+1:]
	}

	return output.String()
}

// AutoEnvVars extends all the specific CLI flags with automatically generated
// env vars by capitalizing the flag, replacing . with _ and prefixing it with
// the specified string.
//
// Note, the prefix should *not* contain the separator underscore, that will be
// added automatically.
func AutoEnvVars(flags []cli.Flag, prefix string) {
	for _, flag := range flags {
		envvar := strings.ToUpper(prefix + "_" + strings.ReplaceAll(strings.ReplaceAll(flag.Names()[0], ".", "_"), "-", "_"))

		switch flag := flag.(type) {
		case *cli.StringFlag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *cli.StringSliceFlag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *cli.BoolFlag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *cli.IntFlag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *cli.Int64Flag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *cli.Uint64Flag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *cli.Float64Flag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *cli.DurationFlag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *cli.PathFlag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *BigFlag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *TextMarshalerFlag:
			flag.EnvVars = append(flag.EnvVars, envvar)

		case *DirectoryFlag:
			flag.EnvVars = append(flag.EnvVars, envvar)
		}
	}
}
