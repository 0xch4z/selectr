# selectr ![test][test-badge]

> Select values from objects/arrays with key-path notation.

## Key-path notation

Key-path notation is a simple format for describing how to traverse data structures. It's very similar to C languages, but with limitations. See the reference for more information.

## Usage

```go
m := map[string]interface{}{
    "foo": map[string]interface{}{
        "bar": []interface{}{
            1,
            2,
            3,
        },
    },
}

sel, _ := selectr.Parse(".foo")
sel.Resolve(m) // => map[string]interface{}{"bar": []interface{}{1, 2, 3}}

sel, _ = selectr.Parse(".foo.bar")
sel.Resolve(m) // => []interface{}{1, 2, 3}

sel, _ = selectr.Parse(".foo.bar[1]")
sel.Resolve(m) // => 2
```

## Use cases

- Referencing a dynamic value in a JSON/YAML file:

  Consider you maintain a program which allows users to reference arbitrary values, and one strategy of resolving a value is referencing a symbol in a JSON file.

  ```json
  # example.json
  {
      "accounts": [{
          "id": 123,
          "name": "main"
      }]
  }
  ```

  The end-user references `example.json` with the selectr `.accounts[0].name`.

  ```go
  jsonBytes, _ := os.Open(fileName)

  var m map[string]interface{}
  json.Unmarshal(jsonBytes, &m)

  sel, _ := selectr.Parse(selectrString)
  val := sel.Resolve(m)
  ```

  The `val` is resolved to `"main"`.

- Import dynamic values from dynamic data files.

[test-badge]: https://github.com/Charliekenney23/selectr/workflows/test/badge.svg
