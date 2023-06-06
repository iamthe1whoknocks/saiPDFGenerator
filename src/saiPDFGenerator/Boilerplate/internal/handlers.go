package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Limpid-LLC/saiService"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

func (is *InternalService) NewHandler() saiService.Handler {
	return saiService.Handler{
		"getPDF": saiService.HandlerElement{
			Name:        "getPDF",
			Description: "get pdf from html file",
			Function: func(data interface{}) (interface{}, int, error) {
				m, ok := data.(map[string]interface{})
				if !ok {
					is.Logger.Error("handlers - getPDF - wrong type of data", zap.Error(fmt.Errorf("Wrong type of data  : %+v\n", reflect.TypeOf(data))))
					return nil, 400, fmt.Errorf("Wrong type of data  : %+v\n", reflect.TypeOf(data))
				}

				htmlIface, ok := m["html"]
				if !ok {
					is.Logger.Error("handlers - getPDF - wrong request - html", zap.Error(fmt.Errorf("Wrong type of data  : %+v\n", reflect.TypeOf(data))))
					return nil, 400, fmt.Errorf("bad request")
				}
				//encoded base64 html bytes array
				encodedHtml, ok := htmlIface.(string)
				if !ok {
					is.Logger.Error("handlers - getPDF - wrong request - html", zap.Error(fmt.Errorf("Wrong type of data  : %+v\n", reflect.TypeOf(data))))
					return nil, 400, fmt.Errorf("bad request")
				}

				//decode base64
				htmlData := make([]byte, base64.StdEncoding.DecodedLen(len(encodedHtml)))
				n, err := base64.StdEncoding.Decode(htmlData, []byte(encodedHtml))
				if err != nil {
					is.Logger.Error("handlers - getPDF - base64.StdEncoding.Decode", zap.Error(err))
					return nil, 400, err
				}

				//validate html
				_, err = html.Parse(bytes.NewReader(htmlData[:n]))
				if err != nil {
					is.Logger.Error("handlers - getPDF - html.Parse", zap.Error(err))
					return nil, 400, err
				}

				var library string
				val, ok := m["library"]
				if ok {
					library, ok = val.(string)
					if !ok {
						library = ""
					}
				} else {
					library = ""
				}

				result, err := is.convert(library, htmlData[:n])
				if err != nil {
					// means that it was s3 error, send html to output
					// todo: error typing
					if strings.Contains(err.Error(), "s3Upload") {
						return result.([]byte), 200, nil
					}
					is.Logger.Error("handlers - getPDF - convert", zap.Error(err))
					return nil, 400, err
				}
				return result.(string), 200, nil

			},
		},
	}
}

func (is *InternalService) get(data interface{}) (string, int, error) {

	return "Get:" + strconv.Itoa(is.Context.GetConfig("common.http.port", 80).(int)), 200, nil
}

type convertRequest struct {
	Link    string `json:"link"`    //link with html
	Library string `json:"library"` //library to use
}
