package transport

import (
	"encoding/json"
	"go-api/internal/todo"
	"log"
	"net/http"
)

type TodoItem struct {
	//ID   string `json:"id"`
	Item string `json:"item"`
}

type Server struct {
	mux *http.ServeMux
}

func NewServer(todoSvc *todo.Service) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		//_, err := w.Write([]byte("Hello World"))
		//if err != nil {
		//	log.Fatal(err)
		//}
		todoItems, err := todoSvc.GetAll()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(todoItems)
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("POST /todo", func(writer http.ResponseWriter, request *http.Request) {
		var t TodoItem
		err := json.NewDecoder(request.Body).Decode(&t)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		err = todoSvc.Add(t.Item)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
		writer.WriteHeader(http.StatusCreated)
		return
	})

	//mux.HandleFunc("DELETE /todo/{id}", func(writer http.ResponseWriter, request *http.Request) {
	//	if request.Method != http.MethodDelete {
	//		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
	//		return
	//	}
	//
	//	idStr := strings.TrimPrefix(request.URL.Path, "/todo/")
	//	id, err := strconv.Atoi(idStr)
	//	if err != nil {
	//		http.Error(writer, "Invalid ID format", http.StatusBadRequest)
	//		return
	//	}
	//
	//	err = todoSvc.Delete(id)
	//	if err != nil {
	//		writer.WriteHeader(http.StatusBadRequest)
	//	}
	//	writer.WriteHeader(http.StatusOK)
	//	return
	//})

	mux.HandleFunc("GET /search", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query().Get("q")
		if query == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		results, err := todoSvc.Search(query)
		if err != nil {
			log.Println(err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(results)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
	})

	return &Server{
		mux: mux,
	}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(":8080", s.mux)
}
