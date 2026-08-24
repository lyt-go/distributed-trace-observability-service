package handler

import (
	"net/http"

	"tracing/internal/model"
	"tracing/pkg/httpx"
)

func (s *Server) registerSpanRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/spans", s.createSpan)
	mux.HandleFunc("GET /api/spans", s.listSpans)
	mux.HandleFunc("GET /api/spans/{id}", s.getSpan)
	mux.HandleFunc("PUT /api/spans/{id}", s.updateSpan)
	mux.HandleFunc("DELETE /api/spans/{id}", s.deleteSpan)
	mux.HandleFunc("POST /api/spans/{id}/finish", s.finishSpan)
}

type createSpanRequest struct {
	TraceID      string `json:"trace_id"`
	ParentSpanID string `json:"parent_span_id"`
	ServiceID    string `json:"service_id"`
	Operation    string `json:"operation"`
	Status       string `json:"status"`
	Tags         string `json:"tags"`
}

func (s *Server) createSpan(w http.ResponseWriter, r *http.Request) {
	var req createSpanRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sp, err := s.svc.CreateSpan(model.Span{
		TraceID: req.TraceID, ParentSpanID: req.ParentSpanID, ServiceID: req.ServiceID,
		Operation: req.Operation, Status: req.Status, Tags: req.Tags,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sp)
}

func (s *Server) listSpans(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SpanFilter{
		TraceID:   r.URL.Query().Get("trace_id"),
		ServiceID: r.URL.Query().Get("service_id"),
		Status:    r.URL.Query().Get("status"),
		Operation: r.URL.Query().Get("operation"),
	}
	items, total, err := s.svc.ListSpans(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSpan(w http.ResponseWriter, r *http.Request) {
	sp, err := s.svc.GetSpan(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sp)
}

func (s *Server) updateSpan(w http.ResponseWriter, r *http.Request) {
	var req createSpanRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sp, err := s.svc.UpdateSpan(r.PathValue("id"), model.Span{Operation: req.Operation, Tags: req.Tags})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sp)
}

func (s *Server) finishSpan(w http.ResponseWriter, r *http.Request) {
	var req finishTraceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sp, err := s.svc.FinishSpan(r.PathValue("id"), req.Status, req.DurationMs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sp)
}

func (s *Server) deleteSpan(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSpan(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
