# Drush

Helper package to execute and store responses from Drush.

Intended for scripting purposes for Drupal developers who're using go.

## Usage

****Running a single command****
```
drush := NewDrushCommand()
drush.Set("myalias", "cc all", false)
_, cmdErr := drush.Output()
```
****Running an infinite amount of commands****
````
commands := NewDrushCommandList()
command1 := NewDrushCommand()
command1.Set("", "cc drush", false)
commands.Add(command1)
_, drushError := commands.Output()
````

## Install

```console
$ go get github.com/fubarhouse/golang-drush
```

## License

MIT