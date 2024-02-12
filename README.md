# runtime-data

This project is a collection of tools to extract runtime data from several interpreters and runtimes.
By runtime data, we mean information about the execution of a program, especially stack unwinding and profiling.

## Usage

### Python

```go
package python

import (
    "fmt"

    "github.com/parca-dev/runtime-data/pkg/python"
)

func main() {
    versions, err := python.GetVersions()
    if err != nil {
        return fmt.Errorf("get python versions: %w", err)
    }

    fmt.Println(versions)
}
```

### Ruby

```go
package ruby

import (
    "fmt"

    "github.com/parca-dev/runtime-data/pkg/ruby"
)

func main() {
    versions, err := ruby.GetVersions()
    if err != nil {
        return fmt.Errorf("get ruby versions: %w", err)
    }

    fmt.Println(versions)
}
```

## Supported runtimes and versions

### Python

- **2.7**: 2.7.x
- **3.x**: 3.3.x, 3.5.x, 3.6.x, 3.7.x, 3.8.x, 3.9.x, 3.10.x, 3.11.x

### Ruby

- **2.6**: 2.6.0, 2.6.3
- **2.7**: 2.7.1, 2.7.4, 2.7.6
- **3.x**: 3.0.0, 3.0.4, 3.1.2, 3.1.3, 3.2.0, 3.2.1

## Tools

Under the `cmd` directory, you can find the following tools:

**structlayout**: Extracts the memory layout using the given map (a struct annotated with certain struct tags).
**mergelayout**: Merges the given layouts into groups of layouts.

### structlayout

[embedmd]:# (tmp/structlayout-help.txt)
```txt
usage: structlayout [flags] <path-to-elf>
e.g: structlayout -m python -v 3.9.5 /usr/bin/python3.9

flags:
  -o string
    	output directory to write the layout file (shorthand)
  -output string
    	output directory to write the layout file
  -r string
    	name of the pre-defined runtime, e.g. python, ruby, libc (shorthand)
  -runtime string
    	name of the pre-defined runtime, e.g. python, ruby, libc
  -v string
    	version of the runtime that the layout to generate, e.g. 3.9.5 (shorthand)
  -version string
    	version of the runtime that the layout to generate, e.g. 3.9.5
```

### mergelayout

[embedmd]:# (tmp/mergelayout-help.txt)
```txt
usage: mergelayout -o outputDir <path-to-layout-files>
e.g: mergelayout -o /tmp/merged /tmp/python/*.yaml

flags:
  -i string
    	path to struct layout files (shorthand)
  -input string
    	path to struct layout files. Can be a glob pattern.
  -o string
    	output directory to write the merged layout file (shorthand)
  -output string
    	output directory to write the merged layout file
```

## Acknowledgments

- [rbperf](https://github.com/javierhonduco/rbperf)
- [rbspy](https://github.com/rbspy/rbspy)
- [py-spy](https://github.com/benfred/py-spy)
- [py-perf](https://github.com/kakkoyun/py-perf)

## License

[Apache 2.0](LICENSE)
