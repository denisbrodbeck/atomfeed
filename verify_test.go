package atomfeed

import (
	"testing"
)

func Test_checkID(t *testing.T) {
	if err := checkID(nil); err == nil {
		t.Error("expected an error on nil ID, got none")
	}
	if err := checkID(&ID{}); err == nil {
		t.Error("expected an error on empty ID, got none")
	}
	if err := checkID(&ID{Value: "example.com"}); err != nil {
		t.Error(err)
	}
}

func Test_checkURI(t *testing.T) {
	if err := checkURI(""); err != nil {
		t.Error(err)
	}
	if err := checkURI("example.com"); err != nil {
		t.Error(err)
	}
	if err := checkURI(":example.com"); err == nil {
		t.Error("expected missing protocol scheme error on invalid uri, got none")
	}
}

func Test_checkEmail(t *testing.T) {
	if err := checkEmail(""); err != nil {
		t.Error(err)
	}
	if err := checkEmail("mail@example.com"); err != nil {
		t.Error(err)
	}
	if err := checkEmail("mailatexample.com"); err == nil {
		t.Error("expected missing @ error on invalid email, got none")
	}
}

func Test_checkPerson(t *testing.T) {
	if err := checkPerson(nil); err != nil {
		t.Error(err)
	}
	if err := checkPerson(&Person{}); err == nil {
		t.Error("should fail on missing name, did not")
	}
	if err := checkPerson(&Person{Name: "Go"}); err != nil {
		t.Error(err)
	}
	if err := checkPerson(&Person{Name: "Go", Email: "wrong.com"}); err == nil {
		t.Error("should fail on invalid email, did not")
	}
}

func Test_checkAuthorsExist(t *testing.T) {
	author := NewPerson("Go", "", "")
	type args struct {
		f *Feed
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "no authors at all",
			args:    args{f: &Feed{}},
			wantErr: true,
		},
		{
			name:    "not enough authors",
			args:    args{f: &Feed{Entries: []Entry{Entry{Author: author}, Entry{Author: nil}}}},
			wantErr: true,
		},
		{
			name:    "only feed author",
			args:    args{f: &Feed{Author: author, Entries: []Entry{Entry{}, Entry{}}}},
			wantErr: false,
		},
		{
			name:    "only entry authors",
			args:    args{f: &Feed{Entries: []Entry{Entry{Author: author}, Entry{Author: author}}}},
			wantErr: false,
		},

		{
			name:    "mixed",
			args:    args{f: &Feed{Author: author, Entries: []Entry{Entry{Author: nil}, Entry{Author: author}}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkAuthorsExist(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("checkAuthorsExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkContent(t *testing.T) {
	type args struct {
		c *Content
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "invalid source",
			args:    args{&Content{Source: ":broken.source"}},
			wantErr: true,
		},
		{
			name:    "source and content",
			args:    args{&Content{Source: "example.org", Value: "not empty"}},
			wantErr: true,
		},
		{
			name:    "invalid source type",
			args:    args{&Content{Source: "example.org", Value: "", Type: "html"}},
			wantErr: true,
		},
		{
			name:    "valid source",
			args:    args{&Content{Source: "example.org", Value: "", Type: "image/png"}},
			wantErr: false,
		},
		{
			name:    "valid content",
			args:    args{&Content{Value: "stuff", Type: "text"}},
			wantErr: false,
		},
		{
			name:    "invalid content mime type",
			args:    args{&Content{Value: "", Type: "gif"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkContent(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("checkContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}