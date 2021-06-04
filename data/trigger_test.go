package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/fxtrigger/data"
)

func TestTrigger_AddTrigger_ValidTrigger_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTrigger := data.Trigger{Name: "Unit test trigger", Description: "Unit test trigger desc", GPIOPin: 23, WebHooks: []data.WebHook{}}

	//	Act
	newTrigger, err := db.AddTrigger(testTrigger.Name, testTrigger.Description, testTrigger.GPIOPin, testTrigger.WebHooks, testTrigger.MinimumSecondsBeforeRetrigger)

	//	Assert
	if err != nil {
		t.Errorf("AddTrigger - Should add trigger without error, but got: %s", err)
	}

	if newTrigger.Created.IsZero() {
		t.Errorf("AddTrigger failed: Should have set an item with the correct datetime: %+v", newTrigger)
	}

	if newTrigger.Enabled != true {
		t.Errorf("AddTrigger failed: Should have enabled the trigger by default: %+v", newTrigger)
	}

}

func TestTrigger_GetTrigger_ValidTrigger_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTrigger1 := data.Trigger{Name: "Trigger 1", Description: "Unit test 1", GPIOPin: 11}
	testTrigger2 := data.Trigger{Name: "Trigger 2", Description: "Unit test 2", GPIOPin: 12}
	testTrigger3 := data.Trigger{Name: "Trigger 3", Description: "Unit test 3", GPIOPin: 13}

	//	Act
	db.AddTrigger(testTrigger1.Name, testTrigger1.Description, testTrigger1.GPIOPin, testTrigger1.WebHooks, testTrigger1.MinimumSecondsBeforeRetrigger)
	newTrigger2, _ := db.AddTrigger(testTrigger2.Name, testTrigger2.Description, testTrigger2.GPIOPin, testTrigger2.WebHooks, testTrigger2.MinimumSecondsBeforeRetrigger)
	db.AddTrigger(testTrigger3.Name, testTrigger3.Description, testTrigger3.GPIOPin, testTrigger3.WebHooks, testTrigger3.MinimumSecondsBeforeRetrigger)

	gotTrigger, err := db.GetTrigger(newTrigger2.ID)

	//	Log the file details:
	t.Logf("Trigger: %+v", gotTrigger)

	//	Assert
	if err != nil {
		t.Errorf("GetTrigger - Should get trigger without error, but got: %s", err)
	}

	if len(gotTrigger.ID) < 2 {
		t.Errorf("GetTrigger failed: Should get valid id but got: %v", gotTrigger.ID)
	}
}

func TestTrigger_GetTrigger_ValidTriggerWithWebhooks_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTrigger1 := data.Trigger{Name: "Trigger 1", Description: "Unit test 1", GPIOPin: 11}
	testTrigger2 := data.Trigger{
		Name:        "Trigger 2",
		Description: "Unit test 2",
		GPIOPin:     12,
		WebHooks: []data.WebHook{
			{URL: "http://www.github.com/webhook1",
				Headers: map[string]string{"header1": "value1", "header2": "value2"}},
			{URL: "http://www.microsoft.com/webhook2"}}}
	testTrigger3 := data.Trigger{Name: "Trigger 3", Description: "Unit test 3", GPIOPin: 13}

	//	Act
	db.AddTrigger(testTrigger1.Name, testTrigger1.Description, testTrigger1.GPIOPin, testTrigger1.WebHooks, testTrigger1.MinimumSecondsBeforeRetrigger)
	newTrigger2, _ := db.AddTrigger(testTrigger2.Name, testTrigger2.Description, testTrigger2.GPIOPin, testTrigger2.WebHooks, testTrigger2.MinimumSecondsBeforeRetrigger)
	db.AddTrigger(testTrigger3.Name, testTrigger3.Description, testTrigger3.GPIOPin, testTrigger3.WebHooks, testTrigger3.MinimumSecondsBeforeRetrigger)

	gotTrigger, err := db.GetTrigger(newTrigger2.ID)

	//	Log the file details:
	t.Logf("Trigger: %+v", gotTrigger)

	//	Assert
	if err != nil {
		t.Errorf("GetTrigger - Should get trigger without error, but got: %s", err)
	}

	if len(gotTrigger.ID) < 2 {
		t.Errorf("GetTrigger failed: Should get valid id but got: %v", gotTrigger.ID)
	}
}

