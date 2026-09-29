/**
 * 回归测试：集群转发必须发往调用方指定的uri（ISSUE-004 / ISSUE-021），
 * 且登录功能关闭时 transferLogin 不得 panic（ISSUE-012）
 * 背景：TransferGlc无视uri参数、硬编码转发到/v1/log/transferAdd，集群模式下
 * 用户保存/改密/删除/登录会话的转发全部被静默丢弃；一旦转发修复激活，
 * 登录关闭的节点被直接或转发命中transferLogin时会因nil会话缓存panic
 *
 * 执行顺序约定（Go按源码顺序执行同包测试）：
 *   1) 守卫测试最前——只依赖init默认态（登录关闭），不写conf、不开存储
 *   2) 删用户测试第二——唯一一次UpdateConfigByEnv发生在任何存储打开之前
 *     （conf全局变量无锁，存储打开后其autoClose协程会并发读取conf，
 *      之后再写conf会构成数据竞争）
 *   3) uri转发测试最后——复用前序测试的conf与已打开的存储，不再写conf
 *
 * 注：本包测试读写固定目录/glogcenter下的leveldb，与ldb包的测试二进制经
 * go test ./...并行运行时会争用.sysmnt文件锁——本包首次打开以recover重试
 * 等待锁释放（有界）；全量测试建议 -p 1 串行各包二进制更稳
 */
package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"glc/conf"
	"glc/gweb"
	"glc/www/cluster"
	"glc/www/service"

	"github.com/gin-gonic/gin"
	"github.com/gotoeasy/glang/cmn"
)

// newRecordPeer 启动记录收到请求路径的伪集群节点
func newRecordPeer(t *testing.T) (*httptest.Server, *sync.Mutex, *[]string) {
	t.Helper()
	var mu sync.Mutex
	var paths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	return srv, &mu, &paths
}

// eventuallySetSysmntItem 等待并写入集群KV：go test ./...并行运行时ldb包的
// 测试二进制会先持有/glogcenter/.sysmnt文件锁数十秒，本包首次打开会因锁
// 报错panic——此处recover重试直至锁释放（有界90秒）
func eventuallySetSysmntItem(t *testing.T, kv *service.KeyValue) {
	t.Helper()
	deadline := time.Now().Add(90 * time.Second)
	for {
		done := func() (ok bool) {
			defer func() { _ = recover() }() // 锁被并行测试二进制持有时panic：重试
			_, err := service.SetSysmntItem(kv)
			return err == nil
		}()
		if done {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("等待/glogcenter/.sysmnt锁释放超时（是否有测试二进制长期持锁？）")
		}
		time.Sleep(time.Second)
	}
}

// eventually 轮询直到条件满足或超时
func eventually(t *testing.T, wait time.Duration, fn func() bool) bool {
	t.Helper()
	deadline := time.Now().Add(wait)
	for {
		if fn() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// 阶段1：登录功能关闭时transferLogin直接忽略，不得panic（ISSUE-012）
// 只依赖init默认态（登录关闭→会话缓存nil），不写conf、不开存储；
// 红灯运行时本测试会因nil会话缓存panic，需单独 -run 执行以取证
func TestUserTransferLoginDisabledNoPanic(t *testing.T) {
	t.Setenv("GLC_ENABLE_LOGIN", "false") // 默认即关闭，显式声明；不调UpdateConfigByEnv以免写conf全局
	if catchSession != nil {              // 还原真实默认态：登录关闭时会话缓存为nil
		catchSession = nil
	}
	t.Cleanup(func() { catchSession = nil })
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/glc/v1/user/transferLogin",
		strings.NewReader(`{"username":"glc","password":"GLogCenter100%666"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	rs := UserTransferLoginController(gweb.NewHttpRequest(c)) // 红灯：此处panic
	if rs == nil || !rs.Success {
		t.Fatalf("登录功能关闭时transferLogin应直接忽略，实际: %+v", rs)
	}
}

// 阶段2：删除用户必须转发到删除接口（ISSUE-004，需与ISSUE-021一同生效）
// 本测试的UpdateConfigByEnv是全套测试中唯一一次conf写入，且发生在任何存储打开之前
func TestUserDelTransfersToDeleteEndpoint(t *testing.T) {
	srv, mu, paths := newRecordPeer(t)
	t.Setenv("GLC_SERVER_URL", "http://127.0.0.1:1")
	t.Setenv("GLC_ENABLE_LOGIN", "true")  // UserDelController要求管理员会话
	t.Setenv("GLC_CLUSTER_MODE", "true") // 转发仅在集群模式下触发
	conf.UpdateConfigByEnv()
	if catchSession == nil { // 包init在登录关闭场景不建会话缓存，测试临时补建
		catchSession = cmn.NewCache(time.Hour)
	}
	t.Cleanup(func() { catchSession = nil }) // 还原，避免影响后续测试

	eventuallySetSysmntItem(t, &service.KeyValue{
		Key:   cluster.KEY_CLUSTER,
		Value: (&cluster.ClusterInfo{MasterUrl: "http://127.0.0.1:1", NodeUrls: srv.URL}).ToJson(),
	})

	const adminToken = "test-admin-token"
	catchSession.Set(adminToken, conf.GetUsername())

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/glc/v1/sysuser/del", strings.NewReader(`{"username":"u1"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("X-Access-Token", adminToken) // 控制器从header取token校验管理员会话
	rs := UserDelController(gweb.NewHttpRequest(c))
	if rs == nil || !rs.Success {
		t.Fatalf("期望删除用户成功，实际: %+v", rs)
	}

	want := conf.GetContextPath() + conf.SysUserTransferDel
	if !eventually(t, 3*time.Second, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(*paths) == 1 && (*paths)[0] == want
	}) {
		mu.Lock()
		defer mu.Unlock()
		t.Fatalf("删除用户应转发到 %s，实际收到 %v（uri或常量有误）", want, *paths)
	}
}

// 阶段3：TransferGlc必须把载荷发往调用方传入的uri（ISSUE-021）
// 复用阶段2写入的conf与已打开的存储，不再写conf（避免与autoClose协程竞争）
func TestTransferGlcForwardsToGivenUri(t *testing.T) {
	srv, mu, paths := newRecordPeer(t)

	eventuallySetSysmntItem(t, &service.KeyValue{
		Key:   cluster.KEY_CLUSTER,
		Value: (&cluster.ClusterInfo{MasterUrl: "http://127.0.0.1:1", NodeUrls: srv.URL}).ToJson(),
	})

	TransferGlc(conf.SysUserTransferDel, `{"username":"u1"}`)

	want := conf.GetContextPath() + conf.SysUserTransferDel
	mu.Lock()
	defer mu.Unlock()
	if len(*paths) != 1 {
		t.Fatalf("期望恰好1条转发请求，实际 %v", *paths)
	}
	if (*paths)[0] != want {
		t.Errorf("TransferGlc被传入uri %s 但实际把载荷发到了 %s（uri参数被无视）", want, (*paths)[0])
	}
}
