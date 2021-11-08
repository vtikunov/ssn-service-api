package internal

//go:generate mockgen -destination=./mocks/api/service_repo_mock.go -package=apimocks github.com/ozonmp/ssn-service-api/internal/repo/subscription/service ServiceRepo
//go:generate mockgen -destination=./mocks/api/event_repo_mock.go -package=apimocks github.com/ozonmp/ssn-service-api/internal/repo/subscription/service EventRepo
//go:generate mockgen -destination=./mocks/api/transactional_session_mock.go -package=apimocks github.com/ozonmp/ssn-service-api/internal/repo TransactionalSession
//go:generate mockgen -destination=./mocks/app/event_repo_mock.go -package=appmocks github.com/ozonmp/ssn-service-api/internal/app/repo EventRepo
//go:generate mockgen -destination=./mocks/app/sender_mock.go -package=appmocks github.com/ozonmp/ssn-service-api/internal/app/sender EventSender
