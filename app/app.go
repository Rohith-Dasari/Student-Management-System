package app

import (
	"database/sql"
	"log"
	"net/http"
	"sms/handlers"
	"sms/middleware"
	gradeRepository "sms/repository/gradesRepository"
	studentsRepository "sms/repository/studentRepository"
	userrepository "sms/repository/userRepository"
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
	studentHandler := handlers.NewStudentHandler(&studentService)
	authHandler := handlers.NewAuthHandler(authSevice)

	mux := http.NewServeMux()

	//auth
	mux.HandleFunc("POST /api/v1/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/signup", authHandler.Signup)

	//student
	mux.Handle("POST /api/v1/students", middleware.JWTAuth(studentHandler.AddStudent))
	mux.Handle("PATCH /api/v1/students/{studentID}", middleware.JWTAuth(studentHandler.UpdateStudent))

	// grades
	mux.Handle("POST /api/v1/grades", middleware.JWTAuth(gradeHandler.AddGrade))
	mux.Handle("GET /api/v1/grades", middleware.JWTAuth(gradeHandler.GetAverageOfClass))
	mux.Handle("GET /api/v1/grades/toppers", middleware.JWTAuth(gradeHandler.GetTopThree))
	mux.Handle("PATCH /api/v1/grades", middleware.JWTAuth(gradeHandler.UpdateGrade))
	return mux
}

func InitDBWithDSN(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Start() {
	DB, error := InitDBWithDSN("sms.db")
	if error != nil {
		log.Fatal(error.Error())
	}
	mux := SetupServer(DB)

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
