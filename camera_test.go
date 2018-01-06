package camera

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitialize(t *testing.T) {

	// Test width < MIN_WIDTH
	if _, err := Initialize(MIN_WIDTH-1, MAX_HEIGHT, false, false, "./"); err == nil {
		t.Error("width < MIN_WIDTH: expected error but was accepted")
	}

	// Test width > MAX_WIDTH
	if _, err := Initialize(MAX_WIDTH+1, MAX_HEIGHT, false, false, "./"); err == nil {
		t.Error("width > MAX_WIDTH: expected error but was accepted")
	}

	// Test height < MIN_HEIGHT
	if _, err := Initialize(MIN_WIDTH, MIN_HEIGHT-1, false, false, "./"); err == nil {
		t.Error("height < MIN_HEIGHT: expected error but was accepted")
	}

	// Test height > MAX_HEIGHT
	if _, err := Initialize(MIN_WIDTH, MAX_HEIGHT+1, false, false, "./"); err == nil {
		t.Error("height > MAX_HEIGHT: expected error but was accepted")
	}

	// Test correct values
	if _, err := Initialize(MIN_WIDTH, MAX_HEIGHT, false, false, "./"); err != nil {
		t.Error(err)
	}
}

func TestPicture(t *testing.T) {
	c, err := Initialize(MAX_WIDTH, MAX_HEIGHT, true, true, "./")
	if err != nil {
		t.Error(err)
	}

	// Test taking picture
	if _, err = c.Picture(); err != nil && err.Error() != "exec: \"raspistill\": executable file not found in $PATH" {
		t.Error(err)
	}
}

func TestVideo(t *testing.T) {
	c, err := Initialize(MAX_WIDTH, MAX_HEIGHT, true, true, "./")
	if err != nil {
		t.Error(err)
	}

	// Test duration < MIN_DURATION
	if _, err = c.Video(0, false); err == nil {
		t.Error("Duration < MIN_DURATION: expected error but was accepted")
	}

	// Test correct values, no convert
	if _, err = c.Video(MIN_DURATION+10, false); err != nil && err.Error() != "exec: \"raspivid\": executable file not found in $PATH" {
		t.Error(err)
	}

	// Test correct values, convert
	if _, err = c.Video(MIN_DURATION+10, true); err != nil && err.Error() != "exec: \"raspivid\": executable file not found in $PATH" {
		t.Error(err)
	}
}

func TestBuildArgs(t *testing.T) {
	c, err := Initialize(MAX_WIDTH, MAX_HEIGHT, true, true, "./")
	if err != nil {
		t.Error(err)
	}

	// Test for Picture
	args := buildArgs(c, "foo.jpg", false, 0)
	ref := fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v",
		STILL, WIDTHF, c.width, HEIGHTF, c.height, NO_PREVIEWF, OUTF, filepath.Join(c.savePath, "foo.jpg"),
		HFLIPF, VFLIPF)
	res := STILL + " " + strings.Join(args, " ")
	if res != ref {
		t.Errorf("Expected: %v, got: %v", ref, res)
	}

	// Test for Video
	args = buildArgs(c, "foo.jpg", true, 50)
	ref = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v %v",
		VID, WIDTHF, c.width, HEIGHTF, c.height, NO_PREVIEWF, OUTF, filepath.Join(c.savePath, "foo.jpg"),
		HFLIPF, VFLIPF, TIMEF, 50*1000)
	res = VID + " " + strings.Join(args, " ")
	if res != ref {
		t.Errorf("Expected: %v, got: %v", ref, res)
	}
}
