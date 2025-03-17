package container

import (
	"log"
	"os/exec"
)

type Docker struct{}

func newDocker() *Docker {
	return &Docker{}
}

// Pull pulls an image from a specified repository using docker command.
func (d *Docker) Pull(imageName string) {
	if err := exec.Command("docker", "pull", imageName).Run(); err != nil {
		log.Fatalf("Failed to pull image %s: %v", imageName, err)
	}
	log.Printf("Image %s pulled successfully.", imageName)
}

// Tag tags an existing image with a new tag using docker command.
func (d *Docker) Tag(imageName string, newTag string) {
	if err := exec.Command("docker", "tag", imageName, newTag).Run(); err != nil {
		log.Fatalf("Failed to tag image %s as %s: %v", imageName, newTag, err)
	}
	log.Printf("Image %s tagged as %s successfully.", imageName, newTag)
}

// Push pushes a tagged image to a specified repository using docker command.
func (d *Docker) Push(imageName string) {
	if err := exec.Command("docker", "push", imageName).Run(); err != nil {
		log.Fatalf("Failed to push image to %s: %v", imageName, err)
	}
	log.Printf("Image pushed to %s successfully.", imageName)
}
