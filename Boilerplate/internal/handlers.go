package internal

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/Limpid-LLC/saiService"
	"go.uber.org/zap"
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

				linkIface, ok := m["link"]
				if !ok {
					is.Logger.Error("handlers - getPDF - wrong request - lin", zap.Error(fmt.Errorf("Wrong type of data  : %+v\n", reflect.TypeOf(data))))
					return nil, 400, fmt.Errorf("bad request")
				}

				link, ok := linkIface.(string)
				if !ok {
					is.Logger.Error("handlers - getPDF - wrong request - lin", zap.Error(fmt.Errorf("Wrong type of data  : %+v\n", reflect.TypeOf(data))))
					return nil, 400, fmt.Errorf("bad request")
				}
				pdfLink, err := is.convert(m["library"].(string), link)
				if err != nil {
					is.Logger.Error("handlers - getPDF - convert", zap.Error(err))
					return nil, 400, err
				}
				return pdfLink, 200, nil

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
