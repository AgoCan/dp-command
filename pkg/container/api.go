package container

// 利用exec进行执行
type Exec interface {
	Push()
	Pull()
	Tag()
}

func New(runType string) Exec {
	var e Exec
	if runType == "containerd" {
		e = newContainerd()
	} else if runType == "docker" {
		e = newDocker()
	} else {

	}

	return e
}
