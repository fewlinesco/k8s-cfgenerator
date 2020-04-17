# cfgenerator

Small tool in charge of reading a [JSONNET][JSONNET] from `STDIN` and output a `JSON` in `STDOUT`.

It accepts a list of paths as parameter and read all the files it can found.
For each file it defines an [`ExtVar`][JSONNET_EXTVAR] based on the file name for the variable name and the content of the file for its value.

It's convenient in a Kubernetes context where we have all the configuration in `ConfigMap` and `Secret` but prefer a file to configure our applications.
We can mount the configurations as volumes and use a JSONNET template as input to produce the desired configuration file.

## Usage

`cfgenerator help` or [read this](/main.go#L18).

Some [examples](/examples) are also available.

## Testing

```
make test
```


## Docker

```
make docker-build
```

[JSONNET]: https://github.com/google/go-jsonnet
[JSONNET_EXTVAR]: https://jsonnet.org/ref/stdlib.html
