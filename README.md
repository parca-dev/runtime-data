# runtime-data

This is a small project to calculate and expose the runtime structs of several interpreters and runtimes.
It heavily relies on existing open-source projects.

## Usage

```
make build
cargo run
```

Check the `out` directory for the generated files.

## Supported version

### Python

- **2.7**: 2.7.x
- **3.x**: 3.3.x, 3.5.x, 3.6.x, 3.7.x, 3.8.x, 3.9.x, 3.10.x, 3.11.x

### Ruby

- **2.6**: 2.6.0, 2.6.3
- **2.7**: 2.7.1, 2.7.4, 2.7.6
- **3.x**: 3.0.0, 3.0.4, 3.1.2, 3.1.3, 3.2.0, 3.2.1

## Acknowledgments

- [Rust](https://github.com/rust-lang)
- [rust-bindgen](https://github.com/rust-lang/rust-bindgen)
- [rbperf](https://github.com/javierhonduco/rbperf)
- [rbspy](https://github.com/rbspy/rbspy)
- [py-spy](https://github.com/benfred/py-spy)

## License

Apache 2
