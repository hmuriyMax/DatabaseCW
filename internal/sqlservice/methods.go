package sqlservice

import (
	"context"
	"fmt"
	sqr "github.com/Masterminds/squirrel"
	"github.com/hmuriyMax/DatabaseCW/internal/sqlservice/entities"
	"html"
	"strconv"

	"sort"
	"strings"
)

type Row = []string

type ReportType string

const (
	SelectRequestsByWorker ReportType = "requests_by_worker"
	SelectUndoneRequests   ReportType = "undone_requests"
	SelectOverdueRequests  ReportType = "overdue_requests"
	GetRequestNumber       ReportType = "request_number"
)

var taskMap = map[ReportType]string{
	SelectUndoneRequests:   "отчет о невыполненных заявках по исполнителю %s",
	SelectRequestsByWorker: "контроль исполнения заявок по исполнителям",
	SelectOverdueRequests:  "отчет о заявках, выполненных с превышением срока",
	GetRequestNumber:       "отчет о количестве заявок типа \"%s\"",
}

type SortFlag int8

const (
	SortByID SortFlag = 1
	SortNone SortFlag = 2
)

type Column struct {
	Name string `db:"column_name"`
	Type string `db:"data_type"`
}

type Header struct {
	Values   []Column
	IDColumn string
}

func (c Header) getIndex(name string) (int, bool) {
	for i, col := range c.Values {
		if col.Name == name {
			return i, true
		}
	}
	return -1, false
}

func (c Header) String() string {
	var sb strings.Builder
	for i, col := range c.Values {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(col.Name)
	}
	return sb.String()
}

type Table struct {
	Name    string
	Columns Header
	Data    []Row
	NextID  int
}

func (s *SQLService) GetTableList(ctx context.Context) (tableList Row, err error) {
	var parseStruct []entities.TableNameStr
	err = s.db.SelectContext(ctx, &parseStruct,
		"select table_name from information_schema.tables where table_schema = 'public'")
	if err != nil {
		return nil, fmt.Errorf("select tables: %v", err)
	}
	for _, el := range parseStruct {
		tableList = append(tableList, el.TableName)
	}
	return
}

func (s *SQLService) SelectAll(ctx context.Context, tableName string) (*Table, error) {
	query := sqr.Select("*").From(tableName)
	query = query.RunWith(s.db)
	tbl, err := s.selectByQuery(ctx, &query, SortByID)
	if err != nil {
		return nil, fmt.Errorf("selectAll: %v", err)
	}
	tbl.Name = tableName
	tbl.Columns.IDColumn = tbl.Columns.Values[0].Name
	tbl.NextID = 1
	if len(tbl.Data) > 0 {
		tbl.NextID, err = strconv.Atoi(tbl.Data[len(tbl.Data)-1][0])
		tbl.NextID++
	}
	return tbl, nil
}

func (s *SQLService) parse(in string) (out string) {
	out = strings.Trim(in, "[]")
	tmp := strings.Fields(out)
	out = strings.Join(tmp, ", ")
	return
}

