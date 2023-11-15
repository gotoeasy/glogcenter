package top.gotoeasy.framework.glc.logback.appender;

import java.util.UUID;

import org.slf4j.MDC;

public class MdcUtil {

    public static final String TRACE_ID = "traceid";
    public static final String CLIENT_IP = "clientip";
    public static final String USER = "user";

    public static void setUser(String user) {
        MDC.put(USER, user);
    }

    public static void setClientIp(String clientip) {
        MDC.put(CLIENT_IP, clientip);
    }

    public static void setTraceId(String traceid) {
        MDC.put(TRACE_ID, traceid);
    }

    public static void removeUser() {
        MDC.remove(USER);
    }

    public static void removeClientIp() {
        MDC.remove(CLIENT_IP);
    }

    public static void removeTraceId() {
        MDC.remove(TRACE_ID);
    }

    public static void clear() {
        MDC.clear();
    }

    public static String generateTraceId() {
        return Util.hash(UUID.randomUUID().toString());
    }

}
