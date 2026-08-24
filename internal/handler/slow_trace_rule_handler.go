package handler

import (
	"net/http"

	"tracing/internal/model"
	"tracing/pkg/httpx"
)

func (s *Server) registerSlowTraceRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/slow-rules", s.createSlowTraceRule)
	mux.HandleFunc("GET /api/slow-rules", s.listSlowTraceRules)
	mux.HandleFunc("GET /api/slow-rules/{id}", s.getSlowTraceRule)
	mux.HandleFunc("PUT /api/slow-rules/{id}", s.updateSlowTraceRule)
	mux.HandleFunc("DELETE /api/slow-rules/{id}", s.deleteSlowTraceRule)
	mux.HandleFunc("GET /api/slow-rules/detect", s.detectSlowTraces)
}

type createSlowTraceRuleRequest struct {
	Name        string `json:"name"`
	ServiceID   string `json:"service_id"`
	ThresholdMs int64  `json:"threshold_ms"`
	Status      string `json:"status"`
}

func (s *Server) createSlowTraceRule(w http.ResponseWriter, r *http.Request) {
	var req createSlowTraceRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.CreateSlowTraceRule(model.SlowTraceRule{
		Name: req.Name, ServiceID: req.ServiceID, ThresholdMs: req.ThresholdMs, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rule)
}

func (s *Server) listSlowTraceRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SlowTraceRuleFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListSlowTraceRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSlowTraceRule(w http.ResponseWriter, r *http.Request) {
	rule, err := s.svc.GetSlowTraceRule(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) updateSlowTraceRule(w http.ResponseWriter, r *http.Request) {
	var req createSlowTraceRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.UpdateSlowTraceRule(r.PathValue("id"), model.SlowTraceRule{
		Name: req.Name, ServiceID: req.ServiceID, ThresholdMs: req.ThresholdMs, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) deleteSlowTraceRule(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSlowTraceRule(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) detectSlowTraces(w http.ResponseWriter, r *http.Request) {
	traces, err := s.svc.DetectSlowTraces()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"items": traces, "total": len(traces)})
}
