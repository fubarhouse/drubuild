# Drush

Helper package to execute and store responses from Drush.

Intended for scripting purposes for Drupal developers who're using go.

## Usage

```
drush := NewDrush(alias, command)
output, error := drush.Output()
```

## Install

```console
$ go get github.com/fubarhouse/golang-drush/
```

## License

MIT