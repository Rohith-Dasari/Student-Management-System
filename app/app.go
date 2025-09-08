package app

import (
	"database/sql"
	"log"
	"net/http"
	"sms/handlers"
	"sms/middleware"
	gradeRepository "sms/respository/gradesRepository"
	studentsRepository "sms/respository/studentRepository"
	userrepository "sms/respository/userRepository"
	"sms/services"

	_ "modernc.org/sqlite"
)

func SetupServer(db *sql.DB) *http.ServeMux {
	//repos
	gradeRepo := gradeRepository.NewGradeRepo(db)
	studentRepo := studentsRepository.NewStudentRepo(db)
	userRepo := userrepository.NewUserRepo(db)

	//services
	gradeService := services.NewGradeService(gradeRepo, studentRepo)
	studentService := services.NewStudentService(studentRepo)
	authSevice := services.NewAuthService(userRepo)

	//handlers
	gradeHandler := handlers.NewGradeHandler(gradeService)
	studentHandler := handlers.NewStudentHandler(studentService)
	authHandler := handlers.NewAuthHandler(authSevice)

	mux := http.NewServeMux()

	//auth
	mux.HandleFunc("POST /api/v1/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/signup", authHandler.Signup)

	//student
	mux.Handle("POST /api/v1/students", middleware.JWTAuth(http.HandlerFunc(studentHandler.AddStudent)))
	mux.Handle("PATCH /api/v1/students/{studentID}", middleware.JWTAuth(http.HandlerFunc(studentHandler.UpdateStudent)))

	// grades
	mux.Handle("POST /api/v1/grades", middleware.JWTAuth(http.HandlerFunc(gradeHandler.AddGrade)))
	mux.Handle("GET /api/v1/grades", middleware.JWTAuth(http.HandlerFunc(gradeHandler.GetAverageOfClass)))
	mux.Handle("GET /api/v1/grades/toppers", middleware.JWTAuth(http.HandlerFunc(gradeHandler.GetTopThree)))
	return mux
}

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "sms.db")
	if err != nil {
		log.Fatal("Could not open db")
	}
	return db
}
func Start() {
	DB := InitDB()
	mux := SetupServer(DB)

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