func TestTrigger_GetAllTriggers_ValidTriggers_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTrigger1 := data.Trigger{Name: "Trigger 1", Description: "Unit test 1", GPIOPin: 11}
	testTrigger2 := data.Trigger{Name: "Trigger 2", Description: "Unit test 2", GPIOPin: 12}
	testTrigger3 := data.Trigger{Name: "Trigger 3", Description: "Unit test 3", GPIOPin: 13}

	//	Act
	db.AddTrigger(testTrigger1.Name, testTrigger1.Description, testTrigger1.GPIOPin, testTrigger1.WebHooks, testTrigger1.MinimumSecondsBeforeRetrigger)
	newTrigger2, _ := db.AddTrigger(testTrigger2.Name, testTrigger2.Description, testTrigger2.GPIOPin, testTrigger2.WebHooks, testTrigger2.MinimumSecondsBeforeRetrigger)
	db.AddTrigger(testTrigger3.Name, testTrigger3.Description, testTrigger3.GPIOPin, testTrigger3.WebHooks, testTrigger3.MinimumSecondsBeforeRetrigger)

	gotTriggers, err := db.GetAllTriggers()

	//	Assert
	if err != nil {
		t.Errorf("GetAllTriggers - Should get all triggers without error, but got: %s", err)
	}

	if len(gotTriggers) < 2 {
		t.Errorf("GetAllTriggers failed: Should get all items but got: %v", len(gotTriggers))
	}

	if gotTriggers[1].Description != newTrigger2.Description {
		t.Errorf("GetAllTriggers failed: Should get an item with the correct details: %+v", gotTriggers[1])
	}
}

func TestTrigger_UpdateTrigger_ValidTriggers_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTrigger1 := data.Trigger{Name: "Trigger 1", Description: "Unit test 1", GPIOPin: 11}
	testTrigger2 := data.Trigger{Name: "Trigger 2", Description: "Unit test 2", GPIOPin: 12}
	testTrigger3 := data.Trigger{Name: "Trigger 3", Description: "Unit test 3", GPIOPin: 13}

	//	Act
	db.AddTrigger(testTrigger1.Name, testTrigger1.Description, testTrigger1.GPIOPin, testTrigger1.WebHooks, testTrigger1.MinimumSecondsBeforeRetrigger)
	newTrigger2, _ := db.AddTrigger(testTrigger2.Name, testTrigger2.Description, testTrigger2.GPIOPin, testTrigger2.WebHooks, testTrigger2.MinimumSecondsBeforeRetrigger)
	db.AddTrigger(testTrigger3.Name, testTrigger3.Description, testTrigger3.GPIOPin, testTrigger3.WebHooks, testTrigger3.MinimumSecondsBeforeRetrigger)
	//	Update the 2nd trigger:
	newTrigger2.Enabled = false
	_, err = db.UpdateTrigger(newTrigger2) //	Update the 2nd trigger

	gotTrigger, _ := db.GetTrigger(newTrigger2.ID) // Refetch to verify

	//	Assert
	if err != nil {
		t.Errorf("UpdateTrigger - Should update trigger without error, but got: %s", err)
	}

	if gotTrigger.Enabled != false {
		t.Errorf("UpdateTrigger failed: Should get an item that has been disabled but got: %+v", gotTrigger)
	}

}

func TestTrigger_DeleteTrigger_ValidTriggers_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTrigger1 := data.Trigger{Name: "Trigger 1", Description: "Unit test 1", GPIOPin: 11}
	testTrigger2 := data.Trigger{Name: "Trigger 2", Description: "Unit test 2", GPIOPin: 12}
	testTrigger3 := data.Trigger{Name: "Trigger 3", Description: "Unit test 3", GPIOPin: 13}

	//	Act
	db.AddTrigger(testTrigger1.Name, testTrigger1.Description, testTrigger1.GPIOPin, testTrigger1.WebHooks, testTrigger1.MinimumSecondsBeforeRetrigger)
	newTrigger2, _ := db.AddTrigger(testTrigger2.Name, testTrigger2.Description, testTrigger2.GPIOPin, testTrigger2.WebHooks, testTrigger2.MinimumSecondsBeforeRetrigger)
	db.AddTrigger(testTrigger3.Name, testTrigger3.Description, testTrigger3.GPIOPin, testTrigger3.WebHooks, testTrigger3.MinimumSecondsBeforeRetrigger)
	err = db.DeleteTrigger(newTrigger2.ID) //	Delete the 2nd trigger

	gotTriggers, _ := db.GetAllTriggers()

	//	Assert
	if err != nil {
		t.Errorf("DeleteTrigger - Should delete trigger without error, but got: %s", err)
	}

	if len(gotTriggers) != 2 {
		t.Errorf("DeleteTrigger failed: Should remove an item but got: %v", len(gotTriggers))
	}

	if gotTriggers[1].Description == newTrigger2.Description {
		t.Errorf("DeleteTrigger failed: Should get an item with different details than the removed item but got: %+v", gotTriggers[1])
	}

}
