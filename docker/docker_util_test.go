package docker

import (
	"testing"
)

//
// 做DockerClient的测试
//

func TestListImage(t *testing.T) {
	dockerClient := NewDockerClient()
	dockerClient.ListImage()
}

func TestLoadImage(t *testing.T) {
	tarPath := "/media/liuxu/data/rattrap/nfs.tar"
	dockerClient := NewDockerClient()
	dockerClient.LoadImage(tarPath)
}

func TestSaveImage(t *testing.T) {
	name := "redis"
	tag := "6.2.6"
	dockerClient := NewDockerClient()
	dockerClient.SaveImage(name, tag)
}
