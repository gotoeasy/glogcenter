package top.gotoeasy.framework.glc.filter;

import java.io.IOException;

import javax.servlet.Filter;
import javax.servlet.FilterChain;
import javax.servlet.FilterConfig;
import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletRequest;

import top.gotoeasy.framework.glc.logback.appender.MdcUtil;

/**
 * 过滤器，自动设定GLC所需的traceid、clientip，优先从请求头中获取
 */
public class GlcFilter implements Filter {

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
    }

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {

        // 设定日志中心相关的的traceid、clientip，若需要设定user，可继承setMdcKeyValues添加相关逻辑
        setMdcKeyValues(request);

        chain.doFilter(request, response);
    }

    @Override
    public void destroy() {
    }

    protected void setMdcKeyValues(ServletRequest request) {
        // 设定日志中心相关的的traceid、clientip
        HttpServletRequest httpServletRequest = (HttpServletRequest)request;
        String traceid = httpServletRequest.getHeader(MdcUtil.TRACE_ID);
        MdcUtil.setTraceId((traceid == null || traceid.length() == 0) ? MdcUtil.generateTraceId() : traceid);
        MdcUtil.setClientIp(getIpAddr(httpServletRequest));
    }

    protected String getIpAddr(HttpServletRequest request) {
        String[] headerNames = { "X-Forwarded-For", "X-Real-IP", "Proxy-Client-IP", "WL-Proxy-Client-IP",
                "HTTP_CLIENT_IP", "HTTP_X_FORWARDED_FOR" };
        for (String headerName : headerNames) {
            String ip = request.getHeader(headerName);
            if (ip != null && ip.length() > 0 && !"unknown".equalsIgnoreCase(ip)) {
                int index = ip.indexOf(',');
                if (index > 0) {
                    ip = ip.substring(0, index);
                }
                return ip;
            }
        }
        return request.getRemoteAddr();
    }

}
