package handlers

import (
	"encoding/json"
	"net/http"
	"sms/middleware"
	"sms/services"
	"sms/utils"
	"strconv"
)

type AddGradeRequest struct {
	StudentID string `json:"studentID"`
	SubjectID string `json:"subjectID"`
	Semester  int    `json:"semester"`
	Grade     int    `json:"grade"`
}

type UpdateGrade struct {
	StudentID string `json:"studentID"`
	SubjectID string `json:"subjectID"`
	NewGrade  int    `json:"new_grade"`
}

type GradeHandler struct {
	gs services.GradeServiceI
}

func NewGradeHandler(gs services.GradeServiceI) *GradeHandler {
	return &GradeHandler{gs}
}

func (gh *GradeHandler) GetAverageOfClass(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.CustomResponseSender(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "faculty" {
		utils.CustomResponseSender(w, http.StatusForbidden, "only faculty can access")
		return
	}
	classID := r.PathValue("classID")
	semester, err := strconv.Atoi(r.PathValue("semester"))
	if err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, "semester must be a number")
		return
	}
	if classID == "" {
		utils.CustomResponseSender(w, http.StatusBadRequest, "invalid classID")
		return
	}

	data, err := gh.gs.GetAverageOfClass(classID, semester)
	if err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.CustomResponseSender(w, http.StatusOK, "ok", data)
}

func (gh *GradeHandler) GetToppers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.CustomResponseSender(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "faculty" {
		utils.CustomResponseSender(w, http.StatusForbidden, "only faculty can access")
		return
	}
	classID := r.PathValue("classID")
	semester, err := strconv.Atoi(r.PathValue("semester"))

	if err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, "semester must be a number")
		return
	}

	query := r.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, "limit must be a number")
		return
	}

	if classID == "" {
		utils.CustomResponseSender(w, http.StatusBadRequest, "invalid classID")
		return
	}

	data, err := gh.gs.GetToppers(classID, semester, limit)
	if err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.CustomResponseSender(w, http.StatusOK, "ok", data)
}

func (gh *GradeHandler) AddGrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.CustomResponseSender(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "faculty" {
		utils.CustomResponseSender(w, http.StatusForbidden, "only faculty can access")
		return
	}
	var req AddGradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err = gh.gs.AddGrades(req.StudentID, req.SubjectID, req.Grade, req.Semester)
	if err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.CustomResponseSender(w, http.StatusCreated, "grade successfully added")
}

func (gh *GradeHandler) UpdateGrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.CustomResponseSender(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "faculty" {
		utils.CustomResponseSender(w, http.StatusForbidden, "only faculty can access")
		return
	}
	var req UpdateGrade
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err = gh.gs.UpdateGrade(req.StudentID, req.SubjectID, req.NewGrade)
	if err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.CustomResponseSender(w, http.StatusOK, "grade updated added")
}
