// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import "code.gitea.io/gitea/modules/structs"

// ExtractKeysFromMapString provides a slice of keys from map
func ExtractKeysFromMapString(in map[string]structs.VisibleType) (keys []string) {
	for k := range in {
		keys = append(keys, k)
	}
	return
}
