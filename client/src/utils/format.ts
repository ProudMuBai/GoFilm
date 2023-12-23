

const fmt = {
    dateFormat( dateStamp:number, format:string='YYYY-mm-dd HH:MM:SS') {
        // 根据时间戳生成当前时间, 单位 毫秒
        let date = new Date(dateStamp*1000)
        const opt = {
            "Y+": date.getFullYear().toString(),        // 年
            "m+": (date.getMonth() + 1).toString(),     // 月
            "d+": date.getDate().toString(),            // 日
            "H+": date.getHours().toString(),           // 时
            "M+": date.getMinutes().toString(),         // 分
            "S+": date.getSeconds().toString()          // 秒
            // 有其他格式化字符需求可以继续添加，必须转化成字符串
        };
        for (let k in opt) {
            // 正则匹配对应的 opt key
            let r = new RegExp("(" + k + ")").exec(format);
            // 如果有匹配成功项
            if (r) {
                format = format.replace(r[1], (r[1].length == 1) ? (opt[k as keyof typeof  opt]) : (opt[k as keyof typeof  opt].padStart(r[1].length, "0")))
            }
        }
        return format;
    }
}

export {fmt}