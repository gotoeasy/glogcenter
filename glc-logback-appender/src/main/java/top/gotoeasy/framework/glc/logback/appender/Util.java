package top.gotoeasy.framework.glc.logback.appender;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.text.CharacterIterator;
import java.text.SimpleDateFormat;
import java.text.StringCharacterIterator;
import java.util.Date;

public class Util {

    private static String serverIp = "";
    private static String serverName = "";

    public static String getServerName() {
        if ("".equals(serverName)) {
            try {
                serverName = InetAddress.getLocalHost().getHostName();
            } catch (UnknownHostException e) {
                // ignore
            }
        }
        return serverName;
    }

    public static String getServerIp() {
        if ("".equals(serverIp)) {
            try {
                serverIp = InetAddress.getLocalHost().getHostAddress();
            } catch (UnknownHostException e) {
                // ignore
            }
        }
        return serverIp;
    }

    public static String getDateString() {
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

    public static String encodeStr(String str) {
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

    private static final char[] hex = "0123456789ABCDEF".toCharArray();

    private static void addUnicode(StringBuilder buf, char c) {
        buf.append("\\u");
        int n = c;
        for (int i = 0; i < 4; ++i) {
            int digit = (n & 0xf000) >> 12;
            buf.append(hex[digit]);
            n <<= 4;
        }
    }

}
