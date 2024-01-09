package top.gotoeasy.framework.glc.logback.appender;

import java.math.BigDecimal;
import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.NetworkInterface;
import java.text.CharacterIterator;
import java.text.SimpleDateFormat;
import java.text.StringCharacterIterator;
import java.util.ArrayList;
import java.util.Date;
import java.util.Enumeration;
import java.util.List;

public class Util {

    private static String serverIp = "";
    private static String serverName = "";

    private static void initServerIp() {

        List<InetAddress> list = new ArrayList<>();
        try {
            Enumeration<NetworkInterface> allNetInterfaces = NetworkInterface.getNetworkInterfaces();
            InetAddress ip;
            while (allNetInterfaces.hasMoreElements()) {
                NetworkInterface netInterface = allNetInterfaces.nextElement();

                if (!netInterface.isLoopback() && !netInterface.isVirtual() && netInterface.isUp()) {
                    Enumeration<InetAddress> addresses = netInterface.getInetAddresses();
                    while (addresses.hasMoreElements()) {
                        ip = addresses.nextElement();
                        if (ip instanceof Inet4Address) {
                            String addr = ip.getHostAddress();
                            if ("eth0".equals(netInterface.getName())) {
                                // 有eth0网卡ip时最优先，直接结束
                                serverIp = addr;
                                return;
                            }
                            list.add(ip);
                        }
                    }
                }
            }
        } catch (Exception e) {
            // ignore
        }

        // 192.* > 172.* > 10.* 
        for (int i = 0; i < list.size(); i++) {
            String ip = list.get(i).getHostAddress();
            if (ip.startsWith("192.")) {
                serverIp = ip;
            } else if (ip.startsWith("17.")) {
                if (!serverIp.startsWith("192.")) {
                    serverIp = ip;
                }
            } else if (ip.startsWith("10.")) {
                if (!serverIp.startsWith("192.") && !serverIp.startsWith("172.")) {
                    serverIp = ip;
                }
            } else {
                if ("".equals(serverIp)) {
                    serverIp = ip;
                }
            }
        }

    }

    public static String getServerName() {
        if ("".equals(serverName)) {
            try {
                serverName = InetAddress.getLocalHost().getHostName();
            } catch (Exception e) {
                // ignore
            }
        }
        return serverName;
    }

    public static String getServerIp() {
        if ("".equals(serverIp)) {
            initServerIp();
        }
        if ("".equals(serverIp)) {
            try {
                serverIp = InetAddress.getLocalHost().getHostAddress();
            } catch (Exception e) {
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

    public static String hash(String str) {
        int rs = 53653;
        int i = (str == null ? 0 : str.length());
        while (i > 0) {
            rs = (rs * 33) ^ str.charAt(--i);
        }
        return new BigDecimal(Long.valueOf(rs & 0x0FFFFFFFFL)).toPlainString();
    }

}
