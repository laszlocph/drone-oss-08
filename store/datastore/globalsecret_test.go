// Copyright 2018 Drone.IO Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datastore

import (
	"testing"

	"github.com/woodpecker-ci/woodpecker/model"
)

func TestGlobalSecretFind(t *testing.T) {
	s := newTest()
	defer func() {
		s.Exec("delete from global_secrets")
		s.Close()
	}()

	err := s.GlobalSecretCreate(&model.GlobalSecret{
		Name:   "password",
		Value:  "correct-horse-battery-staple",
		Images: []string{"golang", "node"},
		Events: []string{"push", "tag"},
	})
	if err != nil {
		t.Errorf("Unexpected error: insert global secret: %s", err)
		return
	}

	secret, err := s.GlobalSecretFind("password")
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := secret.Name, "password"; got != want {
		t.Errorf("Want global secret name %s, got %s", want, got)
	}
	if got, want := secret.Value, "correct-horse-battery-staple"; got != want {
		t.Errorf("Want global secret value %s, got %s", want, got)
	}
	if got, want := secret.Events[0], "push"; got != want {
		t.Errorf("Want global secret event %s, got %s", want, got)
	}
	if got, want := secret.Events[1], "tag"; got != want {
		t.Errorf("Want global secret event %s, got %s", want, got)
	}
	if got, want := secret.Images[0], "golang"; got != want {
		t.Errorf("Want global secret image %s, got %s", want, got)
	}
	if got, want := secret.Images[1], "node"; got != want {
		t.Errorf("Want global secret image %s, got %s", want, got)
	}
}

func TestGlobalSecretList(t *testing.T) {
	s := newTest()
	defer func() {
		s.Exec("delete from global_secrets")
		s.Close()
	}()

	s.GlobalSecretCreate(&model.GlobalSecret{
		Name:  "foo",
		Value: "bar",
	})
	s.GlobalSecretCreate(&model.GlobalSecret{
		Name:  "baz",
		Value: "qux",
	})

	list, err := s.GlobalSecretList()
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(list), 2; got != want {
		t.Errorf("Want %d global secrets, got %d", want, got)
	}
}

func TestGlobalSecretUpdate(t *testing.T) {
	s := newTest()
	defer func() {
		s.Exec("delete from global_secrets")
		s.Close()
	}()

	secret := &model.GlobalSecret{
		Name:  "foo",
		Value: "baz",
	}
	if err := s.GlobalSecretCreate(secret); err != nil {
		t.Errorf("Unexpected error: insert global secret: %s", err)
		return
	}
	secret.Value = "qux"
	if err := s.GlobalSecretUpdate(secret); err != nil {
		t.Errorf("Unexpected error: update global secret: %s", err)
		return
	}
	updated, err := s.GlobalSecretFind("foo")
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := updated.Value, "qux"; got != want {
		t.Errorf("Want global global secret value %s, got %s", want, got)
	}
}

func TestGlobalSecretDelete(t *testing.T) {
	s := newTest()
	defer func() {
		s.Exec("delete from global_secrets")
		s.Close()
	}()

	secret := &model.GlobalSecret{
		Name:  "foo",
		Value: "baz",
	}
	if err := s.GlobalSecretCreate(secret); err != nil {
		t.Errorf("Unexpected error: insert global secret: %s", err)
		return
	}

	if err := s.GlobalSecretDelete(secret); err != nil {
		t.Errorf("Unexpected error: delete global secret: %s", err)
		return
	}
	_, err := s.GlobalSecretFind("foo")
	if err == nil {
		t.Errorf("Expect error: sql.ErrNoRows")
		return
	}
}

func TestGlobalSecretIndexes(t *testing.T) {
	s := newTest()
	defer func() {
		s.Exec("delete from global_secrets")
		s.Close()
	}()

	if err := s.GlobalSecretCreate(&model.GlobalSecret{
		Name:  "foo",
		Value: "bar",
	}); err != nil {
		t.Errorf("Unexpected error: insert secret: %s", err)
		return
	}

	// fail due to duplicate name
	if err := s.GlobalSecretCreate(&model.GlobalSecret{
		Name:  "foo",
		Value: "baz",
	}); err == nil {
		t.Errorf("Unexpected error: dupliate name")
	}
}
