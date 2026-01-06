package top.gotoeasy.framework.glc.logback.appender;

import java.io.DataOutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.core.AppenderBase;
import ch.qos.logback.core.Layout;

/**
 * GLC是glogcenter缩写，一个golang实现的日志中心<br>
 * GlcHttpJsonAppender提供一种http提交Json数据的方式发送日志到GLC，适用于使用logback做日志管理的java项目<br>
 * 仅简易实现为主，若是性能要求高日志量大的场景，应选用其他比如消息队列之类的Appender
 */
public class GlcHttpJsonAppender extends AppenderBase<ILoggingEvent> {

    // 自定义配置，需Getter和Setter方法
    private String glcApiUrl;
    private String glcApiKey;
    private String system = "default";
    private Layout<ILoggingEvent> layout;

    private String headerKey;
    private String headerVal;


    private boolean enableGlc = true;
    private final AtomicInteger errorCount = new AtomicInteger(0); // 连续错误计数

    // 熔断控制
    private volatile long lastFailureTime = 0;
    private static final long PAUSE_PERIOD_MS = 60 * 1000; // 失败后暂停发送60秒
    private static final int QUEUE_SIZE = 5000; // 队列最大缓存日志数
    private static final int CONNECT_TIMEOUT = 2000; // 连接超时缩短为2秒
    private static final int READ_TIMEOUT = 3000;    // 读取超时3秒

    // 线程池
    private ExecutorService executor;


    @Override
    protected void append(ILoggingEvent event) {
        if (!enableGlc) {
            return; // 未启用时跳过
        }

        if (event == null || !isStarted()) {
            return;
        }

        // 1. 熔断检查：如果最近发生过错误，且在冷却时间内，直接丢弃日志，不进入线程池
        if (lastFailureTime > 0) {
            if (System.currentTimeMillis() - lastFailureTime < PAUSE_PERIOD_MS) {
                return; // 处于熔断期，跳过发送
            } else {
                // 过了冷却期，重置状态尝试发送
                lastFailureTime = 0;
            }
        }

        try {
            // 2. 提交任务
            executor.execute(() -> submitToGlogCenter(layout.doLayout(event), event));
        } catch (Exception e) {
            // 线程池满或其他异常，直接忽略，不影响业务
        }
    }
    /**
     * 发送日志到GLC<br>
     * 为不依赖第三方包，仅作java原生包简单实现，性能较差<br>
     * 实际使用时若有性能问题可继承重写实现
     *
     * @param text  日志
     * @param event ILoggingEvent
     */
    protected void submitToGlogCenter(String text, ILoggingEvent event) {
        if (text == null) {
            return; // ignore
        }

        DataOutputStream dos = null;
        String body = null;
        try {
            String traceid = event.getMDCPropertyMap().get(MdcUtil.TRACE_ID);
            String clientip = event.getMDCPropertyMap().get(MdcUtil.CLIENT_IP);
            String user = event.getMDCPropertyMap().get(MdcUtil.USER);

            body = "{\"text\":" + Util.encodeStr(text.trim());
            body += ",\"date\":\"" + Util.getDateString() + "\"";
            body += ",\"system\":" + Util.encodeStr(getSystem());
            body += ",\"servername\":" + Util.encodeStr(Util.getServerName());
            body += ",\"serverip\":" + Util.encodeStr(Util.getServerIp());
            body += ",\"loglevel\":\"" + event.getLevel().toString() + "\"";
            if (traceid != null && !"".equals(traceid)) {
                body += ",\"traceid\":" + Util.encodeStr(traceid);
            }
            if (clientip != null && !"".equals(clientip)) {
                body += ",\"clientip\":" + Util.encodeStr(clientip);
            }
            if (user != null && !"".equals(user)) {
                body += ",\"user\":" + Util.encodeStr(user);
            }
            body += "}";

            URL url = new URL(glcApiUrl);
            HttpURLConnection connection = (HttpURLConnection)url.openConnection();
            // 设置header
            if (headerKey != null && !"".equals(headerVal)) {
                connection.setRequestProperty(headerKey, headerVal);
            }
            connection.setConnectTimeout(CONNECT_TIMEOUT);
            connection.setReadTimeout(READ_TIMEOUT);
            connection.setDoInput(true);
            connection.setDoOutput(true);
            connection.setUseCaches(false);
            connection.setRequestMethod("POST");
            connection.setRequestProperty("Content-Type", "application/json");
            // 发送日志数据
            connection.connect();
            dos = new DataOutputStream(connection.getOutputStream());
            dos.write(body.getBytes("utf-8"));
            dos.flush();
            // 接收响应内筒
            connection.getContent();
            connection.disconnect();
            // 发送成功，重置错误计数
            errorCount.set(0);
        } catch (Exception e) {
            // 记录失败时间，触发熔断
            lastFailureTime = System.currentTimeMillis();
            // 限制控制台报错频率：仅在连续错误的前3次打印，避免刷屏
            if (errorCount.incrementAndGet() <= 3) {
                System.err.println("[GLC日志发送异常][地址：" + glcApiUrl + "]发送日志失败(暂停发送" + (PAUSE_PERIOD_MS/1000) + "秒): " + e.getMessage());
            }
        } finally {
            try {
                if (dos != null) {
                    dos.close();
                }
            } catch (Exception e) {
                // ignore
            }
        }
    }

