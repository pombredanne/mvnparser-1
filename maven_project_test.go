package mvnparser

import (
	"fmt"
	"testing"
)

func TestMavenProject_WalkMavenProject(t *testing.T) {
	mp := NewMavenProject(".", "")
	mp.WalkMavenProject(func(p *MavenProject){
		if p.Build != nil {
			fmt.Println(p.ArtifactId,p.Version, p.Packaging)
		}
	})
}
