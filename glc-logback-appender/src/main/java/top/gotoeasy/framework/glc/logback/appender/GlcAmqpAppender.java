package top.gotoeasy.framework.glc.logback.appender;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;

import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.core.AppenderBase;
import ch.qos.logback.core.Layout;

/**
 * GLC是glogcenter缩写，一个golang实现的日志中心<br>
 * GlcAmqpAppender提供一种发送日志数据到RabbitMQ的方式（GLC则从RabbitMQ接收日志），适用于使用logback做日志管理的java项目
 */
public class GlcAmqpAppender extends AppenderBase<ILoggingEvent> {

    // 自定义配置，需Getter和Setter方法
    private String amqpHost;
    private int amqpPort;
    private String amqpUser;
    private String amqpPassword;
    private String system;
    private Layout<ILoggingEvent> layout;

    private int cnt = 0;
    private boolean enableGlc = true;

    private ExecutorService executor = Executors.newSingleThreadExecutor();

    protected Connection connection = null;
    protected Channel channel = null;

    protected synchronized void initConnectionChannel() throws Exception {
        if (this.channel != null) {
            return;
        }

        Connection conn = null;
        Channel chan = null;
        try {
            ConnectionFactory factory = new ConnectionFactory(); // 创建一个连接工厂
            factory.setHost(amqpHost); // 工厂ip 连接rabbitmq的队列
            factory.setPort(amqpPort); // 端口
            factory.setUsername(amqpUser); // 用户名
            factory.setPassword(amqpPassword); // 密码
            conn = factory.newConnection(); // 创建连接
            chan = conn.createChannel(); // 获取信道

            // 对列名称，
            // durable  是否持久化数据
            // exclusive 排他性，权限私有
            // autoDelete 是否自动删除
            // arguments
            chan.queueDeclare("glc-log-queue", false, false, false, null);
        } finally {
            this.connection = conn;
            this.channel = chan;
        }
    }

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

        // 异步发送日志
        executor.execute(() -> {
            sendToRabbitMQ(layout.doLayout(event), event);
        });
    }

    /**
     * 发送日志到RabbitMQ<br>
     * 
     * @param text 日志
     */
    protected void sendToRabbitMQ(String text, ILoggingEvent event) {
        if (text == null) {
            return; // ignore
        }

        try {
            if (channel == null) {
                initConnectionChannel();
            }

            String traceid = event.getMDCPropertyMap().get(MdcUtil.TRACE_ID);
            String clientip = event.getMDCPropertyMap().get(MdcUtil.CLIENT_IP);
            String user = event.getMDCPropertyMap().get(MdcUtil.USER);

            String body = "{\"text\":" + Util.encodeStr(text.trim());
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

            channel.basicPublish("", "glc-log-queue", null, body.getBytes("utf-8"));
        } catch (Exception e) {
            if (cnt++ < 10) {
                e.printStackTrace();
            }
            resetConnectionChannel();
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
        String host = System.getenv("GLC_AMQP_HOST");
        if (host != null) {
            host = host.trim();
            if (!"".equals(host)) {
                setAmqpHost(host);
            }
        }
        String port = System.getenv("GLC_AMQP_PORT");
        if (port != null) {
            port = port.trim();
            if (!"".equals(port)) {
                setAmqpPort(Integer.valueOf(port));
            }
        }
        String user = System.getenv("GLC_AMQP_USER");
        if (user != null) {
            user = user.trim();
            if (!"".equals(user)) {
                setAmqpUser(user);
            }
        }
        String password = System.getenv("GLC_AMQP_PASSWORD");
        if (password != null) {
            password = password.trim();
            if (!"".equals(password)) {
                setAmqpPassword(password);
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

    public String getAmqpHost() {
        return amqpHost;
    }

    public void setAmqpHost(String amqpHost) {
        this.amqpHost = amqpHost;
    }

    public int getAmqpPort() {
        return amqpPort;
    }

    public void setAmqpPort(int amqpPort) {
        this.amqpPort = amqpPort;
    }

    public String getAmqpUser() {
        return amqpUser;
    }

    public void setAmqpUser(String amqpUser) {
        this.amqpUser = amqpUser;
    }

    public String getAmqpPassword() {
        return amqpPassword;
    }

    public void setAmqpPassword(String amqpPassword) {
        this.amqpPassword = amqpPassword;
    }

    public Layout<ILoggingEvent> getLayout() {
        return layout;
    }

    public void setLayout(Layout<ILoggingEvent> layout) {
        this.layout = layout;
    }

    public void setSystem(String system) {
        this.system = system;
    }

    public String getSystem() {
        return system == null ? "" : system;
    }

    protected synchronized void resetConnectionChannel() {
        try {
            if (channel != null) {
                channel.close();
            }
        } catch (Exception ex) {
            // ignore
        } finally {
            this.channel = null;
        }
        try {
            if (connection != null) {
                connection.close();
            }
        } catch (Exception ex) {
            // ignore
        } finally {
            this.connection = null;
        }
    }

}
