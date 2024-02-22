# runtime-data

This project is a collection of tools to extract runtime data from several interpreters and runtimes.
By runtime data, we mean information about the execution of a program, especially stack unwinding and profiling.

## Requirements

- [Devbox](https://jetpack.io/devbox)
  - Also see [CONTRIBUTING.md](CONTRIBUTING.md) for more information about the development environment.
- [direnv](https://direnv.net/)

```sh
make bootstrap
```

- A container runtime (e.g. Docker, Podman)
  - Ability to run cross-platform containers (e.g. Docker Desktop, Podman with QEMU)

```sh
sudo apt-get install qemu binfmt-support qemu-user-static # Install the qemu packages
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes # This step will execute the registering scripts

docker run --rm -t arm64v8/ubuntu uname -m
```

## Usage

### Python

```go
package python

import (
    "fmt"

    "github.com/parca-dev/runtime-data/pkg/python"
)

func main() {
    layouts, err := python.GetLayouts()
    if err != nil {
        return fmt.Errorf("get python layouts: %w", err)
    }

    fmt.Println(layouts)
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
    layouts, err := ruby.GetLayouts()
    if err != nil {
        return fmt.Errorf("get ruby layouts: %w", err)
    }

    fmt.Println(layouts)
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
e.g: structlayout -r python -v 3.9.5 /usr/bin/python3.9

flags:
  -o string
    	output directory to write the layout file (shorthand)
  -output string
    	output directory to write the layout file
  -r string
    	name of the pre-defined runtime, e.g. python, ruby, libc, musl (shorthand)
  -runtime string
    	name of the pre-defined runtime, e.g. python, ruby, libc, musl
  -v string
    	version of the runtime that the layout to generate, e.g. 3.9.5 (shorthand)
  -version string
    	version of the runtime that the layout to generate, e.g. 3.9.5
```

### mergelayout

[embedmd]:# (tmp/mergelayout-help.txt)
```txt
usage: mergelayout -o outputDir <path-to-layout-files>
e.g: mergelayout -o /tmp/merged '/tmp/python/*.yaml'

flags:
  -o string
    	output directory to write the merged layout file (shorthand)
  -output string
    	output directory to write the merged layout file
```

### debdownload
[embedmd]:# (tmp/debdownload-help.txt)
```txt
NAME
  debdownload

FLAGS
  -o, --output STRING       output directory to write the downloaded deb files
  -t, --temp-dir STRING     temporary directory to download deb files
  -u, --url STRING          URL to download deb files from
  -p, --package STRING      package name to download
  -a, --arch STRING         architectures to download
  -c, --constraint STRING   version constraints to download

```


### apkdownload
[embedmd]:# (tmp/apkdownload-help.txt)
```txt
NAME
  apkdownload

FLAGS
  -o, --output STRING       output directory to write the downloaded apk files (default: tmp/bin)
  -t, --temp-dir STRING     temporary directory to download deb files (default: tmp/apk)
  -u, --url STRING          URL to download apk files from
  -p, --package STRING      package name to download
  -a, --arch STRING         architectures to download
  -c, --constraint STRING   version constraints to download

```

### debuginfofind
[embedmd]:# (tmp/debuginfofind-help.txt)
```txt
NAME
  debdownload

FLAGS
  -d, --debuginfo-dir STRING   directory to write the downloaded debuginfo files

```

## Acknowledgments

- [rbperf](https://github.com/javierhonduco/rbperf)
- [rbspy](https://github.com/rbspy/rbspy)
- [py-spy](https://github.com/benfred/py-spy)
- [py-perf](https://github.com/kakkoyun/py-perf)

## License

[Apache 2.0](LICENSE)
