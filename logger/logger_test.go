package logger

import (
	"encoding/json"
	"log"
	"testing"
)

func TestError(t *testing.T) {
	lc := &Options{
		ErrLog: "./output.log",
		InfoLog: "./output.log",
		MaxSize: 10,
		MaxBackups: 5,
		MaxAge: 7,
		Level: "debug",
		LocalTime: true,
	}

	if err := lc.InitLogger(); err != nil {
		log.Fatal(err.Error())
	}

	Info("test info msg", "this is test info message", map[string]interface{}{"Key": "45222524", "val": "nihaoshijie"})
	if err := json.Unmarshal([]byte(""), ""); err != nil {
		Error("test msg", err.Error(),
			map[string]interface{}{"Key": "123454678", "val": 20, "is_fail": false, "struct": nil, "num": int8(12), "salary": 12.09})
	}
}
