package services

import (
	"errors"
	gradeRepository "sms/repository/gradesRepository"
	studentsRepository "sms/repository/studentRepository"
	"sort"
)

type GradeService struct {
	gr gradeRepository.GradeRepositoryI
	sr studentsRepository.StudentRepositoryI
}
type GradeReponse struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	Grade       int    `json:"grade"`
}

func NewGradeService(gr gradeRepository.GradeRepositoryI, sr studentsRepository.StudentRepositoryI) *GradeService {
	return &GradeService{gr, sr}
}

func (gs *GradeService) GetAverageGradeOfStudent(studentID string, semester int) (int, error) {
	//get grades of student
	grades, err := gs.gr.GetSemesterGrades(studentID, semester)
	if err != nil {
		return 0, err
	}
	if len(grades) == 0 {
		return 0, nil
	}

	//calculate average
	sum := 0
	for _, grade := range grades {
		sum += grade
	}
	avg := sum / len(grades)
	return avg, nil
}

func (gs *GradeService) GetAverageGradesOfEachStudentOfClass(classID string, semester int) (map[string]int, error) {
	//get students of that class

	students, err := gs.sr.GetAllStudentsOfClass(classID, semester)
	if err != nil {
		return nil, err
	}

	//get average of each student

	averageOfEachStudent := make(map[string]int, len(students))

	for _, student := range students {
		averageOfEachStudent[student.StudentID], err = gs.GetAverageGradeOfStudent(student.StudentID, semester)
		if err != nil {
			return nil, err
		}
	}

	return averageOfEachStudent, err

}

func (gs *GradeService) GetAverageOfClass(classID string, semester int) (int, error) {

	averageOfEachStudent, err := gs.GetAverageGradesOfEachStudentOfClass(classID, semester)
	if err != nil {
		return 0, err
	}
	//get average of that average

	sum := 0
	for _, avg := range averageOfEachStudent {
		sum += avg
	}
	averageOfClass := sum / len(averageOfEachStudent)
	return averageOfClass, nil
}

// func (gs *GradeService) GetToppers(classID string, semester) ([]GradeReponse, error) {
// 	averageOfEachStudent, err := gs.GetAverageGradesOfEachStudentOfClass(classID, semester)
// 	if err != nil {
// 		return nil, err
// 	}
// 	//find max 3 and send those

// 	//how to find max 3 from map

// 	var max1_key, max2_key, max3_key string
// 	var max1_value, max2_value, max3_value int

// 	for key, value := range averageOfEachStudent {
// 		if value > max1_value {
// 			max3_key = max2_key
// 			max3_value = max2_value

// 			max2_key = max1_key
// 			max2_value = max1_value

// 			max1_key = key
// 			max1_value = value

// 		} else if value > max2_value {
// 			max3_key = max2_key
// 			max3_value = max2_value

// 			max2_key = key
// 			max2_value = value
// 		} else if value > max3_value {
// 			max3_key = key
// 			max3_value = value
// 		}
// 	}

// 	student1, _ := gs.sr.GetStudentByID(max1_key)
// 	student2, _ := gs.sr.GetStudentByID(max2_key)
// 	student3, _ := gs.sr.GetStudentByID(max3_key)

// 	topThreeStudents := make([]GradeReponse, 3)
// 	topThreeStudents[0] = GradeReponse{
// 		StudentID:   student1.StudentID,
// 		StudentName: student1.Name,
// 		Grade:       max1_value,
// 	}
// 	topThreeStudents[1] = GradeReponse{
// 		StudentID:   student2.StudentID,
// 		StudentName: student2.Name,
// 		Grade:       max2_value,
// 	}
// 	topThreeStudents[2] = GradeReponse{
// 		StudentID:   student3.StudentID,
// 		StudentName: student3.Name,
// 		Grade:       max3_value,
// 	}

// 	return topThreeStudents, nil
// }

type kv struct {
	Key   string
	Value int
}

func (gs *GradeService) GetToppers(classID string, semester int, limit int) ([]GradeReponse, error) {
	averageOfEachStudent, err := gs.GetAverageGradesOfEachStudentOfClass(classID, semester)
	if err != nil {
		return nil, err
	}
	var sorted []kv
	for k, v := range averageOfEachStudent {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	if limit > len(sorted) {
		limit = len(sorted)
	}
	topStudents := make([]GradeReponse, 0, limit)
	for i := 0; i < limit; i++ {
		student, err := gs.sr.GetStudentByID(sorted[i].Key)
		if err != nil {
			return nil, err
		}
		topStudents = append(topStudents, GradeReponse{
			StudentID:   student.StudentID,
			StudentName: student.Name,
			Grade:       sorted[i].Value,
		})
	}

	return topStudents, nil
}

func (gs *GradeService) AddGrades(studentID string, subjectID string, Grade int, semester int) error {
	if Grade < 0 {
		return errors.New("grade can't be negative")
	}
	err := gs.gr.AddGrades(studentID, subjectID, Grade, semester)
	return err
}

func (gs *GradeService) UpdateGrade(studentID string, subjectID string, newGrade int) error {
	if newGrade < 0 {
		return errors.New("grade can't be negative")
	}
	err := gs.gr.UpdateGrade(studentID, subjectID, newGrade)
	return err
}
