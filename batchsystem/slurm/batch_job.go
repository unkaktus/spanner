package slurm

import (
	"fmt"
	"strings"

	"github.com/unkaktus/spanner/batchsystem"
)

func (b *Slurm) JobData(job batchsystem.Job) (string, error) {
	header, err := batchsystem.ExecTemplate(`#!/bin/bash -l
#SBATCH -J {{.Name}}
#SBATCH -o {{.OutputFile}}
#SBATCH -e {{.ErrorFile}}
#SBATCH --mail-type=ALL
#SBATCH --mail-user={{.Email}}
#SBATCH --nodes {{.Nodes}}
#SBATCH --ntasks-per-node {{.TasksPerNode}}
#SBATCH --time={{.Walltime}}
`,
		job)
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	jobData := header

	if job.WorkingDirectory != "" {
		jobData += fmt.Sprintf("cd %s\n", job.WorkingDirectory)
	}

	if len(job.InitScript) > 0 {
		for _, line := range job.InitScript {
			jobData += fmt.Sprintf("%s\n", line)
		}
	}

	task := []string{
		"srun", "spanner", "tent",
	}

	task = append(task, job.Runtime...)
	task = append(task, job.Executable)

	for _, argument := range job.Arguments {
		if strings.Contains(argument, "{{.ConfigFilename}}") {
			argument = strings.ReplaceAll(argument, "{{.ConfigFilename}}", job.ConfigFilename)
		}
		task = append(task, argument)
	}

	jobData += strings.Join(task, " ")
	jobData += "\n"

	return jobData, nil
}
