// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
// limitations under the License.

package main

import (
	_ "steve/servicelauncher/cmd"
	"steve/servicelauncher/launcher"
)

// 为了兼容有些包里有init方法直接在获取exposer，把cmd的init也放到更早的init里
// 目前的真正入口方法是cmd/root.go的init方法
func main() {
	launcher.LoadService()
}
