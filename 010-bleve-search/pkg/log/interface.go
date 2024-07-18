package log

type (
	ILogger interface {
		DPanic(args ...interface{})
		DPanicf(template string, args ...interface{})
		DPanicw(template string, args ...interface{})
		//Print(args ...interface{})
		//Println(args ...interface{})
		//Printf(template string, args ...interface{})
		Debug(args ...interface{})
		Debugf(template string, args ...interface{})
		Debugw(template string, args ...interface{})
		Error(args ...interface{})
		Errorf(template string, args ...interface{})
		Errorw(template string, args ...interface{})
		Fatal(args ...interface{})
		Fatalf(template string, args ...interface{})
		Fatalw(template string, args ...interface{})
		Info(args ...interface{})
		Infof(template string, args ...interface{})
		Infow(template string, args ...interface{})
		Panic(args ...interface{})
		Panicf(template string, args ...interface{})
		Panicw(template string, args ...interface{})
		Sync() error
		Warn(args ...interface{})
		Warnf(template string, args ...interface{})
		Warnw(template string, args ...interface{})
	}
)
