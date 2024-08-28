package service

import (
	"context"
	"fmt"
	"homework/src/fifth-week/webook/internal/repository"
	"homework/src/fifth-week/webook/internal/service/sms"
	"math/rand"
	"time"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

type Config struct {
	MaxRetries    int
	RetryInterval time.Duration
	FlowLimit     int
}

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context,
		biz, phone, inputCode string) (bool, error)
}
type codeService struct {
	repo   repository.CodeRepository
	sms    sms.Service
	config Config
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service, config Config) CodeService {
	return &codeService{
		repo:   repo,
		sms:    smsSvc,
		config: config,
	}
}

func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	// 你在这儿，是不是要开始发送验证码了？
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	err = svc.sms.Send(ctx, codeTplId, []string{code}, phone)
	if err != nil {
		if svc.isFlowLimitExceeded(err) || svc.isServiceProviderCrashed(ctx) {
			svc.repo.DumpFailedRequest(ctx, biz, phone, code)
			go svc.retryFailedRequests(ctx)
			return nil
		}
		return err
	}
	return nil
}

func (svc *codeService) Verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 相当于，我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者，你这个不对
		return false, nil
	}
	return ok, err
}

func (svc *codeService) generate() string {
	// 0-999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func (svc *codeService) isFlowLimitExceeded(err error) bool {
	return err == ErrCodeSendTooMany
}

func (svc *codeService) isServiceProviderCrashed(ctx context.Context) bool {
	err := svc.sms.HealthCheck(ctx)
	return err != nil
}

func (svc *codeService) retryFailedRequests(ctx context.Context) {
	for i := 0; i < svc.config.MaxRetries; i++ {
		requests, err := svc.repo.GetFailedRequests(ctx)
		if err != nil {
			time.Sleep(svc.config.RetryInterval)
			continue
		}

		for _, req := range requests {
			err := svc.sms.Send(ctx, "1877556", []string{req.Code}, req.Phone)
			if svc.isFlowLimitExceeded(err) || svc.isServiceProviderCrashed(ctx) {
				break
				// Remove request from failed requests
			}
		}
		time.Sleep(svc.config.RetryInterval)
	}
}
