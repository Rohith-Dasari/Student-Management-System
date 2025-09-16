package handlers

import (
	"encoding/json"
	"net/http"
	"sms/middleware"
	"sms/models"
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
	ss services.StudentServiceI
}

func NewStudentHandler(ss services.StudentServiceI) *StudentHandler {
	return &StudentHandler{ss: ss}
}

func (sh *StudentHandler) AddStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.CustomResponseSender(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	role, _ := middleware.GetUserRole(r.Context())
	if role != models.Admin {
		utils.CustomResponseSender(w, http.StatusForbidden, "only admin can access")
		return
	}

	var req CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// log.Println("reaching till here")

	student, err := sh.ss.CreateStudent(req.RollNumber, req.Name, req.ClassID, req.Semester)
	if err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, err.Error())
		return
	}
	// log.Println("reaching after db")
	res := CreateStudentResponse{
		StudentID:  student.StudentID,
		RollNumber: student.RollNumber,
		Name:       student.Name,
		ClassID:    student.ClassID,
		Semester:   student.Semester,
	}

	// log.Println("reaching after create response")
	utils.CustomResponseSender(w, http.StatusCreated, "successfully added", res)
}

func (sh *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.CustomResponseSender(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != models.Admin {
		utils.CustomResponseSender(w, http.StatusForbidden, "only admin can access")
		return
	}
	studentID := r.PathValue("studentID")
	if studentID == "" {
		utils.CustomResponseSender(w, http.StatusBadRequest, "invalid studentID")
		return
	}
	var updateStudent UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&updateStudent); err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err = sh.ss.UpdateStudent(studentID, updateStudent.Name, updateStudent.RollNumber, updateStudent.ClassID, updateStudent.Semester)
	if err != nil {
		utils.CustomResponseSender(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.CustomResponseSender(w, http.StatusOK, "updated successfully")
}
