package user

import (
	"github.com/coopernurse/gorp"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"log"
	"strconv"
	"encoding/json"
)

func getUserHandler(w http.ResponseWriter, r *http.Request, dbMap *gorp.DbMap, params martini.Params) {
	id, err := strconv.Atoi(params["id"])

	user, err := dbMap.Get(SecureUser{}, id)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

func registerHandler(user User, r render.Render, dbMap *gorp.DbMap) {
	log.Println(user.Email)
	err := dbMap.Insert(&user)

	if err != nil {
		r.Error(http.StatusNotFound)
	} else {
		r.JSON(http.StatusOK, map[string]interface{}{})
	}
}

func loginHandler(credentials LoginCredentials, r render.Render, redisPool *redis.Pool, dbMap *gorp.DbMap) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	var user User

	/*
	Naive implementation
	 */
	err := dbMap.SelectOne(
		&user,
		"SELECT * FROM users as u WHERE u.email = $1 AND u.user_pass = $2",
		credentials.Email,
		credentials.Password,
	)

	if err != nil {
		panic(credentials)
		r.Error(http.StatusNotFound)
	}

	userSession := NewUserSession(&user)

	responseData, _ := json.Marshal(userSession)

	redisConnection.Do("SET", userSession.Id, responseData)

	r.JSON(http.StatusOK, userSession)
}

func logoutHandler(apiCredentials APICredentials, r render.Render, redisPool *redis.Pool) string {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	userSessionBytes, err := redis.Bytes(redisConnection.Do("GET", apiCredentials.Id))

	if err != nil {
		r.Error(http.StatusNotFound)
		return "Not Found"
	}

	var userSession UserSession
	json.Unmarshal(userSessionBytes, &userSession)

	if userSession.Token == apiCredentials.Token {
		redisConnection.Do("DEL", apiCredentials.Id)
		r.Error(http.StatusOK)
		return "Logged out"
	} else {
		r.Error(http.StatusUnauthorized)
		return "Forbidden"
	}
}
