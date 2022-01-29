package feargreedindex

import (
	"crypto-selling-algorithm/common"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const timeLayout = time.RFC3339

type FearAndGreedDatapoint struct {
	Timestamp        time.Time
	Value            int
	ChangePercentage float64
}

type FearAndGreedProvider interface {
	Get(days int) (fgResponse, error)
	FromFile(csvFile string) ([]*FearAndGreedDatapoint, error)
}

type fearAndGreedProvider struct {
}

func New() FearAndGreedProvider {
	return &fearAndGreedProvider{}
}

type fgResponse struct {
	Name     string     `json:"name"`
	Data     []fgData   `json:"data"`
	Metadata fgMetadata `json:"metadata"`
}

type fgMetadata struct {
	Err string `json:"error"`
}

type fgData struct {
	Value               string `json:"value"`
	ValueClassification string `json:"value_classification"`
	Timestamp           string `json:"timestamp"`
	TimeUntilUpdate     string `json:"time_until_update"`
}

func (p *fearAndGreedProvider) Get(days int) (fgResponse, error) {
	var fgr fgResponse

	c := http.Client{}
	res, err := c.Get(fmt.Sprintf("https://api.alternative.me/fng/?limit=%d", days))
	if err != nil {
		return fgr, err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&fgr)
	if err != nil {
		return fgr, err
	}

	return fgr, nil
}

func (p *fearAndGreedProvider) FromFile(csvFile string) ([]*FearAndGreedDatapoint, error) {
	allLines, err := common.ReadFile(csvFile)
	if err != nil {
		return nil, err
	}

	return p.mapLines(allLines)
}

func (p *fearAndGreedProvider) mapLines(allLines [][]string) ([]*FearAndGreedDatapoint, error) {
	datapoints := make([]*FearAndGreedDatapoint, 0, len(allLines))

	for idx, line := range allLines {

		// skip header
		if idx == 0 {
			continue
		}

		date := line[0]
		v := line[1]

		timestamp, err := time.Parse(timeLayout, date)
		if err != nil {
			return nil, err
		}

		value, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		datapoints = append(datapoints, &FearAndGreedDatapoint{
			Timestamp: timestamp,
			Value:     value,
		})

	}

	return datapoints, nil
}
