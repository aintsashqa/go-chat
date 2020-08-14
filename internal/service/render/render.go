package render

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	WrapperHtmlRenderServiceRenderMethod = "HtmlRenderService.Render"
)

type HtmlRenderService struct {
	directory string
	layout    string

	logger *logrus.Logger
}

func NewHtmlRenderService(directory, layout string, logger *logrus.Logger) *HtmlRenderService {
	return &HtmlRenderService{
		directory: directory,
		layout:    layout,

		logger: logger,
	}
}

func (s *HtmlRenderService) getFilename(name string) string {
	return fmt.Sprintf("%s/%s.html", s.directory, name)
}

func (s *HtmlRenderService) getLayoutFilename() string {
	return s.getFilename(s.layout)
}

func (s *HtmlRenderService) getTempalteFilename(name string) string {
	return s.getFilename(fmt.Sprintf("web/%s", name))
}

func (s *HtmlRenderService) Render(w http.ResponseWriter, name string, data interface{}) {
	layout := s.getLayoutFilename()
	templateFilename := s.getTempalteFilename(name)
	temp := template.Must(template.ParseFiles(layout, templateFilename))
	if err := temp.ExecuteTemplate(w, "layout", data); err != nil {
		s.logger.Error(errors.Wrap(err, WrapperHtmlRenderServiceRenderMethod))
		s.ErrorResponse(w, http.StatusInternalServerError)
	}
}

func (s *HtmlRenderService) ErrorResponse(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}
