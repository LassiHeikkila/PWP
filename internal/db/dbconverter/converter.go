package dbconverter

import (
	"encoding/json"

	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func ConvertOrganization(dborg *db.Organization) types.Organization {
	return types.Organization{
		Name: dborg.Name,
	}
}

func ConvertOrganizationToDB(org *types.Organization) db.Organization {
	return db.Organization{
		Name: org.Name,
	}
}

func ConvertUser(dbuser *db.User) types.User {
	return types.User{
		Name:         dbuser.Name,
		Email:        dbuser.Email,
		Role:         dbuser.Role,
	}
}

func ConvertUserToDB(user *types.User) db.User {
	return db.User{
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

func ConvertMachine(dbmachine *db.Machine) types.Machine {
	return types.Machine{
		Name:        dbmachine.Name,
		Description: dbmachine.Description,
		OS:          dbmachine.OS,
		Arch:        dbmachine.Arch,
	}
}

func ConvertMachineToDB(machine *types.Machine) db.Machine {
	return db.Machine{
		Name:        machine.Name,
		Description: machine.Description,
		OS:          machine.OS,
		Arch:        machine.Arch,
	}
}

func ConvertTask(dbtask *db.Task) types.Task {
	c := make(map[string]interface{})
	_ = json.Unmarshal(dbtask.Content.Bytes, &c)

	return types.Task{
		Name:        dbtask.Name,
		Description: dbtask.Description,
		Content:     c,
	}
}

func ConvertTaskToDB(task *types.Task) db.Task {
	b, _ := json.Marshal(&task.Content)

	return db.Task{
		Name:        task.Name,
		Description: task.Description,
		Content:     db.StringToJSON(string(b)),
	}
}

func ConvertRecord(dbrecord *db.Record) types.Record {
	return types.Record{
		MachineName: dbrecord.Machine.Name,
		TaskName:    dbrecord.Task.Name,
		ExecutedAt:  dbrecord.ExecutedAt,
		Status:      dbrecord.Status,
		Output:      dbrecord.Output,
	}
}

func ConvertRecordToDB(record *types.Record) db.Record {
	return db.Record{
		// cannot set Machine or Task
		ExecutedAt: record.ExecutedAt,
		Status:     record.Status,
		Output:     record.Output,
	}
}

func ConvertSchedule(dbschedule *db.Schedule) types.Schedule {
	c := make(map[string]interface{})
	_ = json.Unmarshal(dbschedule.Content.Bytes, &c)

	return types.Schedule{
		Content: c,
	}
}

func ConvertScheduleToDB(schedule *types.Schedule) db.Schedule {
	b, _ := json.Marshal(&schedule.Content)

	return db.Schedule{
		Content: db.StringToJSON(string(b)),
	}
}

func ConvertUserToken(dbtoken *db.UserToken) types.UserToken {
	var s string
	_ = dbtoken.Value.AssignTo(&s)
	return types.UserToken(s)
}

func ConvertMachineToken(dbtoken *db.MachineToken) types.MachineToken {
	var s string
	_ = dbtoken.Value.AssignTo(&s)
	return types.MachineToken(s)
}
