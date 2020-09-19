/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/jszwec/csvutil"
	"github.com/nandanrao/chance"
	"github.com/nandanrao/go-reloadly/reloadly"
	"github.com/spf13/cobra"
)

func LoadBatchCsv(path string) ([]reloadly.TopupJob, error) {
	var jobs []reloadly.TopupJob
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = csvutil.Unmarshal(b, &jobs)

	if len(jobs) == 0 {
		return jobs, fmt.Errorf("We could not parse the data from the csv. Please ensure it is in the right format and that the required fields (number, amount, country) are present for each row in the csv file.")
	}

	return jobs, nil
}

func WriteBatchCsv(path string, responses []*reloadly.TopupWorkerResponse) error {
	b, err := csvutil.Marshal(responses)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0644)
	return err
}

func BatchTopup(svc *reloadly.Service, numWorkers int, jobs []reloadly.TopupJob) []*reloadly.TopupWorkerResponse {
	tt := make([]interface{}, len(jobs))
	for i := range jobs {
		tt[i] = &jobs[i]
	}

	tasks := chance.Flatten(tt)
	worker := reloadly.TopupWorker(*svc)
	outputs := chance.Pool(numWorkers, tasks, worker.Work)

	res := []*reloadly.TopupWorkerResponse{}
	for r := range outputs {
		rr := r.(*reloadly.TopupWorkerResponse)
		res = append(res, rr)
	}

	return res
}



// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires 2 positional args [input csv] and [output csv]")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := LoadService(cmd)
		if err != nil {
			return err
		}

		input := args[0]
		output := args[1]

		details, err := LoadBatchCsv(input)
		if err != nil {
			return err
		}

		numWorkers, err := cmd.Flags().GetInt("workers")
		if err != nil {
			return err
		}

		responses := BatchTopup(svc, numWorkers, details)
		err = WriteBatchCsv(output, responses)
		if err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("Successfully wrote %v responses from %v rows", len(responses), len(details)))
		return nil
	},
}

func init() {
	topupsCmd.AddCommand(batchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// batchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	batchCmd.Flags().IntP("workers", "w", 12, "Parallelism for http requests")
}
