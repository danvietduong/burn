package convert

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spiermar/burn/types"
)

type Frame struct {
	Name       string
	Line       *int
	ResourceId *int
}

type Sample struct {
	StackId int
}

type Stack struct {
	FrameId  int
	ParentId *int
}

type JSProfile struct {
	Frames    []Frame
	Resources []string
	Samples   []Sample
	Stacks    []Stack
}

func ParseJS(r io.Reader) types.Profile {
	jsProfile := JSProfile{}
	d := json.NewDecoder(r)
	d.Decode(&jsProfile)
	rootNode := types.Node{Name: "root", Value: 0, Children: make(map[string]*types.Node)}
	profile := types.Profile{RootNode: rootNode, Stack: []string{}}

	for i := 0; i < len(jsProfile.Samples); i++ {
		profile.OpenStack()
		sample := jsProfile.Samples[i]
		for stackId := &sample.StackId; stackId != nil; stackId = jsProfile.Stacks[*stackId].ParentId {
			stack := jsProfile.Stacks[*stackId]
			frame := jsProfile.Frames[stack.FrameId]
			name := frame.Name
			if name == "" {
				name = "anonymous"
			}
			location := ""
			if frame.Line != nil && frame.ResourceId != nil {
				location = fmt.Sprintf(" %s:%d", jsProfile.Resources[*frame.ResourceId], *frame.Line)
			}
			out := fmt.Sprintf("%s%s", name, location)
			profile.AddFrame(out)
		}

		profile.CloseStack()
	}

	return profile
}