func (s *SQLService) Insert(ctx context.Context, rows Table) error {
	rowStr := rows.Columns.String()
	for _, row := range rows.Data {
		var tmp []string
		for i, el := range row {
			if el == "" {
				el = "NULL"
			} else if rows.Columns.Values[i].Type != "integer" {
				el = fmt.Sprintf("'%s'", el)
			}
			tmp = append(tmp, fmt.Sprint(el))
		}
		valStr := strings.Join(tmp, ", ")
		query := fmt.Sprintf("insert into %s (%v) values (%v)", rows.Name, rowStr, valStr)
		_, err := s.db.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLService) Delete(ctx context.Context, rows Table) error {
	idColumn := rows.Columns.Values[0].Name
	for _, row := range rows.Data {
		id := row[0]
		query := fmt.Sprintf("delete from %s where %s = %s", rows.Name, idColumn, id)
		_, err := s.db.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLService) Update(ctx context.Context, data Table, id string) interface{} {
	setStr := ""
	for i, el := range data.Data[0] {
		if data.Columns.Values[i].Name != data.Columns.IDColumn {
			var formatStr string
			if el == "" {
				formatStr += "%s=NULL, "
			} else if data.Columns.Values[i].Type == "integer" {
				formatStr += "%s=%s, "
			} else {
				formatStr = "%s='%s', "
			}
			if el == "" {
				setStr += fmt.Sprintf(formatStr, data.Columns.Values[i].Name)
			} else {
				setStr += fmt.Sprintf(formatStr, data.Columns.Values[i].Name, el)
			}
		}
	}
	query := fmt.Sprintf("update %s set %s where %s = %s", data.Name, setStr[:len(setStr)-2],
		data.Columns.IDColumn, id)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLService) selectByQuery(ctx context.Context, query *sqr.SelectBuilder, sortFlag SortFlag) (*Table, error) {
	data := &Table{}
	sql, _, err := query.ToSql()
	if err != nil {
		return data, fmt.Errorf("select by query: %v", err)
	}
	s.lg.Printf("query: %s\n", sql)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return data, fmt.Errorf("select by query: %v", err)
	}

	columns, _ := rows.ColumnTypes()
	if len(columns) == 0 {
		return data, fmt.Errorf("select by query: error parsing columns")
	}

	for _, col := range columns {
		data.Columns.Values = append(data.Columns.Values, Column{
			Name: col.Name(),
			Type: strings.ToLower(col.DatabaseTypeName()),
		})
	}

	for rows.Next() {
		var strRow Row

		anyRow := make([][]byte, len(columns))
		anyRowPointers := make([]any, len(columns))
		for i := range anyRow {
			anyRowPointers[i] = &anyRow[i]
		}
		err = rows.Scan(anyRowPointers...)
		for i, el := range anyRow {
			strEl := string(el)
			if strings.Contains(data.Columns.Values[i].Type, "time") {
				rpl := strings.NewReplacer("T", " ", "Z", "")
				strEl = rpl.Replace(strEl)
			}
			strRow = append(strRow, strEl)
		}
		data.Data = append(data.Data, strRow)
	}
	switch sortFlag {
	case SortByID:
		sort.Slice(data.Data, func(i, j int) bool {
			a, _ := strconv.Atoi(data.Data[i][0])
			b, _ := strconv.Atoi(data.Data[j][0])
			return a < b
		})
	default:
		break
	}
	if err != nil {
		return data, fmt.Errorf("select by query: %v", err)
	}
	return data, nil
}

func (s *SQLService) Report(ctx context.Context, name ReportType, params map[string][]string) (*Table, error) {
	tbl := &Table{
		Name: fmt.Sprint(name),
	}
	query := sqr.Select().RunWith(s.db)
	var textParams []any
	switch name {
	case SelectUndoneRequests:
		query = query.Where("finish_stamp is not null")
		fallthrough
	case SelectRequestsByWorker:
		wName := params["worker_name"][0]
		textParams = append(textParams, wName)
		query = query.Columns("request_id").
			Columns("finish_stamp is null as \"done\"").
			From("repair_request").
			LeftJoin("performer using(worker_id)").
			Where(fmt.Sprintf("name = '%s'", html.EscapeString(wName))).
			OrderBy("performer.name", "request_id")
	case SelectOverdueRequests:
		query = query.Columns("request_id",
			"finish_stamp - start_stamp as \"real_duration\"",
			"duration as \"planned_duration\"").
			From("repair_request").InnerJoin("repair_type using (repair_type_id)").
			Where("finish_stamp - start_stamp > duration")
	case GetRequestNumber:
		typeId := params["type_id"][0]
		textParams = append(textParams, typeId)
		query = query.Columns("count(request_id) as \"number\"").
			From("repair_request").
			InnerJoin("repair_type using (repair_type_id)").
			Where(fmt.Sprintf("name = '%v'", html.EscapeString(typeId))).
			GroupBy("name")
	default:
		return &Table{
			Name: fmt.Sprintf("report_id: %s", name),
		}, fmt.Errorf("такого отчета не существует")
	}
	tbl, err := s.selectByQuery(ctx, &query, SortNone)
	tbl.Name = fmt.Sprintf(taskMap[name], textParams...)
	if len(tbl.Data) == 0 {
		return tbl, fmt.Errorf("нет данных по такому запросу")
	}
	if err != nil {
		return tbl, fmt.Errorf("report: %v", err)
	}
	return tbl, nil
}
