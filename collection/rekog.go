package collection

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

type FaceBox struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (ths *Image) DetectFaces() ([]FaceBox, error) {
	d, err := os.ReadFile(ths.FullPath)
	if err != nil {
		return nil, err
	}

	out, err := ths.collection.Rekog.DetectFaces(context.TODO(), &rekognition.DetectFacesInput{
		Attributes: []types.Attribute{"DEFAULT"},
		Image: &types.Image{
			Bytes: d,
		},
	})
	if err != nil {
		return nil, err
	}

	faces := []FaceBox{}
	for _, face := range out.FaceDetails {
		if *face.Confidence > 90 {
			faces = append(faces, FaceBox{
				X:      int(*face.BoundingBox.Left * 100),
				Y:      int(*face.BoundingBox.Top * 100),
				Width:  int(*face.BoundingBox.Width * 100),
				Height: int(*face.BoundingBox.Height * 100),
			})
		}
	}

	return faces, nil
}
