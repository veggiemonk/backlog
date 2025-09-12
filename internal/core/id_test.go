package core

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"go.yaml.in/yaml/v4"
)

func TestTaskIDNameAndString(t *testing.T) {
	testCases := []struct {
		id       string
		expected string
	}{
		{"T1.2.3", "T01.02.03"},
		{"1.2.3", "T01.02.03"},
		{"T10.20.30", "T10.20.30"},
		{"10.20.30", "T10.20.30"},
		{"10.20.30", "T10.20.30"},
		{"T01", "T01"},
		{"1", "T01"},
	}
	for _, tc := range testCases {
		taskID, err := parseTaskID(tc.id)
		if err != nil {
			t.Errorf("Unexpected error parsing TaskID %q: %v", tc.id, err)
			continue
		}
		if name := taskID.Name(); name != tc.expected {
			t.Errorf("TaskID.Name() for %q = %q; want %q", tc.id, name, tc.expected)
		}
		if str := taskID.String(); str != tc.expected[1:] {
			t.Errorf("TaskID.String() for %q = %q; want %q", tc.id, str, tc.expected[1:])
		}
	}
}

func TestParseTaskID_Error(t *testing.T) {
	_, err := parseTaskID("T1.a.3")
	is.New(t).True(err != nil)
}

func TestTaskID_HasSubTasks(t *testing.T) {
	is := is.New(t)
	taskID, _ := parseTaskID("T1.2")
	is.True(taskID.HasSubTasks())

	taskID, _ = parseTaskID("T1")
	is.True(!taskID.HasSubTasks())
}

func TestTaskID_Less(t *testing.T) {
	is := is.New(t)
	taskID1, _ := parseTaskID("T1.2")
	taskID2, _ := parseTaskID("T1.3")
	is.True(taskID1.Less(taskID2))

	taskID1, _ = parseTaskID("T1.2")
	taskID2, _ = parseTaskID("T1.2.1")
	is.True(taskID1.Less(taskID2))
}

func TestTaskID_Equals(t *testing.T) {
	is := is.New(t)
	taskID1, _ := parseTaskID("T1.2")
	taskID2, _ := parseTaskID("T01.02")
	is.True(taskID1.Equals(taskID2))

	taskID1, _ = parseTaskID("T1.2")
	taskID2, _ = parseTaskID("T1.3")
	is.True(!taskID1.Equals(taskID2))
}

func TestTaskID_Parent(t *testing.T) {
	is := is.New(t)
	taskID, _ := parseTaskID("T1.2.3")
	parent := taskID.Parent()
	expected, _ := parseTaskID("T1.2")
	is.True(parent.Equals(expected))

	taskID, _ = parseTaskID("T1")
	parent = taskID.Parent()
	is.True(parent == nil)
}

func TestTaskID_NextSubTaskID(t *testing.T) {
	is := is.New(t)
	taskID, _ := parseTaskID("T1.2")
	nextSubTaskID := taskID.NextSubTaskID()
	expected, _ := parseTaskID("T1.2.1")
	is.True(nextSubTaskID.Equals(expected))
}

func TestTaskID_NextSiblingID(t *testing.T) {
	is := is.New(t)
	taskID, _ := parseTaskID("T1.2")
	nextSiblingID := taskID.NextSiblingID()
	expected, _ := parseTaskID("T1.3")
	is.True(nextSiblingID.Equals(expected))

	taskID, _ = parseTaskID("T1")
	nextSiblingID = taskID.NextSiblingID()
	expected, _ = parseTaskID("T2")
	is.True(nextSiblingID.Equals(expected))
}

func TestTaskID_JSONMarshalling(t *testing.T) {
	is := is.New(t)
	taskID, _ := parseTaskID("T1.2")
	b, err := json.Marshal(taskID)
	is.NoErr(err)
	is.Equal("\"01.02\"", string(b))

	var unmarshalled TaskID
	err = json.Unmarshal(b, &unmarshalled)
	is.NoErr(err)
	is.True(taskID.Equals(unmarshalled))
}

func TestTaskID_YAMLMarshalling(t *testing.T) {
	is := is.New(t)
	taskID, _ := parseTaskID("T1.2")
	b, err := yaml.Marshal(taskID)
	is.NoErr(err)
	is.Equal("\"01.02\"\n", string(b))

	var unmarshalled TaskID
	err = yaml.Unmarshal(b, &unmarshalled)
	is.NoErr(err)
	is.True(taskID.Equals(unmarshalled))
}

func TestTaskID_TextMarshalling(t *testing.T) {
	is := is.New(t)
	taskID, _ := parseTaskID("T1.2")
	b, err := taskID.MarshalText()
	is.NoErr(err)
	is.Equal("01.02", string(b))

	var unmarshalled TaskID
	err = unmarshalled.UnmarshalText(b)
	is.NoErr(err)
	is.True(taskID.Equals(unmarshalled))
}
