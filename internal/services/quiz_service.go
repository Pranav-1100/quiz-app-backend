package services

import (
    "database/sql"
	"encoding/json" 
	"log"
	"math/rand"  // Add this
    "time" 
    "github.com/Pranav-1100/quiz-app-backend/internal/models"
)

type QuizService struct {
    DB *sql.DB
}

func NewQuizService(db *sql.DB) *QuizService {
    return &QuizService{DB: db}
}

func (s *QuizService) GetNextLevel(currentLevel int) (*models.Level, error) {
    var level models.Level
    err := s.DB.QueryRow("SELECT id, name, difficulty, unlock_threshold FROM levels WHERE id > ? ORDER BY id ASC LIMIT 1", currentLevel).Scan(&level.ID, &level.Name, &level.Difficulty, &level.UnlockThreshold)
    if err == sql.ErrNoRows {
        return nil, nil // No more levels
    }
    if err != nil {
        return nil, err
    }
    return &level, nil
}

func (s *QuizService) GetQuestionsByLevel(levelID int) ([]models.Question, error) {
    rows, err := s.DB.Query("SELECT id, level_id, question, options, correct_answer, difficulty, image_url, explanation FROM questions WHERE level_id = ?", levelID)
    if err != nil {
        log.Printf("Error querying database: %v", err)
        return nil, err
    }
    defer rows.Close()

    var questions []models.Question
    for rows.Next() {
        var q models.Question
        var optionsJSON string
        if err := rows.Scan(&q.ID, &q.LevelID, &q.Question, &optionsJSON, &q.CorrectAnswer, &q.Difficulty, &q.ImageURL, &q.Explanation); err != nil {
            log.Printf("Error scanning row: %v", err)
            return nil, err
        }
        if err := json.Unmarshal([]byte(optionsJSON), &q.Options); err != nil {
            log.Printf("Error unmarshaling options: %v", err)
            return nil, err
        }
        questions = append(questions, q)
    }

    if len(questions) == 0 {
        log.Printf("No questions found for level %d", levelID)
    }

    return questions, nil
}

func (s *QuizService) CheckAnswer(questionID int, userAnswer string) (bool, string, error) {
    var correctAnswer, explanation string
    err := s.DB.QueryRow("SELECT correct_answer, explanation FROM questions WHERE id = ?", questionID).Scan(&correctAnswer, &explanation)
    if err != nil {
        return false, "", err
    }
    return userAnswer == correctAnswer, explanation, nil
}

func (s *QuizService) GetFiftyFiftyOptions(questionID int) ([]string, error) {
    // Implement logic to return two options, including the correct one
    // For example:
    question, err := s.GetQuestionByID(questionID)
    if err != nil {
        return nil, err
    }

    correctAnswer := question.CorrectAnswer
    var incorrectOptions []string
    for _, option := range question.Options {
        if option != correctAnswer {
            incorrectOptions = append(incorrectOptions, option)
        }
    }

    // Randomly select one incorrect option
    rand.Seed(time.Now().UnixNano())
    randomIndex := rand.Intn(len(incorrectOptions))
    return []string{correctAnswer, incorrectOptions[randomIndex]}, nil
}

func (s *QuizService) GetHint(questionID int) (string, error) {
    // Implement logic to return a hint for the question
    // For example:
    question, err := s.GetQuestionByID(questionID)
    if err != nil {
        return "", err
    }

    // You might want to store hints in the database or generate them dynamically
    return "This is a hint for the question: " + question.Question, nil
}

func (s *QuizService) GetExpertAdvice(questionID int) (string, error) {
    // Implement logic to return expert advice for the question
    // For example:
    question, err := s.GetQuestionByID(questionID)
    if err != nil {
        return "", err
    }

    // You might want to store expert advice in the database or generate it dynamically
    return "Expert advice: The correct answer is likely to be " + question.CorrectAnswer, nil
}

func (s *QuizService) GetQuestionByID(questionID int) (*models.Question, error) {
    // Implement logic to fetch a question by its ID from the database
    var question models.Question
    err := s.DB.QueryRow("SELECT id, level_id, question, options, correct_answer, difficulty, image_url, explanation FROM questions WHERE id = ?", questionID).Scan(
        &question.ID, &question.LevelID, &question.Question, &question.Options, &question.CorrectAnswer, &question.Difficulty, &question.ImageURL, &question.Explanation,
    )
    if err != nil {
        return nil, err
    }
    return &question, nil
}

func (s *QuizService) GetAllQuizzes() ([]models.Quiz, error) {
    rows, err := s.DB.Query("SELECT id, name, description, difficulty, unlock_threshold FROM levels ORDER BY difficulty, id")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var quizzes []models.Quiz
    for rows.Next() {
        var q models.Quiz
        if err := rows.Scan(&q.ID, &q.Name, &q.Description, &q.Difficulty, &q.UnlockThreshold); err != nil {
            return nil, err
        }
        quizzes = append(quizzes, q)
    }

    return quizzes, nil
}