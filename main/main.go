package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type PostCard struct {
	Title string `redis:"title"`
	Creator string `redis:"creator"`
	MembershipFee float32 `redis:"membership_fee"`
}

var pool *redis.Pool

//Initializes the connection pool when the program starts
func init() {

	pool = &redis.Pool{
		MaxIdle: 8, // maximum number of idle links
		MaxActive: 0, // indicates the maximum number of links to the database, 0 indicates no limit
		IdleTimeout: 100, // maximum idle time
		Dial: func() (redis. Conn, error) {// initializes the linked code, which IP redis is linked to
			return redis.Dial("tcp", "localhost:6379")
		},
	}

}

func main() {
	conn := pool.Get()

	defer conn.Close()

	_, err := conn.Do(
		"HMSet",
		"podcard",
		"title",
		"tech over tea",
		"creator",
		"kis lupin",
		"membership_fee",
		9.99,
	)
	checkError(err)
	title, err := redis.String(conn.Do("HGet", "podcard", "title"))
	checkError(err)
	fmt.Println("Padcast title", title)
	fee, err := redis.Float64(conn.Do("HGet", "podcard", "membership_fee"))
	checkError(err)
	fmt.Println("Padcast fee", fee)
	value, err := redis.Values(conn.Do("HGetAll", "podcard"))

	var p PostCard
	err = redis.ScanStruct(value, &p)
	checkError(err)
	fmt.Printf("PostCard value: %+v\n", p)

	data := PostCard{
		Title: "test",
		Creator: "lupin",
		MembershipFee: 3.33,
	}
	_, err = conn.Do(
		"HSet",
		redis.Args{}.Add(data.Title).AddFlat(data)...
	)

	checkError(err)
	fmt.Printf("PostCard test value: %+v\n", p)
}