    @Override
    public void start() {
        if (this.layout == null) {
            addError("Layout未被初始化");
            return;
        }
        // 初始化线程池：单线程，有界队列，满载丢弃旧日志
        // DiscardOldestPolicy: 当队列满时，抛弃最老的日志，保证新日志能尝试进入
        // 也可以使用 DiscardPolicy 直接抛弃当前日志
        executor = new ThreadPoolExecutor(
                1, 1,
                0L, TimeUnit.MILLISECONDS,
                new ArrayBlockingQueue<>(QUEUE_SIZE),
                r -> {
                    Thread t = new Thread(r, "GLC-Appender-Thread");
                    t.setDaemon(true); // 设置为守护线程，防止阻碍JVM关闭
                    return t;
                },
                new ThreadPoolExecutor.DiscardPolicy() // 队列满时直接丢弃，不抛异常，保护主业务
        );

        super.start();

        // 优先使用环境变量设定
        String enable = System.getenv("GLC_ENABLE");
        if ("false".equalsIgnoreCase(enable) || "0".equals(enable)) {
            enableGlc = false;
        }
        String apiUrl = System.getenv("GLC_API_URL");
        if (apiUrl != null) {
            apiUrl = apiUrl.trim();
            if (!"".equals(apiUrl)) {
                if (!apiUrl.endsWith("/glc/v1/log/add")) {
                    // 允许省略接口路径，默认自动补足以简化使用
                    if (apiUrl.endsWith("/")) {
                        apiUrl += "glc/v1/log/add";
                    } else {
                        apiUrl += "/glc/v1/log/add";
                    }
                }
                setGlcApiUrl(apiUrl);
            }
        }

        String apiKey = System.getenv("GLC_API_KEY");
        if (apiKey != null) {
            apiKey = apiKey.trim();
            if (!"".equals(apiKey)) {
                setGlcApiKey(apiKey);
            }
        }

        String system = System.getenv("GLC_SYSTEM");
        if (system != null) {
            system = system.trim();
            if (!"".equals(system)) {
                setSystem(system);
            }
        }
    }

    @Override
    public void stop() {
        if (!isStarted()) {
            return;
        }
        // 关闭线程池
        if (executor != null) {
            executor.shutdownNow();
        }
        super.stop();
    }

    public Layout<ILoggingEvent> getLayout() {
        return layout;
    }

    public void setLayout(Layout<ILoggingEvent> layout) {
        this.layout = layout;
    }

    public String GetGlcApiUrl() {
        return glcApiUrl;
    }

    public void setGlcApiUrl(String glcApiUrl) {
        this.glcApiUrl = glcApiUrl;
    }

    public String GetGlcApiKey() {
        return glcApiKey;
    }

    public void setSystem(String system) {
        this.system = system;
    }

    public String getSystem() {
        return system == null ? "" : system;
    }

    public void setGlcApiKey(String glcApiKey) {
        this.glcApiKey = glcApiKey;

        String key = glcApiKey.split(":")[0];
        String value = glcApiKey.substring(key.length() + 1).trim();
        key = key.trim();
        if (!"".equals(key)) {
            headerKey = key;
            headerVal = value;
        }
    }

}
