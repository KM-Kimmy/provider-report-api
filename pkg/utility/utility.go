package utility

import (
	config "provider-report-api/configs"
	"provider-report-api/constant"
	sharedRequest "provider-report-api/internal/modules/shared/dtos/requests"
	"crypto/tls"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type SortInput struct {
	SortField string
	SortOrder string
}

type QueryConfig struct {
	Value     interface{}
	Condition string
}

// SetDefault checks if a pointer is nil and sets it to the default value if it is.
func SetDefaultPtr[T any](p *T, defaultValue T) *T {
	if p == nil {
		return &defaultValue
	}
	return p
}

func CheckDuoRequire[T, S any](data1 *T, data2 *S) bool {
	if data1 == nil && data2 == nil {
		return true
	} else if data1 != nil && data2 != nil {
		return true
	} else {
		return false
	}
}

// ConvertSortInput converts a string input into structured sort information
func ConvertSortInput(str *string, defaultValue string) (*SortInput, error) {
	if str == nil || *str == "" {
		return &SortInput{SortField: defaultValue, SortOrder: "DESC"}, nil
	}
	// Split the string by '-'
	parts := strings.Split(*str, "-")

	// Check if the split parts meet the required conditions
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, errors.New("invalid format")
	}

	// Convert the first part to snake case upper
	sortField := strings.ToUpper(strings.NewReplacer("A", "_A", "B", "_B", "C", "_C", "D", "_D", "E", "_E", "F", "_F", "G", "_G", "H", "_H", "I", "_I", "J", "_J", "K", "_K", "L", "_L", "M", "_M", "N", "_N", "O", "_O", "P", "_P", "Q", "_Q", "R", "_R", "S", "_S", "T", "_T", "U", "_U", "V", "_V", "W", "_W", "X", "_X", "Y", "_Y", "Z", "_Z").Replace(parts[0]))

	// Convert the second part to upper case
	sortOrder := strings.ToUpper(parts[1])

	// Check if the second part is 'ASC' or 'DESC'
	if sortOrder != "ASC" && sortOrder != "DESC" {
		return nil, errors.New("sort order must be 'ASC' or 'DESC'")
	}

	return &SortInput{SortField: sortField, SortOrder: sortOrder}, nil
}

func ProcessQueryConfigs(queryConfigs []QueryConfig, whereClauses *[]string, queryParams *[]interface{}) {
	for _, config := range queryConfigs {
		// Check if the value is a pointer and is nil
		val := reflect.ValueOf(config.Value)
		if val.Kind() == reflect.Ptr && val.IsNil() {
			continue
		}

		switch v := config.Value.(type) {
		case string:
			if v == "" {
				continue
			}
			fmt.Println("Processing string:", v)
			if strings.Contains(config.Condition, "LIKE") {
				*whereClauses = append(*whereClauses, config.Condition)
				count := strings.Count(config.Condition, "?")
				for i := 0; i < count; i++ {
					*queryParams = append(*queryParams, "%"+v+"%")
				}
			} else {
				*whereClauses = append(*whereClauses, config.Condition)
				*queryParams = append(*queryParams, v)
			}
		case *string:
			if v == nil || *v == "" {
				continue
			}
			fmt.Println("Processing string:", *v)
			if strings.Contains(config.Condition, "LIKE") {
				*whereClauses = append(*whereClauses, config.Condition)
				count := strings.Count(config.Condition, "?")
				for i := 0; i < count; i++ {
					*queryParams = append(*queryParams, "%"+*v+"%")
				}
			} else {
				*whereClauses = append(*whereClauses, config.Condition)
				*queryParams = append(*queryParams, *v)
			}
		case *int:
			fmt.Println("Processing integer:", *v)
			if strings.Contains(config.Condition, "LIKE") {
				strValue := strconv.Itoa(*v)
				*whereClauses = append(*whereClauses, config.Condition)
				count := strings.Count(config.Condition, "?")
				for i := 0; i < count; i++ {
					*queryParams = append(*queryParams, "%"+strValue+"%")
				}
			} else {
				*whereClauses = append(*whereClauses, config.Condition)
				*queryParams = append(*queryParams, *v)
			}
		case int:
			fmt.Println("Processing integer:", v)
			if strings.Contains(config.Condition, "LIKE") {
				strValue := strconv.Itoa(v)
				*whereClauses = append(*whereClauses, config.Condition)
				count := strings.Count(config.Condition, "?")
				for i := 0; i < count; i++ {
					*queryParams = append(*queryParams, "%"+strValue+"%")
				}
			} else {
				*whereClauses = append(*whereClauses, config.Condition)
				*queryParams = append(*queryParams, v)
			}
		case *bool:
			fmt.Println("Processing boolean:", *v)
			*whereClauses = append(*whereClauses, config.Condition)
			*queryParams = append(*queryParams, *v)
		default:
			log.Printf("Unsupported type for value: %#v", config.Value)
		}
	}
}

