package logger

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"strings"
)

// ZapLogger 是封装后的 zap.Logger，替代之前的全局变量
type ZapLogger struct {
	zap *zap.Logger
}

// New 创建一个新的 ZapLogger 实例
func New() (*ZapLogger, error) {
	// 自定义日期时间格式
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 设置日期时间格式
	// 自定义日志格式
	encoderConfig.EncodeLevel = customLevelEncoder
	encoderConfig.EncodeCaller = customEncodeCaller
	// 自定义字段宽度
	encoderConfig.ConsoleSeparator = "     "   // 设置字段之间的间隔宽度，这里设置为 5 个空格
	encoderConfig.MessageKey = "msg"           // 设置日志消息的字段名称
	encoderConfig.LevelKey = "level"           // 设置日志级别的字段名称
	encoderConfig.CallerKey = "caller"         // 设置调用者信息的字段名称
	encoderConfig.TimeKey = "time"             // 设置时间字段的名称
	encoderConfig.StacktraceKey = "stacktrace" // 设置堆栈跟踪字段的名称

	// 配置日志输出级别和格式
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel) // 设置日志级别为 Debug
	config.Encoding = "console"                             // 使用控制台输出
	config.OutputPaths = []string{"stdout"}                 // 输出到标准输出
	config.EncoderConfig = encoderConfig                    // 使用自定义的日志格式

	// 创建日志实例
	zapLogger, err := config.Build()
	if err != nil {
		log.Fatal("Failed to initialize zap logger: ", err) // 初始化失败时输出错误并终止
	}

	// 配置日志选项，跳过调用位置（使日志输出时不包含日志生成的文件位置）
	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(3))

	// 确保初始化成功
	fmt.Println("Logger initialized successfully")

	return &ZapLogger{zap: zapLogger}, nil
}

// Writer 返回一个实现了 logx.Writer 接口的 ZapWriter，用于替换 go-zero 默认日志
func (l *ZapLogger) Writer() logx.Writer {
	return NewZapWriter(l.zap)
}

// Sync 调用 zap 的 Sync 方法，刷新缓冲区
func (l *ZapLogger) Sync() error {
	return l.zap.Sync()
}

// 自定义 EncodeLevel 函数，支持颜色和宽度控制
func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	// 将 level 转换为整数类型，避免警告
	levelValue := int(level)
	// 获取 level 的 CapitalString() 输出
	levelString := fmt.Sprintf("%-5s", level.CapitalString()) // 使用 CapitalString()
	// 设置颜色
	var color string
	switch levelValue { // 使用 levelValue 进行整数比较
	case int(zapcore.DebugLevel):
		color = "\x1b[34m" // 蓝色
	case int(zapcore.InfoLevel):
		color = "\x1b[32m" // 绿色
	case int(zapcore.WarnLevel):
		color = "\x1b[33m" // 黄色
	case int(zapcore.ErrorLevel):
		color = "\x1b[31m" // 红色
	case int(zapcore.FatalLevel):
		color = "\x1b[35m" // 紫色
	default:
		color = "\x1b[37m" // 白色
	}

	// 将颜色与宽度设置的级别字符串一起添加到输出
	enc.AppendString(color + levelString + "\x1b[0m") // 末尾加上 "\x1b[0m" 以重置颜色
}

// 自定义 EncodeCaller 函数，用来格式化调用栈信息（保留文件路径、文件名和行号）
func customEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// 获取文件路径（包括目录）
	fileParts := strings.Split(caller.File, "/")
	// 保留文件路径中的最后两部分（路径+文件名）
	fileName := strings.Join(fileParts[len(fileParts)-2:], "/")
	// 拼接文件路径和行号
	callerString := fmt.Sprintf("%s:%d", fileName, caller.Line)
	// 控制拼接后的字符串宽度，这里保证文件路径+行号的宽度为 40
	// 如果拼接后总长度小于40，则填充空格，保证对齐
	callerStringWithWidth := fmt.Sprintf("%-35s", callerString)
	// 将格式化后的 callerString 添加到日志中
	enc.AppendString(callerStringWithWidth)
}

// Infof 格式化输出日志，记录 INFO 级别日志
func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.zap.Info(fmt.Sprintf(format, args...))
}

// Debugf 格式化输出日志，记录 DEBUG 级别日志
func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.zap.Debug(fmt.Sprintf(format, args...))
}

// Warnf 格式化输出日志，记录 WARN 级别日志
func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.zap.Warn(fmt.Sprintf(format, args...))
}

// Errorf 格式化输出日志，记录 ERROR 级别日志
func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.zap.Error(fmt.Sprintf(format, args...))
}

// Fatalf 格式化输出日志，记录 FATAL 级别日志
func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.zap.Fatal(fmt.Sprintf(format, args...))
}

// ZapWriter 实现了 go-zero 的 logx.Writer 接口
type ZapWriter struct {
	logger *zap.Logger
}

// NewZapWriter 创建一个新的 ZapWriter 实例
func NewZapWriter(logger *zap.Logger) logx.Writer {
	return &ZapWriter{
		logger: logger,
	}
}

// Write 拦截所有标准日志输出，实现 io.Writer 接口
func (w *ZapWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p))
	return len(p), nil
}

// Alert 实现 logx.Writer 的 Alert 方法
func (w *ZapWriter) Alert(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

// Close 实现 logx.Writer 的 Close 方法
func (w *ZapWriter) Close() error {
	// zap 日志器通常不需要关闭，如果需要可以在这里调用 w.logger.Sync()
	return nil
}

// Debug 实现 logx.Writer 的 Debug 方法
func (w *ZapWriter) Debug(v interface{}, fields ...logx.LogField) {
	w.logger.Debug(joinMessageWithFields(v, fields))
}

// Error 实现 logx.Writer 的 Error 方法
func (w *ZapWriter) Error(v interface{}, fields ...logx.LogField) {
	w.logger.Error(joinMessageWithFields(v, fields))
}

// Info 实现 logx.Writer 的 Info 方法
func (w *ZapWriter) Info(v interface{}, fields ...logx.LogField) {
	w.logger.Info(joinMessageWithFields(v, fields))
}

// Severe 实现 logx.Writer 的 Severe 方法
func (w *ZapWriter) Severe(v interface{}) {
	w.logger.Fatal(fmt.Sprint(v))
}

// Slow 实现 logx.Writer 的 Slow 方法
func (w *ZapWriter) Slow(v interface{}, fields ...logx.LogField) {
	w.logger.Warn(joinMessageWithFields(v, fields))
}

// Stack 实现 logx.Writer 的 Stack 方法
func (w *ZapWriter) Stack(v interface{}) {
	w.logger.Error(fmt.Sprint(v), zap.Stack("stack"))
}

// Stat 实现 logx.Writer 的 Stat 方法
func (w *ZapWriter) Stat(v interface{}, fields ...logx.LogField) {
	w.logger.Info(joinMessageWithFields(v, fields))
}

// joinMessageWithFields 辅助函数，拼接日志消息和关注的字段信息
func joinMessageWithFields(v interface{}, fields []logx.LogField) string {
	msg := fmt.Sprint(v)

	if len(fields) == 0 {
		return msg
	}

	for _, field := range fields {
		// 只保留你关心的字段
		switch field.Key {
		case "duration", "trace", "span":
			msg += fmt.Sprintf(" %-13s", field.Value)
		}
	}

	return msg
}
