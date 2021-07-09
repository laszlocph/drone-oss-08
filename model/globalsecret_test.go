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

package model

import (
	"testing"

	"github.com/franela/goblin"
)

func TestGlobalSecret(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("GlobalSecret", func() {

		g.It("should match event", func() {
			secret := GlobalSecret{}
			secret.Events = []string{"pull_request"}
			g.Assert(secret.Match("pull_request")).IsTrue()
		})
		g.It("should not match event", func() {
			secret := GlobalSecret{}
			secret.Events = []string{"pull_request"}
			g.Assert(secret.Match("push")).IsFalse()
		})
		g.It("should match when no event filters defined", func() {
			secret := GlobalSecret{}
			g.Assert(secret.Match("pull_request")).IsTrue()
		})
		g.It("should pass validation", func() {
			secret := GlobalSecret{}
			secret.Name = "secretname"
			secret.Value = "secretvalue"
			err := secret.Validate()
			g.Assert(err).Equal(nil)
		})
		g.Describe("should fail validation", func() {
			g.It("when no name", func() {
				secret := GlobalSecret{}
				secret.Value = "secretvalue"
				err := secret.Validate()
				g.Assert(err != nil).IsTrue()
			})
			g.It("when no value", func() {
				secret := GlobalSecret{}
				secret.Name = "secretname"
				err := secret.Validate()
				g.Assert(err != nil).IsTrue()
			})
		})
	})
}
