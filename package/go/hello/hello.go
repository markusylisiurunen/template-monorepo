package hello

import "fmt"

func Say(name string, say func(msg string)) {
	say(fmt.Sprintf("Hello, %s!", name))
}
