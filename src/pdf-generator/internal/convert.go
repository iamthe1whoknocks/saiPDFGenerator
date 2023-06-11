package internal

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

const (
	chromedpLib    = "chromedp"
	wkhtmltopdfLib = "wkhtmltopdf"
	fileOutput     = "file"
	s3Output       = "s3"
)

type Response struct {
	Result string `json:"result"`
	Status string `json:"status"`
}

var (
	fileNum = 1
)

// convert html from link to pdf depending on chosen library
func (is *InternalService) convert(library string, html []byte, output string) (result interface{}, err error) {
	file := []byte{}

	switch library {
	case wkhtmltopdfLib:
		file, err = is.wkhtmltopdfConvert(html)
		if err != nil {
			return nil, fmt.Errorf("convert - wkhtmltopdfConvert - %w", err)
		}
	default:
		file, err = is.chromedpConvert(html)
		if err != nil {
			return nil, fmt.Errorf("convert - chromedpConvert - %w", err)
		}
	}

	switch output {
	case s3Output:
		s3Link, err := is.s3Upload(file)
		if err != nil {
			return nil, fmt.Errorf("convert - s3Upload - %w", err)
		}
		return Response{Result: s3Link, Status: "OK"}, nil
	default:
		// save file to check without s3
		if fileNum > is.FileNum {
			fileNum = 1
		}
		filename := fmt.Sprintf("%s_%s.pdf", "file", strconv.Itoa(fileNum))
		if err := ioutil.WriteFile("files/"+filename, file, 0644); err != nil {
			return nil, fmt.Errorf("convert - file output - ioutil.WriteFile - %w", err)
		}
		host := is.Context.GetConfig("common.http.host", "localhost").(string)
		port := is.Context.GetConfig("common.http.fileserver_port", 8085).(int)

		path := fmt.Sprintf("http://%s:%s/%s", host, strconv.Itoa(port), filename)
		fileNum++
		return Response{Result: path, Status: "OK"}, nil
	}
}

// convert html into pdf using chromedp
func (is *InternalService) chromedpConvert(html []byte) (fileData []byte, err error) {
	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	var pdfBuffer []byte
	if err := chromedp.Run(taskCtx, pdfGrabberFromFile(html, "body", &pdfBuffer)); err != nil {
		is.Logger.Error("handlers - chromedpConvert - chromedp.Run", zap.Error(err))
		return nil, err
	}

	return pdfBuffer, nil
}

// pdfGrabberFromFile func to convert html into pdf
func pdfGrabberFromFile(html []byte, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		actionLoadHTMLContent(html),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

// actionLoadHTMLContent load html from byte slice
func actionLoadHTMLContent(data []byte) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		isSet, isSetLock := false, sync.Mutex{}
		listenerCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		chromedp.ListenTarget(listenerCtx, func(ev interface{}) {
			if _, ok := ev.(*page.EventLoadEventFired); ok {
				// stop listener
				cancel()

				isSetLock.Lock()
				isSet = true
				isSetLock.Unlock()
			}
		})

		frameTree, err := page.GetFrameTree().Do(ctx)
		if err != nil {
			return err
		}

		if err := page.SetDocumentContent(frameTree.Frame.ID, string(data)).Do(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-listenerCtx.Done():
			isSetLock.Lock()
			defer isSetLock.Unlock()
			if isSet {
				return nil
			}
			return listenerCtx.Err()
		}
	}
}

// wkhtmltopdfConvert html into pdf using wkhtmltopdf
func (is *InternalService) wkhtmltopdfConvert(data []byte) (output []byte, err error) {

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	is.Logger.Info("wkhtmltopdfConvert", zap.String("data", string(data)))

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)

	// Create a new input page from an URL
	_page := wkhtmltopdf.NewPageReader(bytes.NewReader(data))

	// Set options for this page
	_page.FooterRight.Set("[page]")
	_page.FooterFontSize.Set(10)
	_page.Zoom.Set(0.95)

	_page.EnableLocalFileAccess.Set(true)

	// Add to document
	pdfg.AddPage(_page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
