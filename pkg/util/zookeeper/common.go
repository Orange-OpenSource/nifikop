// Copyright 2020 Orange SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package apis

package zookeeper

import "strings"

// PrepareConnectionAddress prepares the proper address for Nifi and CC
// The required path for Nifi and CC looks 'example-1:2181/nifi'
func PrepareConnectionAddress(zkAddresse string, zkPath string) string {
	return zkAddresse + zkPath
}

//
func GetHostnameAddress(zkAddresse string) string {
	return strings.Split(zkAddresse, ":")[0]
}

//
func GetPortAddress(zkAddresse string) string {
	return strings.Split(zkAddresse, ":")[1]
}
