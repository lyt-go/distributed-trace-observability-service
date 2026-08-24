package handler

import (
	"net/http"
	"strconv"

	"tracing/internal/model"
	"tracing/pkg/httpx"
)

func (s *Server) registerTraceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/traces", s.createTrace)
	mux.HandleFunc("GET /api/traces", s.listTraces)
	mux.HandleFunc("GET /api/traces/{id}", s.getTrace)
	mux.HandleFunc("PUT /api/traces/{id}", s.updateTrace)
	mux.HandleFunc("DELETE /api/traces/{id}", s.deleteTrace)
	mux.HandleFunc("POST /api/traces/{id}/finish", s.finishTrace)
}

type createTraceRequest struct {
	ServiceID string `json:"service_id"`
	Operation string `json:"operation"`
	Status    string `json:"status"`
}

func (s *Server) createTrace(w http.ResponseWriter, r *http.Request) {
	var req createTraceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTrace(model.Trace{ServiceID: req.ServiceID, Operation: req.Operation, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTraces(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	minDur, _ := strconv.ParseInt(r.URL.Query().Get("min_duration"), 10, 64)
	filter := model.TraceFilter{
		ServiceID:   r.URL.Query().Get("service_id"),
		Status:      r.URL.Query().Get("status"),
		MinDuration: minDur,
		Keyword:     r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListTraces(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTrace(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.GetTrace(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) updateTrace(w http.ResponseWriter, r *http.Request) {
	var req createTraceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTrace(r.PathValue("id"), model.Trace{Operation: req.Operation})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type finishTraceRequest struct {
	Status     string `json:"status"`
	DurationMs int64  `json:"duration_ms"`
}

func (s *Server) finishTrace(w http.ResponseWriter, r *http.Request) {
	var req finishTraceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.FinishTrace(r.PathValue("id"), req.Status, req.DurationMs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTrace(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteTrace(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
