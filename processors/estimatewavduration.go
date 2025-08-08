package processors

import (
	"encoding/binary"
	"errors"
	"io"
	"mime/multipart"
)

func GetWAVDuration(file multipart.File) (float64, error) {

	header := make([]byte, 44)
	_, err := io.ReadFull(file, header)
	if err != nil {
		return 0, err
	}

	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return 0, errors.New("not a valid WAV file")
	}

	// Extract sample rate (bytes 24–27), bits per sample (34–35), channels (22–23), data size (40–43)
	sampleRate := binary.LittleEndian.Uint32(header[24:28])
	bitsPerSample := binary.LittleEndian.Uint16(header[34:36])
	numChannels := binary.LittleEndian.Uint16(header[22:24])
	dataSize := binary.LittleEndian.Uint32(header[40:44])

	if sampleRate == 0 || bitsPerSample == 0 || numChannels == 0 {
		return 0, errors.New("invalid WAV metadata")
	}

	// Calculate duration: size / (sampleRate * numChannels * bitsPerSample / 8)
	bytesPerSecond := float64(sampleRate * uint32(numChannels) * uint32(bitsPerSample) / 8)
	durationSeconds := float64(dataSize) / bytesPerSecond

	return durationSeconds, nil
}

func EstimateMP3Duration(file multipart.File) (float64, error) {
	// Get file size
	fileSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	file.Seek(0, 0)

	// Estimate bitrate — assume 128 kbps (standard)
	const assumedBitrateKbps = 128.0
	bits := float64(fileSize * 8) // bytes to bits
	durationSeconds := bits / (assumedBitrateKbps * 1000)

	return durationSeconds, nil
}
