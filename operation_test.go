package backup

import (
	"reflect"
	"testing"

	"strings"
)

func TestOperationFromBytesGood(t *testing.T) {
	s := []byte(`{
  "destination_root": "/test",
  "sd_card_mount_point": "/Volumes",
  "sd_card_names": [
    "HERA",
    "ZEUS"
  ],
  "folder_mapping": [
    { "source": "DCIM",                      "destination": "DCIM"    },
    { "source": "MP_ROOT",                   "destination": "MP_ROOT" }
  ]
}`)
	expected := &Operation{
		DestinationRoot:  "/test",
		SDCardMountPoint: "/Volumes",
		SDCardNames:      []string{"HERA", "ZEUS"},
		FolderMapping: []folderMapping{
			{Source: "DCIM", Destination: "DCIM"},
			{Source: "MP_ROOT", Destination: "MP_ROOT"},
		},
	}

	op, err := operationFromBytes(s)
	if err != nil {
		t.Fatalf("Unable to read valid config: %s", err)
	}

	if !reflect.DeepEqual(expected, op) {
		t.Errorf("Parsed config string is invalid.\n%#v\n%#v", expected, op)
	}
}

func TestOperationFromBytesBad(t *testing.T) {
	s := []byte("{\"destination_root\"}")

	_, err := operationFromBytes(s)
	if err == nil {
		t.Error("Expected deserialization to fail for invalid config string.")
	}
}

var validationErrors = []struct {
	source     string
	wantPrefix string
}{
	{`{
  "sd_card_mount_point": "/Volumes",
  "sd_card_names": ["HERA", "ZEUS"],
  "folder_mapping": [{"source": "from", "destination": "to"}]
}`,
		"missing `destination_root`"},
	{`{
  "destination_root": "/test",
  "sd_card_names": ["HERA", "ZEUS"],
  "folder_mapping": [{"source": "from", "destination": "to"}]
}`,
		"missing `sd_card_mount_point`"},
	{`{
  "destination_root": "/test",
  "sd_card_mount_point": "/Volumes",
  "folder_mapping": [{"source": "from", "destination": "to"}]
}`,
		"missing `sd_card_names`"},
	{`{
  "destination_root": "/test",
  "sd_card_mount_point": "/Volumes",
  "sd_card_names": [],
  "folder_mapping": [{"source": "from", "destination": "to"}]
}`,
		"empty `sd_card_names`"},
	{`{
  "destination_root": "/test",
  "sd_card_mount_point": "/Volumes",
  "sd_card_names": ["HERA", ""],
  "folder_mapping": [{"source": "from", "destination": "to"}]
}`,
		"contains empty card name"},
	{`{
  "destination_root": "/test",
  "sd_card_mount_point": "/Volumes",
  "sd_card_names": ["HERA", "ZEUS"]
}`,
		"missing `folder_mapping`"},
	{`{
  "destination_root": "/test",
  "sd_card_mount_point": "/Volumes",
  "sd_card_names": ["HERA", "ZEUS"],
  "folder_mapping": []
}`,
		"empty `folder_mapping`"},
	{`{
  "destination_root": "/test",
  "sd_card_mount_point": "/Volumes",
  "sd_card_names": ["HERA", "ZEUS"],
  "folder_mapping": [{"destination": "to"}]
}`,
		"missing `source` in folder mapping"},
	{`{
  "destination_root": "/test",
  "sd_card_mount_point": "/Volumes",
  "sd_card_names": ["HERA", "ZEUS"],
  "folder_mapping": [{"source": "from"}]
}`,
		"missing `destination` in folder mapping"},
}

func TestValidationErrors(t *testing.T) {
	for _, c := range validationErrors {
		_, err := operationFromBytes([]byte(c.source))

		if err == nil {
			t.Errorf("Expected validation error: %#v", c.wantPrefix)
			continue
		}

		if !strings.HasPrefix(err.Error(), c.wantPrefix) {
			t.Errorf("Incorrect prefix.\nPrefix Wanted: %#v\nError Observed: %#v", c.wantPrefix, err.Error())
		}
	}
}
