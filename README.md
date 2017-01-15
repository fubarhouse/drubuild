# Drush

Helper package to execute and store responses from Drush.

Intended for scripting purposes for Drupal developers who're using go.

## Usage

****Running a single command****
```
drush := NewDrush(alias, command)
output, error := drush.Output()
```
****Running an infinite amount of commands****
````
drushList := NewDrushList()
command1 := NewDrush("none", "", false)
command2 := NewDrush("none", "", false)
drushList.Add(command1, command2)
drushList.RemoveIndex(1)
outputArray, errorArray := drushList.Output()
````

## Install

```console
$ go get github.com/fubarhouse/golang-drush/
```

## License

MIT