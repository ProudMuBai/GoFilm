
const COOKIE_KEY_MAP = {
    'FILM_HISTORY':'filmHistory'
}

const cookieUtil =
    {
        /**
         * 设置cookie值
         * cname    string  cookie名称
         * value    any cookie值
         * expire   cookie保存天数
         */
        setCookie(name: string, value: any, expire = 30) {
            let d = new Date();
            d.setTime(d.getTime() + (expire * 24 * 60 * 60 * 1000));
            let expires = "expires=" + d.toUTCString();
            document.cookie = name + "=" + encodeURIComponent(value) + "; " + expires + ': path=/';
        },
        /**
         * 获取cookie值
         * name string  cookie名称
         */
        getCookie(name: string) {
            let cookies = document.cookie.split('; ');
            for (let i = 0; i < cookies.length; i++) {
                let [k,v] = cookies[i].split("=")
                if (k == name) {
                    return decodeURIComponent(v)
                }
            }
            return "";
        },
        /**
         *清除cookie值
         * name   string  cookie名称
         * 将expire 设置为已过期的时间
         */
        clearCookie(name: string) {
            let d = new Date();
            d.setTime(-1);
            let expires = "expires=" + d.toUTCString();
            document.cookie = name + "=''; " + expires;
        },
    }
export  {cookieUtil, COOKIE_KEY_MAP};
