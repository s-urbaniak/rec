package webapp

import "net/http"

// ListenAndServe starts serving the webapp
func ListenAndServe() {
	handleBowerDist := func(name string) {
		h := http.StripPrefix(
			"/"+name+"/",
			FileServerNoReaddir(http.Dir("webapp/bower_components/"+name+"/dist")),
		)

		http.Handle("/"+name+"/", h)
	}

	// bower components
	handleBowerDist("bootstrap")
	handleBowerDist("jquery")
	handleBowerDist("bacon")

	// local assets
	http.Handle("/js/", FileServerNoReaddir(http.Dir("webapp")))
	http.Handle("/", FileServerNoReaddir(http.Dir("webapp/html")))

	http.ListenAndServe(":8080", nil)
}
