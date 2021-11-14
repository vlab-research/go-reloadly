package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-playground/validator/v10"
	"github.com/jszwec/csvutil"
	"github.com/nandanrao/chance"
	"github.com/spf13/cobra"
	"github.com/vlab-research/go-reloadly/reloadly"
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

	validate := validator.New()
	for _, job := range jobs {
		err = validate.Struct(job)
		if err != nil {
			return jobs, err
		}
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

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Make airtime recharges to multiple mobile numbers using a CSV file",
	Long:  "Make airtime recharges to multiple mobile numbers using a CSV file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires 2 positional args [input csv] and [output csv]")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		output := args[1]

		svc, err := LoadService(cmd)
		if err != nil {
			return err
		}

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

	batchCmd.Flags().IntP("workers", "w", 12, "Parallelism for http requests")
}
