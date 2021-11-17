package apiserver

import (
	"calendar/configs"
	"calendar/internal/app/model"
	redisConfig "calendar/internal/app/redis"
	"calendar/internal/app/storage"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type APIserver struct {
	logger *logrus.Logger
	router *mux.Router
	storage *storage.Storage
	redisStorage *redisConfig.RedisStorage
}

func New() *APIserver {
	return &APIserver{
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIserver) Start() error {
	err := s.configureStorage()
	if err != nil {
		s.logger.Fatal(err)
	}
	configuredServer := *s.ConfigureServer()

	if err := configuredServer.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Fatal(err)
	}

	return nil
}

func (s *APIserver) configureStorage() error {
	config := configs.NewConfig()
	store := storage.New(config.Storage)
	redisStorage := redisConfig.New(config.RedisConfig)

	log.Println("Configuring store")
	if err := store.Open(); err != nil {
		return err
	}
	s.storage = store
	s.redisStorage = redisStorage

	return nil
}

// ConfigureServer is used for setting up server's routes and some settings before initialization
func (s *APIserver) ConfigureServer() *http.Server {
	s.router.Handle("/hello", s.Hello())

	s.router.Handle("/api/event/{id:[0-9]+}", auth(http.HandlerFunc(s.EventById))).Methods(http.MethodGet)
	s.router.Handle("/api/events", auth(http.HandlerFunc(s.Events))).Methods(http.MethodGet)
	s.router.Handle("/api/events", auth(http.HandlerFunc(s.CreateEvent))).Methods(http.MethodPost)
	s.router.Handle("/api/events", auth(http.HandlerFunc(s.UpdateEvent))).Methods(http.MethodPut)
	s.router.Handle("/api/events", auth(http.HandlerFunc(s.DeleteEvent))).Methods(http.MethodDelete)

	s.router.Handle("/login", http.HandlerFunc(s.Login)).Methods(http.MethodPost)
	s.router.Handle("/logout", http.HandlerFunc(s.Logout)).Methods(http.MethodPost)

	config := configs.NewConfig()
	server := &http.Server {
		Addr: config.BindAddr,
		Handler: s.router,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return server
}

func (s *APIserver) Hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "Hello")
		if err != nil {
			s.logger.Println(err)
		}
	}
}

// Response is used for preparing content of a request
func (s *APIserver) Response(w http.ResponseWriter, message string, status int) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if status == 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}

	return w
}

// EventById is used for fetching event by ID
func (s *APIserver) EventById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "ID has not appropriate value", http.StatusBadRequest)
		return
	}
	if id <= 0 {
		s.Response(w, "ID has not appropriate value. Negative ID", http.StatusBadRequest)
		return
	}
	event, err := s.storage.Event().Load(id)

	if err != nil {
		s.logger.Println(err)
		s.Response(w, "can't process the request", http.StatusInternalServerError)
		return
	}
	preparedJson, err := json.Marshal(event)
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "can't process the request", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s.Response(w, string(preparedJson), http.StatusOK)
}

func (s *APIserver) Events(w http.ResponseWriter, r *http.Request) {

	events, err := s.storage.Event().All()
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant create an event", http.StatusBadRequest)
	}

	fmt.Println("EVENTS")
	fmt.Println(events)

	result, _ := json.Marshal(events)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s.Response(w, string(result), http.StatusOK)
}

func (s *APIserver) CreateEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		s.Response(w, "Data can not be parsed", http.StatusBadRequest)
		return
	}
	event := &model.Event{}
	err = json.Unmarshal(body, &event)

	if event.Title == "" {
		event.Title = "Birthday"
	}
	//if event.DateFrom == "" {
	//	event.DateFrom = "2021-09-01"
	//}
	//if event.DateTo == "" {
	//	event.DateTo = "2021-09-01"
	//}
	//if event.TimeFrom == "" {
	//	event.TimeFrom = "8:00 PM"
	//}
	//if event.TimeTo == "" {
	//	event.TimeTo = "10:00 PM"
	//}
	defaultTimeZone := "Europe/Riga"
	if event.Timezone == "" {
		event.Timezone = defaultTimeZone
	}

	if event.Time == "" {
		tlayout := "2006-01-02 03-04 AM"
		locTimezone, _ := time.LoadLocation(defaultTimeZone)
		currentTime := time.Now().In(locTimezone).Format(tlayout)

		event.Time = currentTime
		event.Timezone = defaultTimeZone
	}

	if len(event.Notes) == 0 {
		event.Notes = []string{""}
	}
	_, err = s.storage.Event().Create(event)
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant create an event", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s.Response(w, "Successful operation", http.StatusCreated)
}

