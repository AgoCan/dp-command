package container

import (
	"log"
	"os/exec"
)

type Containerd struct{}

func newContainerd() *Containerd {
	return &Containerd{}
}

// Pull pulls an image from a specified repository using ctr command.
func (c *Containerd) Pull(imageName string) (err error) {
	if err := exec.Command("ctr", "images", "pull", imageName).Run(); err != nil {
		log.Fatalf("Failed to pull image %s: %v", imageName, err)
	}
	log.Printf("Image %s pulled successfully.", imageName)
	return nil
}

// Tag tags an existing image with a new tag using ctr command.
func (c *Containerd) Tag(imageName string, newTag string) {
	// Assuming imageName includes the full path including the current tag if any.
	if err := exec.Command("ctr", "images", "tag", imageName, newTag).Run(); err != nil {
		log.Fatalf("Failed to tag image %s as %s: %v", imageName, newTag, err)
	}
	log.Printf("Image %s tagged as %s successfully.", imageName, newTag)

}

// Push pushes a tagged image to a specified repository using ctr command.
func (c *Containerd) Push(imageName string) (err error) {
	if err := exec.Command("ctr", "images", "push", imageName).Run(); err != nil {
		log.Fatalf("Failed to push image to %s: %v", imageName, err)
	}
	log.Printf("Image pushed to %s successfully.", imageName)
	return nil
}
