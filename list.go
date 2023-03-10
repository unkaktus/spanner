package spanner

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func showTable(jobList []Job) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "State", "Queue", "Time", "Nodes", "MPI"})

	for _, job := range jobList {
		timePercentage := int(100 * job.Walltime.Seconds() / job.RequestedWalltime.Seconds())
		table.Append([]string{
			job.Name,
			fmt.Sprintf("%s [%d]", job.State, job.ExitCode),
			job.Queue,
			fmt.Sprintf("%s/%s (%d%%)", job.Walltime, job.RequestedWalltime, timePercentage),
			strconv.Itoa(job.NodeNumber),
			fmt.Sprintf("%d/%d", job.MPIProcessNumber/job.NodeNumber, job.MPIProcessNumber),
		})
	}
	table.Render()
	return nil
}

func ListJobs(bs BatchSystem, state string) error {
	jobList, err := bs.ListJobs()
	if err != nil {
		return fmt.Errorf("query job list: %w", err)
	}

	jobMap := map[string]Job{}
	for _, job := range jobList {
		addedJob, ok := jobMap[job.Name]
		if ok {
			if job.CreationTime.After(addedJob.CreationTime) {
				jobMap[job.Name] = job
			}
		} else {
			jobMap[job.Name] = job
		}
	}

	jobList = []Job{}
	for _, job := range jobMap {
		// Skip the job with other states
		if state != "" && job.State != state {
			continue
		}
		jobList = append(jobList, job)
	}

	sort.Slice(jobList, func(i, j int) bool {
		return jobList[i].Name < jobList[j].Name
	})

	if err := showTable(jobList); err != nil {
		return fmt.Errorf("query list: %w", err)
	}

	return nil
}
