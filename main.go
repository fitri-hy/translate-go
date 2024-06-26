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

const indexHTML = `
<!doctype html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-slate-100 max-w-4xl m-auto pt-28 px-4">
	<nav class="fixed inset-x-0 top-0 z-10 w-full px-4 py-1 bg-white shadow-md border-slate-500 transition duration-700 ease-out">
		<div class="flex justify-between p-4">
			<div class="flex gap-2 items-center text-[2rem] leading-[3rem] tracking-tight font-bold text-black">
				<img class="w-10 h-10" src="https://hy-tech.my.id/images/logo.png">
				<h2 class="hidden lg:flex"><span class="text-indigo-600">Trans</span>Late</h2>
			</div>
			<div class="flex items-center space-x-4 text-lg font-semibold tracking-tight">
				<a href="https://github.com/fitri-hy" class="flex gap-2 items-center px-6 py-2 text-black transition duration-700 ease-out bg-white border border-black rounded-lg hover:bg-gray-200 hover:border">
					<svg width="20" height="20" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
					  <image xlink:href="http://hytech-icons.vercel.app/icons/pro/brands/github.svg" width="24" height="24"/>
					</svg>
					Github
				</a>
				<a href="https://hy-tech.my.id/" class="flex gap-2 items-center px-6 py-2 text-white transition duration-700 ease-out bg-indigo-600 rounded-lg hover:bg-indigo-500">
					<svg width="20" height="20" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
					  <image xlink:href="http://hytech-icons.vercel.app/icons/pro/sharp-solid/globe.svg?color=ffffff" width="24" height="24"/>
					</svg>
					Website
				</a>
			</div>
		</div>
	</nav>
	<form action="/" method="post">
		<div class="flex flex-col gap-4">
			<h2 class="font-bold text-2xl">Masukan Kata</h2>
			<div class="bg-white shadow border rounded-lg p-2">
				<textarea class="w-full p-4" type="text" id="query" rows="6" name="query">{{.Query}}</textarea>
			</div>
			<div class="flex items-center justify-center gap-4 w-auto p-4">
				<select id="target" name="target" class="px-4 py-2 rounded-md bg-indigo-500 text-white shadow">
					<option value="en">English</option>
					<option value="id">Indonesian</option>
					<option value="fr">French</option>
					<option value="es">Spanish</option>
					<option value="de">German</option>
					<option value="it">Italian</option>
					<option value="pt">Portuguese</option>
					<option value="zh">Chinese</option>
					<option value="ja">Japanese</option>
					<option value="ko">Korean</option>
					<option value="ru">Russian</option>
					<option value="ar">Arabic</option>
					<option value="hi">Hindi</option>
					<option value="bn">Bengali</option>
					<option value="ms">Malay</option>
					<option value="vi">Vietnamese</option>
					<option value="th">Thai</option>
					<option value="tl">Filipino</option>
					<option value="tr">Turkish</option>
					<option value="pl">Polish</option>
					<option value="nl">Dutch</option>
					<option value="sv">Swedish</option>
					<option value="fi">Finnish</option>
					<option value="no">Norwegian</option>
					<option value="da">Danish</option>
					<option value="el">Greek</option>
					<option value="hu">Hungarian</option>
					<option value="cs">Czech</option>
					<option value="ro">Romanian</option>
					<option value="uk">Ukrainian</option>
				</select>
				<button type="submit" class="px-4 py-2 rounded-md bg-emerald-500 text-white shadow">Terjemahkan</button>
			</div>
			<h2 class="font-bold text-2xl">Hasil Terjemahan</h2>
			<div class="bg-white shadow border rounded-lg p-2">
				<textarea class="w-full p-4" rows="6" readonly>{{.Result}}</textarea>
			</div>
		</div>
	</form>
	<footer class="bg-white mt-20 mb-4 rounded-lg shadow">
		<div class="container px-6 py-8 mx-auto">
			<div class="flex flex-col items-center sm:flex-row sm:justify-center">
				<p class="text-sm text-gray-500">©2024 <a href="https://hy-tech.my.id/" target="_blank" class="text-blue-500">Hy-Tech Group</a>. All Rights Reserved.</p>
			</div>
		</div>
	</footer>
</body>
</html>
`

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

			tmpl, err := template.New("index").Parse(indexHTML)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}

		tmpl, err := template.New("index").Parse(indexHTML)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server berjalan pada http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