// ValidateIDParam checks if the parameter is present, is a valid integer, and meets the minimum value requirement.
func ValidateIDParam(c *gin.Context, paramName string, minValue int) (int, bool) {
	paramValue := c.Param(paramName)
	if paramValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": paramName + " is required"})
		return 0, false
	}

	id, err := strconv.Atoi(paramValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": paramName + " must be a valid integer"})
		return 0, false
	}

	if id < minValue {
		c.JSON(http.StatusBadRequest, gin.H{"error": paramName + " must be at least " + strconv.Itoa(minValue)})
		return 0, false
	}

	return id, true
}

func DecodeToken(tok string) (jwt.MapClaims, error) {
	// Parse the access token

	token, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
		// You need to provide the secret key used to sign the token for validation
		return []byte("YOUR_SECRET_KEY_HERE"), nil
	})

	// Check for errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid access token")
	}

	// Extract the claims from the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token claims")
	}
}

func GetUsername(c *gin.Context) (string, error) {
	claims, ok := c.Get("claims")
	if !ok {
		return "", errors.New("Unauthorized")
	}
	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Internal Server Error")
	}

	if claimsMap["username"] == nil {
		return "", errors.New("Username not found")
	}

	username := claimsMap["username"].(string)
	return username, nil
}

func ConvertISOToCustomFormat(isoDate string) (string, error) {
	// Parse the ISO string into a time.Time object
	parsedTime, err := time.Parse(time.RFC3339, isoDate)
	if err != nil {
		return "", err
	}

	// Format the time into the desired format
	formattedDate := parsedTime.Format("02:01:2006 15:04")
	return formattedDate, nil
}

func toExcelColumnName(n int) string {
	result := ""
	for n >= 0 {
		result = string(rune('A'+(n%26))) + result
		n = n/26 - 1
	}
	return result
}

func ExportExcelFromSlice(data []map[string]interface{}, headCell []string) (*excelize.File, error) {
	f := excelize.NewFile()

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet := "Sheet1"

	// Create sheet
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}

	// Set headers
	for colIndex, header := range headCell {
		cell := fmt.Sprintf("%s1", toExcelColumnName(colIndex))
		if err := f.SetCellValue(sheet, cell, header); err != nil {
			return nil, err
		}
	}

	// Fill data
	for rowIndex, row := range data {
		for colIndex, header := range headCell {
			cell := fmt.Sprintf("%s%d", toExcelColumnName(colIndex), rowIndex+2)
			value := ""
			if val, ok := row[header]; ok && val != nil {
				value = fmt.Sprintf("%v", val)

				if err := f.SetCellValue(sheet, cell, value); err != nil {
					return nil, err
				}
			}
		}
	}

	// Set active sheet
	f.SetActiveSheet(index)

	return f, nil
}

func Placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.TrimRight(strings.Repeat("?,", n), ",")
}

func DownloadLocalFile(filePath string, c *gin.Context) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open the file: %w", err)
	}
	defer file.Close()

	// Get file information
	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Verify the file has a .xlsx extension
	if filepath.Ext(fileStat.Name()) != ".xlsx" {
		return fmt.Errorf("invalid file extension: %s", filepath.Ext(fileStat.Name()))
	}

	fmt.Println("filename => ", fileStat.Name())

	// Set headers for file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileStat.Name()))

	// Stream the file to the client
	c.File(filePath)

	return nil
}

// ProcessExcelFile: Extract data from excel file
func ProcessExcelFile(fileHeader multipart.FileHeader) ([][]string, error) {
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	// Create a new Excel reader
	xlFile, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}

	// Get rows from the first sheet
	rows, err := xlFile.GetRows(xlFile.GetSheetName(0))
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from Excel sheet: %w", err)
	}

	return rows, nil
}

