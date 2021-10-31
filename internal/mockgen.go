package internal

//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks github.com/ozonmp/ssn-service-api/internal/repo ServiceRepo
//go:generate mockgen -destination=./mocks/event_repo_mock.go -package=mocks github.com/ozonmp/ssn-service-api/internal/app/repo EventRepo
//go:generate mockgen -destination=./mocks/sender_mock.go -package=mocks github.com/ozonmp/ssn-service-api/internal/app/sender EventSender
