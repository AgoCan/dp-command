package container

import (
	"log"
	"os/exec"
)

// 利用exec进行执行
type Exec interface {
	Push(imageName string) error
	Pull(imageName string) error
	Tag(imageName string, newTag string)
}

func New() Exec {
	var e Exec
	if commandExists("docker") {
		e = newDocker()
	} else if commandExists("ctr") {
		e = newContainerd()
	} else {
		log.Fatalf("no found docker or ctr.")
	}

	return e
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
