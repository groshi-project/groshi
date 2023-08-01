# lookupcfg
Easily define and populate your configs from any kind of source using lookup function!

## Installation 
```shell
go get github.com/jieggii/lookupcfg
```

## Usage
It's so simple!

First of all, define your config struct and create config instance:
```go
package main

import "github.com/jieggii/lookupcfg"

func main() {
    type Config struct {
        Host string `env:"HOST"` // use tags to define source name (key) 
                                 // and value name in this source (value)
	Port int `env:"PORT"`
    }
	config := Config{}
}
```

Then choose lookup function which will be used to fill fields of the config struct.
In this example we will define own lookup function, which simply uses `os.LookupEnv` function.
```go
package main

import (
	"github.com/jieggii/lookupcfg"
	"os"
)

func lookup(key string) (value string, found bool) {
	return os.LookupEnv(key)
}

func main() {
	type Config struct {
		Host string `env:"HOST"` // use tags to define source name (key) 
                                         // and value name in this source (value)
		Port int `env:"PORT"`
	}
	config := Config{}
}
```

Let's finally fill our config with some values from environment variables and print the result!
```go
package main

import (
	"fmt"
	"github.com/jieggii/lookupcfg"
	"os"
)

func lookup(key string) (value string, found bool) {
	return os.LookupEnv(key)
}

func main() {
	type Config struct {
		Host string `env:"HOST"` // use tags to define source name (key)
                                         // and value name in this source (value)
		Port int `env:"PORT"`
	}
	config := Config{}
	result := lookupcfg.PopulateConfig(
		"env",   // name of the source defined in struct's field tags
		lookup,  // our lookup function
		&config, // pointer to our config instance
	) // populating our config instance with values from environmental variables

	fmt.Printf("Population result: %+v\n", result)
	// >>> Population result: &{MissingFields:[] IncorrectTypeFields:[]}

	// print our populated config instance
	fmt.Printf("My config: %+v\n", config)
	// >>> My config: {Host:localhost Port:8888}
}
```
Done!

More examples can be found in the [examples](https://github.com/jieggii/lookupcfg/tree/master/examples) directory  (´｡• ◡ •｡`) ♡.
