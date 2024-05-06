package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

// La función main es el punto de entrada del programa
func main() {
	fmt.Println("Server: Iniciando...")

	// Configurar el servidor para servir archivos estáticos desde la carpeta "./static"
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Configurar manejadores de rutas para diferentes URLs
	http.HandleFunc("/", home)
	http.HandleFunc("/info", info)
	http.HandleFunc("/producto", producto)

	http.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/producto", 301)
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fatal error", 501)
	})

	http.HandleFunc("/cabeceras", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("test", "test1")

		w.Header().Set("Content-Type", "application/json; charset=utf-8") 
		fmt.Fprintln(w, "{ \"hola\": 1 }")
	})

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	http.HandleFunc("/template",func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, struct{ Saludo string}{ "Hola mundo!"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Server: On")
	http.ListenAndServe(":8080", nil)
}

// Maneja la ruta principal
func home(w http.ResponseWriter, r *http.Request) {
	html := "<html>"
	html += "<body>"
	html += "<h1>Hola Mundo<h1>"
	html += "<body>"
	html += "<html>"
	w.Write([]byte(html))
}

// Maneja la ruta /info y muestra información sobre la solicitud HTTP
func info(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Host: ", req.Host)
	fmt.Fprintln(w, "URI: ", req.RequestURI)
	fmt.Fprintln(w, "Method: ", req.Method)
	fmt.Fprintln(w, "RemoteAddr: ", req.RemoteAddr)
}

// Maneja la ruta /producto y muestra información sobre los productos
var productos []string

func producto(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	add, okForm := r.Form["add"]

	// Si se envió el parámetro "add" en el formulario, añadir el producto a la lista
	if okForm && len(add) == 1 {
		productos = append(productos, string(add[0]))
		w.Write([]byte("Producto añadido correctamente"))
		return
	}

	prod, ok := r.URL.Query()["prod"]

	// Si se proporcionó el parámetro "prod" en la URL, mostrar información sobre el producto en esa posición
	if ok && len(prod) == 1 {
		pos, err := strconv.Atoi(prod[0])
		if err != nil {
			// Error al convertir a entero, no hacer nada
			return
		}

		html := "<html>"
		html += "<body>"
		html += "<h1> Producto " + productos[pos] + "<h1>"
		html += "<body>"
		html += "<html>"

		w.Write([]byte(html))
		return
	}

	// Mostrar información sobre la cantidad total de productos
	html := "<html>"
	html += "<body>"
	html += "<h1> Total productos " + strconv.Itoa(len(productos)) + "<h1>"
	html += "<body>"
	html += "<html>"

	w.Write([]byte(html))
}