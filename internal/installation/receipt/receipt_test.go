// Copyright 2019 The Kubernetes Authors.
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

package receipt

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"sigs.k8s.io/krew/internal/index/indexscanner"
	"sigs.k8s.io/krew/internal/testutil"
)

func TestStore(t *testing.T) {
	tmpDir, cleanup := testutil.NewTempDir(t)
	defer cleanup()

	testPlugin := testutil.NewPlugin().WithName("some-plugin").WithPlatforms(testutil.NewPlatform().V()).V()
	dest := tmpDir.Path("some-plugin.yaml")

	if err := Store(testPlugin, dest); err != nil {
		t.Fatal(err)
	}

	actual, err := indexscanner.LoadPluginByName(tmpDir.Root(), "some-plugin")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(&testPlugin, &actual); diff != "" {
		t.Fatal(diff)
	}
}

func TestLoad(t *testing.T) {
	tmpDir, cleanup := testutil.NewTempDir(t)
	defer cleanup()

	testPlugin := testutil.NewPlugin().WithName("foo").WithPlatforms(testutil.NewPlatform().V()).V()
	if err := Store(testPlugin, tmpDir.Path("foo.yaml")); err != nil {
		t.Fatal(err)
	}

	gotPlugin, err := Load(tmpDir.Path("foo.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(&gotPlugin, &testPlugin); diff != "" {
		t.Fatal(diff)
	}
}

func TestLoad_preservesNonExistsError(t *testing.T) {
	_, err := Load("non-existing.yaml")
	if !os.IsNotExist(err) {
		t.Fatalf("returned error is not ENOENT: %+v", err)
	}
}
