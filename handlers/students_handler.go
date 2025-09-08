package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sms/middleware"
	"sms/services"
	"sms/utils"
)

type CreateStudentRequest struct {
	RollNumber string `json:"roll_number"`
	Name       string `json:"name"`
	ClassID    string `json:"classID"`
	Semester   int    `json:"semester"`
}

type CreateStudentResponse struct {
	StudentID  string `json:"studentID"`
	RollNumber string `json:"roll_number"`
	Name       string `json:"name"`
	ClassID    string `json:"classID"`
	Semester   int    `json:"semester"`
}

type UpdateStudentRequest struct {
	RollNumber string `json:"roll_number,omitempty"`
	Name       string `json:"name,omitempty"`
	ClassID    string `json:"classID,omitempty"`
	Semester   int    `json:"semester,omitempty"`
}

type StudentHandler struct {
	ss services.StudentService
}

func NewStudentHandler(ss services.StudentService) StudentHandler {
	return StudentHandler{ss: ss}
}

func (sh *StudentHandler) AddStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	role, err := middleware.GetUserRole(r.Context())
	if err != nil {
		utils.CustomError(w, http.StatusInternalServerError, "unable to get user role")
	}
	if role != "admin" {
		utils.CustomError(w, http.StatusForbidden, "only admin can access")
		return
	}

	var req CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.CustomError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	log.Println("reaching till here")

	student, err := sh.ss.CreateStudent(req.RollNumber, req.Name, req.ClassID, req.Semester)
	if err != nil {
		utils.CustomError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("reaching after db")
	res := CreateStudentResponse{
		StudentID:  student.StudentID,
		RollNumber: student.RollNumber,
		Name:       student.Name,
		ClassID:    student.ClassID,
		Semester:   student.Semester,
	}

	log.Println("reaching after create response")
	utils.SendCustomResponse(w, http.StatusOK, "successfully added", res)
}

func (sh *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "admin" {
		utils.CustomError(w, http.StatusForbidden, "only admin can access")
		return
	}
	studentID := r.PathValue("studentID")
	if studentID == "" {
		utils.CustomError(w, http.StatusBadRequest, "invalid studentID")
		return
	}
	var updateStudent UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&updateStudent); err != nil {
		utils.CustomError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err = sh.ss.UpdateStudent(studentID, updateStudent.Name, updateStudent.RollNumber, updateStudent.ClassID, updateStudent.Semester)
	if err != nil {
		utils.CustomError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.CustomError(w, http.StatusOK, "updated successfully")
}
