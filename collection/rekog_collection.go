package collection

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

type RecognitionEngine struct {
	CollectionName string

	path  string
	rekog *rekognition.Client
}

func (ths *RecognitionEngine) Initialize() {
	if ths.CollectionName == "" {
		ths.CollectionName = "imagine"
	}

	_, err := ths.rekog.DescribeCollection(context.TODO(), &rekognition.DescribeCollectionInput{
		CollectionId: &ths.CollectionName,
	})
	if err != nil {
		var rnf *types.ResourceNotFoundException
		if errors.As(err, &rnf) {
			fmt.Println("Creating new collection")
			_, err := ths.rekog.CreateCollection(context.TODO(), &rekognition.CreateCollectionInput{
				CollectionId: &ths.CollectionName,
			})
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	ths.Write()
}

func (ths *RecognitionEngine) AddFace(face image.Image, name string) error {
	var buffer bytes.Buffer
	err := png.Encode(&buffer, face)
	if err != nil {
		return err
	}

	_, err = ths.rekog.IndexFaces(context.TODO(), &rekognition.IndexFacesInput{
		CollectionId: &ths.CollectionName,
		Image: &types.Image{
			Bytes: buffer.Bytes(),
		},
		ExternalImageId: &name,
		MaxFaces:        aws.Int32(1),
	})
	if err != nil {
		return err
	}

	return nil
}

func (ths *RecognitionEngine) Search(face image.Image) ([]string, error) {
	var buffer bytes.Buffer
	err := png.Encode(&buffer, face)
	if err != nil {
		return nil, err
	}

	resp, err := ths.rekog.SearchFacesByImage(context.TODO(), &rekognition.SearchFacesByImageInput{
		CollectionId: &ths.CollectionName,
		Image: &types.Image{
			Bytes: buffer.Bytes(),
		},
		MaxFaces:           aws.Int32(1),
		FaceMatchThreshold: aws.Float32(50),
	})
	if err != nil {
		return nil, err
	}

	matches := []string{}
	for _, match := range resp.FaceMatches {
		matches = append(matches, *match.Face.ExternalImageId)
	}
	return matches, nil
}

func (ths *RecognitionEngine) Write() {
	if data, err := json.MarshalIndent(ths, "", "\t"); err != nil {
		panic(err)
	} else {
		if err := os.WriteFile(ths.path, data, os.FileMode(int(0777))); err != nil {
			panic(err)
		}
	}
}
