package db

import (
	"log"

	"gorm.io/gorm"

	"github.com/jackc/pgtype"
)

type Controller interface {
	LoadModel(model interface{}, id uint) error
	// CRUD
	// Create
	CreateUser(*User) error
	CreateMachine(*Machine) error
	CreateOrganization(*Organization) error
	CreateSchedule(*Schedule) error
	CreateTask(*Task) error
	CreateUserToken(*UserToken) error
	CreateMachineToken(*MachineToken) error
	CreateLoginInfo(*LoginInfo) error
	CreateRecord(*Record) error
	// Read
	ReadUser(name string) (*User, error)
	ReadMachine(name string) (*Machine, error)
	ReadOrganization(name string) (*Organization, error)
	ReadSchedule(machineName string) (*Schedule, error)
	ReadTask(name string) (*Task, error)
	ReadUserToken(value pgtype.UUID) (*UserToken, error)
	ReadMachineToken(value pgtype.UUID) (*MachineToken, error)
	ReadLoginInfo(username string) (*LoginInfo, error)
	ReadRecords(machineName string) ([]Record, error)
	// Update
	UpdateUser(*User) error
	UpdateMachine(*Machine) error
	UpdateOrganization(*Organization) error
	UpdateSchedule(*Schedule) error
	UpdateTask(*Task) error
	UpdateUserToken(*UserToken) error
	UpdateMachineToken(*MachineToken) error
	UpdateLoginInfo(*LoginInfo) error
	UpdateRecord(*Record) error
	// Delete
	DeleteUser(name string) error
	DeleteMachine(name string) error
	DeleteOrganization(name string) error
	DeleteSchedule(machineName string) error
	DeleteTask(name string) error
	DeleteUserToken(value pgtype.UUID) error
	DeleteMachineToken(value pgtype.UUID) error
	DeleteLoginInfo(username string) error
	DeleteRecords(machineName string) error
	DeleteRecord(machineName string, recordID uint64) error
}

type controller struct {
	db *gorm.DB
}

var _ Controller = &controller{}

func NewController(db *gorm.DB) Controller {
	return &controller{
		db: db,
	}
}

type dbError string

func (d dbError) Error() string { return string(d) }

var (
	unimplemented = dbError("unimplemented")
	noDB          = dbError("no db connection")
)

func (c *controller) LoadModel(model interface{}, id uint) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.First(model, id)
	if err := res.Error; err != nil {
		log.Println("error loading model:", err)
		return err
	}
	return nil
}

func (c *controller) CreateUser(user *User) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Create(user)
	if err := res.Error; err != nil {
		log.Println("error creating User:", err)
		return err
	}
	log.Println("inserted user with id:", user.ID)
	return nil
}

func (c *controller) CreateMachine(machine *Machine) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Create(machine)
	if err := res.Error; err != nil {
		log.Println("error creating Machine:", err)
		return err
	}
	log.Println("inserted machine with id:", machine.ID)
	return nil
}

func (c *controller) CreateOrganization(org *Organization) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Create(org)
	if err := res.Error; err != nil {
		log.Println("error creating Organization:", err)
		return err
	}
	log.Println("inserted organization with id:", org.ID)
	return nil
}

func (c *controller) CreateSchedule(schedule *Schedule) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Create(schedule)
	if err := res.Error; err != nil {
		log.Println("error creating Schedule:", err)
		return err
	}
	log.Println("inserted schedule with id:", schedule.ID)
	return nil
}

func (c *controller) CreateTask(task *Task) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Create(task)
	if err := res.Error; err != nil {
		log.Println("error creating Task:", err)
		return err
	}
	log.Println("inserted task with ID:", task.ID)
	return nil
}

func (c *controller) CreateUserToken(userToken *UserToken) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Create(userToken)
	if err := res.Error; err != nil {
		log.Println("error creating UserToken:", err)
		return err
	}
	log.Println("inserted user token with ID:", userToken.ID)
	return nil
}

func (c *controller) CreateMachineToken(machineToken *MachineToken) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Create(machineToken)
	if err := res.Error; err != nil {
		log.Println("error creating MachineToken:", err)
		return err
	}
	log.Println("inserted machine token with ID:", machineToken.ID)
	return nil
}

