package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type TranslationData struct {
	Query      string
	TargetLang string
	Result     string
}

type TranslationResponse struct {
	Translation string `json:"translation"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			query := r.FormValue("query")
			targetLang := r.FormValue("target")

			apiURL := fmt.Sprintf("https://api.hy-tech.my.id/api/translate?text=%s&target=%s", query, targetLang)
			resp, err := http.Get(apiURL)
			if err != nil {
				http.Error(w, "Gagal menghubungi server terjemahan", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, "Gagal membaca respons dari server terjemahan", http.StatusInternalServerError)
				return
			}

			var translationResp TranslationResponse
			err = json.Unmarshal(body, &translationResp)
			if err != nil {
				http.Error(w, "Gagal parsing respons dari server terjemahan", http.StatusInternalServerError)
				return
			}

			data := TranslationData{
				Query:      query,
				TargetLang: targetLang,
				Result:     translationResp.Translation,
			}

			renderTemplate(w, "index.html", data)
			return
		}

		renderTemplate(w, "index.html", nil)
	})

	fmt.Println("Server berjalan pada http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
