package admin

import (
	"net/http"
	"text/template"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("ui/html/admin/dashboard.html"))
	tmpl.Execute(w, nil)
}
