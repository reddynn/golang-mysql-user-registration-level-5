package routes

import (
	"log"
	"net/http"
	"reddynn/controller"
)

func LoadRouters() {
	http.HandleFunc("/", controller.Welcome)
	http.HandleFunc("/signup", controller.Signup)
	http.HandleFunc("/signin", controller.Signin)
	http.HandleFunc("/profile",controller.Profile)
	http.HandleFunc("/signout",controller.Signout)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
