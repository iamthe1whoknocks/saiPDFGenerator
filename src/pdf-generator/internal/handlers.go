package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/Limpid-LLC/saiService"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"reflect"
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
					return nil, 500, fmt.Errorf("Wrong type of data  : %+v\n", reflect.TypeOf(data))
				}

				htmlIface, ok := m["html"]
				if !ok {
					is.Logger.Error("handlers - getPDF - wrong request - html", zap.Error(fmt.Errorf("Wrong type of data  : %+v\n", reflect.TypeOf(data))))
					return nil, 500, fmt.Errorf("bad request")
				}

				//encoded base64 html bytes array
				encodedHtml, ok := htmlIface.(string)
				if !ok {
					is.Logger.Error("handlers - getPDF - wrong request - html", zap.Error(fmt.Errorf("Wrong type of data  : %+v\n", reflect.TypeOf(data))))
					return nil, 500, fmt.Errorf("bad request")
				}

				//decode base64
				htmlData := make([]byte, base64.StdEncoding.DecodedLen(len(encodedHtml)))
				n, err := base64.StdEncoding.Decode(htmlData, []byte(encodedHtml))
				if err != nil {
					is.Logger.Error("handlers - getPDF - base64.StdEncoding.Decode", zap.Error(err))
					return nil, 500, err
				}

				//validate html
				_, err = html.Parse(bytes.NewReader(htmlData[:n]))
				if err != nil {
					is.Logger.Error("handlers - getPDF - html.Parse", zap.Error(err))
					return nil, 500, err
				}

				library := is.mapGetValue(m, "library", is.Context.GetConfig("default_convert_library", "chromedp").(string))
				output := is.mapGetValue(m, "output", is.Context.GetConfig("default_convert_library", "chromedp").(string))

				result, err := is.convert(library, htmlData[:n], output)
				if err != nil {
					is.Logger.Error("handlers - getPDF - convert", zap.Error(err))
					return nil, 500, err
				}

				return result, 200, nil
			},
		},
	}
}
