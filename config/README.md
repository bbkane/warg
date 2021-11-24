`configreader` allows using different types of config files with the app.

the JSON and YAML configreaders duplicate a ton of code

Differences:

- JSON can only unmarshal into `map[string]interface{}` while YAML can only unmarshal into `map[interface{}]interface{}`.
- JSON can only parse float64s, while YAML can also produce ints

# ConfigReader Error Cases

These need to be added to get configs working well enough :)

config not passed -> flag not set
config passed, file doesn't exist -> flag not set  # command should error if a flag isn't set properly
config passed, file exists, can't unmarshall -> ERROR
config passed, file exists, can unmarshall, invalid path-> ERROR
config passed, file exists, can unmarshall, valid path, path not in config -> flag not set
config passed, file exists, can unmarshall, valid path, path in config, value error -> ERROR
config passed, file exists, can unmarshall, valid path, path in config, value created -> flag set
