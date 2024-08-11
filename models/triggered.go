package models

import "github.com/mstgnz/cronjob/pkg/config"

type Triggered struct {
	ScheduleID int `json:"schedule_id"`
}

func (m *Triggered) Create(schedule_id int) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["TRIGGERED_INSERT"])
	if err != nil {
		return err
	}

	_, err = stmt.Exec(schedule_id)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	return nil
}

func (m *Triggered) Delete(schedule_id int) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["TRIGGERED_DELETE"])
	if err != nil {
		return err
	}

	_, err = stmt.Exec(schedule_id)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	return nil
}
