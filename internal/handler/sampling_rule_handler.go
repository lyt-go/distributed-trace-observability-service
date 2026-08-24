package handler

import (
	"net/http"

	"tracing/internal/model"
	"tracing/pkg/httpx"
)

func (s *Server) registerSamplingRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/sampling-rules", s.createSamplingRule)
	mux.HandleFunc("GET /api/sampling-rules", s.listSamplingRules)
	mux.HandleFunc("GET /api/sampling-rules/{id}", s.getSamplingRule)
	mux.HandleFunc("PUT /api/sampling-rules/{id}", s.updateSamplingRule)
	mux.HandleFunc("DELETE /api/sampling-rules/{id}", s.deleteSamplingRule)
	mux.HandleFunc("GET /api/sampling/evaluate", s.evaluateSampling)
}

type createSamplingRuleRequest struct {
	Name           string `json:"name"`
	ServicePattern string `json:"service_pattern"`
	Rate           int    `json:"rate"`
	Status         string `json:"status"`
}

func (s *Server) createSamplingRule(w http.ResponseWriter, r *http.Request) {
	var req createSamplingRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.CreateSamplingRule(model.SamplingRule{
		Name: req.Name, ServicePattern: req.ServicePattern, Rate: req.Rate, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rule)
}

func (s *Server) listSamplingRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SamplingRuleFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListSamplingRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSamplingRule(w http.ResponseWriter, r *http.Request) {
	rule, err := s.svc.GetSamplingRule(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) updateSamplingRule(w http.ResponseWriter, r *http.Request) {
	var req createSamplingRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.UpdateSamplingRule(r.PathValue("id"), model.SamplingRule{
		Name: req.Name, ServicePattern: req.ServicePattern, Rate: req.Rate, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) deleteSamplingRule(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSamplingRule(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) evaluateSampling(w http.ResponseWriter, r *http.Request) {
	rule, sample, err := s.svc.EvaluateSampling(r.URL.Query().Get("service_name"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{
		"sample":       sample,
		"matched_rule": rule,
	})
}
