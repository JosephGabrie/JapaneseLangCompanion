package models

import "time"

type KanaKanji struct {
    KanaKanji_ID   int    `json:"kanakanji_id"`
    Character      string `json:"character"`
    Romanization   string `json:"romanization"`
}

type Progress struct {
    UserID          int       `json:"user_id"`
    KanaKanjiID     int       `json:"kanakanji_id"`
    TimeCompleted   time.Time `json:"time_completed"`
    NextTimeReview  time.Time `json:"next_time_review"`
    MasteryLevel    int       `json:"mastery_level"`
    LastLearned     bool      `json:"last_learned"`
    UserTypedAnswer bool      `json:"user_typed_answer"`
}

type Users struct {
    User_ID   int    `json:"user_id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password  string `json:"password"`
}
