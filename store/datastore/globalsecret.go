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
	"github.com/russross/meddler"
	"github.com/woodpecker-ci/woodpecker/model"
	"github.com/woodpecker-ci/woodpecker/store/datastore/sql"
)

func (db *datastore) GlobalSecretFind(name string) (*model.GlobalSecret, error) {
	stmt := sql.Lookup(db.driver, "global-secret-find-name")
	data := new(model.GlobalSecret)
	err := meddler.QueryRow(db, data, stmt, name)
	return data, err
}

func (db *datastore) GlobalSecretList() ([]*model.GlobalSecret, error) {
	stmt := sql.Lookup(db.driver, "global-secret-find")
	data := []*model.GlobalSecret{}
	err := meddler.QueryAll(db, &data, stmt)
	return data, err
}

func (db *datastore) GlobalSecretCreate(secret *model.GlobalSecret) error {
	return meddler.Insert(db, "global_secrets", secret)
}

func (db *datastore) GlobalSecretUpdate(secret *model.GlobalSecret) error {
	return meddler.Update(db, "global_secrets", secret)
}

func (db *datastore) GlobalSecretDelete(secret *model.GlobalSecret) error {
	stmt := sql.Lookup(db.driver, "global-secret-delete")
	_, err := db.Exec(stmt, secret.ID)
	return err
}
