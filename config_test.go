package backup

import (
	"reflect"
	"testing"
)

func TestOperationFromBytesGood(t *testing.T) {
	s := []byte("{  \"destination_root\": \"/test\",  \"sd_card_mount_point\": \"/Volumes\",  \"sd_card_names\": [    \"HERA\",    \"ZEUS\"  ],  \"folder_mapping\": [    { \"source\": \"DCIM\",                      \"destination\": \"DCIM\"    },    { \"source\": \"MP_ROOT\",                   \"destination\": \"MP_ROOT\" },    { \"source\": \"PRIVATE/AVCHD/BDMV/STREAM\", \"destination\": \"STREAM\"  },    { \"source\": \"PRIVATE/M4ROOT/CLIP\",       \"destination\": \"CLIP\"    }  ]}")
	expected := &Operation{DestinationRoot: "/test", SDCardMountPoint: "/Volumes", SDCardNames: []string{"HERA", "ZEUS"}, FolderMapping: []folderMapping{folderMapping{Source: "DCIM", Destination: "DCIM"}, folderMapping{Source: "MP_ROOT", Destination: "MP_ROOT"}, folderMapping{Source: "PRIVATE/AVCHD/BDMV/STREAM", Destination: "STREAM"}, folderMapping{Source: "PRIVATE/M4ROOT/CLIP", Destination: "CLIP"}}}

	op, err := operationFromBytes(s)
	if err != nil {
		t.Errorf("Unable to read valid config: %s", err)
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

// TODO: enforce that the config data has all expected entries.