func (c *controller) CreateLoginInfo(loginInfo *LoginInfo) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Create(loginInfo)
	if err := res.Error; err != nil {
		log.Println("error creating LoginInfo:", err)
		return err
	}
	log.Println("inserted LoginInfo with ID:", loginInfo.ID)
	return nil
}

func (c *controller) CreateRecord(record *Record) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Create(record)
	if err := res.Error; err != nil {
		log.Println("error creating Record:", err)
		return err
	}
	log.Println("inserted Record with ID:", record.ID)
	return nil
}

func (c *controller) ReadUser(name string) (*User, error) {
	if c == nil || c.db == nil {
		return nil, noDB
	}

	var user User
	res := c.db.First(&user, `name = ?`, name)
	err := res.Error
	if err != nil {
		return nil, err
	}
	log.Println("found User with ID:", user.ID)

	return &user, nil
}

func (c *controller) ReadMachine(name string) (*Machine, error) {
	if c == nil || c.db == nil {
		return nil, noDB
	}

	var machine Machine
	res := c.db.First(&machine, `name = ?`, name)
	err := res.Error
	if err != nil {
		return nil, err
	}
	log.Println("found Machine with ID:", machine.ID)

	return &machine, nil
}

func (c *controller) ReadOrganization(name string) (*Organization, error) {
	if c == nil || c.db == nil {
		return nil, noDB
	}

	var org Organization
	res := c.db.Preload("Users").Preload("Machines").Preload("Tasks").First(&org, `name = ?`, name)
	err := res.Error
	if err != nil {
		return nil, err
	}
	log.Println("found Organization with ID:", org.ID)

	return &org, nil
}

func (c *controller) ReadSchedule(machineName string) (*Schedule, error) {
	if c == nil || c.db == nil {
		return nil, noDB
	}

	machine, err := c.ReadMachine(machineName)
	if err != nil {
		return nil, err
	}

	var schedule Schedule
	res := c.db.Preload("Machine").First(&schedule, `machine_id = ?`, machine.ID)
	err = res.Error
	if err != nil {
		return nil, err
	}
	log.Println("found Schedule with ID:", schedule.ID)

	return &schedule, nil
}

func (c *controller) ReadTask(name string) (*Task, error) {
	if c == nil || c.db == nil {
		return nil, noDB
	}

	var task Task
	res := c.db.First(&task, `name = ?`, name)
	err := res.Error
	if err != nil {
		return nil, err
	}
	log.Println("found Task with ID:", task.ID)

	return &task, nil
}

func (c *controller) ReadUserToken(value pgtype.UUID) (*UserToken, error) {
	if c == nil || c.db == nil {
		return nil, noDB
	}

	var userToken UserToken
	res := c.db.Preload("User").Where(`value = ?`, value).First(&userToken)
	err := res.Error
	if err != nil {
		return nil, err
	}
	log.Println("found UserToken with ID:", userToken.ID)

	return &userToken, nil
}

func (c *controller) ReadMachineToken(value pgtype.UUID) (*MachineToken, error) {
	if c == nil || c.db == nil {
		return nil, noDB
	}

	var machineToken MachineToken
	res := c.db.Preload("Machine").Where(`value = ?`, value).First(&machineToken)
	err := res.Error
	if err != nil {
		return nil, err
	}
	log.Println("found MachineToken with ID:", machineToken.ID)

	return &machineToken, nil
}

func (c *controller) ReadLoginInfo(username string) (*LoginInfo, error) {
	if c == nil || c.db == nil {
		return nil, noDB
	}

	var loginInfo LoginInfo
	res := c.db.Preload("User").Where(`username = ?`, username).First(&loginInfo)
	err := res.Error
	if err != nil {
		return nil, err
	}
	log.Println("found LoginInfo with ID:", loginInfo.ID)

	return &loginInfo, nil
}

func (c *controller) ReadRecords(machineName string) ([]Record, error) {
	if c == nil || c.db == nil {
		return nil, noDB
	}

	machine, err := c.ReadMachine(machineName)
	if err != nil {
		return nil, err
	}

	var records []Record
	res := c.db.Preload("Task").Preload("Machine").Where(`machine_id = ?`, machine.ID).Find(&records)
	err = res.Error
	if err != nil {
		return nil, err
	}

	log.Printf("found %d Record(s) for machine \"%s\"\n", len(records), machineName)

	return records, nil
}

