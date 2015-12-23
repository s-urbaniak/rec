package webapp

import "net/http"

// ListenAndServe starts serving the webapp
func ListenAndServe() {
	http.Handle(
		"/bootstrap/",
		http.StripPrefix("/bootstrap/", FileServerNoReaddir("webapp/bower_components/bootstrap/dist")),
	)

	http.Handle(
		"/jquery/",
		http.StripPrefix("/jquery/", FileServerNoReaddir("webapp/bower_components/jquery/dist")),
	)

	http.Handle("/", http.FileServer(http.Dir("webapp/html")))
	http.ListenAndServe(":8080", nil)
}
