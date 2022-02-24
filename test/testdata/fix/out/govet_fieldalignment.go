//args: -Egovet
//config: linters-settings.govet.enable=fieldalignment
package p

import "log"

type gvfaGood struct {
	y int32
	x byte
	z byte
}

type gvfaBad struct {
	y int32
	x byte
	z byte
}

func test() {
	g := gvfaGood{1, 2, 3}
	b := gvfaBad{1, 2, 3}
	log.Printf("good: %v; bad: %v", g, b)
}
