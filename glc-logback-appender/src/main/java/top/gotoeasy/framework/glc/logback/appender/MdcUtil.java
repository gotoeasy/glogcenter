package top.gotoeasy.framework.glc.logback.appender;

import java.util.UUID;

import org.slf4j.MDC;

public class MdcUtil {

    public static final String TRACE_ID = "traceid";
    public static final String CLIENT_IP = "clientip";
    public static final String USER = "user";

    public static void setUser(String user) {
        try {
            MDC.put(USER, user);
        } catch (Exception e) {
            // ignore
        }
    }

    public static String getUser() {
        try {
            return MDC.get(USER);
        } catch (Exception e) {
            return "";
        }
    }

    public static void setClientIp(String clientip) {
        try {
            MDC.put(CLIENT_IP, clientip);
        } catch (Exception e) {
            // ignore
        }
    }

    public static String getClientIp() {
        try {
            return MDC.get(CLIENT_IP);
        } catch (Exception e) {
            return "";
        }
    }

    public static void setTraceId(String traceid) {
        try {
            MDC.put(TRACE_ID, traceid);
        } catch (Exception e) {
            // ignore
        }
    }

    public static String getTraceId() {
        try {
            return MDC.get(TRACE_ID);
        } catch (Exception e) {
            return "";
        }
    }

    public static void removeUser() {
        try {
            MDC.remove(USER);
        } catch (Exception e) {
            // ignore
        }
    }

    public static void removeClientIp() {
        try {
            MDC.remove(CLIENT_IP);
        } catch (Exception e) {
            // ignore
        }
    }

    public static void removeTraceId() {
        try {
            MDC.remove(TRACE_ID);
        } catch (Exception e) {
            // ignore
        }
    }

    public static void clear() {
        MDC.clear();
    }

    public static String generateTraceId() {
        return Util.hash(UUID.randomUUID().toString());
    }

}
