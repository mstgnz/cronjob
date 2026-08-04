package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mstgnz/cronjob/pkg/config"
)

// Triggered is the distributed lock guarding a schedule's execution. One row means
// one instance is running that schedule right now; the primary key on schedule_id
// is what keeps every other instance out.
type Triggered struct {
	ScheduleID int    `json:"schedule_id"`
	InstanceID string `json:"instance_id"`
}

// Acquire tries to take the lock for the given schedule.
// It reports true only when this instance now holds it, which happens when the
// schedule was idle or when the previous holder's lease had already run out.
// Contention is not an error: losing the race simply means another instance is
// running this schedule, so false is the normal result for every replica but one.
func (m *Triggered) Acquire(scheduleID int, instanceID string, lease time.Duration) (bool, error) {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["TRIGGERED_ACQUIRE"])
	if err != nil {
		return false, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	var acquired int
	err = stmt.QueryRow(scheduleID, instanceID, lease.Seconds()).Scan(&acquired)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	m.ScheduleID = scheduleID
	m.InstanceID = instanceID
	return true, nil
}

// Release gives the lock back. It only deletes this instance's own row, so an
// instance whose lease expired and was taken over cannot cut the new holder short.
func (m *Triggered) Release(scheduleID int, instanceID string) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["TRIGGERED_RELEASE"])
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(scheduleID, instanceID)
	return err
}
