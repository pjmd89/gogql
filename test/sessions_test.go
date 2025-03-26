package test

// import (
// 	"os"
// 	"testing"

// 	"github.com/pjmd89/gogql/lib/http/sessions"
// )

// func TestFileProvider_Init_Set_Get_Count_Destroy(t *testing.T) {
// 	tempDir := t.TempDir()

// 	sessMngr := sessions.NewSessionManager(sessions.FILE_PROVIDER, tempDir, 2)

// 	sessionID := "test_session"
// 	sessionData := struct {
// 		Name string
// 		Age  int
// 	}{
// 		Name: "John Doe",
// 		Age:  30,
// 	}
// 	sessMngr.SessionProvider.Init(sessionID, sessionData)

// 	// Prueba Set
// 	newData := struct {
// 		Name string
// 		Age  int
// 		City string
// 	}{
// 		Name: "Jane Doe",
// 		Age:  25,
// 		City: "New York",
// 	}
// 	err := sessMngr.SessionProvider.Set(sessionID, newData)
// 	if err != nil {
// 		t.Errorf("Set() error = %v", err)
// 	}

// 	// Prueba Get
// 	var retrievedData struct {
// 		Name string
// 		Age  int
// 		City string
// 	}
// 	_, err = sessMngr.SessionProvider.Get(sessionID, &retrievedData)
// 	if err != nil {
// 		t.Errorf("Get() error = %v", err)
// 	}

// 	if retrievedData.Name != newData.Name || retrievedData.Age != newData.Age || retrievedData.City != newData.City {
// 		t.Errorf("Get() returned incorrect data: got %+v, want %+v", retrievedData, newData)
// 	}

// 	// Prueba Count
// 	count, err := sessMngr.SessionProvider.Count()
// 	if err != nil {
// 		t.Errorf("Count() error = %v", err)
// 	}
// 	if count != 1 {
// 		t.Errorf("Count() should return 1, got %d", count)
// 	}

// 	// Prueba Destroy
// 	err = sessMngr.SessionProvider.Destroy(sessionID)
// 	if err != nil {
// 		t.Errorf("Destroy() error = %v", err)
// 	}

// 	count, err = sessMngr.SessionProvider.Count()
// 	if err != nil {
// 		t.Errorf("Count() error = %v", err)
// 	}

// 	if count != 0 {
// 		t.Errorf("Count() should return 0 after Destroy(), got %d", count)
// 	}
// }

// func TestFileProvider_Set_NilData(t *testing.T) {
// 	tempDir := t.TempDir()
// 	sessMngr := sessions.NewSessionManager(sessions.FILE_PROVIDER, tempDir, 2)
// 	sessionID := "test_session"

// 	err := sessMngr.SessionProvider.Set(sessionID, nil)
// 	if err != nil {
// 		t.Errorf("Set() with nil data should not return an error, got %v", err)
// 	}

// 	_, err = os.Stat(tempDir + sessionID + ".gob")
// 	if !os.IsNotExist(err) {
// 		t.Errorf("Set() with nil data should not create a file")
// 	}
// }

// func TestFileProvider_Get_NilReceiver(t *testing.T) {
// 	tempDir := t.TempDir()
// 	sessMngr := sessions.NewSessionManager(sessions.FILE_PROVIDER, tempDir, 2)
// 	sessionID := "test_session"
// 	sessionData := struct{ Name string }{Name: "Test"}

// 	sessMngr.SessionProvider.Set(sessionID, sessionData)
// 	sess := struct{ Name string }{}
// 	_, err := sessMngr.SessionProvider.Get(sessionID, &sess)
// 	if err != nil {
// 		t.Errorf("Get() with nil receiver should not return an error, got %v", err)
// 	}
// }
