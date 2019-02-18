package error

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"runtime"
)
var (
	ServerError                       = GenerateError("Something went wrong! Please try again later")
	ProjectDoesNotExistError          = GenerateError("project does not exists")
	QuestionDoesNotExistError         = GenerateError("question does not exists")
	AnswerDoesNotExistError           = GenerateError("answer does not exists")
	UserNotFoundError                 = GenerateError("user not found")
	UserNotAuthorizedToUnpublishError = GenerateError("you are not authorized to unpublish this project")
	AlreadyUnpublishedError           = GenerateError("project is already unpublished")
	ReviewDoesNotExistsError          = GenerateError("review does not exists")
	ReviewCommentDoesNotExistsError   = GenerateError("review comment does not exists")
	TimeStampError                    = GenerateError("time should be a unix timestamp")
	InternalServerError               = GenerateError("internal server error")
	InvalidStructureIDError           = GenerateError("invalid structure id")
)

func GenerateError(err string) error {
	return errors.New(err)
}
func IsForeignKeyError(err error) bool {
	pgErr := err.(*pq.Error);
	if pgErr.Code == "23503" {
		return true
	}
	return false
}

func DebugPrintf(err_ error, args ...interface{}) string {
	programCounter, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(programCounter)
	msg := fmt.Sprintf("[%s: %s %d] %s, %s", file, fn.Name(), line, err_, args)
	log.Println(msg)
	return msg
}