func GetInsurerCode(insurerId int) (string, error) {
	db := config.GetDB()
	var insurerCode string

	query := `
		SELECT INSURER_CODE 
		FROM tpacaredb.CLIENT.INSURER
		WHERE INSURER_ID = ?
	`

	err := db.QueryRow(query, insurerId).Scan(&insurerCode)
	if err != nil {
		return "", fmt.Errorf("query error: %w", err)
	}

	return insurerCode, nil
}

func FindAgent(agentCode string) (*int, error) {
	db := config.GetDB()

	query := `
		SELECT INSURER_AGENT_ID
		FROM tpacaredb.CLIENT.INSURER_AGENT
		WHERE AGENT_CODE = ?
	`

	var id int
	err := db.QueryRow(query, agentCode).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func SaveImportMonitorLog(data any) (string, error) {
	// Define UTC+7 time zone
	location := time.FixedZone("UTC+7", 7*60*60)

	// Get current time in UTC+7
	currentTime := time.Now().In(location)

	json, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	payload := strings.NewReader(string(json))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: os.Getenv("VERIFY_CERT") == "true"},
		},
	}
	req, err := http.NewRequest(http.MethodPost, os.Getenv("LOGSTASH_URL"), payload)

	if err != nil {
		log.Print("error: ", err)

		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Print("error: ", err)

		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	fmt.Println("body => ", string(body))

	if err != nil || strings.Contains(string(body), "invalid") {
		log.Println("error: ", err)

		return "", err
	}

	return currentTime.Format("2006-01-02T15:04"), nil
}

func SaveAuditLog(action string, currentUser string, data any, method, level, subModule string) error {
	auditLog, err := json.Marshal(sharedRequest.AuditLog{
		Id:        uuid.NewString(),
		Action:    action,
		Username:  currentUser,
		Data:      data,
		Timestamp: time.Now().Format("2006-01-02T15:04"),
		Level:     level,
		SubModule: subModule,
		Service:   constant.CLIENT_MGMT_INDEX,
	})

	if err != nil {
		return err
	}

	payload := strings.NewReader(string(auditLog))

	log.Printf("payload  %+v : ", payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, os.Getenv("LOGSTASH_URL"), payload)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	fmt.Println("body => ", string(body))

	if err != nil || strings.Contains(string(body), "invalid") {
		return err
	}

	return nil
}
func CompareStructs(a, b interface{}) map[string][2]interface{} {
	differences := make(map[string][2]interface{})

	valA := reflect.ValueOf(a)
	valB := reflect.ValueOf(b)
	typeA := reflect.TypeOf(a)

	if typeA != reflect.TypeOf(b) {
		fmt.Println("Structs are of different types")
		return differences
	}

	if valA.Kind() == reflect.Ptr {
		valA = valA.Elem()
	}
	if valB.Kind() == reflect.Ptr {
		valB = valB.Elem()
	}
	if typeA.Kind() == reflect.Ptr {
		typeA = typeA.Elem()
	}

	for i := 0; i < valA.NumField(); i++ {
		field := typeA.Field(i)
		label := field.Tag.Get("label")
		if label == "" {
			label = field.Name
		}

		valueA := valA.Field(i).Interface()
		valueB := valB.Field(i).Interface()

		differences[label] = [2]interface{}{valueA, valueB}
	}

	return differences
}

func deref(v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		return rv.Elem().Interface()
	}
	return v
}

func stringify(v interface{}) string {
	if v == nil {
		return ""
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return ""
		}
		v = rv.Elem().Interface()
	}

	switch val := v.(type) {
	case string:
		if val == "" {
			return ""
		}
		return val
	case time.Time:
		if val.IsZero() {
			return ""
		}
		return val.Format("2006-01-02") // แสดงเฉพาะวัน
	default:
		return fmt.Sprintf("%v", v)
	}
}

