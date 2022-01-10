/*
Copyright © (2020-2022) Bizfly Cloud

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
	"fmt"
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"log"
)

var (
	scheduledVolumeBackupHeader = []string{"ID", "Volume ID", "Frequency", "Size", "Scheduled Hour", "Next Run",
		"Billing Plan", "Created At"}
	volumeId  string
	frequency string
	size      string
	hour      int
)

var scheduledVolumeBackupCmd = &cobra.Command{
	Use:   "schedule-volume-backup",
	Short: "Bizfly Cloud Scheduled Volume Backup",
	Long:  `Bizfly Cloud Scheduled Volume Backup Action: Create, List, Get, Delete, Update`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("schedule-volume-backup called")
	},
}

var scheduledVolumeBackupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List scheduled volume backup",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		backups, err := client.ScheduledVolumeBackup.List(ctx)
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, key := range backups {
			data = append(data, []string{
				key.ID,
				key.ResourceID,
				key.Options.Frequency,
				key.Options.Size,
				fmt.Sprintf("%d", key.ScheduledHour),
				key.NextRunAt,
				key.BillingPlan,
				key.CreatedAt,
			})
		}
		formatter.Output(scheduledVolumeBackupHeader, data)
	},
}

var scheduledVolumeBackupGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get scheduled volume backup",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) == 0 {
			log.Fatal("Please specify backup id")
		} else if len(args) > 1 {
			log.Fatal("Too many arguments")
		}
		backup, err := client.ScheduledVolumeBackup.Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{
			backup.ID,
			backup.ResourceID,
			backup.Options.Frequency,
			backup.Options.Size,
			fmt.Sprintf("%d", backup.ScheduledHour),
			backup.NextRunAt,
			backup.BillingPlan,
			backup.CreatedAt,
		})
		formatter.Output(scheduledVolumeBackupHeader, data)
	},
}

var scheduleVolumeBackupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create scheduled volume backup",
	Long:  "Create scheduled volume backup: bizfly schedule-volume-backup create <volume_id> --frequency=<frequency> --size=<size> --hour=<hour>",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) == 0 {
			log.Fatal("Please specify volume id")
		} else if len(args) > 1 {
			log.Fatal("Too many arguments")
		}
		volumeId = args[0]
		if len(frequency) == 0 {
			log.Fatal("Please specify frequency")
		}
		if len(size) == 0 {
			log.Fatal("Please specify size")
		}
		if hour == -1 {
			hour = 0
		} else if hour < 0 || hour > 23 {
			log.Fatal("Invalid hour")
		}
		payload := &gobizfly.CreateBackupPayload{
			ResourceID: volumeId,
			Frequency:  frequency,
			Size:       size,
			Hour:       hour,
		}
		backup, err := client.ScheduledVolumeBackup.Create(ctx, payload)
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{
			backup.ID,
			backup.ResourceID,
			backup.Options.Frequency,
			backup.Options.Size,
			fmt.Sprintf("%d", backup.ScheduledHour),
			backup.NextRunAt,
			backup.BillingPlan,
			backup.CreatedAt,
		})
		formatter.Output(scheduledVolumeBackupHeader, data)
	},
}

var scheduledVolumeBackupDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete scheduled volume backup",
	Long:  "Delete backup: bizfly schedule-volume-backup delete <backup_id>",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) == 0 {
			log.Fatal("Please specify backup id")
		} else if len(args) > 1 {
			log.Fatal("Too many arguments")
		}
		err := client.ScheduledVolumeBackup.Delete(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Backup deleted")
	},
}

var scheduledVolumeBackupUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update scheduled volume backup",
	Long:  "Update backup: bizfly scheduled-volume-backup update <backup_id> --frequency=<frequency> --size=<size> --hour=<hour>",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) == 0 {
			log.Fatal("Please specify backup id")
		} else if len(args) > 1 {
			log.Fatal("Too many arguments")
		}
		backup, err := client.ScheduledVolumeBackup.Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		payload := gobizfly.UpdateBackupPayload{}
		if frequency != "" {
			payload.Frequency = frequency
		}
		if size != "" {
			payload.Size = size
		}
		if hour != -1 {
			payload.Hour = hour
		} else {
			payload.Hour = backup.ScheduledHour
		}

		backup, err = client.ScheduledVolumeBackup.Update(ctx, args[0], &payload)
		if err != nil {
			log.Fatal(err)
		}

		var data [][]string
		data = append(data, []string{
			backup.ID,
			backup.ResourceID,
			backup.Options.Frequency,
			backup.Options.Size,
			fmt.Sprintf("%d", backup.ScheduledHour),
			backup.NextRunAt,
			backup.BillingPlan,
			backup.CreatedAt,
		})
		formatter.Output(scheduledVolumeBackupHeader, data)
	},
}

func init() {
	rootCmd.AddCommand(scheduledVolumeBackupCmd)
	scheduledVolumeBackupCmd.AddCommand(scheduledVolumeBackupListCmd)
	scheduledVolumeBackupCmd.AddCommand(scheduledVolumeBackupGetCmd)
	scheduledVolumeBackupCmd.AddCommand(scheduledVolumeBackupDeleteCmd)

	scheduledVolumeBackupCmd.AddCommand(scheduleVolumeBackupCreateCmd)
	bcpf := scheduleVolumeBackupCreateCmd.PersistentFlags()
	bcpf.IntVar(&hour, "hour", -1, "Schedule hour")
	bcpf.StringVar(&size, "size", "", "The size of keeping snapshot")
	bcpf.StringVar(&frequency, "frequency", "", "Frequency of backup snapshot")
	_ = cobra.MarkFlagRequired(bcpf, "size")
	_ = cobra.MarkFlagRequired(bcpf, "frequency")

	scheduledVolumeBackupCmd.AddCommand(scheduledVolumeBackupUpdateCmd)
	bupf := scheduledVolumeBackupUpdateCmd.PersistentFlags()
	bupf.IntVar(&hour, "hour", -1, "Schedule hour")
	bupf.StringVar(&size, "size", "", "The size of keeping snapshot")
	bupf.StringVar(&frequency, "frequency", "", "Frequency of backup snapshot")
}
