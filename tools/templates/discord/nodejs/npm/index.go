package npm

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/abdfnx/botway/constants"
	"github.com/abdfnx/botway/tools/templates"
	"github.com/abdfnx/botway/tools/templates/discord/nodejs"
	"github.com/abdfnx/looker"
	"github.com/tidwall/sjson"
)

func DiscordNodejsNpm(botName string) {
	npm, err := looker.LookPath("npm")

	if err != nil {
		fmt.Print(constants.FAIL_BACKGROUND.Render("ERROR"))
		fmt.Println(constants.FAIL_FOREGROUND.Render(" npm is not installed"))
	} else {
		npmInit := npm + " init -y"

		cmd := exec.Command("bash", "-c", npmInit)

		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell.exe", npmInit)
		}

		cmd.Dir = botName
		err = cmd.Run()

		if err != nil {
			log.Printf("error: %v\n", err)
		}

		packageJson, err := ioutil.ReadFile(filepath.Join(botName, "package.json"))

		if err != nil {
			log.Printf("error: %v\n", err)
		}

		version, _ := sjson.Set(string(packageJson), "version", "0.1.0")
		description, _ := sjson.Delete(version, "description")
		keywords, _ := sjson.Delete(description, "keywords")
		license, _ := sjson.Delete(keywords, "license")
		main, _ := sjson.Set(string(license), "main", "src/index.js")
		author, _ := sjson.Delete(string(main), "author")
		final, _ := sjson.Delete(author, "scripts")

		newPackageJson := ioutil.WriteFile(filepath.Join(botName, "package.json"), []byte(final), 0644)

		if newPackageJson != nil {
			log.Printf("error: %v\n", newPackageJson)
		}

		indexFile := os.WriteFile(filepath.Join(botName, "src", "index.js"), []byte(nodejs.IndexJSContent()), 0644)
		dockerFile := os.WriteFile(filepath.Join(botName, "Dockerfile"), []byte(DockerfileContent(botName)), 0644)
		procFile := os.WriteFile(filepath.Join(botName, "Procfile"), []byte("process: node ./src/index.js"), 0644)
		resourcesFile := os.WriteFile(filepath.Join(botName, "resources.md"), []byte(nodejs.Resources()), 0644)

		if resourcesFile != nil {
			log.Fatal(resourcesFile)
		}

		if indexFile != nil {
			log.Fatal(indexFile)
		}

		if dockerFile != nil {
			log.Fatal(dockerFile)
		}

		if procFile != nil {
			log.Fatal(procFile)
		}

		npmInstall := npm + " install " + nodejs.Packages

		installCmd := exec.Command("bash", "-c", npmInstall)

		if runtime.GOOS == "windows" {
			installCmd = exec.Command("powershell.exe", npmInstall)
		}

		installCmd.Dir = botName
		installCmd.Stdin = os.Stdin
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		err = installCmd.Run()

		if err != nil {
			log.Printf("error: %v\n", err)
		}

		if runtime.GOOS == "windows" {
			installWindowsBuildTools := exec.Command("powershell.exe", npm + " install --global --production --add-python-to-path windows-build-tools")

			installWindowsBuildTools.Dir = botName
			installWindowsBuildTools.Stdin = os.Stdin
			installWindowsBuildTools.Stdout = os.Stdout
			installWindowsBuildTools.Stderr = os.Stderr
			err = installWindowsBuildTools.Run()

			if err != nil {
				log.Printf("error: %v\n", err)
			}
		}

		templates.CheckProject(botName, "discord")
	}
}
