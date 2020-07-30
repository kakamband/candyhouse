package main

import (
	"bytes"
	"candyHouse/models/db"
	"candyHouse/models/viewmodels"
	"candyHouse/routers"
	"candyHouse/utils"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benweissmann/memongo"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type SingupResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

//testSetup ...
func testSetup() (*memongo.Server, *gin.Engine) {
	binPath := utils.EnvVar("IN_MEMORY_MONGO_DB_PATH")

	mongoServer, err := memongo.StartWithOptions(&memongo.Options{
		MongodBin:      binPath,
		StartupTimeout: 500000,
	})

	if err != nil {
		log.Println(err)
	}

	db.ConnectToDatabase(mongoServer.URI(), memongo.RandomDatabase())
	router := routers.InitRoute()
	return mongoServer, router
}

//setupTestAccount ...
func setupTestAccount(router *gin.Engine) {
	talentRegisterModel := &viewmodels.TalentRegister{
		Email:     "shuvo009@yahoo.com",
		FirstName: "shuvo",
		LastName:  "hasan",
		Password:  "123456",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(talentRegisterModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/talent/register", buf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
}

//TestTalentRegister /talent/register/
func TestTalentRegister(t *testing.T) {
	mongoServer, router := testSetup()

	talentRegisterModel := &viewmodels.TalentRegister{
		Email:     "shuvo009@yahoo.com",
		FirstName: "shuvo",
		LastName:  "hasan",
		Password:  "123456",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(talentRegisterModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/talent/register", buf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	singupResponse := SingupResponse{}
	body, _ := ioutil.ReadAll(w.Body)

	json.Unmarshal([]byte(string(body)), &singupResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "talent is created", singupResponse.Message)
	mongoServer.Stop()
}

//TestDuplicateTalentRegister /talent/register/
func TestDuplicateTalentRegister(t *testing.T) {
	mongoServer, router := testSetup()

	setupTestAccount(router)

	talentRegisterModel := &viewmodels.TalentRegister{
		Email:     "shuvo009@yahoo.com",
		FirstName: "shuvo",
		LastName:  "hasan",
		Password:  "123456",
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(talentRegisterModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/talent/register", buf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	singupResponse := SingupResponse{}
	body, _ := ioutil.ReadAll(w.Body)

	json.Unmarshal([]byte(string(body)), &singupResponse)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "Already Exist", singupResponse.Error)
	mongoServer.Stop()
}

//TestTalentLogin /talent/login/
func TestTalentLogin(t *testing.T) {
	mongoServer, router := testSetup()

	setupTestAccount(router)

	loginModel := &viewmodels.LoginModel{
		Username: "shuvo009@yahoo.com",
		Password: "123456",
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(loginModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/talent/login", buf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	loginResponse := LoginResponse{}
	body, _ := ioutil.ReadAll(w.Body)

	json.Unmarshal([]byte(string(body)), &loginResponse)

	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "", loginResponse.Token)
	mongoServer.Stop()
}

//TestTalentInvalidLogin /talent/login/
func TestTalentInvalidUsernameLogin(t *testing.T) {
	mongoServer, router := testSetup()

	setupTestAccount(router)

	loginModel := &viewmodels.LoginModel{
		Username: "shuvo009_wrong@yahoo.com",
		Password: "123456",
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(loginModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/talent/login", buf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	loginResponse := LoginResponse{}
	body, _ := ioutil.ReadAll(w.Body)

	json.Unmarshal([]byte(string(body)), &loginResponse)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "Username or password is invalid", loginResponse.Error)
	mongoServer.Stop()
}

//TestTalentInvalidPasswordLogin ...
func TestTalentInvalidPasswordLogin(t *testing.T) {
	mongoServer, router := testSetup()

	setupTestAccount(router)

	loginModel := &viewmodels.LoginModel{
		Username: "shuvo009@yahoo.com",
		Password: "1234568",
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(loginModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/talent/login", buf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	loginResponse := LoginResponse{}
	body, _ := ioutil.ReadAll(w.Body)

	json.Unmarshal([]byte(string(body)), &loginResponse)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "Username or password is invalid", loginResponse.Error)
	mongoServer.Stop()
}
