/**
 * 回归测试：日志仓关闭后，其内部接收协程应退出
 * 背景：LogDataStorage.readyGo中select内的break只会退出select不会退出for，
 * 导致每次关闭日志仓都泄漏一个协程，内存随运行时间增长（关联上游 issue #70）
 * 覆盖两条退出路径：忙时关闭（select case分支收到nil）与闲置自动关闭（内层阻塞接收收到nil）
 * 注：测试放在ldb包内与其他存储测试同进程运行，避免并行测试进程争用固定目录
 * /glogcenter下leveldb的文件锁
 */
package ldb

import (
	"glc/conf"
	"glc/ldb/storage/logdata"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestLogDataStorageCloseNoGoroutineLeak(t *testing.T) {
	t.Setenv("GLC_MAX_IDLE_TIME", "1") // 闲置1秒即自动关闭，按生产路径加速测试
	conf.UpdateConfigByEnv()

	name := "leakverify-" + strconv.Itoa(int(time.Now().UnixNano()))
	t.Cleanup(func() { // 清理测试产生的日志仓目录，避免残留在/glogcenter
		os.RemoveAll(filepath.Join(conf.GetStorageRoot(), name+"-warm"))
		os.RemoveAll(filepath.Join(conf.GetStorageRoot(), name+"-busy"))
		os.RemoveAll(filepath.Join(conf.GetStorageRoot(), name))
	})

	// 预热：吸收惰性初始化的协程，待预热仓的各定时器协程退出后再计基线
	w := logdata.NewLogDataStorage(name + "-warm")
	w.Add(&logdata.LogDataModel{Text: "warmup log line"})
	time.Sleep(6 * time.Second)
	before := runtime.NumGoroutine()

	// 阶段1：忙时关闭——消费协程正在处理日志时收到退出信号（走select case分支的return）
	s1 := logdata.NewLogDataStorage(name + "-busy")
	s1.Add(&logdata.LogDataModel{Text: "log line busy"})
	s1.Close() // 立即关闭，nil信号在消费协程处理完日志前就已入队
	// 先睡过件数保存协程开仓后首个5秒tick（否则其未退出会干扰计数判断），再静置
	time.Sleep(6 * time.Second)
	settleGoroutines(t)
	mid := runtime.NumGoroutine()
	t.Logf("busy close: before=%d after=%d delta=%+d", before, mid, mid-before)
	if mid > before {
		t.Fatalf("忙时关闭日志仓后仍有%+d个协程存活：readyGo消费协程未退出（协程泄漏）", mid-before)
	}

	// 阶段2：闲置自动关闭——消费协程空闲停驻时收到退出信号（走内层阻塞接收分支的return）
	s2 := logdata.NewLogDataStorage(name)
	for i := 0; i < 3; i++ {
		s2.Add(&logdata.LogDataModel{Text: "log line " + strconv.Itoa(i)})
	}
	time.Sleep(6 * time.Second) // 覆盖闲置自动关闭（约2秒）与件数保存协程首个5秒tick
	settleGoroutines(t)         // 等待残余协程自行退出
	runtime.GC()
	after := runtime.NumGoroutine()
	t.Logf("idle close: before=%d after=%d delta=%+d", before, after, after-before)
	if after > before {
		t.Errorf("日志仓闲置自动关闭后仍有%+d个协程存活：readyGo消费协程未退出（协程泄漏）", after-before)
	}
}

// 静置至goroutine数量稳定（最多约15秒），避免慢机器上的偶发误差
func settleGoroutines(t *testing.T) {
	t.Helper()
	stable, prev := 0, runtime.NumGoroutine()
	for i := 0; i < 30; i++ {
		time.Sleep(500 * time.Millisecond)
		cur := runtime.NumGoroutine()
		if cur == prev {
			stable++
			if stable >= 3 { // 连续1.5秒无变化视为稳定
				return
			}
		} else {
			stable = 0
		}
		prev = cur
	}
}
