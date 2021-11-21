package internal

//go:generate mockgen -destination=./mocks/api/service_repo_mock.go -package=apimocks github.com/ozonmp/ssn-service-api/internal/repo/subscription/service ServiceRepo
//go:generate mockgen -destination=./mocks/api/event_repo_mock.go -package=apimocks github.com/ozonmp/ssn-service-api/internal/repo/subscription/service EventRepo
//go:generate mockgen -destination=./mocks/api/transactional_session_mock.go -package=apimocks github.com/ozonmp/ssn-service-api/internal/repo TransactionalSession
//go:generate mockgen -destination=./mocks/retranslator/event_repo_mock.go -package=retranslatormocks github.com/ozonmp/ssn-service-api/internal/retranslator/repo EventRepo
//go:generate mockgen -destination=./mocks/retranslator/sender_mock.go -package=retranslatormocks github.com/ozonmp/ssn-service-api/internal/retranslator/sender EventSender
//go:generate mockgen -destination=./mocks/facade/service_repo_mock.go -package=facademocks github.com/ozonmp/ssn-service-api/internal/facade/repo/subscription/service ServiceRepo
//go:generate mockgen -destination=./mocks/facade/transactional_session_mock.go -package=facademocks github.com/ozonmp/ssn-service-api/internal/facade/repo TransactionalSession
