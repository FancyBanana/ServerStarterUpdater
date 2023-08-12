package main

import (
	"MinecraftServerManager/config"
	"MinecraftServerManager/cursefoge"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

const configPath = "./manager-config.yaml"
const serverConfigPath = "./server-setup-config.yaml"
const serverConfigBakPath = "./server-setup-config.yaml.BAK"

func getModpackLatestFileLink(config *config.Config) (string, error) {
	key := config.ApiKey

	var modId int
	var err error
	client := cursefoge.New(key)

	if config.ModpackSlug != "" {
		fmt.Println("Mod slug defined, searching mod by slug")
		modId, err = client.FindModIdBySlug(config.ModpackSlug)
		if err != nil {
			return "", errors.New("Could not get id for mod " + config.ModpackSlug + " , " + err.Error())
		}
	} else {
		modId = config.ModpackId
	}
	fmt.Println("Searching for latest file for mod " + strconv.Itoa(modId))
	lin, err := client.GetModFile(modId)
	if err != nil {
		return "", errors.New("Could not get mod file list " + err.Error())
	}

	return lin, nil
}

func copyFile(destination string, source string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	copied, err := io.Copy(dst, src)
	if err != nil || copied < 1 {
		return errors.New("Failed to create server-setup-config.yaml backup")
	}

	return nil
}

func UpdateServerSetupConfig(fileLink string) error {
	_, err := os.Stat(serverConfigPath)
	if os.IsNotExist(err) {
		return errors.New("Could not find server-setup-config.yaml\n" + err.Error())
	}

	fmt.Println("Doing backup of " + serverConfigPath)
	err = copyFile(serverConfigBakPath, serverConfigPath)
	if err != nil {
		return err
	}
	fmt.Println("Backup made at" + serverConfigBakPath)

	file, err := os.ReadFile(serverConfigPath)
	if err != nil {
		return err
	}

	var conf map[string]interface{}

	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return err
	}

	switch install := conf["install"].(type) {
	case map[string]interface{}:
		install["modpackUrl"] = fileLink
	default:
		return errors.New("Failed to modify server-setup-config.yaml")
	}

	new_conf, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	err = os.WriteFile(serverConfigPath, new_conf, 0666)

	return nil

}

func main() {
	fmt.Println("Reading configuration")
	config, err := config.ReadConfig(configPath)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Configuration loaded")

	fmt.Println("Fetching mod download link")
	fileLink, err := getModpackLatestFileLink(config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Found latest file link: " + fileLink)

	fmt.Println("Updating server-setup-config.yaml")
	err = UpdateServerSetupConfig(fileLink)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("server-setup-config.yaml updated")
}
