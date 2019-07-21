package connection

func Init(groupCount, maxConnCount int) error {
	err := InitPool(groupCount, maxConnCount)
	if err != nil {
		return err
	}
	InitEpoller()
	return nil
}
