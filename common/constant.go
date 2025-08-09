package common

import (
	"bytes"
	"fmt"
	"mime/multipart"
)

const (
	SUCCESS   = "S"
	ERROR     = "E"
	COMPLETED = "C"
	FAILED    = "F"
)

func BytesToMultipartFile(pData []byte, pFilename string) (multipart.File, *multipart.FileHeader, error) {

	lBody := &bytes.Buffer{}
	lWriter := multipart.NewWriter(lBody)

	lPart, lErr := lWriter.CreateFormFile("file", pFilename)
	if lErr != nil {
		return nil, nil, lErr
	}

	_, lErr = lPart.Write(pData)
	if lErr != nil {
		return nil, nil, lErr
	}

	lWriter.Close()

	lReader := multipart.NewReader(bytes.NewReader(lBody.Bytes()), lWriter.Boundary())
	lForm, lErr := lReader.ReadForm(int64(len(lBody.Bytes())))
	if lErr != nil {
		return nil, nil, lErr
	}

	lFiles := lForm.File["file"]
	if len(lFiles) == 0 {
		return nil, nil, fmt.Errorf("no file found")
	}

	lFileHeader := lFiles[0]
	lFile, lErr := lFileHeader.Open()
	if lErr != nil {
		return nil, nil, lErr
	}

	return lFile, lFileHeader, nil
}
