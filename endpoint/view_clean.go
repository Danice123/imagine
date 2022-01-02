package endpoint

type CleanData struct {
	RemovedImages []string
}

// func Clean(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
// 	tags, err := collection.New(ths.Root)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	data := CleanData{
// 		RemovedImages: []string{},
// 	}
// 	for path := range tags.Mapping {
// 		if _, err := os.Stat(filepath.Join(ths.Root, path)); err != nil {
// 			data.RemovedImages = append(data.RemovedImages, path)
// 		}
// 	}

// 	var temp = template.New("Clean")
// 	if html, err := os.ReadFile(filepath.Join("templates", "clean.html")); err != nil {
// 		panic(err.Error())
// 	} else {
// 		if _, err := temp.Parse(string(html)); err != nil {
// 			panic(err.Error())
// 		} else {
// 			temp.Execute(w, data)
// 		}
// 	}
// }
