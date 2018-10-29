// Graeme Ferguson | ggf221 | 09/25/18
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Review struct holds entry data within website
type Review struct {
	ID     int
	Album  string
	Artist string
	Rating int
	Body   string
}

var reviewMap = map[int]*Review{}
var idCount = 0

/* Http request handler for "/" page. Displays items and links for
creation, editing, and deletion. */
func handler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("main.html"))
	t.Execute(w, reviewMap)
}

/* Http request handler for "/save/" pages. Receives form value from either
/edit/ or /create/ and saves item info accordingly (creating a new one
if necessary) */
func saveHandler(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.URL.Path[len("/save/"):])
	Album := r.PostFormValue("albumName")
	Artist := r.PostFormValue("artistName")
	Rating, _ := strconv.Atoi(r.PostFormValue("rating"))
	Body := r.PostFormValue("bodyText")
	if ID == 0 {
		createReview(Album, Artist, Rating, Body)
	} else if review, exist := reviewMap[ID]; exist {
		review.Album = Album
		review.Artist = Artist
		review.Rating = Rating
		review.Body = Body
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

/* Http request handler for "/create/" pages. Allows users to enter form
values and then redirects form submissions to "/save/" */
func createHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./create.html")
}

/* Http request handler for "/edit/" pages. Similar to "/create/" with
the caveat that the forms have values filled in from the item to be
edited. Upon form submission, redirects to "/save/" */
func editHandler(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.URL.Path[len("/edit/"):])
	review, exist := reviewMap[ID]
	if !exist {
		http.Redirect(w, r, "/create/", http.StatusFound)
		return
	}
	t := template.Must(template.ParseFiles("edit.html"))
	t.Execute(w, review)
}

/* Http request handler for "/delete/" pages. Receives id of item and
deletes info from map and server if it exists */
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.URL.Path[len("/delete/"):])
	_, exist := reviewMap[ID]
	if exist {
		delete(reviewMap, ID)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

/* Function for creating a new item and instantiating it with correct values
and a unique id. */
func createReview(Album string, Artist string, Rating int, Body string) {
	idCount = idCount + 1
	reviewMap[idCount] = &Review{ID: idCount, Album: Album,
		Artist: Artist, Rating: Rating, Body: Body}
}

func main() {
	portPtr := flag.Int("listen", 8080, "Specify port"+
		" for web server to listen on.")
	flag.Parse()
	portStr := ":" + strconv.Itoa(*portPtr)

	createReview("Grace", "Jeff Buckley", 9, "David Bowie's favorite album!")
	createReview("Exmilitary", "Death Grips", 10, "Zach Hill is good drummer!")
	createReview("Q: Are We Not Men?", "DEVO", 8, "Those are some funny hats!")

	http.HandleFunc("/", handler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/delete/", deleteHandler)

	fmt.Println("Server spun up!")
	log.Fatal(http.ListenAndServe(portStr, nil))
}
