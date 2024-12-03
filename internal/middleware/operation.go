package middleware

import (
	"LinuxOnM/docs"
	"LinuxOnM/internal/api/services"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/copier"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// OperationLog is a middleware function that logs operation details.
// It skips logging for GET requests or requests with "search" in the URL path.
// For other requests, it extracts relevant information from the request and response,
// formats it, and saves it as an operation log in the database.
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for GET requests or requests with "search" in the URL path
		if strings.Contains(c.Request.URL.Path, "search") || c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		// Load log information from the request path
		source := loadLogInfo(c.Request.URL.Path)
		// Initialize an OperationLog struct with basic information
		record := models.OperationLog{
			Source:    source,
			IP:        c.ClientIP(),
			Method:    strings.ToLower(c.Request.Method),
			Path:      strings.ReplaceAll(c.Request.URL.Path, "/api/handler", ""),
			UserAgent: c.Request.UserAgent(),
		}
		var (
			// Variable to hold the unmarshalled Swagger JSON data
			swagger swaggerJson
			// Variable to hold the operation details from Swagger JSON's x-panel-log section
			operationDic operationJson
		)
		// Unmarshal the Swagger JSON data
		if err := json.Unmarshal(docs.SwaggerJson, &swagger); err != nil {
			c.Next()
			return
		}
		// Get the path information from the Swagger JSON data based on the request path
		path, hasPath := swagger.Paths[record.Path]
		if !hasPath {
			c.Next()
			return
		}
		// Assert the type of the path information to a map[string]interface{}
		methodMap, isMethodMap := path.(map[string]interface{})
		if !isMethodMap {
			c.Next()
			return
		}
		// Check if the POST method data exists in the method map
		dataMap, hasPost := methodMap["post"]
		if !hasPost {
			c.Next()
			return
		}
		// Assert the type of the POST method data to a map[string]interface{}
		data, isDataMap := dataMap.(map[string]interface{})
		if !isDataMap {
			c.Next()
			return
		}
		// Check if the x-panel-log data exists in the POST method data
		xlog, hasXlog := data["x-panel-log"]
		if !hasXlog {
			c.Next()
			return
		}
		// Copy the x-panel-log data to the operationDic struct
		if err := copier.Copy(&operationDic, xlog); err != nil {
			c.Next()
			return
		}

		// If the Chinese format string in operationDic is empty, skip logging
		if len(operationDic.FormatZH) == 0 {
			c.Next()
			return
		}

		// Create a map to hold the formatted data
		formatMap := make(map[string]interface{})
		// If there are body keys specified in operationDic, process the request body
		if len(operationDic.BodyKeys) != 0 {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
			bodyMap := make(map[string]interface{})
			_ = json.Unmarshal(body, &bodyMap)
			for _, key := range operationDic.BodyKeys {
				if _, ok := bodyMap[key]; ok {
					formatMap[key] = bodyMap[key]
				}
			}
		}
		// If there are before functions specified in operationDic, process them
		if len(operationDic.BeforeFunctions) != 0 {
			for _, funcs := range operationDic.BeforeFunctions {
				for key, value := range formatMap {
					if funcs.InputValue == key {
						var names []string
						if funcs.IsList {
							query := fmt.Sprintf("SELECT %s FROM %s WHERE %s in (?)", funcs.OutputColumn, funcs.DB, funcs.InputColumn)
							_ = global.DB.Raw(query, value).Scan(&names)
						} else {
							query := fmt.Sprintf("SELECT %s FROM %s WHERE %s =?", funcs.OutputColumn, funcs.DB, funcs.InputColumn)
							_ = global.DB.Raw(query, value).Scan(&names)
						}
						formatMap[funcs.OutputValue] = strings.Join(names, ",")
						break
					}
				}
			}
		}
		// Replace the placeholders in the format strings with the actual values from formatMap
		for key, value := range formatMap {
			if strings.Contains(operationDic.FormatEN, "["+key+"]") {
				t := reflect.TypeOf(value)
				if t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
					operationDic.FormatZH = strings.ReplaceAll(operationDic.FormatZH, "["+key+"]", fmt.Sprintf("[%v]", value))
					operationDic.FormatEN = strings.ReplaceAll(operationDic.FormatEN, "["+key+"]", fmt.Sprintf("[%v]", value))
				} else {
					val := reflect.ValueOf(value)
					length := val.Len()

					var elements []string
					for i := 0; i < length; i++ {
						element := val.Index(i).Interface().(string)
						elements = append(elements, element)
					}
					operationDic.FormatZH = strings.ReplaceAll(operationDic.FormatZH, "["+key+"]", fmt.Sprintf("[%v]", strings.Join(elements, ",")))
					operationDic.FormatEN = strings.ReplaceAll(operationDic.FormatEN, "["+key+"]", fmt.Sprintf("[%v]", strings.Join(elements, ",")))
				}
			}
		}
		// Set the English and Chinese detail strings in the record after formatting
		record.DetailEN = strings.ReplaceAll(operationDic.FormatEN, "[]", "")
		record.DetailZH = strings.ReplaceAll(operationDic.FormatZH, "[]", "")

		// Create a custom response body writer to capture the response body
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		datas := writer.body.Bytes()
		// If the content encoding of the request is gzip, decompress the response body
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			buf := bytes.NewReader(writer.body.Bytes())
			reader, err := gzip.NewReader(buf)
			if err != nil {
				record.Status = constant.StatusFailed
				record.Message = fmt.Sprintf("gzip new reader failed, err: %v", err)
				latency := time.Since(now)
				record.Latency = latency

				if err := services.NewILogService().CreateOperationLog(record); err != nil {
					global.LOG.Errorf("create operation record failed, err: %v", err)
				}
				return
			}
			defer reader.Close()
			datas, _ = io.ReadAll(reader)
		}
		var res response
		_ = json.Unmarshal(datas, &res)
		// Set the status and message in the record based on the response code
		if res.Code == 200 {
			record.Status = constant.StatusSuccess
		} else {
			record.Status = constant.StatusFailed
			record.Message = res.Message
		}

		latency := time.Since(now)
		record.Latency = latency

		// Save the operation log record to the database
		if err := services.NewILogService().CreateOperationLog(record); err != nil {
			global.LOG.Errorf("create operation record failed, err: %v", err)
		}
	}
}

