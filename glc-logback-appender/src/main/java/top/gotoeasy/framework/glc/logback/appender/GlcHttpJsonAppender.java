package top.gotoeasy.framework.glc.logback.appender;

import java.io.DataOutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.text.CharacterIterator;
import java.text.SimpleDateFormat;
import java.text.StringCharacterIterator;
import java.util.Date;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

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

    private int cnt = 0;
    private boolean enableGlc = true;

    private ExecutorService executor = Executors.newSingleThreadExecutor();

    @Override
    protected void append(ILoggingEvent event) {
        if (!enableGlc) {
            return; // 未启用时跳过
        }

        if (event == null || !isStarted()) {
            if (cnt++ < 10) {
                System.err.println("日志事件为空或该Appender未被初始化");
            }
            return;
        }

        // 异步发送日志到GLC
        executor.execute(() -> {
            submitToGlogCenter(layout.doLayout(event));
        });
    }

    /**
     * 发送日志到GLC<br>
     * 为不依赖第三方包，仅作java原生包简单实现，性能较差<br>
     * 实际使用时若有性能问题可继承重写实现
     * 
     * @param text 日志
     */
    protected void submitToGlogCenter(String text) {
        if (text == null) {
            return; // ignore
        }

        DataOutputStream dos = null;
        String body = null;
        try {
            body = "{" + encodeStr("text") + ":" + encodeStr(text.trim());
            body += "," + encodeStr("date") + ":" + encodeStr(getDateString());
            body += "," + encodeStr("system") + ":" + encodeStr(getSystem());
            body += "}";

            URL url = new URL(glcApiUrl);
            HttpURLConnection connection = (HttpURLConnection)url.openConnection();
            // 设置header
            if (headerKey != null && !"".equals(headerVal)) {
                connection.setRequestProperty(headerKey, headerVal);
            }
            connection.setConnectTimeout(5000);
            connection.setReadTimeout(5000);
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
        } catch (Exception e) {
            if (cnt++ < 10) {
                System.err.println("[GLC日志发送异常][地址：" + glcApiUrl + "][异常信息：" + e.getMessage() + "]");
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
            System.err.println("Layout未被初始化");
        }
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

    private String encodeStr(String str) {
        StringBuilder buf = new StringBuilder();
        buf.append('"');
        CharacterIterator it = new StringCharacterIterator(str);
        for (char c = it.first(); c != CharacterIterator.DONE; c = it.next()) {
            if (c == '"')
                buf.append("\\\"");
            else if (c == '\\')
                buf.append("\\\\");
            else if (c == '/')
                buf.append("\\/");
            else if (c == '\b')
                buf.append("\\b");
            else if (c == '\f')
                buf.append("\\f");
            else if (c == '\n')
                buf.append("\\n");
            else if (c == '\r')
                buf.append("\\r");
            else if (c == '\t')
                buf.append("\\t");
            else if (Character.isISOControl(c)) {
                addUnicode(buf, c);
            } else {
                buf.append(c);
            }
        }
        buf.append('"');
        return buf.toString();
    }

    static final char[] hex = "0123456789ABCDEF".toCharArray();

    private void addUnicode(StringBuilder buf, char c) {
        buf.append("\\u");
        int n = c;
        for (int i = 0; i < 4; ++i) {
            int digit = (n & 0xf000) >> 12;
            buf.append(hex[digit]);
            n <<= 4;
        }
    }

    private static String getDateString() {
        SimpleDateFormat sdf = getSimpleDateFormat("yyyy-MM-dd HH:mm:ss.SSS");
        return sdf.format(new Date());
    }

    private static ThreadLocal<SimpleDateFormat> threadLocal = new ThreadLocal<SimpleDateFormat>();
    private static Object lockObject = new Object();

    private static SimpleDateFormat getSimpleDateFormat(String format) {
        SimpleDateFormat simpleDateFormat = threadLocal.get();
        if (simpleDateFormat == null) {
            synchronized (lockObject) {
                if (simpleDateFormat == null) {
                    simpleDateFormat = new SimpleDateFormat(format);
                    simpleDateFormat.setLenient(false);
                    threadLocal.set(simpleDateFormat);
                }
            }
        }
        simpleDateFormat.applyPattern(format);
        return simpleDateFormat;
    }

}
