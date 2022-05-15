package removebg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

var APIENDPOINT = "https://api.remove.bg/v1.0/removebg"

type RemoveBG struct {
	APIKey string
}

func NewRemoveBG(apiKey string) *RemoveBG {
	return &RemoveBG{
		APIKey: apiKey,
	}
}

type RemoveOption struct {
	// Size of the output image (`'auto'` = highest available resolution, `'preview'`|`'small'`|`'regular'` = 0.25 MP, `'medium'` = 1.5 MP, `'hd'` = 4 MP, `'full'`|`'4k'` = original size)
	Size string
	// Type foreground object (`'auto'` = autodetect, `'person'`, `'product'`, `'car'`)
	Type string
	// TypeLevel classification level of the foreground object (`'none'` = no classification, `'1'` = coarse classification (e.g. `'car'`), `'2'` = specific classification (e.g. `'car_interior'`), `'latest'` = latest classification)
	TypeLevel string
	// Format image format (`'auto'` = autodetect, `'png'`, `'jpg'`, `'zip'`)
	Format string
	// ROI region of interest, where to look for foreground object (x1, y1, x2, y2) in px or relative (%)
	ROI string
	// Crop px or relative, single val = all sides, two vals = top/bottom, left/right, four vals = top, right, bottom, left
	Crop string
	// Scale image scale relative to the total image size
	Scale string
	// Position `'center'`, `'original'`, single val = horizontal and vertical, two vals = horizontal, vertical
	Position string
	// Channels request the finalized image (`'rgba'`) or an alpha mask (`'alpha'`)
	Channels string
	// Shadow whether to add an artificial shadow (some types aren't supported)
	Shadow bool
	// Semitransparency  for windows or glass objects (some types aren't supported)
	Semitransparency bool
	// Background (`None` = no background, path, url, color hex code (e.g. `'81d4fa'`, `'fff'`), color name (e.g. `'green'`))
	Background string
	// Background type (`None` = no background, `'path'`, `'url'`, `'color'`)
	BackgroundType string
	NewFileName    string
}

func NewRemoveOption() *RemoveOption {
	rb := &RemoveOption{
		Size:             "regular",
		Type:             "auto",
		TypeLevel:        "none",
		Format:           "auto",
		ROI:              "0 0 100% 100%",
		Position:         "original",
		Channels:         "rgba",
		Shadow:           false,
		Semitransparency: true,
		NewFileName:      "no-bg.png",
	}
	return rb
}

func isStringInSlice(haystack []string, s string) bool {
	for _, e := range haystack {
		if e == s {
			return true
		}
	}
	return false
}

// check if opts are valid.
func (opt *RemoveOption) check() error {
	sizes := []string{"auto", "preview", "small", "regular", "medium", "hd", "full", "4k"}
	if opt.Size == "" {
		opt.Size = "regular"
	}
	if !isStringInSlice(sizes, opt.Size) {
		return errors.New("size argument wrong")
	}
	typies := []string{"auto", "person", "product", "animal", "car", "car_interior", "car_part", "transportation", "graphics", "other"}
	if opt.Type == "" {
		opt.Type = "auto"
	}
	if !isStringInSlice(typies, opt.Type) {
		return errors.New("type argument wrong")
	}
	typeLevies := []string{"none", "latest", "1", "2"}
	if opt.TypeLevel == "" {
		opt.TypeLevel = "none"
	}
	if !isStringInSlice(typeLevies, opt.TypeLevel) {
		return errors.New("type-level argument wrong")
	}

	formates := []string{"jpg", "zip", "png", "auto"}
	if opt.Format == "" {
		opt.Format = "auto"
	}
	if !isStringInSlice(formates, opt.Format) {
		return errors.New("format argument wrong")
	}
	channels := []string{"rgba", "alpha"}
	if opt.Channels == "" {
		opt.Channels = "rgba"
	}
	if !isStringInSlice(channels, opt.Channels) {
		return errors.New("channels argument wrong")
	}
	if opt.ROI == "" {
		opt.ROI = `0 0 100% 100%`
	}
	if opt.Position == "" {
		opt.Position = "original"
	}
	if opt.NewFileName == "" {
		opt.NewFileName = "no-bg.png"
	}
	return nil
}

func (rb *RemoveBG) saveFile(resp *http.Response, fileName string) error {
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return ioutil.WriteFile(fileName, body, os.ModePerm)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	r := map[string][]map[string]string{}
	err := json.Unmarshal(body, &r)
	if err != nil {
		return err
	}
	return errors.New(r["errors"][0]["title"])

}