func (s *APIserver) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	timezone, _ := time.Now().Zone()
	//notes := ""//[]string{"notes", "shmots", "enots"}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant process credentials", http.StatusBadRequest)
		return
	}
	fmt.Println("EVENT ", string(body))

	event := &model.Event{}
	err = json.Unmarshal(body, &event)

	if event.Timezone == "" {
		event.Timezone = timezone
	}

	fmt.Println("EVENT ", event)

	_, err = s.storage.Event().Update(event)
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant create an event", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s.Response(w, "Successful operation", http.StatusOK)
}

func (s *APIserver) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant process credentials", http.StatusBadRequest)
		return
	}

	type DeleteEvent struct {
		Id int `json:"id"`
	}

	deleteEvent := &DeleteEvent{}
	err = json.Unmarshal(body, &deleteEvent)

	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant process credentials", http.StatusBadRequest)
		return
	}

	if int(deleteEvent.Id) <= 0 {
		s.logger.Println(err)
		s.Response(w, "cant process credentials", http.StatusBadRequest)
		return
	}

	result, err := s.storage.Event().Delete(deleteEvent.Id)
	if err != nil {
		s.logger.Println(err)
	}
	if !result {
		s.Response(w, "event can't be deleted", http.StatusInternalServerError)
		return
	}

	s.Response(w, "", http.StatusNoContent)
}

func (s *APIserver) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant process credentials", http.StatusBadRequest)
		return
	}

	credentials := &model.Credentials{}
	err = json.Unmarshal(body, &credentials)

	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant process credentials", http.StatusBadRequest)
		return
	}

	User, err := Authorize(*credentials)
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant process credentials", http.StatusInternalServerError)
		return
	}

	type AccessToken struct {
		Token string `json:"token"`
	}

	token, err := GenerateJWT(User.Username, User.Role)
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant process credentials", http.StatusInternalServerError)
		return
	}
	at := AccessToken{
		Token: token,
	}

	accessTokenResult, err := json.Marshal(at)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s.Response(w, string(accessTokenResult), http.StatusOK)
}

func (s *APIserver) Logout(w http.ResponseWriter, r *http.Request) {
	type TokenParameter struct {
		Token string `json:"token"`
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant process credentials", http.StatusBadRequest)
		return
	}

	tokenParameter := &TokenParameter{}
	err = json.Unmarshal(body, &tokenParameter)

	if err != nil {
		s.logger.Println(err)
		s.Response(w, "cant process token", http.StatusBadRequest)
		return
	}

	if err := Logout(tokenParameter.Token); err != nil {
		s.logger.Println(err)
		s.Response(w, "cant process token", http.StatusBadRequest)
		return
	}

	s.Response(w, "successfully logged out", http.StatusOK)
}

func (s *APIserver) Migrate() error {
	config := configs.NewConfig()
	log.Println(config.Storage.DatabaseURL)
	dirPath, err := os.Getwd()
	dirPath = dirPath + "/" + "migrations/"

	m, err := migrate.New(
		"file:" + dirPath,
		config.Storage.DatabaseURL,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}

	return nil
}

func (s *APIserver) CreateMigrationFiles(migrationFileName string) error {
	if len(migrationFileName) == 0 {
		return errors.New("migration name cant be empty")
	}

	prefixes := []string {
		"up",
		"down",
	}

	for _, prefix := range prefixes {
		fileHandler, err := s.createMigrationFile(migrationFileName, prefix)

		if err != nil {
			return err
		}

		contentCreationEventMigration := ""
		if migrationFileName == "event" && prefix == "up" {
			contentCreationEventMigration =
				`CREATE TABLE event
			(
				id          serial,
				title       text,
				description text,
				time        text,
				timezone    text,
				duration    integer,
				notes       text
			);`
		}
		if migrationFileName == "event" && prefix == "down" {
			contentCreationEventMigration = "DROP TABLE event;"
		}

		defer fileHandler.Close()
		if _, err := fileHandler.Write([]byte(contentCreationEventMigration)); err != nil {
			return err
		}
	}

	return nil
}

// Create file
func (s *APIserver) createMigrationFile(name string, prefix string) (*os.File, error) {
	_, e := os.Stat(name)
	if e == nil {
		return nil, os.ErrExist
	}
	rootFolder, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	dateAsString := currentTime.Format("20060102150405")
	fileName := dateAsString + "_" + "create_" + strings.ToLower(name) + "." + strings.ToLower(prefix) + ".sql"

	fullPath := rootFolder + "/migrations/" + fileName
	return os.Create(fullPath)
}