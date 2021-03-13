### Usage

[Generic](doc/usage.html) to both *gen* and *dump* commands.

[Specific to *gen*](doc/gen-usage.html).

If generating program source code, two difficulties may appear:
 1. For a satisfactory experience when debugging stack traces,
template expansions must match the number of lines in the template source code.
Workaround examples may be found in [```_test_src/```](./_test_src/).
 2. [gofmt](https://golang.org/cmd/gofmt/) is confused by template syntax e.g. "{{...}}".

"gemp" is a portmanteau of "Go-tEMPlate".

### See also

Other Go-based templating utilities, targeting somewhat different use cases:
 - [stringer](https://pkg.go.dev/golang.org/x/tools@v0.1.0/cmd/stringer)
 - [Kubernetes templates](https://pkg.go.dev/k8s.io/kubernetes/pkg/kubectl/util/templates)
