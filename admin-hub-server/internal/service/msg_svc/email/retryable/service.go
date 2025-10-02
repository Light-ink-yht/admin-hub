package retryable

// 小心并发问题
//type CodeService struct {
//	svc *tencent.EmailService
//	// 重试
//	retryCnt int
//}
//
//func (c CodeService) Send(ctx context.Context, email string, subject, body string) error {
//	err := c.svc.Send(ctx, email, subject, body)
//	for err != nil && c.retryCnt < 10 {
//		err = c.svc.Send(ctx, email, subject, body)
//		c.retryCnt++
//		return
//	}
//}
