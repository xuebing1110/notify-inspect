package redis

import (
	"testing"
	"time"

	"github.com/xuebing1110/notify-inspect/pkg/plugin"
	"github.com/xuebing1110/notify-inspect/pkg/plugin/storage"
	_ "github.com/xuebing1110/notify-inspect/pkg/plugin/storage/redis"
	"github.com/xuebing1110/notify-inspect/pkg/schedule"
	"github.com/xuebing1110/notify-inspect/pkg/schedule/cron"
)

func TestScheduler(t *testing.T) {
	record := &plugin.PluginRecord{
		Id:       "1",
		UserId:   "admin",
		PluginId: "test",
		Disable:  "False",
		Cron: &cron.CronTaskSetting{
			Interval:        "1m",
			ClockLimitStart: "08:00",
			ClockLimitEnd:   "22:00",
			WeekLimit:       "notHoliday",
		},
		Data: []plugin.PluginData{
			{"filed1", "字段1", "value1"},
		},
	}

	err := storage.GlobalStorage.SavePluginRecord(record)
	if err != nil {
		t.Fatal(err)
	}

	curtime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-01-15 08:30:00", time.Local)

	task := record.GetCronTask()
	err = schedule.DefaultScheduler.PutTask(task, curtime)
	if err != nil {
		t.Fatal(err)
	}

	tasks := make([]*cron.CronTask, 0)
	for task := range schedule.DefaultScheduler.FetchTasks(curtime.Add(time.Minute)) {
		tasks = append(tasks, task)
	}
	if len(tasks) <= 0 {
		t.Fatalf("expect a task, but get nothing!")
	}
}