// swaggerJson represents the structure of the Swagger JSON data related to paths.
type swaggerJson struct {
	Paths map[string]interface{} `json:"paths"`
}

// operationJson represents the structure of the operation details from Swagger JSON's x-panel-log section.
type operationJson struct {
	API             string         `json:"api"`
	Method          string         `json:"method"`
	BodyKeys        []string       `json:"bodyKeys"`
	ParamKeys       []string       `json:"paramKeys"`
	BeforeFunctions []functionInfo `json:"beforeFunctions"`
	FormatZH        string         `json:"FormatZH"`
	FormatEN        string         `json:"formatEN"`
}

// functionInfo represents the information about a before function used for formatting operation details.
type functionInfo struct {
	InputColumn  string `json:"input_column"`
	InputValue   string `json:"input_value"`
	IsList       bool   `json:"isList"`
	DB           string `json:"db"`
	OutputColumn string `json:"output_column"`
	OutputValue  string `json:"output_value"`
}

// response represents the structure of the response data.
type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// responseBodyWriter is a custom writer that captures the response body.
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write method of responseBodyWriter captures the written data to the buffer and then writes it to the original writer.
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// loadLogInfo extracts the log source information from the request path.
func loadLogInfo(path string) string {
	path = strings.ReplaceAll(path, "/api/handler", "")
	if !strings.Contains(path, "/") {
		return ""
	}
	pathArrays := strings.Split(path, "/")
	if len(pathArrays) < 2 {
		return ""
	}
	return pathArrays[1]
}