func (c *controller) UpdateUser(user *User) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Save(user)
	err := res.Error
	if err != nil {
		return err
	}
	log.Println("Saved User with ID:", user.ID)

	return nil
}

func (c *controller) UpdateMachine(machine *Machine) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Save(machine)
	err := res.Error
	if err != nil {
		return err
	}
	log.Println("Saved Machine with ID:", machine.ID)

	return nil
}

func (c *controller) UpdateOrganization(org *Organization) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Save(org)
	err := res.Error
	if err != nil {
		return err
	}
	log.Println("Saved Organization with ID:", org.ID)

	return nil
}

func (c *controller) UpdateSchedule(schedule *Schedule) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Save(schedule)
	err := res.Error
	if err != nil {
		return err
	}
	log.Println("Saved Schedule with ID:", schedule.ID)

	return nil
}

func (c *controller) UpdateTask(task *Task) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Save(task)
	err := res.Error
	if err != nil {
		return err
	}
	log.Println("Saved Task with ID:", task.ID)

	return nil
}

func (c *controller) UpdateUserToken(userToken *UserToken) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Save(userToken)
	err := res.Error
	if err != nil {
		return err
	}
	log.Println("Saved UserToken with ID:", userToken.ID)

	return nil
}

func (c *controller) UpdateMachineToken(machineToken *MachineToken) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Save(machineToken)
	err := res.Error
	if err != nil {
		return err
	}
	log.Println("Saved MachineToken with ID:", machineToken.ID)

	return nil
}

func (c *controller) UpdateLoginInfo(loginInfo *LoginInfo) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Save(loginInfo)
	err := res.Error
	if err != nil {
		return err
	}
	log.Println("Saved LoginInfo with ID:", loginInfo.ID)

	return nil
}

func (c *controller) UpdateRecord(record *Record) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Save(record)
	err := res.Error
	if err != nil {
		return err
	}
	log.Println("Saved Record with ID:", record.ID)

	return nil
}

func (c *controller) DeleteUser(name string) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Where(`name = ?`, name).Delete(&User{})
	err := res.Error
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) DeleteMachine(name string) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Where(`name = ?`, name).Delete(&Machine{})
	err := res.Error
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) DeleteOrganization(name string) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Where(`name = ?`, name).Delete(&Organization{})
	err := res.Error
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) DeleteSchedule(machineName string) error {
	if c == nil || c.db == nil {
		return noDB
	}

	machine, err := c.ReadMachine(machineName)
	if err != nil {
		return err
	}

	res := c.db.Where(`machine_id = ?`, machine.ID).Delete(&Schedule{})
	err = res.Error
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) DeleteTask(name string) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Where(`name = ?`, name).Delete(&Task{})
	err := res.Error
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) DeleteUserToken(value pgtype.UUID) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Where(`value = ?`, value).Delete(&UserToken{})
	err := res.Error
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) DeleteMachineToken(value pgtype.UUID) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Where(`value = ?`, value).Delete(&MachineToken{})
	err := res.Error
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) DeleteLoginInfo(username string) error {
	if c == nil || c.db == nil {
		return noDB
	}

	res := c.db.Where(`username = ?`, username).Delete(&LoginInfo{})
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (c *controller) DeleteRecords(machineName string) error {
	if c == nil || c.db == nil {
		return noDB
	}

	machine, err := c.ReadMachine(machineName)
	if err != nil {
		return err
	}

	res := c.db.Where(`machine_id = ?`, machine.ID).Delete(&Record{})
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (c *controller) DeleteRecord(machineName string, recordID uint64) error {
	if c == nil || c.db == nil {
		return noDB
	}

	machine, err := c.ReadMachine(machineName)
	if err != nil {
		return err
	}

	res := c.db.Where(`machine_id = ? and id = ?`, machine.ID, recordID).Delete(&Record{})
	if err := res.Error; err != nil {
		return err
	}
	return nil
}
