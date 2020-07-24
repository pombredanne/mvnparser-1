package mvnparser

import (
	"fmt"
	"os"
	"testing"
)

func TestMavenProject_WalkMavenProject(t *testing.T) {
	pwd,_ := os.Getwd()
	fmt.Println("current dir: ", pwd)
	mp := NewMavenProject(".", "")
	mp.WalkMavenProject(func(p *MavenProject){
		if p.IsNeededToBuild() {
			fmt.Println(p.ArtifactId, p.Version, p.Packaging, p.Parent.RelativePath, p.Parent.ArtifactId)
		}
	})
}
