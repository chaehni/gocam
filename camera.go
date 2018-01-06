package camera

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// commands and flags
const (
	STILL       = "raspistill"
	VID         = "raspivid"
	WIDTHF      = "-w"
	HEIGHTF     = "-h"
	HFLIPF      = "-hf"
	VFLIPF      = "-vf"
	TIMEF       = "-t"
	NO_PREVIEWF = "-n"
	OUTF        = "-o"
)

// output characteristics
const (
	MIN_WIDTH    = 64
	MAX_WIDTH    = 1920
	MIN_HEIGHT   = 64
	MAX_HEIGHT   = 1080
	MIN_DURATION = 1
	TYPE_S       = ".jpg"
	RAW_TYPE_V   = ".h264"
	TYPE_V       = ".mp4"
	TIME_STAMP   = "02.01.06_15-04-05"
)

type camera struct {
	width    int
	height   int
	hflip    bool
	vflip    bool
	savePath string
}

// Initialize creates a new camera with the provided characteristics
func Initialize(width, height int, hflip, vflip bool, path string) (*camera, error) {
	if width < MIN_WIDTH || width > MAX_WIDTH || height < MIN_HEIGHT || height > MAX_HEIGHT {
		return nil, fmt.Errorf("width must be between %v and %v, height must be between %v and %v",
			MIN_WIDTH, MAX_WIDTH, MIN_HEIGHT, MAX_HEIGHT)
	}
	return &camera{width, height, hflip, vflip, path}, nil
}

// Picture takes a picture with the given Camera and returns the path to the picture taken
func (c *camera) Picture() (string, error) {
	filename := time.Now().Format(TIME_STAMP) + TYPE_S

	args := buildArgs(c, filename, false, 0)

	cmd := exec.Command(STILL, args...)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return filepath.Join(c.savePath, filename), nil
}

// Video captures a video with the given Camera and the specified duration and converts it to TYPE_V if desired
// it returns the path to the captured video
func (c *camera) Video(duration int, conv bool) (string, error) {

	if duration < MIN_DURATION {
		return "", fmt.Errorf("duration must be at least %v", MIN_DURATION)
	}

	filename := time.Now().Format(TIME_STAMP) + RAW_TYPE_V

	args := buildArgs(c, filename, true, duration)

	cmd := exec.Command(VID, args...)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	if conv {
		path, err := convert(filepath.Join(c.savePath, filename))
		if err != nil {
			return "", err
		}
		return path, nil
	}

	return filepath.Join(c.savePath, filename), nil
}

func convert(path string) (string, error) {
	cmd := exec.Command("MP4Box", "-add", path, strings.Replace(path, RAW_TYPE_V, TYPE_V, 1))
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	cmd = exec.Command("rm", path)
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.Replace(path, RAW_TYPE_V, TYPE_V, 1), nil
}

func buildArgs(c *camera, filename string, video bool, duration int) []string {
	args := []string{
		WIDTHF, strconv.Itoa(c.width),
		HEIGHTF, strconv.Itoa(c.height),
		NO_PREVIEWF,
		OUTF, filepath.Join(c.savePath, filename),
	}

	if c.hflip {
		args = append(args, HFLIPF)
	}
	if c.vflip {
		args = append(args, VFLIPF)
	}
	if video {
		args = append(args, TIMEF)
		args = append(args, strconv.Itoa(duration*1000))
	}

	return args
}
