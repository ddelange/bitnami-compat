/*
Copyright © 2022 ZCube <zcubekr@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list bitnami dockerfiles",
	Long:  `list bitnami dockerfiles`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var dockerfiles []string
		if len(app) > 0 {
			dockerfiles, err = doublestar.FilepathGlob(fmt.Sprintf("bitnami-dockers/bitnami-docker-%v/**/Dockerfile", app))
		} else {
			dockerfiles, err = doublestar.FilepathGlob(fmt.Sprintf("bitnami-dockers/bitnami-docker-*/**/Dockerfile"))
		}
		if err != nil {
			log.Panic(err)
		}

		for i := range dockerfiles {
			if appInfo, err := InspectDockerfile(dockerfiles[i]); err == nil {
				// fmt.Println(appInfo.Dependencies)

				patchFound := false
				var err error
				var patchs []PatchInfo
				if patchs, err = FindPatchs(appInfo); err == nil {
					for _, patch := range patchs {
						patchFound = (patch.BashPatch != "" ||
							patch.DockerFromPatch != "" ||
							patch.DockerInstallPatch != "" ||
							patch.GolangBuild != "")
						if patchFound == false {
							break
						}
					}
					if len(patchs) == 0 {
						patchFound = len(patchs) == 0
					}
				} else {
					patchFound = len(patchs) == 0
				}

				if patchFound {
					fmt.Println(fmt.Sprintf("(o) %v:%v", appInfo.Name, appInfo.Version.Original()))
				} else {
					fmt.Println(fmt.Sprintf("(x) %v:%v", appInfo.Name, appInfo.Version.Original()))
					for _, patch := range patchs {
						patchFound = (patch.BashPatch != "" ||
							patch.DockerFromPatch != "" ||
							patch.DockerInstallPatch != "" ||
							patch.GolangBuild != "")
						if !patchFound {
							fmt.Println(fmt.Sprintf(" * %v.%v patch needed", patch.PackageInfo.Name, patch.PackageInfo.Version.Original()))
						}
					}
					if len(patchs) == 0 {
						patchFound = len(patchs) == 0
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringVar(&app, "app", "", "app")
}
