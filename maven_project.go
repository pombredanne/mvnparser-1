package mvnparser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// Represent a POM file
type MavenProject struct {
	XMLName      *xml.Name     `xml:"project"`
	ModelVersion string       `xml:"modelVersion"`
	Parent	     *Parent		  `xml:"parent"`
	GroupId      string       `xml:"groupId"`
	ArtifactId   string       `xml:"artifactId"`
	Version      string       `xml:"version"`
	Packaging    string       `xml:"packaging"`
	Name         string       `xml:"name"`

	Modules      []string     `xml:"modules>module"`
	Build        *Build                `xml:"build"`

	RelativePath string      `xml:"-"`
	modules      []*MavenProject `xml:"-"`
}

type Parent struct {
	GroupId      string       `xml:"groupId"`
	ArtifactId   string       `xml:"artifactId"`
	Version      string       `xml:"version"`
	RelativePath string 	  `xml:"relativePath"`
}

type Build struct {
	FinalName string `xml:"finalName"`
	Plugins   []*Plugin `xml:plugins>plugin`
}

type Plugin struct {
	Configuration *Configuration `xml:"configuration"`
}

type Configuration struct {
	MainClass string `xml:"mainClass"`
}

func (mp *MavenProject)  AddSubModules(subModule *MavenProject) {
	mp.modules = append(mp.modules, subModule)
}

func (mp *MavenProject) IsNeededToBuild() bool {
	if mp.Packaging != "" && mp.Packaging != "pom" && mp.Build.FinalName != "" {
		return true
	}
	return false
}

type WalkFunc func(project *MavenProject)
func (mp *MavenProject) WalkMavenProject(wf WalkFunc) {
	wf(mp)
	if len(mp.modules) != 0 {
		for _, m := range mp.modules {
			m.WalkMavenProject(wf)
		}
	}
}

func NewMavenProject(relativePath, version string) *MavenProject {
	mp := &MavenProject{
		RelativePath: relativePath,
	}

	// 加载pom文件到mp
	data,err:=ioutil.ReadFile(filepath.Join(relativePath, "pom.xml"))
	if err != nil {
		panic(errors.New(fmt.Sprintf("%s at relativePath: %s", err, relativePath)))
	}
	if err:=xml.Unmarshal(data, mp);err!=nil {
		panic(err)
	}

	if mp.Version == "" {
		mp.Version = version
	}

	if len(mp.Modules) != 0 {
		for _, m := range mp.Modules {
			mp.AddSubModules(NewMavenProject(filepath.Join(relativePath, m), mp.Version))
		}
	}

	return mp
}
