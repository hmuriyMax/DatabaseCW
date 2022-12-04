package httpservice

import (
	"CW/pkg/sqlservice"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strings"
	"time"
)

const HTMLPath = "./web/"

func (s *HTTPService) addHandlers() {
	fileServer := http.FileServer(http.Dir("./web"))
	s.mux.PathPrefix("/res/").Handler(http.StripPrefix("/res/", fileServer))
	s.mux.HandleFunc("/", s.indexHandler)
	s.mux.HandleFunc("/table/{tableName}", s.tableHandler)
	s.mux.HandleFunc("/insert/{tableName}", s.insertHandler)
	s.mux.HandleFunc("/update/{tableName}", s.updateHandler)
	s.mux.HandleFunc("/delete/{tableName}", s.deleteHandler)
	s.mux.HandleFunc("/reports/{reportName}", s.reportHandler)
}

func (s *HTTPService) indexHandler(writer http.ResponseWriter, request *http.Request) {
	indexPath := HTMLPath + "index.html"
	pageTemplate := template.Must(template.ParseFiles(indexPath))

	ctx, cancelFunc := context.WithTimeout(request.Context(), 1*time.Second)
	defer cancelFunc()

	tableList, err := s.db.GetTableList(ctx)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["tableList"] = tableList
	err = pageTemplate.Execute(writer, data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *HTTPService) tableHandler(writer http.ResponseWriter, request *http.Request) {
	indexPath := HTMLPath + "table.html"
	pageTemplate := template.Must(template.ParseFiles(indexPath))
	tableName := mux.Vars(request)["tableName"]

	if request.FormValue("showid") == "1" {
		http.SetCookie(writer, &http.Cookie{
			Name:   "showid",
			Value:  request.FormValue("showid"),
			MaxAge: int(10 * time.Hour.Seconds()),
		})
		http.Redirect(writer, request, request.URL.Path, http.StatusFound)
		return
	}
	if request.FormValue("showid") == "-1" {
		http.SetCookie(writer, &http.Cookie{
			Name:   "showid",
			Value:  request.FormValue("showid"),
			MaxAge: -1,
		})
		http.Redirect(writer, request, request.URL.Path, http.StatusFound)
		return
	}
	data := make(map[string]interface{})
	defer func() {
		err := pageTemplate.Execute(writer, data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}()

	tableList, err := s.db.GetTableList(request.Context())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	data["tableList"] = tableList

	if tableName != "" {
		ctx, cancelFunc := context.WithTimeout(request.Context(), 1*time.Second)
		defer cancelFunc()

		table, err := s.db.SelectAll(ctx, tableName)
		if err != nil {
			data["Error"] = template.HTML(
				fmt.Sprintf("Не удалось открыть <span class=\"mono\">%s</span>:", tableName))
			data["SubError"] = err.Error()
			return
		}
		data["Table"] = *table
		inErr := request.FormValue("inError")
		data["inError"] = inErr

		_, err = request.Cookie("showid")
		data["showID"] = err != http.ErrNoCookie
		data["Controls"] = true
		data["editingID"] = strings.Split(strings.Trim(request.FormValue("row"), "[] "), " ")[0]
		if !s.IDEditable {
			data["IDro"] = 1
		} else {
			data["IDro"] = 0
		}
	} else {
		data["Error"] = "Таблица не выбрана"
	}
}

func (s *HTTPService) insertHandler(writer http.ResponseWriter, request *http.Request) {
	ctx, cancelFunc := context.WithTimeout(request.Context(), 1*time.Second)
	defer cancelFunc()
	tableName := mux.Vars(request)["tableName"]

	var data sqlservice.Table
	data.Name = tableName
	data.Data = make([][]string, 1)
	for key, val := range request.URL.Query() {
		var col sqlservice.Column
		err := json.Unmarshal([]byte(key), &col)
		if err != nil {
			url := fmt.Sprintf("/table/%s?inError=%s", tableName, err.Error())
			http.Redirect(writer, request, url, http.StatusFound)
			return
		}
		data.Columns.Values = append(data.Columns.Values, col)
		data.Data[0] = append(data.Data[0], val[0])
	}
	err := s.db.Insert(ctx, data)
	if err != nil {
		url := fmt.Sprintf("/table/%s?inError=%s", tableName, err.Error())
		http.Redirect(writer, request, url, http.StatusFound)
		return
	}
	http.Redirect(writer, request, "/table/"+tableName, http.StatusFound)
}

func (s *HTTPService) deleteHandler(writer http.ResponseWriter, request *http.Request) {
	ctx, cancelFunc := context.WithTimeout(request.Context(), 1*time.Second)
	defer cancelFunc()
	tableName := mux.Vars(request)["tableName"]

	var data sqlservice.Table
	data.Name = tableName
	data.Data = make([][]string, 1)
	row := strings.Split(strings.Trim(request.URL.Query().Get("row"), "[] "), " ")
	cols := strings.Split(request.URL.Query().Get("cols"), ", ")

	cols = strings.Split(cols[0], " ")
	data.Columns.Values = make([]sqlservice.Column, len(cols))
	for i, col := range cols {
		data.Columns.Values[i] = sqlservice.Column{
			Name: strings.Trim(col, " []{}"),
		}
	}
	data.Data[0] = row
	err := s.db.Delete(ctx, data)
	if err != nil {
		s.logger.Println(err)
		url := fmt.Sprintf("/table/%s?inError=%s", tableName, err.Error())
		http.Redirect(writer, request, url, http.StatusFound)
		return
	}
	http.Redirect(writer, request, "/table/"+tableName, http.StatusFound)
}

func (s *HTTPService) updateHandler(writer http.ResponseWriter, request *http.Request) {
	_, cancelFunc := context.WithTimeout(request.Context(), 1*time.Second)
	defer cancelFunc()
	tableName := mux.Vars(request)["tableName"]

	var data sqlservice.Table
	data.Name = tableName
	data.Data = make([][]string, 1)
	id := request.URL.Query().Get("id")
	data.Columns.Values = []sqlservice.Column{{
		Name: request.URL.Query().Get("idCol"),
		Type: "integer",
	}}
	data.Columns.IDColumn = data.Columns.Values[0].Name
	for key, val := range request.URL.Query() {
		if key == "idCol" || key == "id" {
			continue
		}
		var col sqlservice.Column
		err := json.Unmarshal([]byte(key), &col)
		if err != nil {
			url := fmt.Sprintf("/table/%s?inError=%s", tableName, err.Error())
			http.Redirect(writer, request, url, http.StatusFound)
			return
		}
		if col.Name != data.Columns.IDColumn {
			data.Columns.Values = append(data.Columns.Values, col)
		}
		data.Data[0] = append(data.Data[0], val[0])
	}

	err := s.db.Update(request.Context(), data, id)
	if err != nil {
		s.logger.Println(err)
		url := fmt.Sprintf("/table/%s?inError=%s", tableName, err)
		http.Redirect(writer, request, url, http.StatusFound)
	}
	url := fmt.Sprintf("/table/%s", tableName)
	http.Redirect(writer, request, url, http.StatusFound)
}

func (s *HTTPService) reportHandler(writer http.ResponseWriter, request *http.Request) {
	_, cancelFunc := context.WithTimeout(request.Context(), 1*time.Hour)
	defer cancelFunc()
	indexPath := HTMLPath + "table.html"
	pageTemplate := template.Must(template.ParseFiles(indexPath))
	reportName := mux.Vars(request)["reportName"]

	data := make(map[string]interface{})
	defer func() {
		err := pageTemplate.Execute(writer, data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}()
	table, err := s.db.Report(request.Context(), sqlservice.ReportType(reportName), request.URL.Query())
	if err != nil {
		data["Error"] = template.HTML(
			fmt.Sprintf("Не удалось открыть <span class=\"mono\">%s</span>:", table.Name))
		data["SubError"] = err.Error()
		return
	}
	data["Table"] = *table
	inErr := request.FormValue("inError")
	data["inError"] = inErr
	data["Controls"] = false
	data["showID"] = "1"
}