func (rb *RemoveBG) RemoveFromFile(filePath string, opt *RemoveOption) error {
	if err := opt.check(); err != nil {
		log.Println(err)
		return err
	}
	params := map[string]string{
		"size":       opt.Size,
		"type":       opt.Type,
		"type_level": opt.TypeLevel,
		"format":     opt.Format,
		"roi":        opt.ROI,
		"scale":      opt.Scale,
		"position":   opt.Position,
		"channels":   opt.Channels,
	}
	if opt.Crop != "" {
		params["crop"] = "true"
		params["crop_margin"] = opt.Crop
	} else {
		params["crop"] = "false"
	}

	if opt.Shadow {
		params["add_shadow"] = "true"
	} else {
		params["add_shodow"] = "false"
	}

	if opt.Semitransparency {
		params["semitransparency"] = "true"
	} else {
		params["semitransparency"] = "false"
	}

	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)
	for k, v := range params {
		multiPartWriter.WriteField(k, v)
	}
	if opt.BackgroundType == "path" {
		bgFileWriter, err := multiPartWriter.CreateFormFile("bg_image_file", path.Base(opt.Background))
		if err != nil {
			return err
		}
		bgFile, err := os.Open(opt.Background)
		if err != nil {
			return err
		}
		defer bgFile.Close()
		_, err = io.Copy(bgFileWriter, bgFile)
		if err != nil {
			return err
		}
	} else if opt.BackgroundType == "color" {
		multiPartWriter.WriteField("bg_color", opt.Background)
	} else if opt.BackgroundType == "url" {
		multiPartWriter.WriteField("bg_color", opt.Background)

	}
	imgFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer imgFile.Close()
	imgFileWriter, err := multiPartWriter.CreateFormFile("image_file", path.Base(filePath))
	if err != nil {
		return err
	}
	_, err = io.Copy(imgFileWriter, imgFile)
	if err != nil {
		return err
	}
	multiPartWriter.Close()
	request, err := http.NewRequest("POST", APIENDPOINT, &requestBody)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", multiPartWriter.FormDataContentType())
	request.Header.Add("X-Api-Key", rb.APIKey)
	// request.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	return rb.saveFile(resp, opt.NewFileName)

}

func (rb *RemoveBG) RemoveFromURL(url string, opt *RemoveOption) error {
	if opt == nil {
		opt = NewRemoveOption()
	}
	if err := opt.check(); err != nil {
		log.Println(err)
		return err
	}
	params := map[string]string{
		"image_url":  url,
		"size":       opt.Size,
		"type":       opt.Type,
		"type_level": opt.TypeLevel,
		"format":     opt.Format,
		"roi":        opt.ROI,
		"scale":      opt.Scale,
		"position":   opt.Position,
		"channels":   opt.Channels,
	}
	if opt.Crop != "" {
		params["crop"] = "true"
		params["crop_margin"] = opt.Crop
	} else {
		params["crop"] = "false"
	}
	if opt.Shadow {
		params["add_shadow"] = "true"
	} else {
		params["add_shodow"] = "false"
	}

	if opt.Semitransparency {
		params["semitransparency"] = "true"
	} else {
		params["semitransparency"] = "false"
	}
	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)
	for k, v := range params {
		multiPartWriter.WriteField(k, v)
	}
	if opt.BackgroundType == "path" {
		bgFileWriter, err := multiPartWriter.CreateFormFile("bg_image_file", path.Base(opt.Background))
		if err != nil {
			return err
		}
		bgFile, err := os.Open(opt.Background)
		if err != nil {
			return err
		}
		defer bgFile.Close()
		_, err = io.Copy(bgFileWriter, bgFile)
		if err != nil {
			return err
		}
	} else if opt.BackgroundType == "color" {
		multiPartWriter.WriteField("bg_color", opt.Background)
	} else if opt.BackgroundType == "url" {
		multiPartWriter.WriteField("bg_color", opt.Background)

	}

	multiPartWriter.Close()
	request, err := http.NewRequest("POST", APIENDPOINT, &requestBody)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", multiPartWriter.FormDataContentType())
	request.Header.Add("X-Api-Key", rb.APIKey)
	// request.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	return rb.saveFile(resp, opt.NewFileName)
}

func (rb *RemoveBG) RemoveFromBase64(base64 string, opt *RemoveOption) error {
	if opt == nil {
		opt = NewRemoveOption()
	}
	if err := opt.check(); err != nil {
		log.Println(err)
		return err
	}
	params := map[string]string{
		"image_file_b64": base64,
		"size":           opt.Size,
		"type":           opt.Type,
		"type_level":     opt.TypeLevel,
		"format":         opt.Format,
		"roi":            opt.ROI,
		"scale":          opt.Scale,
		"position":       opt.Position,
		"channels":       opt.Channels,
	}
	if opt.Crop != "" {
		params["crop"] = "true"
		params["crop_margin"] = opt.Crop
	} else {
		params["crop"] = "false"
	}
	if opt.Shadow {
		params["add_shadow"] = "true"
	} else {
		params["add_shodow"] = "false"
	}

	if opt.Semitransparency {
		params["semitransparency"] = "true"
	} else {
		params["semitransparency"] = "false"
	}

	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)
	for k, v := range params {
		multiPartWriter.WriteField(k, v)
	}
	if opt.BackgroundType == "path" {
		bgFileWriter, err := multiPartWriter.CreateFormFile("bg_image_file", path.Base(opt.Background))
		if err != nil {
			return err
		}
		bgFile, err := os.Open(opt.Background)
		if err != nil {
			return err
		}
		defer bgFile.Close()
		_, err = io.Copy(bgFileWriter, bgFile)
		if err != nil {
			return err
		}
	} else if opt.BackgroundType == "color" {
		multiPartWriter.WriteField("bg_color", opt.Background)
	} else if opt.BackgroundType == "url" {
		multiPartWriter.WriteField("bg_color", opt.Background)

	}

	multiPartWriter.Close()
	request, err := http.NewRequest("POST", APIENDPOINT, &requestBody)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", multiPartWriter.FormDataContentType())
	request.Header.Add("X-Api-Key", rb.APIKey)
	// request.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	return rb.saveFile(resp, opt.NewFileName)
}
