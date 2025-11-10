package utils

import (
	"context"
	crypto "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/G4GANCHAUDHARY/encryption-service/global"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	LETTERBYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func CreateServerContext(rq *http.Request) context.Context {
	requestId := rq.Header.Get("X-Request-ID")
	url := rq.URL.Path
	ctx := context.Background()
	if requestId == "" {
		requestId = generateUuidWithoutFail()
	}
	ctx = context.WithValue(ctx, "rid", requestId)
	ctx = context.WithValue(ctx, "source", url)
	return ctx
}

func generateUuidWithoutFail() string {
	uuid, err := GenerateUuidV4()
	if err != nil {
		uuid = GenerateUnsafeUuid()
	}
	return uuid
}

func GenerateUuidV4() (string, error) {
	b := make([]byte, 16)
	_, err := crypto.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return strings.ToLower(uuid), nil
}

func GenerateUnsafeUuid() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 16)
	for i := range b {
		b[i] = LETTERBYTES[rand.Intn(len(LETTERBYTES))]
	}
	return string(b)
}

type HandlePanicRequest struct {
	Tx         *gorm.DB
	FuncName   string
	Err        *error
	CallbackFn func()
}

// HandlePanic accepts all optional parameters and does
// log the panic
// rollback
// sets err value to default error
// and executes the callback function
func HandlePanic(rq HandlePanicRequest) {
	if rec := recover(); rec != nil {
		if rq.Tx != nil {
			rq.Tx.Rollback()
		}
		if rq.Err != nil {
			*rq.Err = getErrFromRecover(rec)
		}
		if rq.CallbackFn != nil {
			rq.CallbackFn()
		}
	}
}

func getErrFromRecover(rec interface{}) error {
	switch rec.(type) {
	case string:
		return errors.New(rec.(string))
	default:
		return global.Panic
	}
}

type IPagination interface {
	GetOffSet() int
	GetLimit() int
}

type Pagination struct {
	OffSet int
	Limit  int
}

func (p Pagination) GetOffSet() int {
	return p.OffSet
}

func (p Pagination) GetLimit() int {
	return p.Limit
}

type ISortObject interface {
	GetSortBy() string
	GetSortOrder() string
	GetQueryString() string
}

type SortObject struct {
	SortBy    string
	SortOrder string
}

func (s SortObject) GetSortBy() string {
	return s.SortBy
}

func (s SortObject) GetSortOrder() string {
	return s.SortOrder
}

func (s SortObject) GetQueryString() string {
	return fmt.Sprintf("%s %s", s.SortBy, s.SortOrder)
}

func HTTPFailWithXxx(errorMessage string, httpStatusCode int, rw http.ResponseWriter) {
	responseJson, err := json.Marshal(struct {
		Success      bool   `json:"success"`
		ErrorMessage string `json:"error_message"`
	}{
		false,
		errorMessage,
	})
	if err != nil {
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(httpStatusCode)
	rw.Write(responseJson)
}

func HTTPSuccessWith200(data interface{}, rw http.ResponseWriter) {
	responseJson, err := json.Marshal(struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}{
		true,
		data,
	})
	if err != nil {
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(responseJson)
}

func GetSortParams(rq *http.Request) (string, string, error) {
	vars := rq.URL.Query()
	orderBy := vars.Get("order_by")
	if strings.TrimSpace(orderBy) == "" {
		orderBy = "created_at"
	}
	orderType := vars.Get("order_type")
	if strings.TrimSpace(orderType) == "" {
		orderType = "desc"
	}
	return orderBy, orderType, nil
}
