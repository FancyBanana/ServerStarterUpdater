package cursefoge

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	resty "github.com/go-resty/resty/v2"
)

const minecraftGameId = "432"
const serverAddress = "https://api.curseforge.com"

type curseForgeClient struct {
	client resty.Client
}

func New(apiKey string) curseForgeClient {
	c := *resty.New().SetHeader("Accept", "application/json")
	c.SetHeader("x-api-key", apiKey)
	return curseForgeClient{c}
}

// Finds mod/modpack ID by slug. Since slugs are unique we assume the request returns a single result
func (client *curseForgeClient) FindModIdBySlug(slug string) (int, error) {

	resp, err := client.client.R().SetQueryParams(
		map[string]string{
			"gameId": minecraftGameId,
			"slug":   slug,
		}).
		Get(serverAddress + "/v1/mods/search")
	if err != nil {
		return -1, err
	}

	var result struct {
		Data []struct {
			Id int
		}
	}
	err = json.Unmarshal(resp.Body(), &result)

	if err != nil {
		fmt.Println(err)
	}

	if len(result.Data) == 0 {
		return -1, errors.New("Could not find mod by slug:" + slug)
	}

	return result.Data[0].Id, nil

}

// Gets the most recent file for the mod/modpack.
// Since file ID's are incremental we take the most recent file by selecting the highest ID
func (client *curseForgeClient) GetModFile(modId int) (string, error) {

	resp, err := client.client.R().
		SetPathParam("modId", strconv.Itoa(modId)).
		Get(serverAddress + "/v1/mods/{modId}/files")
	if err != nil {
		return "", err
	}

	var result struct {
		Data []struct {
			Id          int
			DownloadUrl string
		}
	}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	fileId, fileLink := 0, ""
	for i := 0; i < len(result.Data); i++ {
		file := result.Data[i]
		if file.Id > fileId {
			fileId = file.Id
			fileLink = file.DownloadUrl
		}
	}

	return fileLink, nil
}