func FormatAuditLogEdit(differences map[string][2]interface{}, isAdd bool) string {
	if len(differences) == 0 {
		return ""
	}

	var sb strings.Builder
	i := 1

	keys := make([]string, 0, len(differences))
	for k := range differences {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, field := range keys {
		diff := differences[field]
		from := stringify(deref(diff[0]))
		to := stringify(deref(diff[1]))

		if isAdd || from == "" {
			sb.WriteString(fmt.Sprintf("%d) value added to %s - value : %s\n", i, strings.ToUpper(field), to))
		} else {
			sb.WriteString(fmt.Sprintf("%d) %s - From: %s To: %s\n", i, strings.ToUpper(field), from, to))
		}
		i++
	}

	return sb.String()
}

func GetEsClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTICSEARCH_URL"),
		},
		Username: os.Getenv("ELASTICSEARCH_USERNAME"),
		Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: os.Getenv("VERIFY_CERT") == "true",
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fmt.Println("Connected to Elasticsearch successfully")
	return es, nil
}

func SplitSort(str string) (string, string, error) {
	parts := strings.Split(str, "-")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid sort format: expected two parts separated by '-', got %d parts", len(parts))
	}
	if parts[0] == "" || parts[1] == "" {
		return "", "", errors.New("invalid sort format: field name or direction is missing")
	}

	return parts[0], parts[1], nil
}

// IntPtr returns a pointer to the given integer.
func IntPtr(i int) *int {
	return &i
}

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}

// WriteLocalCSV creates a CSV file and writes the provided records to it.
func WriteLocalCSV(filePath string, records [][]string) error {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create a directory: %w", err)
	}

	// Create and write to the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create a file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record to CSV: %w", err)
		}
	}

	return nil
}

func ConvertToTimePtr(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil // return nil without error if input is empty
	}
	layout := "2/1/2006" // expected input format: DD/MM/YY

	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected DD/MM/YYYY: %w", err)
	}

	utcTime := parsedTime.UTC()
	return &utcTime, nil
}

func GetLookupDetailId(code string, lookupSetupId int) (float64, error) {
	db := config.GetDB()

	var lookupDetailId float64
	query := `
		SELECT LOOKUPDETAIL_ID
		FROM tpacaredb.MAINTAIN.LOOKUPDETAIL
		WHERE LOOKUPSETUP_ID = ? AND LOOKUP_VALUE2 = ?
	`

	err := db.QueryRow(query, lookupSetupId, code).Scan(&lookupDetailId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no record found for categoryCode: %s", code)
		}
		return 0, fmt.Errorf("query error: %w", err)
	}

	return lookupDetailId, nil
}

func FindPricing(providerCode string, categoryId int) (*int, error) {
	db := config.GetDB()

	query := `
		SELECT PROVIDER_PRICING_ID
		FROM tpacaredb.CLIENT.PROVIDER_PRICING
		WHERE PROVIDER_CODE = ? AND PRICING_CATEGORY_ID = ?
	`

	var id int
	err := db.QueryRow(query, providerCode, categoryId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &id, nil
}

func GetProviderCode(providerId int) (string, error) {
	db := config.GetDB()
	var providerCode string

	query := `
		SELECT PROVIDER_CODE 
		FROM tpacaredb.CLIENT.PROVIDER
		WHERE PROVIDER_ID = ?
	`

	err := db.QueryRow(query, providerId).Scan(&providerCode)
	if err != nil {
		return "", fmt.Errorf("query error: %w", err)
	}

	return providerCode, nil
}

func FindDiscount(providerCode string, categoryId int) (*int, error) {
	db := config.GetDB()

	query := `
		SELECT PROVIDER_DISCOUNT_ID
		FROM tpacaredb.CLIENT.PROVIDER_DISCOUNT
		WHERE PROVIDER_CODE = ? AND DISCOUNT_CATEGORY_ID = ?
	`

	var id int
	err := db.QueryRow(query, providerCode, categoryId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &id, nil
}
func ParseIDs(ids []string) ([]int, error) {
	var parsedNumbers []int
	for _, WantedStr := range ids {
		fmt.Println(WantedStr)
		// Split the string by commas
		numbers := strings.Split(WantedStr, ",")

		// Parse each substring to an integer

		for _, numStr := range numbers {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				fmt.Printf("Error parsing %s: %v\n", numStr, err)
				// Handle parsing error
			}
			parsedNumbers = append(parsedNumbers, num)
		}
	}
	return parsedNumbers, nil
}
