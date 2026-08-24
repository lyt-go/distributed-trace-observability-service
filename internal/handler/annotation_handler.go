package handler

import (
	"net/http"

	"tracing/internal/model"
	"tracing/pkg/httpx"
)

func (s *Server) registerAnnotationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/annotations", s.createAnnotation)
	mux.HandleFunc("GET /api/annotations", s.listAnnotations)
	mux.HandleFunc("GET /api/annotations/{id}", s.getAnnotation)
	mux.HandleFunc("PUT /api/annotations/{id}", s.updateAnnotation)
	mux.HandleFunc("DELETE /api/annotations/{id}", s.deleteAnnotation)
}

type createAnnotationRequest struct {
	SpanID  string `json:"span_id"`
	TraceID string `json:"trace_id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (s *Server) createAnnotation(w http.ResponseWriter, r *http.Request) {
	var req createAnnotationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAnnotation(model.Annotation{
		SpanID: req.SpanID, TraceID: req.TraceID, Name: req.Name, Message: req.Message,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAnnotations(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AnnotationFilter{
		SpanID:  r.URL.Query().Get("span_id"),
		TraceID: r.URL.Query().Get("trace_id"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListAnnotations(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAnnotation(w http.ResponseWriter, r *http.Request) {
	a, err := s.svc.GetAnnotation(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) updateAnnotation(w http.ResponseWriter, r *http.Request) {
	var req createAnnotationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateAnnotation(r.PathValue("id"), model.Annotation{Name: req.Name, Message: req.Message})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAnnotation(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteAnnotation(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
