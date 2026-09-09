package foundry

func (f *FoundryApi) background(fn func()) {
	f.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				f.Transport.Logger.Error("Program panicked", "err", err.(error).Error())
			}
		}()

		fn()
	})
}
